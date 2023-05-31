package store

import (
	"context"
	"database/sql"
	"time"
)

// WorkflowInfo represents the information of a workflow
type WorkflowInfo struct {
	Name                  string `json:"name"`
	RAM                   string `json:"ram"`
	Core                  string `json:"core"`
	Policy                string `json:"policy"`
	ExpectedExecutionTime string `json:"expectedExecutionTime"`
	ActualExecutionTime   string `json:"actualExecutionTime"`
	AssignedVM            string `json:"assignedVM"`
	AssignedAt            string `json:"assignedAt"`
	CompletedAt           string `json:"completedAt"`
	Type                  string `json:"type"`
	Status                string `json:"status"`
	LastUpdated           time.Time `json:"lastUpdated"`

}

// GetWorkflowInfo retrieves workflow information from the database
func (s *MySQLStore) GetWorkflowInfo() ([]WorkflowInfo, error) {
	// TODO: Implement the actual retrieval of workflow info from the database
	// For this example, let's assume the data is retrieved from the database successfully

	// Simulated pseudo data in case of error
	pseudoData := []WorkflowInfo{
		{
			Name:                  "WF1",
			RAM:                   "8GB",
			Core:                  "4",
			Policy:                "Some Policy",
			ExpectedExecutionTime: "2 hours",
			ActualExecutionTime:   "1 hour",
			AssignedVM:            "VM1",
			AssignedAt:            "2023-05-31 10:00:00",
			CompletedAt:           "2023-05-31 11:00:00",
			Type:                  "DNA",
			Status:                "Pending",
			LastUpdated:           time.Now().Add(-time.Minute*2),
		},
		{
			Name:                  "WF2",
			RAM:                   "16GB",
			Core:                  "8",
			Policy:                "Another Policy",
			ExpectedExecutionTime: "3 hours",
			ActualExecutionTime:   "2.5 hours",
			AssignedVM:            "VM2",
			AssignedAt:            "2023-05-31 12:00:00",
			CompletedAt:           "2023-05-31 14:30:00",
			Type:                  "RNA",
			Status:                "Completed",
			LastUpdated:           time.Now().Add(-time.Hour),
		},
	}

	// Return the pseudo data in case of error
	return pseudoData, nil
}


type Workflow struct {
	ID                  int     `json:"id"`
	Name                string    `json:"name"`
	Type                string    `json:"type"`
	Duration            int     `json:"duration"`
	ReceivedAt          time.Time `json:"received_at"`
	StartedExecutionAt  sql.NullTime `json:"started_execution_at,omitempty"`
	CompletedAt         sql.NullTime `json:"completed_at,omitempty"`
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
		return nil,err
	}

	return w,nil
}


func (s *MySQLStore) StartWorkflow(ctx context.Context, workflowID int) error {
	query := "UPDATE workflow SET started_execution_at = NOW() WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, workflowID)
	return err
}

func (s *MySQLStore) CompleteWorkflow(ctx context.Context, id int) error {
	query := "UPDATE workflow SET completed_at = NOW() WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}
