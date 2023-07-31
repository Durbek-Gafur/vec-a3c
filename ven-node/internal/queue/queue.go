package queue

import (
	"context"
	"fmt"
	"log"
	"time"

	s "vec-node/internal/store"
	wf "vec-node/internal/workflow"

	"github.com/pkg/errors"
)

//go:generate mockgen -destination=mocks/queue_mock.go -package=queueu_mock vec-node/internal/queue Queue

// Queue handles operations on queue
type Queue interface {
	// Enqueue(ctx context.Context, workflowID int) (int, error)
	// Peek(ctx context.Context) (*s.Queue, error)
	// GetQueueStatus(ctx context.Context) ([]s.Queue, error)
	// ProcessWorkflowInQueue(ctx context.Context, workflowID int) error
	// CompleteWorkflowInQueue(ctx context.Context, id int) error
	StartPeriodicCheck(ctx context.Context)
}

type Service struct {
	queueStore s.QueueStore
	workflow   wf.Workflow
}

func NewService(store s.QueueStore, workflow wf.Workflow) *Service {
	return &Service{
		queueStore: store,
		workflow:   workflow,
	}
}

// StartPeriodicCheck starts a periodic check every 2 minutes
func (s *Service) StartPeriodicCheck(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				err := s.runCheckAndExecute(ctx)
				log.Printf("Starting periodic check.")
				if err != nil {
					fmt.Printf("An error occurred during periodic check: %v\n", err)
					// Handle error appropriately
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

// runCheckAndExecute checks for any running workflows, if none, it peeks a workflow and executes it
func (s *Service) runCheckAndExecute(ctx context.Context) error {
	// Check if queue is empty
	isEmpty, err := s.queueStore.IsEmpty(ctx)
	if err != nil {
		return errors.Wrap(err, "runCheckAndExecute: error occurred while checking if queue is empty")
	}

	// If queue is empty, return early
	if isEmpty {
		log.Printf("Queue is empty.")
		return nil
	}

	// Check if any workflow is currently running
	runningWorkflow, err := s.queueStore.PeekInProcess(ctx)
	if err != nil {
		return errors.Wrap(err, "runCheckAndExecute: error occurred while checking running workflow")
	}

	// If no workflow is running, peek the next workflow and start execution
	if runningWorkflow == nil {
		log.Printf("No running workflow, starting a new one.")
		return s.peekAndExecute(ctx)
	}

	// Check if the running workflow has completed its execution
	isComplete, err := s.workflow.IsComplete(ctx, runningWorkflow.WorkflowID)
	if err != nil {
		return errors.Wrap(err, "runCheckAndExecute: error occurred while checking if workflow is complete")
	}

	// If the running workflow has completed, mark it as complete and start the next workflow
	if isComplete {
		log.Printf("Running workflow has completed, marking as complete and starting a new one.")
		if err := s.queueStore.CompleteWorkflowInQueue(ctx, runningWorkflow.WorkflowID); err != nil {
			return errors.Wrap(err, "runCheckAndExecute: error occurred while marking workflow as complete in queue")
		}

		duration, err := s.workflow.GetScriptDuration()
		if err != nil {
			return errors.Wrap(err, "runCheckAndExecute: error occurred while checking duration")
		}

		err = s.workflow.Complete(ctx, runningWorkflow.ID, duration)
		if err != nil {
			return errors.Wrap(err, "runCheckAndExecute: error occurred while Complete workflow")
		}

		// Informing master node about completed workflow, this part might be different based on your implementation
		// if err := s.InformMasterNode(ctx, runningWorkflow.WorkflowID); err != nil {
		// 	return errors.Wrap(err, "runCheckAndExecute: error occurred while informing master node about completed workflow")
		// }

		return s.peekAndExecute(ctx)
	}

	log.Printf("Workflow is still running.")
	return nil
}

// peekAndExecute peeks the next workflow and starts execution
func (s *Service) peekAndExecute(ctx context.Context) error {
	nextWorkflow, err := s.queueStore.PeekQueued(ctx)
	if err != nil {
		return errors.Wrap(err, "peekAndExecute: error occurred while trying to peek next workflow")
	}
	if nextWorkflow != nil {
		// StartExecution
		log.Printf("starting wf %d", int(nextWorkflow.ID))
		err := s.workflow.StartExecution(ctx, nextWorkflow.ID)
		if err != nil {
			log.Println(err.Error())
		}
		err = s.queueStore.ProcessWorkflowInQueue(ctx, nextWorkflow.WorkflowID)
		if err != nil {
			return errors.Wrap(err, "peekAndExecute: error occurred while trying to process next workflow")
		}
	}
	return nil
}
