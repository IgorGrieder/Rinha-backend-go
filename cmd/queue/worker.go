package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/IgorGrieder/Rinha-backend-go/internal/config"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type paymentInput struct {
	correlationId uuid.UUID
	amount        float32
	requestedAt   string
}

func StartPaymentQueue(cfg *config.Config, redis *redis.Client) {
	ctx := context.Background()

	for {
		data, err := redis.BLPop(ctx, 0, cfg.QUEUE).Result()
		if err != nil {
			return
		}

		json, err := json.Marshal(data[1])
		res, err := http.Post(cfg.REDIS_ADDR, "application/json", bytes.NewBuffer(json))
		if err != nil {
			return
		}
	}
}
