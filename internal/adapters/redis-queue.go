package adapters

import (
	"context"

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

func (q *PaymentQueue) Enqueue(ctx context.Context, payment *domain.Payment) {
	return
}
