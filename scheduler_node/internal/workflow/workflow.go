package workflow

/*
import (
	"context"
	s "scheduler-node/internal/store"
	"time"
)

// ggggo:generate mockgen -destination=mocks/workflow_mock.go -package=workflow_mock scheduler-node/internal/workflow Workflow

// Workflow handles operations on Workflow sizes
type Workflow interface {
	StartExecution(ctx context.Context, id int) error
	Complete(ctx context.Context, id int) error
	UpdateWorkflow(ctx context.Context, wf *s.Workflow) (*s.Workflow,error)
}


// NewWorkflow returns a new Workflow instance
func NewWorkflow(name, wType string, duration int) *s.Workflow {
	return &s.Workflow{
		Name:       name,
		Type:       wType,
		Duration:   duration,
		ReceivedAt: time.Now(),
	}
}

type Service struct {
	workflowStore s.WorkflowStore
}

func NewService(store s.WorkflowStore) *Service {
	return &Service{
		workflowStore: store,
	}
}

func (s *Service) StartExecution(ctx context.Context, workflowID int) error {
	return s.workflowStore.StartWorkflow(ctx, workflowID)
}

func (s *Service) Complete(ctx context.Context, id int) error {
	return s.workflowStore.CompleteWorkflow(ctx, id)
}

func (s *Service) UpdateWorkflow(ctx context.Context, wf *s.Workflow) (*s.Workflow,error) {
	return s.workflowStore.UpdateWorkflow(ctx, wf)
}

*/
