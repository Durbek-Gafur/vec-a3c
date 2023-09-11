package store

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/stretchr/testify/assert"
)

var testStore *MySQLStore
var ctx context.Context

func TestMain(m *testing.M) {
	runTests := func() int {
		host := os.Getenv("MYSQL_HOST")
		user := os.Getenv("MYSQL_USER")
		password := os.Getenv("MYSQL_PASSWORD")
		database := os.Getenv("MYSQL_DBNAME")
		port := os.Getenv("MYSQL_PORT")

		dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + "?parseTime=true"
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Failed to connect to MySQL: %v", err)
		}
		defer db.Close()

		// Create test database
		testDBName := database + "_test"
		_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + testDBName)
		if err != nil {
			log.Fatalf("Failed to create test database: %v", err)
		}

		defer func() {
			_, err := db.Exec("DROP DATABASE IF EXISTS " + testDBName)
			if err != nil {
				log.Printf("Failed to drop test database: %v", err)
			} else {
				log.Printf("Dropped test database: %v", testDBName)
			}
		}()

		// Run migrations on the test database

		newDsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + testDBName + "?multiStatements=true"
		newDB, err := sql.Open("mysql", newDsn)
		if err != nil {
			log.Fatalf("Failed to connect to MySQL: %v", err)
		}
		defer newDB.Close()

		migrationsDir := "file:///app/migrations"
		driver, err := mysql.WithInstance(newDB, &mysql.Config{})
		if err != nil {
			log.Fatalf("Failed to create MySQL driver for migrations: %v", err)
		}
		migrationInstance, err := migrate.NewWithDatabaseInstance(migrationsDir, testDBName, driver)
		if err != nil {
			log.Fatalf("Failed to create migration instance: %v", err)
		}
		err = migrationInstance.Up()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run migrations: %v", err)
		}

		// Create store instance for tests
		testStore, err = NewMySQLStore(user + ":" + password + "@tcp(" + host + ":" + port + ")/" + testDBName + "?parseTime=true")
		if err != nil {
			log.Fatalf("Failed to create MySQL store: %v", err)
		}

		ctx = context.Background()
		return m.Run()
	}

	code := runTests() // The deferred functions will run when this function returns.
	os.Exit(code)
}

func createWorkflow() *WorkflowInfo {

	createdAt := time.Now().UTC()
	assignedAt := time.Now().UTC()
	completedAt := time.Now().UTC()

	newWorkflow := WorkflowInfo{
		Name:                  "Test WorkflowInfo",
		Type:                  "Sequential",
		RAM:                   "16GB",
		Core:                  "4",
		Policy:                "First-Come-First-Serve",
		ExpectedExecutionTime: sql.NullString{String: "2 hours", Valid: true},
		ActualExecutionTime:   sql.NullString{String: "1 hour", Valid: true},
		AssignedVM:            sql.NullString{String: "VM1", Valid: true},
		SubmittedBy:           sql.NullString{String: "TestUser", Valid: true},
		Status:                "pending",
		CreatedAt:             sql.NullTime{Time: createdAt, Valid: true},
		AssignedAt:            sql.NullTime{Time: assignedAt, Valid: true},
		CompletedAt:           sql.NullTime{Time: completedAt, Valid: true},
	}

	createdWorkflow, err := testStore.SaveWorkflow(ctx, &newWorkflow)
	if err != nil {
		panic(err)
	}

	return createdWorkflow
}

func deleteWorkflow(workflow *WorkflowInfo) {
	// Clean up the workflow from the database. This may need to be more complex in a real application.
	_, err := testStore.db.ExecContext(ctx, "DELETE FROM workflow_info WHERE id = ?", workflow.ID)
	if err != nil {
		panic(err)
	}
}

func TestSaveWorkflow(t *testing.T) {
	t.Parallel()
	newWorkflow := createWorkflow()
	defer deleteWorkflow(newWorkflow)

	assert.NotEqual(t, 0, newWorkflow.ID)
}

