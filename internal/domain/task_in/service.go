package taskin

import (
	"backoffice/internal/adapter/broker"
	"backoffice/internal/domain/task_in/dto"
	"backoffice/internal/exception"
	"context"
	"encoding/json"
	"time"
)

type UseCase interface {
	Process(ctx context.Context, ct *dto.Task) error
}

type service struct {
	brokerProvider broker.Broker
}

func New(brokerProvider broker.Broker) UseCase {
	return &service{
		brokerProvider: brokerProvider,
	}
}

func (s *service) Process(ctx context.Context, ct *dto.Task) error {
	if ct.ProcessingTimeMS > 2000 {
		return exception.ErrProcessingTimeout
	}
	taskOut := dto.TaskOut{
		TaskID:      ct.TaskID,
		Status:      "processed",
		ProcessedAt: time.Now(),
	}
	payload, err := json.Marshal(&taskOut)
	if err != nil {
		return err
	}
	return s.brokerProvider.Publish(ctx, payload)
}
