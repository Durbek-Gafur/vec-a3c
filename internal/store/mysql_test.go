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


// package store_test

// import (
// 	"context"
// 	"net/http"
// 	"testing"
// 	"time"

// 	"vec-node/internal/store"
// 	"vec-node/internal/workflow"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/stretchr/testify/assert"
// )

// func TestMySQLStore_GetQueueSize(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)

// 	defer db.Close()

// 	rows := sqlmock.NewRows([]string{"size"}).AddRow(10)
// 	mock.ExpectQuery("SELECT size FROM queue_size WHERE id = 1").WillReturnRows(rows)

// 	mysqlStore := store.NewMySQLStoreWithDB(db)

// 	size, err := mysqlStore.GetQueueSize(context.Background())
// 	assert.NoError(t, err)
// 	assert.Equal(t, 10, size)

// 	rows = sqlmock.NewRows([]string{"size"})
// 	mock.ExpectQuery("SELECT size FROM queue_size WHERE id = 1").WillReturnRows(rows)

// 	size, err = mysqlStore.GetQueueSize(context.Background())
// 	assert.Error(t, err)
// 	assert.Equal(t, http.StatusNotFound, err.(*store.StoreError).StatusCode)
// 	assert.Equal(t, store.ErrNotFound, err.(*store.StoreError).Err)
// }

// func TestMySQLStore_SetQueueSize(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)

// 	defer db.Close()

// 	mock.ExpectExec("INSERT INTO queue_size").WithArgs(10, 10).WillReturnResult(sqlmock.NewResult(1, 1))

// 	mysqlStore := store.NewMySQLStoreWithDB(db)

// 	err = mysqlStore.SetQueueSize(context.Background(), 10)
// 	assert.NoError(t, err)

// 	mock.ExpectExec("INSERT INTO queue_size").WithArgs(20, 20).WillReturnResult(sqlmock.NewResult(1, 1))
// 	mock.ExpectExec("UPDATE queue_size").WithArgs(20).WillReturnResult(sqlmock.NewResult(1, 1))

// 	err = mysqlStore.SetQueueSize(context.Background(), 20)
// 	assert.NoError(t, err)

// 	err = mysqlStore.SetQueueSize(context.Background(), 20)
// 	assert.NoError(t, err)
// }

// func TestMySQLStore_UpdateQueueSize(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)

// 	defer db.Close()

// 	mock.ExpectExec("UPDATE queue_size").WithArgs(10).WillReturnResult(sqlmock.NewResult(1, 1))

// 	mysqlStore := store.NewMySQLStoreWithDB(db)

// 	err = mysqlStore.UpdateQueueSize(context.Background(), 10)
// 	assert.NoError(t, err)
// }

// func TestMySQLStore_GetWorkflowByID(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)

// 	defer db.Close()

// 	wf := &workflow.Workflow{
// 		ID:                  1,
// 		Name:                "workflow-1",
// 		Type:                "type-1",
// 		Duration:            1 * time.Hour,
// 		ReceivedAt:          time.Now(),
// 		StartedExecutionAt:  time.Now(),
// 		CompletedAt:         time.Now(),
// 	}

// 	rows := sqlmock.NewRows([]string{"id", "name", "type", "duration", "received_at", "started_execution_at", "completed_at"}).
// 		AddRow(wf.ID, wf.Name, wf.Type, wf.Duration, wf.ReceivedAt, wf.StartedExecutionAt, wf.CompletedAt)
// 	mock.ExpectQuery("SELECT id, name, type, duration, received_at, started_execution_at, completed_at FROM workflows WHERE id = ?").WithArgs(wf.ID).WillReturnRows(rows)
// 	mysqlStore := store.NewMySQLStoreWithDB(db)

// 	result, err := mysqlStore.GetWorkflowByID(wf.ID)
// 	assert.NoError(t, err)
// 	assert.Equal(t, wf, result)

