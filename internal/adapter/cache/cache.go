package cache

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, payload any, duration time.Duration) error
}
