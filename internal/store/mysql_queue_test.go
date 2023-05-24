package store

import (
	"context"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)



func TestEnqueue(t *testing.T) {
	t.Cleanup(func() {
		if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE workflow;"); err != nil {
			t.Fatalf("Failed to clean up test database: %v", err)
		}
		if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE queue;"); err != nil {
			t.Fatalf("Failed to clean up test database: %v", err)
		}
	})

	// Preparing a new workflow to be used for the enqueue operation
	wf := &Workflow{
		Name:       "test-workflow",
		Type:       "type1",
		Duration:   1,
		ReceivedAt: time.Now(),
	}
	wf_saved, err := testStore.SaveWorkflow(ctx, wf)
	if err != nil {
		t.Fatalf("CreateWorkflow failed: %v", err)
	}

	// Test Enqueue
	queueID, err := testStore.Enqueue(ctx, wf_saved.ID)
	if err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}

	// Test GetQueueStatus
	queues, err := testStore.GetQueueStatus(ctx)
	if err != nil {
		t.Fatalf("GetQueueStatus failed: %v", err)
	}

	var found bool
	for _, q := range queues {
		if q.ID == queueID && q.WorkflowID == wf_saved.ID && q.Status == "pending" {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("Enqueued item not found in the queue")
	}
}


func TestGetQueueStatus(t *testing.T) {
	t.Cleanup(func() {
		if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE workflow;"); err != nil {
			t.Fatalf("Failed to clean up test database: %v", err)
		}
		if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE queue;"); err != nil {
			t.Fatalf("Failed to clean up test database: %v", err)
		}
	})

	// Prepare some sample queue data
	queueData := []Queue{
		{
			WorkflowID:  1,
			// Status:      "pending",
		},
		{
			WorkflowID:  2,
			// Status:      "processing",
		},
		{
			WorkflowID:  3,
			// Status:      "done",
		},
		{
			WorkflowID:  4,
			// Status:      "pending",
		},
		{
			WorkflowID:  5,
			// Status:      "processing",
		},
		{
			WorkflowID:  6,
			// Status:      "done",
		},
		{
			WorkflowID:  7,
			// Status:      "pending",
		},
		{
			WorkflowID:  8,
			// Status:      "processing",
		},
		{
			WorkflowID:  9,
			// Status:      "done",
		},
		{
			WorkflowID:  10,
			// Status:      "pending",
		},
		{
			WorkflowID:  11,
			// Status:      "processing",
		},
		{
			WorkflowID:  12,
			// Status:      "done",
		},
	}



	// Call GetQueueStatus on Empty 
	queues, err := testStore.GetQueueStatus(context.Background())
	if err != nil {
		t.Fatalf("GetQueueStatus failed: %v", err)
	}

	// Verify the results
	if len(queues) != 0{
		t.Fatalf("Expected %d queues, got %d", len(queueData), len(queues))
	}

	// Enqueue
	for _,q := range queueData{
		_, err := testStore.Enqueue(ctx, q.WorkflowID)
		if err != nil {
			t.Fatalf("Enqueue failed: %v", err)
		}
		time.Sleep(1*time.Second)
	}

	// Call GetQueueStatus
	queues, err = testStore.GetQueueStatus(context.Background())
	if err != nil {
		t.Fatalf("GetQueueStatus failed: %v", err)
	}

	// Verify the results
	actualLen := len(queues)
	if actualLen != 10 { //actualLen != len(queueData) &&  
		t.Fatalf("Expected %d queues, got %d", len(queueData), len(queues))
	}

	for i, q := range queues {
		expected := queueData[i]
		if q.WorkflowID != expected.WorkflowID   { //|| q.Status != expected.Status
			t.Errorf("Mismatch in queue data at index %d (worflowID %d != expected %d)", i,q.WorkflowID, expected.WorkflowID)
		}
	}
}

func TestProcessWorkflowInQueue(t *testing.T) {
	t.Cleanup(func() {
		if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE workflow;"); err != nil {
			t.Fatalf("Failed to clean up test database: %v", err)
		}
		if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE queue;"); err != nil {
			t.Fatalf("Failed to clean up test database: %v", err)
		}
	})
	// Preparing test data
	wf := &Workflow{
		Name:       "test-workflow",
		Type:       "type1",
		Duration:   1,
		ReceivedAt: time.Now(),
	}
	wf, err := testStore.SaveWorkflow(ctx, wf)
	if err != nil {
		t.Fatalf("CreateWorkflow failed: %v", err)
	}

	// Enqueue the workflow
	_, err = testStore.Enqueue(ctx, wf.ID)
	if err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}

	// Test UpdateStatus
	newStatus := "processing"
	err = testStore.ProcessWorkflowInQueue(ctx, wf.ID)
	if err != nil {
		t.Fatalf("UpdateStatus failed: %v", err)
	}

	// Retrieve the updated status
	queues, err := testStore.GetQueueStatus(ctx)
	if err != nil {
		t.Fatalf("GetQueueStatus failed: %v", err)
	}

	// Find the updated item in the queue
	var found bool
	for _, q := range queues {
		if q.WorkflowID == wf.ID && q.Status == newStatus {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("Updated item not found in the queue")
	}
}


func TestCompleteWorkflowInQueue(t *testing.T) {
	t.Cleanup(func() {
		if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE workflow;"); err != nil {
			t.Fatalf("Failed to clean up test database: %v", err)
		}
		if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE queue;"); err != nil {
			t.Fatalf("Failed to clean up test database: %v", err)
		}
	})
	// Preparing test data
	wf := &Workflow{
		Name:       "test-workflow",
		Type:       "type1",
		Duration:   1,
		ReceivedAt: time.Now(),
	}
	wf, err := testStore.SaveWorkflow(ctx, wf)
	if err != nil {
		t.Fatalf("CreateWorkflow failed: %v", err)
	}

	// Enqueue the workflow
	_, err = testStore.Enqueue(ctx, wf.ID)
	if err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}

	newStatus := "done"
	err = testStore.ProcessWorkflowInQueue(ctx, wf.ID)
	if err != nil {
		t.Fatalf("UpdateStatus failed: %v", err)
	}

	// Retrieve the updated status
	queues, err := testStore.GetQueueStatus(ctx)
	if err != nil {
		t.Fatalf("GetQueueStatus failed: %v", err)
	}

	// Try to find the updated item in the queue
	for _, q := range queues {
		if q.WorkflowID == wf.ID && q.Status == newStatus {
			t.Fatalf("Completed item found in the queue")
		}
	}


}
