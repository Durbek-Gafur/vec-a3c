package store

import (
	"context"
)


func (s *MySQLStore) Enqueue(ctx context.Context, workflowID int64) (int64, error) {
	query := "INSERT INTO queue (workflow_id, status) VALUES (?, 'pending')"
	res, err := s.db.ExecContext(ctx, query, workflowID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *MySQLStore) Dequeue(ctx context.Context) (*Queue, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Select the next 'pending' job
	query := "SELECT id, workflow_id FROM queue WHERE status = 'pending' ORDER BY enqueued_at ASC LIMIT 1 FOR UPDATE"
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

func (s *MySQLStore) GetQueueStatus(ctx context.Context) ([]Queue, error) {
	query := "SELECT id, workflow_id, status, enqueued_at FROM queue"
	rows, err := s.db.QueryContext(ctx, query)
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

func (s *MySQLStore) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := "UPDATE queue SET status = ? WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, status, id)
	return err
}
