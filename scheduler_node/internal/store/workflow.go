package store

import (
	"context"
	"database/sql"
	"time"
)

// WorkflowInfo represents the information of a WorkflowInfo
type WorkflowInfo struct {
	ID                    int           `json:"id"`
	Name                  string        `json:"name"`
	RAM                   string        `json:"ram"`
	Core                  string        `json:"core"`
	Policy                string        `json:"policy"`
	ExpectedExecutionTime sql.NullString  `json:"expectedExecutionTime"`
	ActualExecutionTime   sql.NullString  `json:"actualExecutionTime"`
	AssignedVM            sql.NullString  `json:"assignedVM"`
	ProcessingStartedAt   sql.NullTime  `json:"processingStartedAt"`
	AssignedAt            sql.NullTime  `json:"assignedAt"`
	CompletedAt           sql.NullTime  `json:"completedAt"`
	Type                  string        `json:"type"`
	Status                string        `json:"status"`
	SubmittedBy           sql.NullString  `json:"submittedBy"`
	LastUpdated           time.Time     `json:"lastUpdated"`
}







func (s *MySQLStore) GetWorkflowByID(ctx context.Context, id int) (*WorkflowInfo, error) {
	wf := &WorkflowInfo{}
	err := s.db.QueryRowContext(ctx,
		"SELECT id,name, type, ram, core, policy, expected_execution_time, actual_execution_time, assigned_vm, assigned_at, completed_at, submitted_by, status, last_updated FROM workflow_info WHERE id = ?",
		id,
	).Scan(&wf.ID,&wf.Name, &wf.Type, &wf.RAM, &wf.Core, &wf.Policy, &wf.ExpectedExecutionTime, &wf.ActualExecutionTime, &wf.AssignedVM, &wf.AssignedAt, &wf.CompletedAt, &wf.SubmittedBy, &wf.Status, &wf.LastUpdated)

	if err != nil {
		return nil, err
	}

	return wf, nil
}

// pseudoGetWorkflowInfo retrieves WorkflowInfo information from the database
func (s *MySQLStore) pseudoGetWorkflowInfo(ctx context.Context) ([]WorkflowInfo, error) {
	// TODO: Implement the actual retrieval of WorkflowInfo info from the database
	// For this example, let's assume the data is retrieved from the database successfully
	const layout = "2006-01-02 15:04:05"
	assignedAt, err := time.Parse(layout, "2023-05-31 10:00:00")
	if err != nil {
		panic(err)
	}
	assignedAt = assignedAt.UTC()

	completedAt, err := time.Parse(layout, "2023-05-31 10:00:00")
	if err != nil {
		panic(err)
	}
	completedAt = completedAt.UTC()
	// Simulated pseudo data in case of error
	pseudoData := []WorkflowInfo{
		{
			Name:                  "WF1",
			RAM:                   "8GB",
			Core:                  "4",
			Policy:                "Some Policy",
			ExpectedExecutionTime: sql.NullString{String:"2 hours", Valid: true},
			ActualExecutionTime:   sql.NullString{String:"2 hours", Valid: true},
			AssignedVM:            sql.NullString{String:"VM1", Valid: true},
			AssignedAt:            sql.NullTime{Time:assignedAt, Valid: true},
			CompletedAt:           sql.NullTime{Time:completedAt, Valid: true},
			Type:                  "DNA",
			Status:                "Pending",
			SubmittedBy:           sql.NullString{String:"User A", Valid: true},
			LastUpdated:           time.Now().Add(-time.Minute*2),
		},
		{
			Name:                  "WF2",
			RAM:                   "16GB",
			Core:                  "8",
			Policy:                "Another Policy",
			ExpectedExecutionTime: sql.NullString{String:"3 hours", Valid: true},
			ActualExecutionTime:   sql.NullString{String:"2.5 hours", Valid: true},
			AssignedVM:            sql.NullString{String:"VM2", Valid: true},
			AssignedAt:            sql.NullTime{Time:assignedAt, Valid: true},
			CompletedAt:           sql.NullTime{Time:completedAt, Valid: true},
			Type:                  "RNA",
			Status:                "Completed",
			SubmittedBy:           sql.NullString{String:"User A", Valid: true},
			LastUpdated:           time.Now().Add(-time.Hour),
		},
	}

	// Return the pseudo data in case of error
	return pseudoData, nil
}

func (s *MySQLStore) GetWorkflows(ctx context.Context) ([]WorkflowInfo, error) {
	query := `SELECT id, name, type, ram, core, policy, 
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

/*
func (s *MySQLStore) GetWorkflows(ctx context.Context) ([]WorkflowInfo, error) {
	return s.pseudoGetWorkflowInfo(ctx)
	
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, type, ram, core, policy, expected_execution_time, actual_execution_time, assigned_vm, assigned_at, completed_at, submitted_by, status, last_updated FROM workflow_info")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workflows := []WorkflowInfo{}
	for rows.Next() {
		var wf WorkflowInfo
		err := rows.Scan(&wf.ID,&wf.Name, &wf.Type, &wf.RAM, &wf.Core, &wf.Policy, &wf.ExpectedExecutionTime, &wf.ActualExecutionTime, &wf.AssignedVM, &wf.AssignedAt, &wf.CompletedAt, &wf.SubmittedBy, &wf.Status, &wf.LastUpdated)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, wf)
	}

	return workflows, nil

}
	*/
	
func (s *MySQLStore) SaveWorkflow(ctx context.Context, w *WorkflowInfo) (*WorkflowInfo, error) {
	res, err := s.db.ExecContext(ctx,
		"INSERT INTO workflow_info (name, type, ram, core, policy, expected_execution_time, actual_execution_time, assigned_vm, submitted_by, status, assigned_at, completed_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		w.Name, w.Type, w.RAM, w.Core, w.Policy, w.ExpectedExecutionTime, w.ActualExecutionTime, w.AssignedVM, w.SubmittedBy, w.Status,  w.AssignedAt, w.CompletedAt,
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
		return nil,err
	}

	return w,nil
}

func (s *MySQLStore) AssignWorkflow(ctx context.Context, workflowID int) error {
	query := "UPDATE workflow_info SET status = 'assigned', assigned_at = NOW(), last_updated = NOW() WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, workflowID)
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

