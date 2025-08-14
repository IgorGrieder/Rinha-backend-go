package ports

import (
	"context"
	"time"
)

type Repository interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}
