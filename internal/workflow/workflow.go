package workflow

import (
	"context"
	"time"
	s "vec-node/internal/store"
)

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

func (s *Service) StartExecution(ctx context.Context, id int) error {
	return s.workflowStore.StartWorkflow(ctx, id)
}

func (s *Service) Complete(ctx context.Context, id int) error {
	return s.workflowStore.CompleteWorkflow(ctx, id)
}

func (s *Service) UpdateWorkflow(ctx context.Context, wf *s.Workflow) (*s.Workflow,error) {
	return s.workflowStore.UpdateWorkflow(ctx, wf)
}
