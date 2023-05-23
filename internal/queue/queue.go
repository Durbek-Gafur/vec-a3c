package queue

import (
	"context"
	s "vec-node/internal/store"
)



type Service struct {
	store s.Store
}


func NewService(store s.Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) Enqueue(ctx context.Context, workflowID int64) (int64, error) {
	return s.store.Enqueue(ctx, workflowID)
}

func (s *Service) Dequeue(ctx context.Context) (*Queue, error) {
	return s.store.Dequeue(ctx)
}

func (s *Service) GetQueueStatus(ctx context.Context) ([]Queue, error) {
	return s.store.GetQueueStatus(ctx)
}

func (s *Service) UpdateStatus(ctx context.Context, id int64, status string) error {
	return s.store.UpdateStatus(ctx, id, status)
}
