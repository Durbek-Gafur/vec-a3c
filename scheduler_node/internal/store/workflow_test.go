package store

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

	newWorkflow := WorkflowInfo{
		Name:                  "Test WorkflowInfo",
		Type:                  "Sequential",
		RAM:                   "16GB",
		Core:                  "4",
		Policy:                "First-Come-First-Serve",
		ExpectedExecutionTime: "2 hours",
		ActualExecutionTime:   "1 hour",
		AssignedVM:            "VM1",
		SubmittedBy:           "TestUser",
		Status:                "pending",
		AssignedAt:            assignedAt,
		CompletedAt:           completedAt,
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
	assert.WithinDuration(t, time.Now(), startedWorkflow.AssignedAt, time.Second)
}

func TestCompleteWorkflow(t *testing.T) {
	t.Parallel()
	newWorkflow := createWorkflow()
	defer deleteWorkflow(newWorkflow)

	err := testStore.CompleteWorkflow(ctx, newWorkflow.ID)

	assert.NoError(t, err)

	completedWorkflow, err := testStore.GetWorkflowByID(ctx, newWorkflow.ID)

	assert.NoError(t, err)
	assert.WithinDuration(t, time.Now(), completedWorkflow.CompletedAt, time.Second)
}

func TestPeek(t *testing.T) {

	// Creating multiple workflows.
	newWorkflow1 := createWorkflow()
	defer deleteWorkflow(newWorkflow1)

	newWorkflow2 := createWorkflow()
	defer deleteWorkflow(newWorkflow2)

	newWorkflow3 := createWorkflow()
	defer deleteWorkflow(newWorkflow3)
	// Completing a workflow.
	err := testStore.CompleteWorkflow(ctx, newWorkflow1.ID)
	assert.NoError(t, err)

	// We expect newWorkflow2 to be the earliest non-completed workflow.
	workflow, err := testStore.Peek(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, workflow)
	assert.Equal(t, newWorkflow2.ID, workflow.ID)
}

func TestGetQueue(t *testing.T) {
	t.Parallel()

	// Creating multiple workflows.
	newWorkflow1 := createWorkflow()
	defer deleteWorkflow(newWorkflow1)

	newWorkflow2 := createWorkflow()
	defer deleteWorkflow(newWorkflow2)

	newWorkflow3 := createWorkflow()
	defer deleteWorkflow(newWorkflow3)

	// Completing a workflow.
	err := testStore.CompleteWorkflow(ctx, newWorkflow1.ID)
	assert.NoError(t, err)

	// Retrieving the queue.
	queue, err := testStore.GetQueue(ctx)

	assert.NoError(t, err)
	assert.NotEmpty(t, queue)

	// Verifying that the completed workflow is not in the queue.
	for _, workflow := range queue {
		assert.NotEqual(t, newWorkflow1.ID, workflow.ID)
	}
}

