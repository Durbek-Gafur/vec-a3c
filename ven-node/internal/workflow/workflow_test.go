package workflow

import (
	"context"
	"testing"
	"time"
	s "vec-node/internal/store"
	store_mock "vec-node/internal/store/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// func TestStartExecution(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockStore := store_mock.NewMockWorkflowStore(ctrl)
// 	service := NewService(mockStore, nil)

// 	ctx := context.TODO()
// 	mockStore.EXPECT().StartWorkflow(ctx, 123).Return(nil)

// 	err := service.StartExecution(ctx, 123)
// 	assert.NoError(t, err)
// }

func TestComplete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store_mock.NewMockWorkflowStore(ctrl)
	service := NewService(mockStore, nil)

	ctx := context.TODO()
	mockStore.EXPECT().GetWorkflowByID(ctx, 456).Return(&s.Workflow{ID: 456}, nil)
	mockStore.EXPECT().UpdateWorkflow(ctx, gomock.Any()).Return(nil, nil)
	mockStore.EXPECT().CompleteWorkflow(ctx, 456).Return(nil)

	err := service.Complete(ctx, 456, 4)
	assert.NoError(t, err)
}

func TestUpdateWorkflow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store_mock.NewMockWorkflowStore(ctrl)
	service := NewService(mockStore, nil)

	ctx := context.TODO()
	wf := &s.Workflow{
		Name:       "test",
		Type:       "type1",
		Duration:   1,
		ReceivedAt: time.Now(),
	}
	mockStore.EXPECT().UpdateWorkflow(ctx, wf).Return(wf, nil)

	updatedWf, err := service.UpdateWorkflow(ctx, wf)
	assert.NoError(t, err)
	assert.Equal(t, wf, updatedWf)
}
