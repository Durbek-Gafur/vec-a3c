package store

import "time"

type Queue struct {
	ID         int     `json:"id"`
	WorkflowID int     `json:"workflow_id"`
	Status     string    `json:"status"`
	EnqueuedAt time.Time `json:"enqueued_at"`
}