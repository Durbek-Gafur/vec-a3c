package queue

import (
	"context"
	s "vec-node/internal/store"
	wf "vec-node/internal/workflow"
)

//go:generate mockgen -destination=mocks/queue_mock.go -package=queueu_mock vec-node/internal/queue Queue

// Queue handles operations on queue
type Queue interface {
	Enqueue(ctx context.Context, workflowID int) (int, error)
	Peek(ctx context.Context) (*s.Queue, error)
	GetQueueStatus(ctx context.Context) ([]s.Queue, error)
	ProcessWorkflowInQueue(ctx context.Context, workflowID int) error
	CompleteWorkflowInQueue(ctx context.Context, id int) error
}


type Service struct {
	queueStore s.QueueStore
	workflow wf.Workflow
}


func NewService(store s.QueueStore, workflow wf.Workflow) *Service {
	return &Service{
		queueStore: store,
		workflow: workflow,
	}
}

func (s *Service) Enqueue(ctx context.Context, workflowID int) (int, error) {
	return s.queueStore.Enqueue(ctx, workflowID)
}

func (s *Service) Peek(ctx context.Context) (*s.Queue, error) {
	return s.queueStore.Peek(ctx)
}

func (s *Service) GetQueueStatus(ctx context.Context) ([]s.Queue, error) {
	return s.queueStore.GetQueueStatus(ctx)
}

func (s *Service) ProcessWorkflowInQueue(ctx context.Context, workflowID int) error {
	err := s.workflow.StartExecution(ctx,workflowID)
	if err!=nil{
		return err
	}
	return s.queueStore.ProcessWorkflowInQueue(ctx, workflowID)
}

func (s *Service) CompleteWorkflowInQueue(ctx context.Context, id int) error {
	err := s.workflow.Complete(ctx,id)
	if err!=nil{
		return err
	}
	return s.queueStore.CompleteWorkflowInQueue(ctx, id)
}

