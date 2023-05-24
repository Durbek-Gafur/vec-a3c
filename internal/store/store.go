package store

import (
	"context"
	"errors"
)

// Store is an interface that defines the methods for interacting with the data store.
type Store interface {
	GetQueueSize(ctx context.Context) (int, error)
	SetQueueSize(ctx context.Context, size int) error
	UpdateQueueSize(ctx context.Context, size int) error

	// Workflow related methods
	GetWorkflowByID(ctx context.Context,id int) (*Workflow, error)
	GetWorkflows(ctx context.Context,filter *WorkflowFilter) ([]Workflow, error)
	SaveWorkflow(ctx context.Context,workflow *Workflow) (*Workflow, error)
	UpdateWorkflow(ctx context.Context, w *Workflow) (*Workflow, error) 
	StartWorkflow(ctx context.Context, id int) error
	CompleteWorkflow(ctx context.Context, id int) error 

	// Queue related methods
	Enqueue(ctx context.Context, workflowID int) (int, error)
	Dequeue(ctx context.Context) (*Queue, error)
	GetQueueStatus(ctx context.Context) ([]Queue, error)
	ProcessWorkflowInQueue(ctx context.Context, id int) error
	CompleteWorkflowInQueue(ctx context.Context, id int) error
	IsSpaceAvailable(ctx context.Context) (bool, error) 

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

