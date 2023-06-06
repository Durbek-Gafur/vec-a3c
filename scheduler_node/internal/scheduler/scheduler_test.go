package scheduler

import (
	"context"
	"database/sql"
	"net/http"
	"scheduler-node/internal/store"
	store_mock "scheduler-node/internal/store/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestSchedulerService_SubmitWorkflow(t *testing.T) {
	// Create a new instance of the mock object
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Initialize the mock QueueStore
	mockQueueStore := store_mock.NewMockQueueStore(ctrl)

	// Create SchedulerService with the mocked QueueStore
	s := NewSchedulerService(mockQueueStore)

	ctx := context.TODO()
	ven := store.VENInfo{URL: "http://durbek.ga"}

	workflow := store.WorkflowInfo{
		ID:   1,
		Name: "TestWorkflow",
		Type: "typeA",
		ExpectedExecutionTime: sql.NullString{
			String: "10s",
			Valid:  true,
		},
	}

	// Mock the AssignWorkflow call
	mockQueueStore.EXPECT().
		AssignWorkflow(ctx, workflow.ID, ven.Name).
		Return(nil)

	// Activate the httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock the HTTP POST request
	httpmock.RegisterResponder("POST", ven.URL,
		httpmock.NewStringResponder(http.StatusOK, ""))

	// Call the method under test
	err := s.SubmitWorkflow(ctx, ven, workflow)

	// Assert
	assert.NoError(t, err)
}