// 	rows = sqlmock.NewRows([]string{"id", "name", "type", "duration", "received_at", "started_execution_at", "completed_at"})
// 	mock.ExpectQuery("SELECT id, name, type, duration, received_at, started_execution_at, completed_at FROM workflows WHERE id = ?").WithArgs(wf.ID).WillReturnRows(rows)

// 	result, err = mysqlStore.GetWorkflowByID(wf.ID)
// 	assert.Error(t, err)
// 	assert.Nil(t, result)

// }

// func TestMySQLStore_GetWorkflows(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	wf1 := workflow.Workflow{
// 		ID:         1,
// 		Name:       "workflow-1",
// 		Type:       "type-1",
// 		Duration:   1 * time.Hour,
// 		ReceivedAt: time.Now(),
// 	}
// 	wf2 := workflow.Workflow{
// 		ID:         2,
// 		Name:       "workflow-2",
// 		Type:       "type-2",
// 		Duration:   2 * time.Hour,
// 		ReceivedAt: time.Now().Add(-2 * time.Hour),
// 	}
// 	wf3 := workflow.Workflow{
// 		ID:         3,
// 		Name:       "workflow-3",
// 		Type:       "type-2",
// 		Duration:   3 * time.Hour,
// 		ReceivedAt: time.Now().Add(-4 * time.Hour),
// 	}

// 	rows := sqlmock.NewRows([]string{"id", "name", "type", "duration", "received_at", "started_execution_at", "completed_at"}).
// 		AddRow(wf1.ID, wf1.Name, wf1.Type, wf1.Duration, wf1.ReceivedAt, wf1.StartedExecutionAt, wf1.CompletedAt).
// 		AddRow(wf2.ID, wf2.Name, wf2.Type, wf2.Duration, wf2.ReceivedAt, wf2.StartedExecutionAt, wf2.CompletedAt).
// 		AddRow(wf3.ID, wf3.Name, wf3.Type, wf3.Duration, wf3.ReceivedAt, wf3.StartedExecutionAt, wf3.CompletedAt)

// 	mock.ExpectQuery("SELECT id, name, type, duration, received_at, started_execution_at, completed_at FROM workflows WHERE 1").WillReturnRows(rows)

// 	mysqlStore := store.NewMySQLStoreWithDB(db)

// 	filter := &workflow.WorkflowFilter{
// 		Type:       "type-2",
// 		StartTime:  time.Now().Add(-5 * time.Hour),
// 		EndTime:    time.Now(),
// 	}

// 	result, err := mysqlStore.GetWorkflows(filter)
// 	assert.NoError(t, err)
// 	assert.Len(t, result, 2)
// 	assert.Contains(t, result, wf2)
// 	assert.Contains(t, result, wf3)
// }

// func TestMySQLStore_SaveWorkflow(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	mock.ExpectExec("INSERT INTO workflows").WithArgs("workflow-1", "type-1", 1*time.Hour, mock.AnythingOfType("time.Time")).WillReturnResult(sqlmock.NewResult(1, 1))

// 	mysqlStore := store.NewMySQLStoreWithDB(db)

// 	wf := &workflow.Workflow{
// 		Name:       "workflow-1",
// 		Type:       "type-1",
// 		Duration:   1 * time.Hour,
// 		ReceivedAt: time.Now(),
// 	}

// 	err = mysqlStore.SaveWorkflow(wf)
// 	assert.NoError(t, err)
// 	assert.NotZero(t, wf.ID)

// }

// func TestNewMySQLStore(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	mock.ExpectPing()

// 	mysqlStore, err := store.NewMySQLStoreWithDB(db)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, mysqlStore)

// }

// func TestNewMySQLStore_WithInvalidDSN(t *testing.T) {
// 	mysqlStore, err := store.NewMySQLStore("invalid-dsn")
// 	assert.Error(t, err)
// 	assert.Nil(t, mysqlStore)
// }

// // func TestNewMySQLStore_WithNilDB(t *testing.T) {
// // 	mysqlStore, err := store.NewMySQLStore(nil)
// // 	assert.Error(t, err)
// // 	assert.Nil(t, mysqlStore)
// // }


