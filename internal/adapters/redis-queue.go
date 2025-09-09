package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"math"

	"github.com/IgorGrieder/Rinha-backend-go/internal/domain"
	"github.com/IgorGrieder/Rinha-backend-go/internal/ports"
	"github.com/redis/go-redis/v9"
)

type PaymentQueue struct {
	redisClient *redis.Client
}

func NewQueue(c *redis.Client) ports.Queue {
	return &PaymentQueue{redisClient: c}
}

func (q *PaymentQueue) Enqueue(ctx context.Context, queueName string, payment *domain.InternalPayment) error {
	json, err := json.Marshal(payment)
	if err != nil {
		log.Println("FATAL: error while encoding json to append to the queue")
		err := fmt.Errorf("Error in enqueue enqueue process in the queue: %s for value %v", queueName, payment)

		return err
	}

	q.redisClient.RPush(ctx, queueName, string(json))
	return nil
}

func (q *PaymentQueue) Dequeue(ctx context.Context, queueName string) []string {
	// Dequeu safe, with backoff logic
	const maxRetries = 5
	const initialBackoff = 5 * time.Second

	for range maxRetries {
		data, err := q.redisClient.BLPop(ctx, 0, queueName).Result()
		
		if err != nil {
			return data
		}

		// exponential retry with backoff
		expoRetry := math.Pow(2, float64(maxRetries))
		time.Sleep(time.Duration(expoRetry))

	}

	return nil
}
