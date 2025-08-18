package adapters

import (
	"context"
	"time"

	"github.com/IgorGrieder/Rinha-backend-go/internal/ports"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	redisClient *redis.Client
}

func NewRepository(c *redis.Client) ports.Repository {
	return &Repository{redisClient: c}
}

func (r *Repository) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return nil
}

func (r *Repository) Get(ctx context.Context, key string) (string, error) {
	return "", nil
}
