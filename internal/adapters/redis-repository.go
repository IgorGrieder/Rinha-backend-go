package adapters

import (
	"context"
	"fmt"
	"math"
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

func (r *Repository) SetValue(ctx context.Context, key string, value int64) error {
	const maxRetryes = 5
	const initialBackoff = 1 * time.Second
	redisKey := fmt.Sprintf(
		"%s:%s",
		r.hashPrefix,
		key,
	)

	for range maxRetryes {
		err := r.redisClient.IncrBy(ctx, redisKey, value).Err()
		if err == nil {
			return nil
		}

		// exponential retry with backoff
		expoRetry := math.Pow(2, float64(maxRetryes))
		time.Sleep(time.Duration(expoRetry))
	}

	// send the error, we won't store it in a dead letter queue
	err := fmt.Errorf("Error while inserting to the %s key the value %d", redisKey, value)
	return err
}

func (r *Repository) GetValue(ctx context.Context, key string) (string, error) {
	// this route will be used for returning the processed ammounts in hh
	return "", nil
}
