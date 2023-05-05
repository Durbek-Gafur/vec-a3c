package store

import (
	"context"
	"errors"
	w "vec-node/internal/workflow"
)

// Store is an interface that defines the methods for interacting with the data store.
type Store interface {
	GetQueueSize(ctx context.Context) (int, error)
	SetQueueSize(ctx context.Context, size int) error
	UpdateQueueSize(ctx context.Context, size int) error

	// Workflow related methods
	GetWorkflowByID(ctx context.Context,id int64) (*w.Workflow, error)
	GetWorkflows(ctx context.Context,filter *w.WorkflowFilter) ([]w.Workflow, error)
	SaveWorkflow(ctx context.Context,workflow *w.Workflow) (*w.Workflow, error)
	UpdateWorkflow(ctx context.Context, w *w.Workflow) (*w.Workflow, error)  

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

