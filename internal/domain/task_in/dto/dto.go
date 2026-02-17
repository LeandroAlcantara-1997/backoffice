package dto

import "time"

type Task struct {
	TaskID           string `json:"task_id"`
	Payload          string `json:"payload"`
	ProcessingTimeMS int    `json:"processing_time_ms"`
}

type TaskOut struct {
	TaskID      string    `json:"task_id"`
	Status      string    `json:"status"`
	ProcessedAt time.Time `json:"processed_at"`
}
