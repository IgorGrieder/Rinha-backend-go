package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/IgorGrieder/Rinha-backend-go/internal/config"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type paymentInput struct {
	correlationId uuid.UUID
	amount        float32
	requestedAt   string
}

func StartPaymentQueue(cfg *config.Config, redisClient *redis.Client) {
	log.Println("Starting payment queue worker...")
	queueName := cfg.QUEUE
	ctx := context.Background()

	for {
		retry := true
		data, err := redisClient.BLPop(ctx, 0, queueName).Result()
		if err != nil {
			log.Printf("ERROR: Failed to pop from Redis queue '%s': %v", queueName, err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Printf("INFO: Received new payment job: %s", data[1])

		var jsonToSend bytes.Buffer
		if err := json.NewEncoder(&jsonToSend).Encode(data[1]); err != nil {
			log.Printf("ERROR: Failed to encode JSON: %v", err)
			continue
		}

		for retry {
			r, err := http.Post(decideServer(cfg), "application/json", &jsonToSend)
			if err != nil {
				log.Printf("ERROR: Failed to POST payment data: %v", err)
				continue
			}
			defer r.Body.Close()

			if r.StatusCode != http.StatusOK {
				log.Printf("WARN: Server responded with non-200 status: %s", r.Status)
			} else {
				log.Printf("INFO: Successfully processed payment job.")
				retry = false
			}
		}
	}
}

func decideServer(cfg *config.Config) string {
	// Make an post and see if teh service is available
	return "http://localhost:8080/process_payment"
}
