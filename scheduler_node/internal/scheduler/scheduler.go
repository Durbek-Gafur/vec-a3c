package scheduler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"scheduler-node/internal/store"
	"strings"
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
	wfStore    store.WorkflowStore
}

func NewSchedulerService(queueStore store.QueueStore, venStore store.VENStore, wfStore store.WorkflowStore) *SchedulerService {
	return &SchedulerService{
		queueStore: queueStore,
		venStore:   venStore,
		wfStore:    wfStore,
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
	if resp.StatusCode == http.StatusUnavailableForLegalReasons {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	// Check the response status
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	// If the response is OK, enqueue the workflow
	err = s.wfStore.AssignWorkflow(ctx, workflow.Name, ven.Name)
	if err != nil {
		log.Printf("Failed to AssignWorkflow: %v", err)
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
		for i, ven := range vens {
			// TODO add logic if queue size is not 0 or less than 1
			if ven.CurrentQueueSize == "0" {
				break
			}
			if err := s.SubmitWorkflow(ctx, ven, wf); err != nil {
				// log error and perhaps continue or break based on logic
				if strings.Contains(err.Error(), fmt.Sprint(http.StatusUnavailableForLegalReasons)) {
					log.Printf("Queue is full for ven %s", ven.Name)
					vens[i].CurrentQueueSize = "0"
					break
				}
				log.Printf("Failed to submit wf %s with ven %s: %v", wf.Name, ven.Name, err)
				continue
			} else {
				// Successfully scheduled, move to next workflow
				log.Printf("Successfully scheduled wf %s with ven %s", wf.Name, ven.Name)
				break
			}
		}
	}
	log.Printf("Successfully completed schedule cycle")
	return nil
}
