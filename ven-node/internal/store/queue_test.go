package store

import (
	"context"
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/stretchr/testify/assert"
)

func TestEnqueue(t *testing.T) {
	t.Cleanup(func() {
		if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE workflow;"); err != nil {
			t.Fatalf("Failed to clean up test database: %v", err)
		}
		if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE queue;"); err != nil {
			t.Fatalf("Failed to clean up test database: %v", err)
		}
		if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE queue_size;"); err != nil {
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
	t.Log(queues)
	var found bool
	for _, q := range queues {
		if q.ID == queueID && q.WorkflowID == wf_saved.ID && q.Status == WorkflowStatusQueued {
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
		if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE queue_size;"); err != nil {
			t.Fatalf("Failed to clean up test database: %v", err)
		}
	})

	// Prepare some sample queue data
	queueData := []Queue{
		{
			WorkflowID: 1,
			// Status:      "pending",
		},
		{
			WorkflowID: 2,
			// Status:      "processing",
		},
		{
			WorkflowID: 3,
			// Status:      "done",
		},
		{
			WorkflowID: 4,
			// Status:      "pending",
		},
		{
			WorkflowID: 5,
			// Status:      "processing",
		},
		{
			WorkflowID: 6,
			// Status:      "done",
		},
		{
			WorkflowID: 7,
			// Status:      "pending",
		},
		{
			WorkflowID: 8,
			// Status:      "processing",
		},
		{
			WorkflowID: 9,
			// Status:      "done",
		},
		{
			WorkflowID: 10,
			// Status:      "pending",
		},
		{
			WorkflowID: 11,
			// Status:      "processing",
		},
		{
			WorkflowID: 12,
			// Status:      "done",
		},
	}

	// Call GetQueueStatus on Empty
	queues, err := testStore.GetQueueStatus(context.Background())
	if err != nil {
		t.Fatalf("GetQueueStatus failed: %v", err)
	}

	// Verify the results
	if len(queues) != 0 {
		t.Fatalf("Expected %d queues, got %d", len(queueData), len(queues))
	}

	// Enqueue
	for _, q := range queueData {
		_, err := testStore.Enqueue(ctx, q.WorkflowID)
		if err != nil {
			t.Fatalf("Enqueue failed: %v", err)
		}
		time.Sleep(1 * time.Second)
	}

	// Call GetQueueStatus
	queues, err = testStore.GetQueueStatus(context.Background())
	if err != nil {
		t.Fatalf("GetQueueStatus failed: %v", err)
	}

	// Verify the results
	queue_size_original, err := testStore.GetQueueSizeFromDBorENV(ctx)
	if err != nil {
		fmt.Println("Error in GetQueueSizeFromDBorENV in test")
	}
	actualLen := len(queues)
	if actualLen != queue_size_original { //actualLen != len(queueData) &&
		t.Fatalf("Expected %d queues, got %d", len(queueData), len(queues))
	}

	for i, q := range queues {
		expected := queueData[i]
		if q.WorkflowID != expected.WorkflowID { //|| q.Status != expected.Status
			t.Errorf("Mismatch in queue data at index %d (worflowID %d != expected %d)", i, q.WorkflowID, expected.WorkflowID)
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
		if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE queue_size;"); err != nil {
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
	newStatus := WorkflowStatusInProcess
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
		if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE queue_size;"); err != nil {
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

	newStatus := WorkflowStatusComplete
	err = testStore.CompleteWorkflowInQueue(ctx, wf.ID)
	if err != nil {
		t.Fatalf("CompleteWorkflowInQueue failed: %v", err)
	}

	// Retrieve the updated status
	queues, err := testStore.GetQueueStatus(ctx)
	if err != nil {
		t.Fatalf("GetQueueStatus failed: %v", err)
	}

	// Try to find the updated item in the queue
	for _, q := range queues {
		if q.WorkflowID == wf.ID || q.Status == newStatus {
			t.Fatalf("Completed item found in the queue")
		}
	}
}

func TestIsSpaceAvailable(t *testing.T) {
	t.Cleanup(func() {
		if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE workflow;"); err != nil {
			t.Fatalf("Failed to clean up test database: %v", err)
		}
		if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE queue;"); err != nil {
			t.Fatalf("Failed to clean up test database: %v", err)
		}
		if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE queue_size;"); err != nil {
			t.Fatalf("Failed to clean up test database: %v", err)
		}
	})

	// Prepare sample queue data
	queueData := []Queue{
		{
			WorkflowID: 1,
		},
		{
			WorkflowID: 2,
		},
	}

	// Enqueue queueData
	for _, q := range queueData {
		_, err := testStore.Enqueue(ctx, q.WorkflowID)
		assert.NoError(t, err, "Enqueue failed")
	}

	// Call IsSpaceAvailable when space is available
	available, err := testStore.IsSpaceAvailable(context.Background())
	assert.NoError(t, err, "IsSpaceAvailable failed")
	assert.True(t, available, "Expected space to be available")

	// Fill up the queue to the maximum size
	for i := len(queueData); i < 15; i++ {
		_, err := testStore.Enqueue(ctx, i+1)
		assert.NoError(t, err, "Enqueue failed")
	}

	// Call IsSpaceAvailable when space is not available
	available, err = testStore.IsSpaceAvailable(context.Background())
	assert.NoError(t, err, "IsSpaceAvailable failed")
	assert.False(t, available, "Expected space to be unavailable")
}

