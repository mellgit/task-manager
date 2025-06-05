package worker

import "github.com/google/uuid"

type TaskPayload struct {
	TaskID uuid.UUID `json:"task_id"`
	UserID uuid.UUID `json:"user_id"`
}
