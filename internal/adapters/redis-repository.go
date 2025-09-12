package adapters

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/IgorGrieder/Rinha-backend-go/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	redisClient  *redis.Client
	hashDefault  string
	hashFallback string
}

func NewRepository(c *redis.Client, hashDefault string, hashFallback string) *Repository {
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

func (r *Repository) GetPayments(startScore, endScore float64) ([]domain.InternalPayment, error) {
	const maxRetries = 5
	const initialBackoff = 5 * time.Second

	for range maxRetries {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Use ZRANGEBYSCORE to get all payment IDs within the date range
		paymentIDs, err := r.redisClient.ZRangeByScore(ctx, "payments:by:date", &redis.ZRangeBy{
			Min: fmt.Sprintf("%f", startScore),
			Max: fmt.Sprintf("%f", endScore),
		}).Result()

		if err == nil {
			continue
		}

		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Use a pipeline to send all HGETALL commands at once.
		pipe := r.redisClient.Pipeline()
		cmds := make(map[string]*redis.MapStringStringCmd)
		for _, id := range paymentIDs {
			cmds[id] = pipe.HGetAll(ctx, id)
		}

		// Execute the pipeline.
		if _, err := pipe.Exec(ctx); err != nil {
			fmt.Printf("Pipeline execution failed: %v\n", err)
			return nil, nil
		}

		// Create a slice to store the final payments.
		payments := make([]domain.InternalPayment, 0, len(paymentIDs))

		// Iterate through the executed commands to get their results.
		for id, cmd := range cmds {
			data, err := cmd.Result()
			if err != nil {
				fmt.Printf("Failed to get data for ID %s: %v\n", id, err)
				continue
			}

			var payment domain.InternalPayment
			payment.Id, err = uuid.Parse(id)
			if err != nil {
				return nil, err
			}

			if amount, err := strconv.ParseFloat(data["amount"], 32); err == nil {
				payment.Amount = float32(amount)
			}

			if requestedAt, err := time.Parse(time.RFC3339, data["requestedAt"]); err == nil {
				payment.RequestedAt = requestedAt
			}

			payments = append(payments, payment)
		}

		return payments, nil
	}

	fmt.Println("Erro while getting all payments")

	return make([]domain.InternalPayment, 0), nil
}
