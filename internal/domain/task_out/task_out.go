package taskin

import (
	"backoffice/internal/adapter/cache"
	"backoffice/internal/domain/task_out/dto"
	"context"
	"encoding/json"
	"time"
)

type UseCase interface {
	Process(ctx context.Context, ct *dto.TaskOut) error
}

type service struct {
	cacheService cache.Cache
}

func New(cacheService cache.Cache) UseCase {
	return &service{
		cacheService: cacheService,
	}
}

func (s *service) Process(ctx context.Context, ct *dto.TaskOut) error {
	taskOut := dto.TaskOut{
		TaskID:      ct.TaskID,
		Status:      "processed",
		ProcessedAt: time.Now(),
	}
	payload, err := json.Marshal(&taskOut)
	if err != nil {
		return err
	}
	return s.cacheService.Set(ctx, "taskout:"+ct.TaskID, payload, 0)
}
