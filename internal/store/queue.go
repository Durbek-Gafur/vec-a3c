package store

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Queue struct {
	ID         int     `json:"id"`
	WorkflowID int     `json:"workflow_id"`
	Status     string    `json:"status"`
	EnqueuedAt time.Time `json:"enqueued_at"`
}

func (s *MySQLStore) Enqueue(ctx context.Context, workflowID int) (int, error) {
	query := "INSERT INTO queue (workflow_id, status) VALUES (?, 'pending')"
	res, err := s.db.ExecContext(ctx, query, workflowID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *MySQLStore) Dequeue(ctx context.Context) (*Queue, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Select the next 'pending' or 'processing' job
	query := "SELECT id, workflow_id FROM queue WHERE status <> 'done' ORDER BY enqueued_at ASC LIMIT 1 FOR UPDATE"
	row := tx.QueryRowContext(ctx, query)

	q := &Queue{}
	err = row.Scan(&q.ID, &q.WorkflowID)
	if err != nil {
		return nil, err
	}

	// Update its status to 'processing'
	updateQuery := "UPDATE queue SET status = 'processing' WHERE id = ?"
	_, err = tx.ExecContext(ctx, updateQuery, q.ID)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return q, nil
}

func (s *MySQLStore) Peek(ctx context.Context) (*Queue, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Select the next 'pending' or 'processing' job
	query := "SELECT id, workflow_id, status, enqueued_at FROM queue WHERE status <> 'done' ORDER BY enqueued_at ASC LIMIT 1"
	row := tx.QueryRowContext(ctx, query)

	q := &Queue{}
	err = row.Scan(&q.ID, &q.WorkflowID, &q.Status, &q.EnqueuedAt)
	if err != nil {
		return nil, err
	}

	return q, nil
}


func (s *MySQLStore) GetQueueStatus(ctx context.Context) ([]Queue, error) {
	// Fetch the queue size limit from the queue_size table
	queueSize, err := s.GetQueueSizeFromDBorENV(ctx)
	if err != nil {
		return nil, err
	}

	query := "SELECT id, workflow_id, status, enqueued_at FROM queue WHERE status <> 'done' ORDER BY enqueued_at ASC LIMIT ?;"
	rows, err := s.db.QueryContext(ctx, query, queueSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var qs []Queue
	for rows.Next() {
		var q Queue
		err = rows.Scan(&q.ID, &q.WorkflowID, &q.Status, &q.EnqueuedAt)
		if err != nil {
			return nil, err
		}
		qs = append(qs, q)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return qs, nil
}
func (s *MySQLStore) getQueueSizeFromDB(ctx context.Context) (int, error) {
	
	queueSizeQuery := "SELECT size FROM queue_size LIMIT 1;"
	var queueSize int
	err := s.db.QueryRowContext(ctx, queueSizeQuery).Scan(&queueSize)
	if err != nil {
		return 0, err
	}
	// fmt.Printf("queue_size %d from DB\n", queueSize)
	return queueSize, nil
}

func (s *MySQLStore) getQueueSizeFromEnv() (int, error) {
	queueSizeEnv := os.Getenv("QUEUE_SIZE")
	if queueSizeEnv == "" {
		fmt.Println("QUEUE_SIZE not set, using default value 10")
		return 10, nil //errors.New("QUEUE_SIZE not set, using default value 10")
	}
	queueSize, err := strconv.Atoi(queueSizeEnv)
	if err != nil {
		return 0, fmt.Errorf("invalid value for QUEUE_SIZE: %s", queueSizeEnv)
	}
	// fmt.Printf("QUEUE_SIZE from ENV %d", queueSize)
	return queueSize, nil
}

func (s *MySQLStore) GetQueueSizeFromDBorENV(ctx context.Context) (int, error) {
	queueSize, err := s.getQueueSizeFromDB(ctx)
	if err == nil {
		return queueSize, nil
	}

	queueSize, err = s.getQueueSizeFromEnv()
	if err != nil {
		return 0, err
	}

	return queueSize, nil
}

func (s *MySQLStore) IsSpaceAvailable(ctx context.Context) (bool, error) {
	queueSize, err := s.GetQueueSizeFromDBorENV(ctx)
	if err != nil {
		return false, err
	}

	query := "SELECT COUNT(*) FROM queue WHERE status <> 'done';"
	var currentSize int
	err = s.db.QueryRowContext(ctx, query).Scan(&currentSize)
	if err != nil {
		return false, err
	}

	availableSpace := currentSize < queueSize

	return availableSpace, nil
}

func (s *MySQLStore) GetAvailableSpace(ctx context.Context) (int, error) {
    queueSize, err := s.GetQueueSizeFromDBorENV(ctx)
    if err != nil {
        return 0, fmt.Errorf("failed to get queue size: %w", err)
    }

    query := "SELECT COUNT(*) FROM queue WHERE status <> 'done';"
    stmt, err := s.db.PrepareContext(ctx, query)
    if err != nil {
        return 0, fmt.Errorf("failed to prepare query: %w", err)
    }
    defer stmt.Close()

    var currentSize int
    err = stmt.QueryRowContext(ctx).Scan(&currentSize)
    if err != nil {
        return 0, fmt.Errorf("failed to query queue size: %w", err)
    }

    availableSpace := queueSize - currentSize
    if availableSpace < 0 {
        return 0, fmt.Errorf("negative available space: queueSize=%d, currentSize=%d", queueSize, currentSize)
    }

    return availableSpace, nil
}


func (s *MySQLStore) updateStatus(ctx context.Context, id int, status string) error {
	query := "UPDATE queue SET status = ? WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, status, id)
	return err
}

func (s *MySQLStore) ProcessWorkflowInQueue(ctx context.Context, id int) error {
	return s.updateStatus(ctx,id,"processing")
}

func (s *MySQLStore) CompleteWorkflowInQueue(ctx context.Context, id int) error {
	return s.updateStatus(ctx,id,"done")
}