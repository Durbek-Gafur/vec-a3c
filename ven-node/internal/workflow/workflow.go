package workflow

import (
	"context"
	"strconv"

	// "errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"

	"vec-node/internal/store"
	s "vec-node/internal/store"
)

//go:generate mockgen -destination=mocks/workflow_mock.go -package=workflow_mock vec-node/internal/workflow Workflow

// Workflow handles operations on Workflow sizes
type Workflow interface {
	StartExecution(ctx context.Context, id int) error
	Complete(ctx context.Context, id int, duration int) error
	UpdateWorkflow(ctx context.Context, wf *s.Workflow) (*s.Workflow, error)
	IsComplete(ctx context.Context, id int) (bool, error)
	IsScriptComplete() (bool, error)
	GetScriptDuration() (int, error)
	GetWorkflowByID(ctx context.Context, workflowID int) (*store.Workflow, error)
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

type service struct {
	workflowStore s.WorkflowStore
	cmd           *exec.Cmd
	logFile       *os.File
}

func NewService(store s.WorkflowStore, logFile *os.File) Workflow {
	return &service{
		workflowStore: store,
		logFile:       logFile,
	}
}

const (
	fileName = "/app/workflow/data/generated/result.txt"
)

func (s *service) GetWorkflowByID(ctx context.Context, workflowID int) (*store.Workflow, error) {
	return s.workflowStore.GetWorkflowByID(ctx, workflowID)
}

func (s *service) StartExecution(ctx context.Context, workflowID int) error {
	// Create a new context with a 5 minute timeout
	log.Println("Preparing to execute the script...")

	filename := "demo.fastq" // This should be set to whatever the desired filename is.
	// Prepare to execute the script
	s.cmd = exec.Command("bash", "-c", "/app/workflow/rna.sh "+filename)

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
	err = s.workflowStore.StartWorkflow(ctx, workflowID)
	if err != nil {
		log.Printf("Failed to save workflow start time: %v", err)
		return fmt.Errorf("failed to save workflow start time: %w", err)
	}
	log.Println("Script started successfully. Setting up a ticker to monitor the script...")

	return nil
}

func (s *service) Complete(ctx context.Context, id int, duration int) error {
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

func (s *service) UpdateWorkflow(ctx context.Context, wf *s.Workflow) (*s.Workflow, error) {
	return s.workflowStore.UpdateWorkflow(ctx, wf)
}

func (s *service) IsComplete(ctx context.Context, id int) (bool, error) {
	// not implemented
	complete, err := s.IsScriptComplete()
	if err != nil {
		return false, err

	}

	if complete {
		return true, nil
	}
	return false, nil
}

// GetScriptDuration reads the script duration from a file
func (s *service) GetScriptDuration() (int, error) {

	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	duration := strings.TrimSpace(string(contents))
	if duration == "" {
		return 0, errors.New("Script is still running")
	}
	durationInt, err := strconv.Atoi(duration)
	if err != nil {
		return 0, errors.New("duration is not a valid integer")
	}
	return durationInt, nil
}

// IsScriptComplete reads the script duration from a file
func (s *service) IsScriptComplete() (bool, error) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		// If the file does not exist, return false
		return false, nil
	}

	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return false, err
	}

	durationStr := strings.TrimSpace(string(contents))

	if durationStr == "" {
		log.Println("script is still running")
		return false, nil
	}

	if strings.Contains(durationStr, "Error") || strings.Contains(durationStr, "Killed") {
		log.Print("The workflow failed")
		return false, errors.New("The workflow failed")
	}

	return true, nil
}
