package workflow

import (
	"context"
	"time"
	s "vec-node/internal/store"
)

// NewWorkflow returns a new Workflow instance
func NewWorkflow(name, wType string, duration int64) *s.Workflow {
	return &s.Workflow{
		Name:       name,
		Type:       wType,
		Duration:   duration,
		ReceivedAt: time.Now(),
	}
}

type Service struct {
	store s.Store
}

func NewService(store s.Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) StartExecution(ctx context.Context, id int) error {
	return s.store.StartWorkflow(ctx, id)
}

func (s *Service) Complete(ctx context.Context, id int) error {
	return s.store.CompleteWorkflow(ctx, id)
}

func (s *Service) UpdateWorkflow(ctx context.Context, id int, wf *s.Workflow) (*s.Workflow,error) {
	return s.store.UpdateWorkflow(ctx, wf)
}