func TestPeek(t *testing.T) {
	ctx := context.Background()

	t.Cleanup(func() {
		// Clean up test database
		tables := []string{"workflow", "queue", "queue_size"}
		for _, table := range tables {
			if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE "+table+";"); err != nil {
				t.Fatalf("Failed to clean up table %s: %v", table, err)
			}
		}
	})

	// Prepare some sample queue data
	queueData := []Queue{
		{
			WorkflowID: 1,
		},
		{
			WorkflowID: 2,
		},
		{
			WorkflowID: 3,
		},
		{
			WorkflowID: 4,
		},
		{
			WorkflowID: 5,
		},
		{
			WorkflowID: 6,
		},
		{
			WorkflowID: 7,
		},
		{
			WorkflowID: 8,
		},
		{
			WorkflowID: 9,
		},
		{
			WorkflowID: 10,
		},
		{
			WorkflowID: 11,
		},
		{
			WorkflowID: 12,
		},
	}

	// Call Peek on Empty Queue
	wf, err := testStore.Peek(ctx, WorkflowStatusQueued)
	if wf != nil || err == nil {
		t.Fatalf("Expected nil, got %d", wf.WorkflowID)
	}

	// Enqueue
	for _, q := range queueData {
		_, err := testStore.Enqueue(ctx, q.WorkflowID)
		if err != nil {
			t.Fatalf("Enqueue failed: %v", err)
		}
		time.Sleep(1 * time.Second)
	}

	// Call Peek
	wf, err = testStore.Peek(ctx, WorkflowStatusQueued)
	if err != nil {
		t.Fatalf("Peek failed: %v", err)
	}
	if wf.WorkflowID != queueData[0].WorkflowID {
		t.Fatalf("Expected workflow ID %d, got %d", queueData[0].WorkflowID, wf.WorkflowID)
	}

	// Test each queue item
	for i := range queueData {
		wf, err := testStore.Peek(ctx, WorkflowStatusQueued)
		expectedID := i + 1
		if err != nil {
			t.Fatalf("Peek failed: %v", err)
		}
		if wf.WorkflowID != expectedID {
			t.Errorf("Mismatch in queue data at index %d (expected workflow ID %d, got %d)", i, expectedID, wf.WorkflowID)
		}
		err = testStore.CompleteWorkflowInQueue(ctx, wf.ID)
		if err != nil {
			t.Fatalf("CompleteWorkflowInQueue failed: %v", err)
		}
	}
}

func TestPeekInProcess(t *testing.T) {
	ctx := context.Background()

	t.Cleanup(func() {
		// Clean up test database
		tables := []string{"workflow", "queue", "queue_size"}
		for _, table := range tables {
			if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE "+table+";"); err != nil {
				t.Fatalf("Failed to clean up table %s: %v", table, err)
			}
		}
	})

	// Prepare some sample queue data
	queueData := []Queue{
		{
			WorkflowID: 1,
		},
		{
			WorkflowID: 2,
		},
		{
			WorkflowID: 3,
		},
	}

	// Enqueue
	for _, q := range queueData {
		_, err := testStore.Enqueue(ctx, q.WorkflowID)
		if err != nil {
			t.Fatalf("Enqueue failed: %v", err)
		}
		time.Sleep(1 * time.Second)
	}

	// Set wf to inprocess item
	for range queueData {
		wf, err := testStore.PeekQueued(ctx)
		if err != nil {
			t.Fatalf("Peek failed: %v", err)
		}

		err = testStore.ProcessWorkflowInQueue(ctx, wf.ID)
		if err != nil {
			t.Fatalf("ProcessWorkflowInQueue failed: %v", err)
		}
	}

	// Test each queue item
	for i := range queueData {
		wf, err := testStore.Peek(ctx, WorkflowStatusInProcess)
		expectedID := i + 1
		if err != nil {
			t.Fatalf("Peek failed: %v", err)
		}
		if wf.WorkflowID != expectedID {
			t.Errorf("Mismatch in queue data at index %d (expected workflow ID %d, got %d)", i, expectedID, wf.WorkflowID)
		}
		err = testStore.CompleteWorkflowInQueue(ctx, wf.ID)
		if err != nil {
			t.Fatalf("CompleteWorkflowInQueue failed: %v", err)
		}
	}
}

