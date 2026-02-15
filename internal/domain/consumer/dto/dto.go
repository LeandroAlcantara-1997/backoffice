package dto

type Task struct {
	TaskID           string `json:"task_id"`
	Payload          string `json:"payload"`
	ProcessingTimeMS int    `json:"processing_time_ms"`
}
