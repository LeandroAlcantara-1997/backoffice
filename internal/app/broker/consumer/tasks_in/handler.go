package tasksin

import (
	taskin "backoffice/internal/domain/task_in"
	"backoffice/internal/domain/task_in/dto"
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type handler struct {
	taskin.UseCase
}

func NewTasksIn(c taskin.UseCase) *handler {
	return &handler{
		UseCase: c,
	}
}
func (h *handler) HandleFunc(ctx context.Context, d amqp.Delivery) error {
	var consumerContract dto.Task
	if err := json.Unmarshal(d.Body, &consumerContract); err != nil {
		return err
	}

	println(consumerContract.TaskID)
	if err := h.UseCase.Process(ctx, &consumerContract); err != nil {
		return err
	}
	return nil
}
