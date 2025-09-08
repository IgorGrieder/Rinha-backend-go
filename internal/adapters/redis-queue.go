package adapters

import (
	"context"
	"encoding/json"
	"log"

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
		return err
	}

	q.redisClient.RPush(ctx, queueName, string(json))
	return nil
}

func (q *PaymentQueue) Dequeue(ctx context.Context, queueName string) *redis.StringSliceCmd {
	return q.redisClient.BLPop(ctx, 0, queueName)
}
