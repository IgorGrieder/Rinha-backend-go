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

func (q *PaymentQueue) Enqueue(ctx context.Context, queueName string, payment *domain.Payment) {
	json, err := json.Marshal(payment)
	if err != nil {
		log.Println("FATAL: error while encoding json to append to the queue")
		return
	}

	q.redisClient.RPush(ctx, queueName, string(json))
}
