package taskin

import (
	"backoffice/internal/domain/task_in/dto"
	"backoffice/internal/exception"
	"backoffice/internal/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	ctx        context.Context
	brokerMock *mocks.BrokerMock
	srv        *service
}

func (s *ServiceTestSuite) SetupSubTest() {
	s.brokerMock = mocks.NewBrokerMock(s.T())
	s.srv = &service{
		brokerProvider: s.brokerMock,
	}
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) TestProcess() {
	s.Run("must return success when pubish message", func() {
		s.brokerMock.EXPECT().Publish(s.ctx, mock.Anything).Return(nil)
		err := s.srv.Process(s.ctx, &dto.Task{
			TaskID:           "1",
			Payload:          "processed",
			ProcessingTimeMS: 500,
		})
		s.NoError(err)
	})

	s.Run("must return error when broker call fails", func() {
		s.brokerMock.EXPECT().Publish(s.ctx, mock.Anything).Return(errors.New("unexpcted error"))
		err := s.srv.Process(s.ctx, &dto.Task{
			TaskID:           "1",
			Payload:          "processed",
			ProcessingTimeMS: 500,
		})
		s.Error(err)
	})

	s.Run("must return error when processing ms is bigger than 2000 ms", func() {
		err := s.srv.Process(s.ctx, &dto.Task{
			TaskID:           "1",
			Payload:          "processed",
			ProcessingTimeMS: 2001,
		})
		s.Error(err)
		s.ErrorIs(err, exception.ErrProcessingTimeout)
	})
}
