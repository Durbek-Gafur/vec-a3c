package store

import (
	"time"
)

type Workflow struct {
	ID                  int64     `json:"id"`
	Name                string    `json:"name"`
	Type                string    `json:"type"`
	Duration            int64     `json:"duration"`
	ReceivedAt          time.Time `json:"received_at"`
	StartedExecutionAt  time.Time `json:"started_execution_at,omitempty"`
	CompletedAt         time.Time `json:"completed_at,omitempty"`
}

type WorkflowFilter struct {
	Type      string
	StartTime time.Time
	EndTime   time.Time
}
