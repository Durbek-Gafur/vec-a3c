package scheduler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"scheduler-node/internal/store"
	"time"
)

//go:generate mockgen -destination=mocks/scheduler_mock.go -package=scheduler_mock scheduler-node/internal/scheduler Scheduler

// Workflow handles operations on Workflow sizes
type Scheduler interface {
	SubmitWorkflow(ctx context.Context, ven store.VENInfo, workflow store.WorkflowInfo) error
}


type SchedulerService struct{
	queueStore store.QueueStore
}

func NewSchedulerService(queueStore store.QueueStore) *SchedulerService {
	return &SchedulerService{
		queueStore: queueStore,
	}
}





type WorkflowToBeSent struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Duration int    `json:"duration"`
}

func (s *SchedulerService) SubmitWorkflow(ctx context.Context, ven store.VENInfo, workflow store.WorkflowInfo) error {
	// Prepare the workflow to submit
	submitWorkflow := WorkflowToBeSent{
		Name:     workflow.Name,
		Type:     workflow.Type,
		Duration: 0,
	}

	// If ExpectedExecutionTime is a valid time string, convert it to seconds as duration
	if workflow.ExpectedExecutionTime.Valid {
		duration, err := time.ParseDuration(workflow.ExpectedExecutionTime.String)
		if err == nil {
			submitWorkflow.Duration = int(duration.Seconds())
		}
	}

	// Marshal the workflow into JSON
	body, err := json.Marshal(submitWorkflow)
	if err != nil {
		return err
	}

	// Send the POST request
	resp, err := http.Post(ven.URL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	// If the response is OK, enqueue the workflow
	err = s.queueStore.AssignWorkflow(ctx, workflow.ID,ven.Name)
	if err != nil {
		return err
	}

	return nil
}
