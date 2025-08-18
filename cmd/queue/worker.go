package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type paymentInput struct {
	correlationId uuid.UUID
	amount        float32
	requestedAt   string
}

func StartPaymentQueue(redis *redis.Client) {
	ctx := context.Background()
	queue := "payment-queue"

	for {
		data, err := redis.BLPop(ctx, 0, queue).Result()
		if err != nil {
			return
		}

		json, err := json.Marshal(data[1])
		req, err := http.NewRequest("POST", "http://payment-processor-default:8080", bytes.NewBuffer(json))
		if err != nil {
			return
		}

		req.Header.Set("Content-Type", "application/json")

	}
}
