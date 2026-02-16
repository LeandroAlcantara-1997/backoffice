package tasksin

import (
	"backoffice/internal/adapter/broker"
	taskin "backoffice/internal/domain/task_in"
	"backoffice/internal/domain/task_in/dto"
	"context"
	"encoding/json"
)

type handler struct {
	taskin.UseCase
}

func NewTasksIn(c taskin.UseCase) *handler {
	return &handler{
		UseCase: c,
	}
}
func (h *handler) HandleFunc(ctx context.Context, d broker.Delivery) error {
	var consumerContract dto.Task
	if err := json.Unmarshal(d.Body(), &consumerContract); err != nil {
		return err
	}

	if err := h.UseCase.Process(ctx, &consumerContract); err != nil {
		return err
	}
	return nil
}
