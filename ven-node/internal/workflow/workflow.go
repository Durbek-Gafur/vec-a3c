package workflow

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	s "vec-node/internal/store"
)

//go:generate mockgen -destination=mocks/workflow_mock.go -package=workflow_mock vec-node/internal/workflow Workflow

// Workflow handles operations on Workflow sizes
type Workflow interface {
	StartExecution(ctx context.Context, id int) error
	Complete(ctx context.Context, id int, duration int) error
	UpdateWorkflow(ctx context.Context, wf *s.Workflow) (*s.Workflow, error)
	IsComplete(ctx context.Context, id int) (bool, error)
}

// NewWorkflow returns a new Workflow instance
func NewWorkflow(name, wType string, duration int) *s.Workflow {
	return &s.Workflow{
		Name:       name,
		Type:       wType,
		Duration:   duration,
		ReceivedAt: time.Now(),
	}
}

type Service struct {
	workflowStore s.WorkflowStore
	cmd           *exec.Cmd
	logFile       *os.File
}

func NewService(store s.WorkflowStore, logFile *os.File) *Service {
	return &Service{
		workflowStore: store,
		logFile:       logFile,
	}
}

func (s *Service) StartExecution(ctx context.Context, workflowID int) error {
	// Create a new context with a 5 minute timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	log.Println("Preparing to execute the script...")

	// Prepare to execute the script
	s.cmd = exec.Command("bash", "-c", "/app/workflow/rna.sh")

	// Create a pipe for stdout and stderr, and wrap it with a log writer
	s.cmd.Stdout = s.logFile
	s.cmd.Stderr = s.logFile

	log.Println("Starting the script in the background...")

	// Run the script in the background
	err := s.cmd.Start()
	if err != nil {
		log.Printf("Failed to start the script: %v", err)
		return fmt.Errorf("failed to start the script: %w", err)
	}

	log.Println("Script started successfully. Setting up a ticker to monitor the script...")

	// Set up a ticker to check if the script has finished every 10 seconds
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		err := s.cmd.Wait()
		if err != nil {
			s.logFile.WriteString(fmt.Sprintf("script execution failed: %v", err))
		}

		// ticker.Stop() // If you want to stop the ticker here
	}()

	// Goroutine to monitor the script
	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// if err := s.cmd.Process.Signal(syscall.Signal(0)); err != nil {
				// 	log.Println("Script has finished executing.")
				// 	duration := s.cmd.ProcessState.UserTime()
				// 	log.Printf("Duration %s", duration)
				// 	err := s.Complete(ctx, workflowID, int(duration.Seconds()))
				// 	if err != nil {
				// 		log.Printf("Failed to complete the workflow: %v", err)
				// 	}
				// 	return
				// }

				if s.cmd.ProcessState != nil && s.cmd.ProcessState.Exited() {
					log.Println("Script has finished executing.")
					duration := s.cmd.ProcessState.UserTime()
					log.Printf("Duration %s", duration)
					err := s.Complete(ctx, workflowID, int(duration.Seconds()))
					if err != nil {
						log.Printf("Failed to complete the workflow: %v", err)
					}
					return
				}
				// TODO context is cancelled and the request dies
				// the call above doesn't work because ctx has died.

				// GET duration by reading output_workflow_id.txt or think of smth else

				// case <-ctx.Done():
				// 	log.Println("Context has been cancelled, stopping the ticker...")
				// 	ticker.Stop()
				// 	return

				// default:
				// 	log.Println("waiting ... ")
				// 	continue
			}
		}
	}()

	log.Println("Script is now being monitored. Execution continues in the background.")

	return nil
}

func (s *Service) Complete(ctx context.Context, id int, duration int) error {
	// Update the workflow with the duration
	wf, err := s.workflowStore.GetWorkflowByID(ctx, id)
	if err != nil {
		return err
	}

	wf.Duration = duration
	_, err = s.UpdateWorkflow(ctx, wf)
	if err != nil {
		return err
	}

	return s.workflowStore.CompleteWorkflow(ctx, id)
}

func (s *Service) UpdateWorkflow(ctx context.Context, wf *s.Workflow) (*s.Workflow, error) {
	return s.workflowStore.UpdateWorkflow(ctx, wf)
}

func (s *Service) IsComplete(ctx context.Context, id int) (bool, error) {
	// not implemented
	return false, nil
}
