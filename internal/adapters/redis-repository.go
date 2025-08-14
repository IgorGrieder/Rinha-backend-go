package adapters

import (
	"context"
	"time"

	"github.com/IgorGrieder/Rinha-backend-go/internal/ports"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	redisClient *redis.Client
}

func NewRepository(c *redis.Client) ports.Repository {
	return &RedisRepository{redisClient: c}
}

func (r *RedisRepository) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return nil
}

func (r *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	return "", nil
}
