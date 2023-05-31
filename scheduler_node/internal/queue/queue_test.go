package queue

import (
	"context"
	"errors"
	"testing"
	"time"
	"scheduler-node/internal/store"
	store_mock "scheduler-node/internal/store/mocks"
	workflow_mock "scheduler-node/internal/workflow/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)



func TestEnqueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store_mock.NewMockQueueStore(ctrl)
	mockWorkflow := workflow_mock.NewMockWorkflow(ctrl)
	service := NewService(mockStore, mockWorkflow)

	ctx := context.TODO()
	workflowID := 123
	queueID := 456
	mockStore.EXPECT().Enqueue(ctx, workflowID).Return(queueID, nil)

	id, err := service.Enqueue(ctx, workflowID)
	assert.NoError(t, err)
	assert.Equal(t, queueID, id)
}

func TestPeek(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store_mock.NewMockQueueStore(ctrl)
	mockWorkflow := workflow_mock.NewMockWorkflow(ctrl)
	service := NewService(mockStore, mockWorkflow)

	ctx := context.TODO()
	expectedQueue := &store.Queue{
		ID:             123,
		WorkflowID:     456,
		Status:         "pending",
		EnqueuedAt:     time.Now(),
	}
	mockStore.EXPECT().Peek(ctx).Return(expectedQueue, nil)

	queue, err := service.Peek(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedQueue, queue)
}

func TestGetQueueStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store_mock.NewMockQueueStore(ctrl)
	mockWorkflow := workflow_mock.NewMockWorkflow(ctrl)
	service := NewService(mockStore, mockWorkflow)

	ctx := context.TODO()
	expectedQueues := []store.Queue{
		{ID: 123, WorkflowID: 456, Status: "pending", EnqueuedAt: time.Now()},
		{ID: 789, WorkflowID: 987, Status: "processing", EnqueuedAt: time.Now()},
	}
	mockStore.EXPECT().GetQueueStatus(ctx).Return(expectedQueues, nil)

	queues, err := service.GetQueueStatus(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedQueues, queues)
}

func TestProcessWorkflowInQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store_mock.NewMockQueueStore(ctrl)
	mockWorkflow := workflow_mock.NewMockWorkflow(ctrl)
	service := NewService(mockStore, mockWorkflow)

	ctx := context.TODO()
	workflowID := 456

	mockWorkflow.EXPECT().StartExecution(ctx, workflowID).Return(nil)
	mockStore.EXPECT().ProcessWorkflowInQueue(ctx, workflowID).Return(nil)

	err := service.ProcessWorkflowInQueue(ctx, workflowID)
	assert.NoError(t, err)
}

func TestProcessWorkflowInQueueWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store_mock.NewMockQueueStore(ctrl)
	mockWorkflow := workflow_mock.NewMockWorkflow(ctrl)
	service := NewService(mockStore, mockWorkflow)

	ctx := context.TODO()
	workflowID := 456
	expectedError := errors.New("failed to start execution")

	mockWorkflow.EXPECT().StartExecution(ctx, workflowID).Return(expectedError)
	mockStore.EXPECT().ProcessWorkflowInQueue(gomock.Any(), gomock.Any()).Times(0)

	err := service.ProcessWorkflowInQueue(ctx, workflowID)
	assert.Equal(t, expectedError, err)
}


func TestCompleteWorkflowInQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store_mock.NewMockQueueStore(ctrl)
	mockWorkflow := workflow_mock.NewMockWorkflow(ctrl)
	service := NewService(mockStore, mockWorkflow)

	ctx := context.TODO()
	workflowID := 456

	mockWorkflow.EXPECT().Complete(ctx, workflowID).Return(nil)
	mockStore.EXPECT().CompleteWorkflowInQueue(ctx, workflowID).Return(nil)

	err := service.CompleteWorkflowInQueue(ctx, workflowID)
	assert.NoError(t, err)
}
