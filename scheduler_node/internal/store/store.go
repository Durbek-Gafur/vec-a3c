package store

import (
	"context"
	"database/sql"
	"errors"
)

type Store interface {
	GetVENInfo() ([]VENInfo, error)
	GetWorkflowInfo() ([]WorkflowInfo, error)
}


//go:generate mockgen -destination=mocks/store_mock.go -package=store_mock scheduler-node/internal/store WorkflowStore,QueueStore,QueueSizeStore
// WorkflowStore handles operations on workflows
type WorkflowStore interface {
	GetWorkflowByID(ctx context.Context,id int) (*WorkflowInfo, error)
	GetWorkflows(ctx context.Context,filter *WorkflowFilter) ([]WorkflowInfo, error)
	SaveWorkflow(ctx context.Context,WorkflowInfo *WorkflowInfo) (*WorkflowInfo, error)
	UpdateWorkflow(ctx context.Context, w *WorkflowInfo) (*WorkflowInfo, error) 
	StartWorkflow(ctx context.Context, id int) error
	CompleteWorkflow(ctx context.Context, id int) error 
}

// QueueStore handles operations on queues
type QueueStore interface {
	Peek(ctx context.Context) (*Queue, error) 
	GetQueue(ctx context.Context) ([]Queue, error)

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

func NewMySQLStore(dsn string) (*MySQLStore, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &MySQLStore{db: db}, nil
}