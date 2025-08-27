package adapters

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/IgorGrieder/Rinha-backend-go/internal/domain"
	"github.com/IgorGrieder/Rinha-backend-go/internal/ports"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	redisClient  *redis.Client
	hashDefault  string
	hashFallback string
}

func NewRepository(c *redis.Client, hashDefault string, hashFallback string) ports.Repository {
	return &Repository{
		redisClient:  c,
		hashDefault:  hashDefault,
		hashFallback: hashFallback,
	}
}

func (r *Repository) SetValue(ctx context.Context, key string, value domain.InternalPayment, isDefault bool) error {
	const maxRetryes = 5
	const initialBackoff = 1 * time.Second
	var hash string

	if isDefault {
		hash = r.hashDefault
	} else {

		hash = r.hashFallback
	}
	redisKey := fmt.Sprintf(
		"%s:%s",
		hash,
		key,
	)

	for range maxRetryes {
		err := r.redisClient.Set(ctx, redisKey, value, 0).Err()
		if err == nil {
			return nil
		}

		// exponential retry with backoff
		expoRetry := math.Pow(2, float64(maxRetryes))
		time.Sleep(time.Duration(expoRetry))
	}

	// send the error, we won't store it in a dead letter queue
	err := fmt.Errorf("Error while inserting to the %s key the value %v", redisKey, value)
	return err
}

func (r *Repository) GetValue(ctx context.Context, key string) (string, error) {
	// this route will be used for returning the processed ammounts in the application
	return "", nil
}