func TestPeekQueued(t *testing.T) {
	ctx := context.Background()

	t.Cleanup(func() {
		// Clean up test database
		tables := []string{"workflow", "queue", "queue_size"}
		for _, table := range tables {
			if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE "+table+";"); err != nil {
				t.Fatalf("Failed to clean up table %s: %v", table, err)
			}
		}
	})

	// Prepare some sample queue data
	queueData := []Queue{
		{
			WorkflowID: 1,
		},
		{
			WorkflowID: 2,
		},
		{
			WorkflowID: 3,
		},
	}

	// Enqueue
	for _, q := range queueData {
		_, err := testStore.Enqueue(ctx, q.WorkflowID)
		if err != nil {
			t.Fatalf("Enqueue failed: %v", err)
		}
		time.Sleep(1 * time.Second)
	}

	// Test each queue item
	for i := range queueData {
		wf, err := testStore.PeekQueued(ctx)
		expectedID := i + 1
		if err != nil {
			t.Fatalf("Peek failed: %v", err)
		}
		if wf.WorkflowID != expectedID {
			t.Errorf("Mismatch in queue data at index %d (expected workflow ID %d, got %d)", i, expectedID, wf.WorkflowID)
		}
		err = testStore.CompleteWorkflowInQueue(ctx, wf.ID)
		if err != nil {
			t.Fatalf("CompleteWorkflowInQueue failed: %v", err)
		}
	}
}

func TestGetAvailableSpace(t *testing.T) {

	t.Cleanup(func() {
		if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE queue;"); err != nil {
			t.Fatalf("Failed to clean up test database: %v", err)
		}
		if _, err := testStore.db.ExecContext(ctx, "TRUNCATE TABLE queue_size;"); err != nil {
			t.Fatalf("Failed to clean up test database: %v", err)
		}
	})

	// Prepare sample queue data
	queueData := []Queue{
		{
			WorkflowID: 1,
		},
		{
			WorkflowID: 2,
		},
	}

	// Enqueue queueData
	for _, q := range queueData {
		_, err := testStore.Enqueue(ctx, q.WorkflowID)
		assert.NoError(t, err, "Enqueue failed")
	}

	// Call GetAvailableSpace when space is available
	availableSpace, err := testStore.GetAvailableSpace(ctx)
	assert.NoError(t, err, "GetAvailableSpace failed")
	fmt.Println(availableSpace)
	assert.True(t, availableSpace > 0, "Expected available space to be more than 0")

	// Fill up the queue to the maximum size
	queue_size_original, err := testStore.GetQueueSizeFromDBorENV(ctx)
	if err != nil {
		fmt.Println("Error in GetQueueSizeFromDBorENV in test")
	}
	for i := len(queueData); i < queue_size_original; i++ {
		_, err := testStore.Enqueue(ctx, i+1)
		assert.NoError(t, err, "Enqueue failed")
	}

	// Call GetAvailableSpace when space is not available
	availableSpace, err = testStore.GetAvailableSpace(ctx)
	assert.NoError(t, err, "GetAvailableSpace failed")
	assert.Equal(t, 0, availableSpace, "Expected available space to be 0")

	// Check if the function correctly handles cases where the current size is greater than the queue size
	_, err = testStore.Enqueue(ctx, queue_size_original+1)
	assert.NoError(t, err, "Enqueue failed")

	_, err = testStore.GetAvailableSpace(ctx)
	assert.Error(t, err, "Expected GetAvailableSpace to return an error")
	assert.Contains(t, err.Error(), "negative available space", "Error message should contain 'negative available space'")
}
