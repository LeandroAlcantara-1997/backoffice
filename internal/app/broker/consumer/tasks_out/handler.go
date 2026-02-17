package taskout

import (
	"backoffice/internal/adapter/broker"
	taskout "backoffice/internal/domain/task_out"
	"backoffice/internal/domain/task_out/dto"
	"context"
	"encoding/json"
)

type handler struct {
	taskout.UseCase
}

func NewTasksOut(c taskout.UseCase) *handler {
	return &handler{
		UseCase: c,
	}
}
func (h *handler) HandleFunc(ctx context.Context, d broker.Delivery) error {
	var consumerContract dto.TaskOut
	if err := json.Unmarshal(d.Body(), &consumerContract); err != nil {
		return err
	}
	if err := h.UseCase.Process(ctx, &consumerContract); err != nil {
		return err
	}
	return nil
}
