package workflow

import (
	"time"
)

// Workflow represents a workflow with its properties
type Workflow struct {
	ID                  int64     `json:"id"`
	Name                string    `json:"name"`
	Type                string    `json:"type"`
	Duration            int64     `json:"duration"`
	ReceivedAt          time.Time `json:"received_at"`
	StartedExecutionAt  time.Time `json:"started_execution_at,omitempty"`
	CompletedAt         time.Time `json:"completed_at,omitempty"`
}

// WorkflowFilter is a struct to filter workflows during retrieval
type WorkflowFilter struct {
	Type      string
	StartTime time.Time
	EndTime   time.Time
}

// NewWorkflow returns a new Workflow instance
func NewWorkflow(name, wType string, duration int64) *Workflow {
	return &Workflow{
		Name:       name,
		Type:       wType,
		Duration:   duration,
		ReceivedAt: time.Now(),
	}
}

// StartExecution sets the StartedExecutionAt field to the current time
func (w *Workflow) StartExecution() {
	w.StartedExecutionAt = time.Now()
}

// Complete sets the CompletedAt field to the current time
func (w *Workflow) Complete() {
	w.CompletedAt = time.Now()
}
