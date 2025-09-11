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

func (r *Repository) SetValue(key string, payment domain.InternalPayment, isDefault bool) error {
	const maxRetries = 5
	const initialBackoff = 5 * time.Second
	var hash string
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

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

	for range maxRetries {
		err := r.redisClient.HSet(
			ctx,
			redisKey,
			"paymentId", payment.Id.String(),
			"amount", fmt.Sprintf("%f", payment.Amount),
			"requestedAt", payment.RequestedAt.Format(time.RFC3339),
		).Err()
		if err == nil {
			continue
		}

		// Now we will store the value to a sorted set
		score := float64(payment.RequestedAt.Unix())
		err = r.redisClient.ZAdd(ctx, "payments:by:date", redis.Z{
			Score:  score,
			Member: payment.Id.String(),
		}).Err()

		// exponential retry with backoff
		expoRetry := math.Pow(2, float64(maxRetries))
		time.Sleep(time.Duration(expoRetry))
	}

	// send the error, we won't store it in a dead letter queue
	err := fmt.Errorf("Error while inserting to the payment %+v", payment)
	return err
}

func (r *Repository) GetValue(key string) (string, error) {
	// this route will be used for returning the processed ammounts in the application
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return "", nil
}
