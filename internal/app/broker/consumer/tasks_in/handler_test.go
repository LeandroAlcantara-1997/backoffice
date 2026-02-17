package tasksin

import (
	"backoffice/internal/adapter/broker"
	"backoffice/internal/domain/task_in/dto"
	"backoffice/internal/mocks"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

func TestHandleFunc(t *testing.T) {
	t.Run("must return success when read payload", func(t *testing.T) {
		var (
			ctx     = context.Background()
			useCase = mocks.NewTaskInProcessUseCaseMock(t)
			h       = handler{
				UseCase: useCase,
			}
			ct = &dto.Task{
				TaskID:           "1",
				Payload:          "processing",
				ProcessingTimeMS: 500,
			}
		)
		useCase.EXPECT().Process(ctx, ct).Return(nil)
		payload, err := json.Marshal(ct)
		assert.NoError(t, err)

		err = h.HandleFunc(ctx, broker.Delivery{D: amqp091.Delivery{
			Body: payload,
		}})

		assert.NoError(t, err)
	})

	t.Run("must return error from use case", func(t *testing.T) {
		var (
			ctx     = context.Background()
			useCase = mocks.NewTaskInProcessUseCaseMock(t)
			h       = handler{
				UseCase: useCase,
			}
			ct = &dto.Task{
				TaskID:           "1",
				Payload:          "processing",
				ProcessingTimeMS: 500,
			}
		)
		useCase.EXPECT().Process(ctx, ct).Return(errors.New("error"))
		payload, err := json.Marshal(ct)
		assert.NoError(t, err)

		err = h.HandleFunc(ctx, broker.Delivery{D: amqp091.Delivery{
			Body: payload,
		}})

		assert.Error(t, err)
	})

	t.Run("must return error when sent a invalid body", func(t *testing.T) {
		var (
			ctx     = context.Background()
			useCase = mocks.NewTaskInProcessUseCaseMock(t)
			h       = handler{
				UseCase: useCase,
			}
		)
		payload, err := json.Marshal("potato")
		assert.NoError(t, err)

		err = h.HandleFunc(ctx, broker.Delivery{D: amqp091.Delivery{
			Body: payload,
		}})

		assert.Error(t, err)
	})
}