func TestGetWorkflowByID(t *testing.T) {
	t.Parallel()
	newWorkflow := createWorkflow()
	defer deleteWorkflow(newWorkflow)

	workflowByID, err := testStore.GetWorkflowByID(ctx, newWorkflow.ID)

	assert.NoError(t, err)
	assert.Equal(t, newWorkflow.ID, workflowByID.ID)
}

func TestUpdateWorkflow(t *testing.T) {
	t.Parallel()
	newWorkflow := createWorkflow()
	defer deleteWorkflow(newWorkflow)

	newWorkflow.Policy = "Round-Robin"
	_, err := testStore.UpdateWorkflow(ctx, newWorkflow)

	assert.NoError(t, err)

	updatedWorkflow, err := testStore.GetWorkflowByID(ctx, newWorkflow.ID)

	assert.NoError(t, err)
	assert.Equal(t, "Round-Robin", updatedWorkflow.Policy)
}

func TestStartWorkflow(t *testing.T) {
	t.Parallel()
	newWorkflow := createWorkflow()
	defer deleteWorkflow(newWorkflow)

	err := testStore.StartWorkflow(ctx, newWorkflow.ID)

	assert.NoError(t, err)

	startedWorkflow, err := testStore.GetWorkflowByID(ctx, newWorkflow.ID)

	assert.NoError(t, err)
	assert.WithinDuration(t, time.Now(), startedWorkflow.AssignedAt.Time, time.Second)
}

func TestCompleteWorkflow(t *testing.T) {
	t.Parallel()
	newWorkflow := createWorkflow()
	defer deleteWorkflow(newWorkflow)

	err := testStore.CompleteWorkflow(ctx, newWorkflow.ID)

	assert.NoError(t, err)

	completedWorkflow, err := testStore.GetWorkflowByID(ctx, newWorkflow.ID)

	assert.NoError(t, err)
	assert.WithinDuration(t, time.Now(), completedWorkflow.CompletedAt.Time, time.Second)
}

// Add this function to create multiple workflows
func createWorkflows(num int) []*WorkflowInfo {
	workflows := make([]*WorkflowInfo, num)
	for i := 0; i < num; i++ {
		workflows[i] = createWorkflow()
	}
	return workflows
}

// Adjust this function to delete all workflows
func deleteAllWorkflows() {
	// Clean up all workflows from the database. This may need to be more complex in a real application.
	_, err := testStore.db.ExecContext(ctx, "DELETE FROM workflow_info")
	if err != nil {
		panic(err)
	}
}

// Add this function to test the GetWorkflows function
func TestGetWorkflows(t *testing.T) {

	// Create 5 workflows in the database
	workflows := createWorkflows(5)
	defer deleteAllWorkflows()

	// Get the workflows from the database
	dbWorkflows, err := testStore.GetWorkflows(ctx)
	if err != nil {
		t.Fatalf("GetWorkflows failed with error: %s", err)
	}

	// Check the number of workflows returned
	if len(dbWorkflows) != 5 {
		t.Fatalf("expected 5 workflows, but got %d", len(dbWorkflows))
	}

	// Check the workflows returned
	for i, wf := range dbWorkflows {
		if wf.ID != workflows[i].ID {
			t.Errorf("expected workflow ID %d, but got %d", workflows[i].ID, wf.ID)
		}
	}
}

func TestAssignWorkflow(t *testing.T) {
	t.Parallel()

	// 1. Create a workflow in your test database.
	newWorkflow := createWorkflow()
	defer deleteWorkflow(newWorkflow)

	venName := "TestVEN"

	// 2. Assign the workflow using the AssignWorkflow function.
	err := testStore.AssignWorkflow(ctx, newWorkflow.Name, venName, 0)
	assert.NoError(t, err)

	// 3. Retrieve the workflow from the database and check updates.
	assignedWorkflow, err := testStore.GetWorkflowByID(ctx, newWorkflow.ID)
	assert.NoError(t, err)

	assert.Equal(t, "assigned", assignedWorkflow.Status)
	assert.Equal(t, venName, assignedWorkflow.AssignedVM.String)

	// Ensure the `AssignedAt` time is close to the current time, within a tolerance of 1 second.
	assert.WithinDuration(t, time.Now(), assignedWorkflow.AssignedAt.Time, time.Second)
}
