package dto

import "time"

type TaskOut struct {
	TaskID      string    `json:"task_id"`
	Status      string    `json:"status"`
	ProcessedAt time.Time `json:"processed_at"`
}
