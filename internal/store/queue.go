package store

import "time"

type Queue struct {
	ID         int64     `json:"id"`
	WorkflowID int64     `json:"workflow_id"`
	Status     string    `json:"status"`
	EnqueuedAt time.Time `json:"enqueued_at"`
}