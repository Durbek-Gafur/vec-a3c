package store

import (
	"context"
	"os"
	"testing"
	"time"
	w "vec-node/internal/workflow"

	_ "github.com/go-sql-driver/mysql"
)

var testStore *MySQLStore
var ctx context.Context


func TestMain(m *testing.M) {
	host := os.Getenv("MYSQL_HOST")
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	database := os.Getenv("MYSQL_DBNAME")
	port := os.Getenv("MYSQL_PORT")
	// log.Printf(database)

	dsn := user+ ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?parseTime=true"
	
	store, err := NewMySQLStore(dsn)
	if err != nil {
		panic(err)
	}
	testStore = store
	ctx = context.Background()
	code := m.Run()
	os.Exit(code)
}

func TestQueueSize(t *testing.T) {
	// Test SetQueueSize
	err := testStore.SetQueueSize(ctx,5)
	if err != nil {
		t.Fatalf("SetQueueSize failed: %v", err)
	}

	// Test GetQueueSize
	size, err := testStore.GetQueueSize(ctx)
	if err != nil {
		t.Fatalf("GetQueueSize failed: %v", err)
	}

	if size != 5 {
		t.Fatalf("Expected queue size 5, got %d", size)
	}

	// Test UpdateQueueSize
	err = testStore.UpdateQueueSize(ctx,10)
	if err != nil {
		t.Fatalf("UpdateQueueSize failed: %v", err)
	}

	// Test GetQueueSize after update
	size, err = testStore.GetQueueSize(ctx)
	if err != nil {
		t.Fatalf("GetQueueSize failed: %v", err)
	}

	if size != 10 {
		t.Fatalf("Expected queue size 10, got %d", size)
	}
}

func TestWorkflow(t *testing.T) {
	// Test CreateWorkflow
	newWorkflow := w.Workflow{
		Name:              "Test Workflow",
		Type:              "Sequential",
		Duration:          10,
		ReceivedAt:        time.Now(),
		StartedExecutionAt: time.Now(),
		CompletedAt:       time.Now(),
	}

	createdWorkflow, err := testStore.SaveWorkflow(ctx,&newWorkflow)
	if err != nil {
		t.Fatalf("CreateWorkflow failed: %v", err)
	}
	t.Log(newWorkflow)
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
	_,err = testStore.UpdateWorkflow(ctx, workflowByID)
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
}



