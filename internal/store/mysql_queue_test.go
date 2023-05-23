package store

import (
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)



func TestEnqueue(t *testing.T) {
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
