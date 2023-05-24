package store

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
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

        log.Println(database)
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
