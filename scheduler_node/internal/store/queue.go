package store

import (
	"context"
	"fmt"
	"os"
	"strconv"
)

// Peek gets the first workflow in the queue (that is not 'done') and does not remove it from the queue
func (s *MySQLStore) Peek(ctx context.Context) (*WorkflowInfo, error) {
	wf := &WorkflowInfo{}
	err := s.db.QueryRowContext(ctx,
		"SELECT id, created_at, name, type, ram, core, policy, expected_execution_time, actual_execution_time, assigned_vm, assigned_at, completed_at, submitted_by, status, last_updated FROM workflow_info WHERE status <> 'done' ORDER BY id ASC LIMIT 1",
	).Scan(&wf.ID,&wf.CreatedAt,&wf.Name, &wf.Type, &wf.RAM, &wf.Core, &wf.Policy, &wf.ExpectedExecutionTime, &wf.ActualExecutionTime, &wf.AssignedVM, &wf.AssignedAt, &wf.CompletedAt, &wf.SubmittedBy, &wf.Status, &wf.LastUpdated)

	if err != nil {
		return nil, err
	}

	return wf, nil
}

// GetQueue gets all workflows that are not 'done'
func (s *MySQLStore) GetQueue(ctx context.Context) ([]WorkflowInfo, error) {
	// get QUEUE_SIZE from environment variable
	queueSizeStr := os.Getenv("QUEUE_SIZE")
	if queueSizeStr == "" {
		return nil, fmt.Errorf("QUEUE_SIZE environment variable not set")
	}
	
	// convert it to integer
	queueSize, err := strconv.Atoi(queueSizeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid QUEUE_SIZE value: %v", err)
	}
	
	// use as a limit in SQL query
	rows, err := s.db.QueryContext(ctx, "SELECT name,created_at, type, ram, core, policy, expected_execution_time, actual_execution_time, assigned_vm, assigned_at, completed_at, submitted_by, status, last_updated FROM workflow_info WHERE status != 'done' ORDER BY id ASC LIMIT ?", queueSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workflows := []WorkflowInfo{}
	for rows.Next() {
		var wf WorkflowInfo
		err := rows.Scan(&wf.Name,&wf.CreatedAt, &wf.Type, &wf.RAM, &wf.Core, &wf.Policy, &wf.ExpectedExecutionTime, &wf.ActualExecutionTime, &wf.AssignedVM, &wf.AssignedAt, &wf.CompletedAt, &wf.SubmittedBy, &wf.Status, &wf.LastUpdated)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, wf)
	}

	return workflows, nil
}





