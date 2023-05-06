package store

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLStore struct {
	db *sql.DB
}

func NewMySQLStore(dsn string) (*MySQLStore, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &MySQLStore{db: db}, nil
}

func (s *MySQLStore) GetQueueSize(ctx context.Context) (int, error) {
	var size int
	err := s.db.QueryRowContext(ctx, "SELECT size FROM queue_size WHERE id = 1").Scan(&size)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, &StoreError{Err: ErrNotFound, StatusCode: http.StatusNotFound}
		}
		return 0, err
	}

	return size, nil
}

func (s *MySQLStore) SetQueueSize(ctx context.Context, size int) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO queue_size (id, size) VALUES (1, ?) ON DUPLICATE KEY UPDATE size = ?", size, size)
	if err != nil {
		return err
	}
	return nil
}

func (s *MySQLStore) UpdateQueueSize(ctx context.Context, size int) error {
	_, err := s.db.ExecContext(ctx, "UPDATE queue_size SET size = ? WHERE id = 1", size)
	if err != nil {
		return err
	}
	return nil
}

func (s *MySQLStore) GetWorkflowByID(ctx context.Context, id int64) (*Workflow, error) {
	wf := &Workflow{}
	err := s.db.QueryRowContext(ctx,
		"SELECT id, name, type, duration, received_at, started_execution_at, completed_at FROM workflow WHERE id = ?",
		id,
	).Scan(&wf.ID, &wf.Name, &wf.Type, &wf.Duration, &wf.ReceivedAt, &wf.StartedExecutionAt, &wf.CompletedAt)

	if err != nil {
		return nil, err
	}

	return wf, nil
}

func (s *MySQLStore) GetWorkflows(ctx context.Context, filter *WorkflowFilter) ([]Workflow, error) {
	query := "SELECT id, name, type, duration, received_at, started_execution_at, completed_at FROM workflow WHERE 1"

	if filter.Type != "" {
		query += " AND type = ?"
	}
	if !filter.StartTime.IsZero() {
		query += " AND received_at >= ?"
	}
	if !filter.EndTime.IsZero() {
		query += " AND received_at <= ?"
	}

	rows, err := s.db.QueryContext(ctx, query, filter.Type, filter.StartTime, filter.EndTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workflows := []Workflow{}
	for rows.Next() {
		var wf Workflow
		err := rows.Scan(&wf.ID, &wf.Name, &wf.Type, &wf.Duration, &wf.ReceivedAt, &wf.StartedExecutionAt, &wf.CompletedAt)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, wf)
	}

	return workflows, nil
}

func (s *MySQLStore) SaveWorkflow(ctx context.Context, w *Workflow) (*Workflow, error) {
	res, err := s.db.ExecContext(ctx,
		"INSERT INTO workflow (name, type, duration, received_at) VALUES (?, ?, ?, ?)",
		w.Name, w.Type, w.Duration, w.ReceivedAt,
	)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	w.ID = id
	return w, nil
}

func (s *MySQLStore) UpdateWorkflow(ctx context.Context, w *Workflow) (*Workflow, error) {
	_, err := s.db.ExecContext(ctx,
		"UPDATE workflow SET name = ?, type = ?, duration = ?, received_at = ?, started_execution_at = ?, completed_at = ? WHERE id = ?",
		w.Name, w.Type, w.Duration, w.ReceivedAt, w.StartedExecutionAt, w.CompletedAt, w.ID,
	)
	if err != nil {
		return nil,err
	}

	return w,nil
}


func (s *MySQLStore) StartWorkflow(ctx context.Context, id int) error {
	query := "UPDATE workflow SET started_execution_at = NOW() WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

func (s *MySQLStore) CompleteWorkflow(ctx context.Context, id int) error {
	query := "UPDATE workflow SET completed_at = NOW() WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}
