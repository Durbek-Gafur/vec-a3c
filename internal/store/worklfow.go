package store

import (
	"database/sql"
	"time"
)

type Workflow struct {
	ID                  int64     `json:"id"`
	Name                string    `json:"name"`
	Type                string    `json:"type"`
	Duration            int64     `json:"duration"`
	ReceivedAt          time.Time `json:"received_at"`
	StartedExecutionAt  sql.NullTime `json:"started_execution_at,omitempty"`
	CompletedAt         sql.NullTime `json:"completed_at,omitempty"`
}

type WorkflowFilter struct {
	Type      string
	StartTime time.Time
	EndTime   time.Time
}
