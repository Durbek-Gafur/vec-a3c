package store

import (
	"context"
	"database/sql"
	"time"
)

// WorkflowInfo represents the information of a WorkflowInfo
type WorkflowInfo struct {
	ID                    int            `json:"id"`
	CreatedAt             sql.NullTime   `json:"createdAt"`
	Name                  string         `json:"name"`
	RAM                   string         `json:"ram"`
	Core                  string         `json:"core"`
	Policy                string         `json:"policy"`
	ExpectedExecutionTime sql.NullString `json:"expectedExecutionTime"`
	ActualExecutionTime   sql.NullString `json:"actualExecutionTime"`
	AssignedVM            sql.NullString `json:"assignedVM"`
	ProcessingStartedAt   sql.NullTime   `json:"processingStartedAt"`
	AssignedAt            sql.NullTime   `json:"assignedAt"`
	CompletedAt           sql.NullTime   `json:"completedAt"`
	Type                  string         `json:"type"`
	Status                string         `json:"status"`
	SubmittedBy           sql.NullString `json:"submittedBy"`
	LastUpdated           time.Time      `json:"lastUpdated"`
}

func (s *MySQLStore) GetWorkflowByID(ctx context.Context, id int) (*WorkflowInfo, error) {
	wf := &WorkflowInfo{}
	err := s.db.QueryRowContext(ctx,
		"SELECT id,created_at,name, type, ram, core, policy, expected_execution_time, actual_execution_time, assigned_vm, assigned_at, completed_at, submitted_by, status, last_updated FROM workflow_info WHERE id = ?",
		id,
	).Scan(&wf.ID, &wf.CreatedAt, &wf.Name, &wf.Type, &wf.RAM, &wf.Core, &wf.Policy, &wf.ExpectedExecutionTime, &wf.ActualExecutionTime, &wf.AssignedVM, &wf.AssignedAt, &wf.CompletedAt, &wf.SubmittedBy, &wf.Status, &wf.LastUpdated)

	if err != nil {
		return nil, err
	}

	return wf, nil
}

func (s *MySQLStore) GetWorkflows(ctx context.Context) ([]WorkflowInfo, error) {
	query := `SELECT id,created_at, name, type, ram, core, policy, 
                      expected_execution_time, actual_execution_time, 
                      assigned_vm, processing_started_at, assigned_at, 
                      completed_at, status, submitted_by, last_updated 
               FROM workflow_info`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []WorkflowInfo
	for rows.Next() {
		var wf WorkflowInfo
		err := rows.Scan(
			&wf.ID,
			&wf.CreatedAt,
			&wf.Name,
			&wf.Type,
			&wf.RAM,
			&wf.Core,
			&wf.Policy,
			&wf.ExpectedExecutionTime,
			&wf.ActualExecutionTime,
			&wf.AssignedVM,
			&wf.ProcessingStartedAt,
			&wf.AssignedAt,
			&wf.CompletedAt,
			&wf.Status,
			&wf.SubmittedBy,
			&wf.LastUpdated,
		)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, wf)
	}

	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return workflows, nil
}

func (s *MySQLStore) SaveWorkflow(ctx context.Context, w *WorkflowInfo) (*WorkflowInfo, error) {
	res, err := s.db.ExecContext(ctx,
		"INSERT INTO workflow_info (name, created_at, type, ram, core, policy, expected_execution_time, actual_execution_time, assigned_vm, submitted_by, status, assigned_at, completed_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		w.Name, w.CreatedAt, w.Type, w.RAM, w.Core, w.Policy, w.ExpectedExecutionTime, w.ActualExecutionTime, w.AssignedVM, w.SubmittedBy, w.Status, w.AssignedAt, w.CompletedAt,
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

func (s *MySQLStore) UpdateWorkflow(ctx context.Context, w *WorkflowInfo) (*WorkflowInfo, error) {
	_, err := s.db.ExecContext(ctx,
		"UPDATE workflow_info SET name = ?, type = ?, ram = ?, core = ?, policy = ?, expected_execution_time = ?, actual_execution_time = ?, assigned_vm = ?, assigned_at = ?, completed_at = ?, submitted_by = ?, status = ? WHERE id = ?",
		w.Name, w.Type, w.RAM, w.Core, w.Policy, w.ExpectedExecutionTime, w.ActualExecutionTime, w.AssignedVM, w.AssignedAt, w.CompletedAt, w.SubmittedBy, w.Status, w.ID,
	)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (s *MySQLStore) UpdateWorkflowByName(ctx context.Context, w *WorkflowInfo) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE workflow_info SET  actual_execution_time = ?, processing_started_at = ?, completed_at = ?, status = ?  WHERE name = ?",
		w.ActualExecutionTime, w.ProcessingStartedAt, w.CompletedAt, w.Status, w.Name)
	if err != nil {
		return err
	}

	return nil
}

func (s *MySQLStore) AssignWorkflow(ctx context.Context, workflowName string, venName string, expectedTime float64) error {
	query := `UPDATE workflow_info 
			  SET status = 'assigned', 
			      assigned_at = NOW(), 
			      assigned_vm = ?,
			      last_updated = NOW(),
			      expected_execution_time = ?
			  WHERE name = ?`
	_, err := s.db.ExecContext(ctx, query, venName, expectedTime, workflowName)
	return err
}

func (s *MySQLStore) StartWorkflow(ctx context.Context, workflowID int) error {
	query := "UPDATE workflow_info SET status = 'processing', processing_started_at = NOW(), last_updated = NOW() WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, workflowID)
	return err
}

func (s *MySQLStore) CompleteWorkflow(ctx context.Context, id int) error {
	query := "UPDATE workflow_info SET status = 'done', completed_at = NOW(), last_updated = NOW() WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

func (store *MySQLStore) CountWorkflows() (int, error) {
	var count int
	err := store.db.QueryRow("SELECT COUNT(*) FROM workflow_info").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (store *MySQLStore) InsertWorkflow(wf *WorkflowInfo) error {
	_, err := store.db.Exec(`
		INSERT INTO workflow_info 
		(name, type, ram, core, policy, submitted_by, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		wf.Name,
		wf.Type,
		wf.RAM,
		wf.Core,
		wf.Policy,
		wf.SubmittedBy,
		wf.CreatedAt,
	)
	return err
}
