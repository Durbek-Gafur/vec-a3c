package scheduler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"scheduler-node/internal/store"
	"time"
)

//go:generate mockgen -destination=mocks/scheduler_mock.go -package=scheduler_mock scheduler-node/internal/scheduler Scheduler

// Workflow handles operations on Workflow sizes
type Scheduler interface {
	SubmitWorkflow(ctx context.Context, ven store.VENInfo, workflow store.WorkflowInfo) error
}

type SchedulerService struct {
	queueStore store.QueueStore
	venStore   store.VENStore
}

func NewSchedulerService(queueStore store.QueueStore, venStore store.VENStore) *SchedulerService {
	return &SchedulerService{
		queueStore: queueStore,
		venStore:   venStore,
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
		Name: workflow.Name,
		Type: workflow.Type,
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
	resp, err := http.Post(ven.URL+"/workflow", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	// If the response is OK, enqueue the workflow
	err = s.queueStore.AssignWorkflow(ctx, workflow.ID, ven.Name)
	if err != nil {
		return err
	}

	return nil
}

const scanInterval = 10 * time.Second // adjust as necessary

func (s *SchedulerService) StartScheduling(ctx context.Context) {
	ticker := time.NewTicker(scanInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := s.processUnscheduledWorkflows(ctx)
			if err != nil {
				log.Printf("Failed to processUnscheduledWorkflows: %v", err)
				return
			}
		}
	}
}

func (s *SchedulerService) processUnscheduledWorkflows(ctx context.Context) error {
	// Assume you've methods in your store to get unscheduled workflows and available VENs
	workflows, err := s.queueStore.GetPendingQueue(ctx)
	if err != nil {
		log.Printf("Failed to GetPendingQueue: %v", err)
		return err
	}

	vens, err := s.venStore.GetAvailableVEN()
	if err != nil {
		log.Printf("Failed to GetVENInfos: %v", err)
		return err
	}

	for _, wf := range workflows {
		for _, ven := range vens {
			if err := s.SubmitWorkflow(ctx, ven, wf); err != nil {
				// log error and perhaps continue or break based on logic
			} else {
				// Successfully scheduled, move to next workflow
				break
			}
		}
	}
	return nil
}
