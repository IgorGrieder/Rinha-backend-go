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
	def := cfg.DEFAULT_ADDR
	fallback := cfg.FALLBACK_ADDR

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
			r, err := http.Post(decideServer(def, fallback), "application/json", &jsonToSend)
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

func decideServer(def string, fallback string) string {
	// using an WaitGrop to request in parallel
	var wg sync.WaitGroup
	urls := []string{def, fallback}
	resultsChan := make(chan result, len(urls))

	for idx, url := range urls {

		wg.Add(1)
		go func(idx int, url string) {
			defer wg.Done()

			r, err := http.Get(url)
			if err != nil {
				result := result{
					Err:        err,
					StatusCode: r.StatusCode,
				}
				resultsChan <- result
				return
			}

			var res response
			err = json.NewDecoder(r.Body).Decode(&res)
			defer r.Body.Close()
			log.Printf("RESPONSE: Server responded with: %v", res)

			result := result{
				Url:        url,
				Status:     res.Status,
				Err:        nil,
				TimeLimit:  res.TimeLimit,
				StatusCode: r.StatusCode,
			}
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
	// when receiving information in the channel append to an slice
	// when the wait group finishes the await and the channel closes it will end the for loop
	for result := range resultsChan {
		allResults = append(allResults, result)
	}

	r1 := allResults[0]
	r2 := allResults[1]
	t1 := r1.StatusCode
	t2 := r2.StatusCode

	// in case we receive 429 from both services we will trust the default
	if t1 == http.StatusTooManyRequests && t2 == http.StatusTooManyRequests {
		return def
	}

	// simple case, just going with the option with the least timeout received from the backends
	if r1.TimeLimit > r2.TimeLimit {
		log.Printf("CHOOSEN: Server choose the: %s", r1.Url)
		return r1.Url
	} else {
		log.Printf("CHOOSEN: Server choose the: %s", r2.Url)
		return r2.Url
	}
}
