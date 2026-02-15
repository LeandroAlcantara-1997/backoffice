package consumer

import (
	"backoffice/internal/domain/consumer/dto"
	"context"
)

type UseCase interface {
	Process(ctx context.Context, ct *dto.Task)
}

type service struct {
}

func New() UseCase {
	return nil
}

func (s *service) Process(ctx context.Context, ct *dto.Task) {}
