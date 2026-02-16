package redis

import (
	"context"
	"fmt"
	"time"

	redisotel "github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	*redis.Client
}

func New(ctx context.Context, host, port, pass string,
	readTimeout, writeTimeout time.Duration) (*redis.Client, error) {
	redisClient := redis.NewClient(
		&redis.Options{
			Addr:         fmt.Sprintf("%s:%s", host, port),
			Password:     pass,
			DB:           0,
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
		},
	)
	redisotel.InstrumentTracing(redisClient)
	cmd := redisClient.Ping(ctx)
	if cmd.Err() != nil {
		return nil, fmt.Errorf("ping -> %w", cmd.Err())
	}
	return redisClient, nil
}
