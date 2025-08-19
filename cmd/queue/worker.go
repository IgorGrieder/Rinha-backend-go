package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
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

type response struct {
	Status    bool `json:"failing"`
	TimeLimit int  `json:"minResponseTime"`
}

type result struct {
	Url        string
	StatusCode int
	Status     bool
	TimeLimit  int
	Err        error
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
	// using an WaitGroup to request in parallel
	var wg sync.WaitGroup
	urls := []string{cfg.DEFAULT_ADDR, cfg.FALLBACK_ADDR}
	resultsChan := make(chan result, len(urls))

	for idx, url := range urls {

		wg.Add(1)
		go func(idx int, url string) {
			defer wg.Done()

			r, err := http.Get(url)
			if err != nil {
				result := result{Err: err, StatusCode: r.StatusCode}
				resultsChan <- result

				return
			}

			var res response
			err = json.NewDecoder(r.Body).Decode(&res)
			defer r.Body.Close()

			result := result{Url: url, Status: res.Status, Err: nil, TimeLimit: res.TimeLimit, StatusCode: r.StatusCode}
			resultsChan <- result

		}(idx, url)

	}

	// Go routine to close the channle
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	wg.Wait()

	var allResults []result
	for result := range resultsChan {
		allResults = append(allResults, result)
	}

	// tenho que comparar os dois timeouts
	// vou setar como default sempre chamar o default mesmo
	// in case we receive 429 from both services we will trust the default
	if allResults[0].StatusCode == http.StatusTooManyRequests && allResults[1].StatusCode == http.StatusTooManyRequests {
		return cfg.DEFAULT_ADDR
	}

	// If
	biggestTimeout := biggestTimeout(allResults)

	return cfg.DEFAULT_ADDR
}

func biggestTimeout(results []result) result {
	timeOut1 := results[0]
	timeOut2 := results[1]
	return results[0]
}
