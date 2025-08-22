package adapters

import (
	"context"
	"fmt"

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

func (r *Repository) SetValue(ctx context.Context, key string, value int64) error {
	err := r.redisClient.IncrBy(ctx, r.hashPrefix, value).Err()
	if err != nil {
		fmt.Println("ERROR: writing to the hash in redis")
		return err
	}
	return nil
}

func (r *Repository) GetValue(ctx context.Context, key string) (string, error) {
	return "", nil
}
