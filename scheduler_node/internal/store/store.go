package store

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

//go:generate mockgen -destination=mocks/store_mock.go -package=store_mock scheduler-node/internal/store WorkflowStore,QueueStore,VENStore
// WorkflowStore handles operations on workflows
type WorkflowStore interface {
	GetWorkflowByID(ctx context.Context,id int) (*WorkflowInfo, error)
	GetWorkflows(ctx context.Context) ([]WorkflowInfo, error)
	SaveWorkflow(ctx context.Context,WorkflowInfo *WorkflowInfo) (*WorkflowInfo, error)
	UpdateWorkflow(ctx context.Context, w *WorkflowInfo) (*WorkflowInfo, error) 
	AssignWorkflow(ctx context.Context, workflowID int, venName string) error
	StartWorkflow(ctx context.Context, id int) error
	CompleteWorkflow(ctx context.Context, id int) error 
}

// VENStore handles operations on workflows
type VENStore interface {
	GetVENInfos() ([]VENInfo, error)
}


// QueueStore handles operations on queues
type QueueStore interface {
	Peek(ctx context.Context) (*WorkflowInfo, error) 
	GetQueue(ctx context.Context) ([]WorkflowInfo, error)

	AssignWorkflow(ctx context.Context, workflowID int, venName string) error
	StartWorkflow(ctx context.Context, id int) error
	CompleteWorkflow(ctx context.Context, id int) error 
}


// StoreError is a custom error type for store-related errors. It includes the original error and a status code.
type StoreError struct {
	Err        error
	StatusCode int
}

// Error returns the error message of the wrapped error.
func (e *StoreError) Error() string {
	return e.Err.Error()
}

// Unwrap returns the original wrapped error.
func (e *StoreError) Unwrap() error {
	return e.Err
}

// ErrNotFound is a sentinel error for when a requested item is not found in the store.
var ErrNotFound = errors.New("not found")


type MySQLStore struct {
	db *sql.DB
}

func (m *MySQLStore) GetDB() *sql.DB {
	return m.db 
}



func NewMySQLStore(dsn string) (*MySQLStore, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &MySQLStore{db: db}, nil
}