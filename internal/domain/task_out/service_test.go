package taskout

import (
	"backoffice/internal/domain/task_out/dto"
	"backoffice/internal/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	ctx        context.Context
	brokerMock *mocks.BrokerMock
	cacheMock  *mocks.CacheMock
	srv        *service
}

func (s *ServiceTestSuite) SetupSubTest() {
	s.cacheMock = mocks.NewCacheMock(s.T())
	s.srv = &service{
		cacheService: s.cacheMock,
	}
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) TestProcess() {
	s.Run("must return error when cache call fails", func() {
		s.cacheMock.EXPECT().Set(s.ctx, "taskout:1", mock.Anything, time.Duration(0)).Return(errors.New("unexpcted error"))
		err := s.srv.Process(s.ctx, &dto.TaskOut{
			TaskID:      "1",
			Status:      "must processing",
			ProcessedAt: time.Now(),
		})
		s.Error(err)
	})
	s.Run("must return success when save in cache", func() {
		s.cacheMock.EXPECT().Set(s.ctx, "taskout:1", mock.Anything, time.Duration(0)).Return(nil)
		err := s.srv.Process(s.ctx, &dto.TaskOut{
			TaskID:      "1",
			Status:      "must processing",
			ProcessedAt: time.Now(),
		})
		s.NoError(err)
	})
}
