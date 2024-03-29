package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// type WorkflowProvider interface {
// 	GetWorkflowByID(ctx context.Context, id int) (*Workflow, error)
// 	GetWorkflows(ctx context.Context, filter *WorkflowFilter) ([]Workflow, error)
// 	SaveWorkflow(ctx context.Context, w *Workflow) (*Workflow, error)
// 	UpdateWorkflow(ctx context.Context, w *Workflow) (*Workflow, error)
// }

type Workflow struct {
	ID                 int          `json:"id"`
	Name               string       `json:"name"`
	Type               string       `json:"type"`
	Duration           int          `json:"duration"`
	ReceivedAt         time.Time    `json:"received_at"`
	StartedExecutionAt sql.NullTime `json:"started_execution_at,omitempty"`
	CompletedAt        sql.NullTime `json:"completed_at,omitempty"`
}

type WorkflowFilter struct {
	Type      string
	StartTime time.Time
	EndTime   time.Time
}

func (s *MySQLStore) GetWorkflowByID(ctx context.Context, id int) (*Workflow, error) {
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

	w.ID = int(id)
	return w, nil
}

func (s *MySQLStore) UpdateWorkflow(ctx context.Context, w *Workflow) (*Workflow, error) {
	_, err := s.db.ExecContext(ctx,
		"UPDATE workflow SET name = ?, type = ?, duration = ?, received_at = ?, started_execution_at = ?, completed_at = ? WHERE id = ?",
		w.Name, w.Type, w.Duration, w.ReceivedAt, w.StartedExecutionAt, w.CompletedAt, w.ID,
	)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (s *MySQLStore) StartWorkflow(ctx context.Context, workflowID int) error {
	// Check if started_execution_at is already set
	var startedAt sql.NullTime
	checkQuery := "SELECT started_execution_at FROM workflow WHERE id = ?"
	err := s.db.QueryRowContext(ctx, checkQuery, workflowID).Scan(&startedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no workflow found with ID: %d", workflowID)
		}
		return err
	}

	if startedAt.Valid {
		// Log that the value for started_execution_at already exists
		log.Printf("The workflow with ID %d already has a started_execution_at value of %s", workflowID, startedAt.Time)
		return nil // Exit without making any update
	}

	// If started_execution_at is not set, then update it
	query := "UPDATE workflow SET started_execution_at = NOW() WHERE id = ?"
	_, err = s.db.ExecContext(ctx, query, workflowID)
	return err
}

func (s *MySQLStore) CompleteWorkflow(ctx context.Context, id int) error {
	query := "UPDATE workflow SET completed_at = NOW() WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}
