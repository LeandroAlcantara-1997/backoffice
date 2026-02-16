package container

import (
	env "backoffice/config"
	"backoffice/internal/adapter/broker"
	"backoffice/internal/adapter/cache"
	taskin "backoffice/internal/domain/task_in"
	taskout "backoffice/internal/domain/task_out"
	"backoffice/internal/infra/rabbitmq"
	"backoffice/internal/infra/redis"
	"context"
	"time"
)

type Container struct {
	UseCases
	*components
}

type components struct {
	TaskInBroker  broker.Broker
	TaskOutBroker broker.Broker
	TaskDQLBroker broker.Broker
	cacheService  cache.Cache
}
type UseCases struct {
	TaskInUseCase  taskin.UseCase
	TaskOutUseCase taskout.UseCase
}

func New(ctx context.Context) (*Container, error) {
	var err error
	if err = env.LoadEnv(); err != nil {
		return nil, err
	}

	tasksIn, err := setupBroker(
		env.Env.TasksInBrokerConfig.QueueName,
		env.Env.TasksInBrokerConfig.URL)
	if err != nil {
		return nil, err
	}

	tasksOut, err := setupBroker(
		env.Env.TasksOutBrokerConfig.QueueName,
		env.Env.TasksOutBrokerConfig.URL)
	if err != nil {
		return nil, err
	}

	tasksDLQ, err := setupBroker(env.Env.TasksDLQBrokerConfig.QueueName,
		env.Env.TasksDLQBrokerConfig.URL)
	if err != nil {
		return nil, err
	}

	redisClient, err := redis.New(ctx, env.Env.RedisCache.Host,
		env.Env.RedisCache.Port, env.Env.RedisCache.Pass,
		time.Second*time.Duration(env.Env.ReadTimeout),
		time.Second*time.Duration(env.Env.WriteTimeout),
	)
	if err != nil {
		return nil, err
	}
	return &Container{
		components: &components{
			TaskInBroker:  tasksIn,
			TaskOutBroker: tasksOut,
			TaskDQLBroker: tasksDLQ,
			cacheService:  cache.NewRedis(redisClient),
		},
		UseCases: UseCases{
			TaskInUseCase:  taskin.New(tasksOut),
			TaskOutUseCase: taskout.New(cache.NewRedis(redisClient)),
		},
	}, nil
}

func setupBroker(queueName, url string) (broker.Broker, error) {
	conn, ch, err := rabbitmq.New(url)
	if err != nil {
		return nil, err
	}

	return broker.NewBroker(queueName, ch, conn), nil
}

func (ctn *Container) CloseConnections() error {
	if err := ctn.components.TaskDQLBroker.Close(); err != nil {
		return err
	}

	if err := ctn.components.TaskInBroker.Close(); err != nil {
		return err
	}

	if err := ctn.components.TaskOutBroker.Close(); err != nil {
		return err
	}

	return nil
}
