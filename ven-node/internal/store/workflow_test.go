package store

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)

func TestWorkflow(t *testing.T) {
	// Test CreateWorkflow
	newWorkflow := Workflow{
		Name:       "Test Workflow",
		Type:       "Sequential",
		Duration:   10,
		ReceivedAt: time.Now(),
		StartedExecutionAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		CompletedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}

	createdWorkflow, err := testStore.SaveWorkflow(ctx, &newWorkflow)
	if err != nil {
		t.Fatalf("CreateWorkflow failed: %v", err)
	}
	if createdWorkflow.ID == 0 {
		t.Fatalf("Expected a valid ID, got %d", createdWorkflow.ID)
	}

	// Test GetWorkflowByID
	workflowByID, err := testStore.GetWorkflowByID(ctx, createdWorkflow.ID)
	if err != nil {
		t.Fatalf("GetWorkflowByID failed: %v", err)
	}

	if workflowByID.ID != createdWorkflow.ID {
		t.Fatalf("Expected workflow ID %d, got %d", createdWorkflow.ID, workflowByID.ID)
	}

	// Test UpdateWorkflow
	workflowByID.Duration = 20
	_, err = testStore.UpdateWorkflow(ctx, workflowByID)
	if err != nil {
		t.Fatalf("UpdateWorkflow failed: %v", err)
	}

	// Test GetWorkflowByID after update
	updatedWorkflow, err := testStore.GetWorkflowByID(ctx, createdWorkflow.ID)
	if err != nil {
		t.Fatalf("GetWorkflowByID failed: %v", err)
	}

	if updatedWorkflow.Duration != 20 {
		t.Fatalf("Expected workflow duration 20, got %d", updatedWorkflow.Duration)
	}

	// Test StartWorkflow when started_execution_at is NULL
	err = testStore.StartWorkflow(ctx, createdWorkflow.ID)
	if err != nil {
		t.Fatalf("StartWorkflow failed: %v", err)
	}

	// Test GetWorkflowByID after StartWorkflow
	workflowAfterStart, err := testStore.GetWorkflowByID(ctx, createdWorkflow.ID)
	if err != nil {
		t.Fatalf("GetWorkflowByID after StartWorkflow failed: %v", err)
	}

	if !workflowAfterStart.StartedExecutionAt.Valid {
		t.Fatalf("Expected started_execution_at to be set after StartWorkflow")
	}

	initialStartTime := workflowAfterStart.StartedExecutionAt.Time

	// Test StartWorkflow again, expecting started_execution_at to remain unchanged
	err = testStore.StartWorkflow(ctx, createdWorkflow.ID)
	if err != nil {
		t.Fatalf("StartWorkflow failed on the second call: %v", err)
	}

	// Fetch again to check
	workflowAfterSecondStart, err := testStore.GetWorkflowByID(ctx, createdWorkflow.ID)
	if err != nil {
		t.Fatalf("GetWorkflowByID after second StartWorkflow failed: %v", err)
	}

	if workflowAfterSecondStart.StartedExecutionAt.Time != initialStartTime {
		t.Fatalf("Expected started_execution_at to remain unchanged after the second StartWorkflow")
	}
}
