package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	*redis.Client
}

func NewRedis(r *redis.Client) Cache {
	return &redisClient{Client: r}
}
func (r *redisClient) Set(ctx context.Context, key string, payload any, duration time.Duration) error {
	return r.Client.Set(ctx, key, payload, duration).Err()
}
