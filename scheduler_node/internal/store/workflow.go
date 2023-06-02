package store

import (
	"context"
	"time"
)

// WorkflowInfo represents the information of a WorkflowInfo
type WorkflowInfo struct {
	ID                    int    `json:"id"`
	Name                  string `json:"name"`
	RAM                   string `json:"ram"`
	Core                  string `json:"core"`
	Policy                string `json:"policy"`
	ExpectedExecutionTime string `json:"expectedExecutionTime"`
	ActualExecutionTime   string `json:"actualExecutionTime"`
	AssignedVM            string `json:"assignedVM"`
	AssignedAt            time.Time `json:"assignedAt"`
	CompletedAt           time.Time `json:"completedAt"`
	Type                  string `json:"type"`
	Status                string `json:"status"`
	SubmittedBy           string    `json:"submittedBy"`
	LastUpdated           time.Time `json:"lastUpdated"`

}

// GetWorkflowInfo retrieves WorkflowInfo information from the database
func (s *MySQLStore) GetWorkflowInfo() ([]WorkflowInfo, error) {
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
			ExpectedExecutionTime: "2 hours",
			ActualExecutionTime:   "1 hour",
			AssignedVM:            "VM1",
			AssignedAt:            assignedAt,
			CompletedAt:           completedAt,
			Type:                  "DNA",
			Status:                "Pending",
			SubmittedBy:           "User A",
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
			AssignedAt:            assignedAt,
			CompletedAt:           completedAt,
			Type:                  "RNA",
			Status:                "Completed",
			SubmittedBy:           "User B",
			LastUpdated:           time.Now().Add(-time.Hour),
		},
	}

	// Return the pseudo data in case of error
	return pseudoData, nil
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

func (s *MySQLStore) GetWorkflows(ctx context.Context) ([]WorkflowInfo, error) {
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



func (s *MySQLStore) StartWorkflow(ctx context.Context, workflowID int) error {
	query := "UPDATE workflow_info SET status = 'processing', assigned_at = NOW(), last_updated = NOW() WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, workflowID)
	return err
}


func (s *MySQLStore) CompleteWorkflow(ctx context.Context, id int) error {
	query := "UPDATE workflow_info SET status = 'done', completed_at = NOW(), last_updated = NOW() WHERE id = ?"
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

