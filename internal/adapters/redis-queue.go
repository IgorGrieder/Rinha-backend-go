package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

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

func (q *PaymentQueue) Enqueue(queueName string, payment *domain.InternalPayment) error {
	// Enqueue safe, with backoff logic
	const maxRetries = 5
	const initialBackoff = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	json, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("FATAL: error while encoding json to append to the queue: %w", err)
	}

	for range maxRetries {

		if err = q.redisClient.RPush(ctx, queueName, string(json)).Err(); err != nil {
			return nil
		}

		log.Println("FATAL: error while writing to the queue")

		// exponential retry with backoff
		expoRetry := math.Pow(2, float64(maxRetries))
		time.Sleep(time.Duration(expoRetry))
	}

	return fmt.Errorf("Error in enqueue process in the queue: %s for value %+v", queueName, payment)
}

func (q *PaymentQueue) Dequeue(queueName string) []string {
	// Dequeu safe, with backoff logic
	const maxRetries = 5
	const initialBackoff = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for range maxRetries {
		data, err := q.redisClient.BLPop(ctx, 0, queueName).Result()
		fmt.Printf("An error happened in the BLPop in redis %v", err)

		if err == nil {
			return data
		}

		// exponential retry with backoff
		expoRetry := math.Pow(2, float64(maxRetries))
		time.Sleep(time.Duration(expoRetry))
	}

	return nil
}
