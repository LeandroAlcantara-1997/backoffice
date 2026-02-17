package taskout

import (
	"backoffice/internal/adapter/broker"
	"backoffice/internal/domain/task_out/dto"
	"backoffice/internal/mocks"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleFunc(t *testing.T) {
	t.Run("must return success when read payload", func(t *testing.T) {
		var (
			ctx     = context.Background()
			useCase = mocks.NewTaskOutProcessUseCaseMock(t)
			h       = handler{
				UseCase: useCase,
			}
			ct = &dto.TaskOut{
				TaskID:      "1",
				Status:      "processed",
				ProcessedAt: time.Now(),
			}
		)
		useCase.EXPECT().Process(ctx, mock.Anything).Return(nil)
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
			useCase = mocks.NewTaskOutProcessUseCaseMock(t)
			h       = handler{
				UseCase: useCase,
			}
			ct = &dto.TaskOut{
				TaskID:      "1",
				Status:      "processed",
				ProcessedAt: time.Now(),
			}
		)
		useCase.EXPECT().Process(ctx, mock.Anything).Return(errors.New("error"))
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
			useCase = mocks.NewTaskOutProcessUseCaseMock(t)
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
