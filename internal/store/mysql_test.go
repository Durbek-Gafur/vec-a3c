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
)

var testStore *MySQLStore
var ctx context.Context

func TestMain(m *testing.M) {
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

	
	log.Println(database)
    // Run migrations on the test database

	newDsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" +testDBName+ "?multiStatements=true"
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
    code := m.Run()
    os.Exit(code)
}


func TestQueueSize(t *testing.T) {
	// Test SetQueueSize
	err := testStore.SetQueueSize(ctx, 5)
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
	err = testStore.UpdateQueueSize(ctx, 10)
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
}
