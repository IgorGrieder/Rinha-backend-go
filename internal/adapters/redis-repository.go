package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/IgorGrieder/Rinha-backend-go/internal/ports"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	redisClient *redis.Client
	hashPrefix  string
}

func NewRepository(c *redis.Client, hashPrefix string) ports.Repository {
	return &Repository{
		redisClient: c,
		hashPrefix:  hashPrefix,
	}
}

func (r *Repository) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	err := r.Set(ctx, r.hashPrefix, value, ttl)
	if err != nil {
		fmt.Println("ERROR: writing to the hash in redis")
		return err
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, key string) (string, error) {
	return "", nil
}
