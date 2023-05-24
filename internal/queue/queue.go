package queue

import (
	"context"
	s "vec-node/internal/store"
	wf "vec-node/internal/workflow"
)



type Service struct {
	store s.Store
	workflow wf.Service
}


func NewService(store s.Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) Enqueue(ctx context.Context, workflowID int) (int, error) {
	return s.store.Enqueue(ctx, workflowID)
}

func (s *Service) Dequeue(ctx context.Context) (*s.Queue, error) {
	return s.store.Dequeue(ctx)
}

func (s *Service) GetQueueStatus(ctx context.Context) ([]s.Queue, error) {
	return s.store.GetQueueStatus(ctx)
}

func (s *Service) ProcessWorkflowInQueue(ctx context.Context, id int, status string) error {
	err := s.workflow.StartExecution(ctx,id)
	if err!=nil{
		return err
	}
	return s.store.ProcessWorkflowInQueue(ctx, id, status)
}

func (s *Service) CompleteWorkflowInQueue(ctx context.Context, id int, status string) error {
	err := s.workflow.Complete(ctx,id)
	if err!=nil{
		return err
	}
	return s.store.CompleteWorkflowInQueue(ctx, id, status)
}

// TODO write unit test for these