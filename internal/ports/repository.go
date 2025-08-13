package ports

import (
	"context"
	"time"
)

type Repository interface {
	Ping(ctx context.Context) error
	Close() error

	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}
