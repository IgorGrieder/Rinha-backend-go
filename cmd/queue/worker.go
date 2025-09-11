package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/IgorGrieder/Rinha-backend-go/internal/domain"
	"github.com/IgorGrieder/Rinha-backend-go/internal/ports"
)

type Worker struct {
	repository   ports.Repository
	queue        ports.Queue
	queueName    string
	fallbackAddr string
	defaultAddr  string
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

func NewWorker(r ports.Repository, q ports.Queue, queue string, fallbackAddr string, defaultAddr string) *Worker {
	return &Worker{
		repository:   r,
		queue:        q,
		queueName:    queue,
		fallbackAddr: fallbackAddr,
		defaultAddr:  defaultAddr,
	}
}

func (w *Worker) StartPaymentQueue(workerId int) {
	log.Printf("Starting payment queue worker %d...", workerId)
	client := &http.Client{Timeout: 1 * time.Second}

	for {
		data := w.queue.Dequeue(w.queueName)

		if data == nil {
			log.Printf("ERROR: Failed to pop from Redis queue with backoff '%s'", w.queueName)
			time.Sleep(1 * time.Second)
			continue
		}

		log.Printf("INFO: Received new payment job: %s", data[1])

		var job domain.InternalPayment
		dataInBytes := []byte(data[1])

		if err := json.Unmarshal(dataInBytes, &job); err != nil {
			// In this case we have a parsing error, so we won't retry because of bad json

			log.Printf("ERROR: Failed to parse JSON from queue: %+v", err)
			continue
		}

		for {

			urlToCall := decideServer(w.defaultAddr, w.fallbackAddr, client)
			r, err := client.Post(urlToCall, "application/json", bytes.NewBuffer(dataInBytes))

			if err != nil {
				log.Printf("ERROR: Failed to POST payment data: %v", err)
				continue
			}

			isDefault := urlToCall == w.defaultAddr

			isErrorStatus := r.StatusCode != http.StatusOK
			if isErrorStatus {
				log.Printf("WARN: Server responded with non-200 status: %s", r.Status)
			} else {
				log.Printf("INFO: Successfully processed payment job.")

				// Writting to redis the actual value, if we get an error we must queue again the job
				if err = w.repository.SetValue(job.Id.String(), job, isDefault); err != nil {
					// We could send to and DLQ but in this case I will just schedule an go routine
					// for some time to provide the proper write

					go func() {
						// We won't care about the error or not in this situation
						time.Sleep(1 * time.Minute)
						w.queue.Enqueue(w.queueName, &job)
					}()

				}

				break
			}

			r.Body.Close()
		}
	}
}

func decideServer(def string, fallback string, client *http.Client) string {
	// using an WaitGrop to request in parallel
	var wg sync.WaitGroup

	urls := []string{def, fallback}
	resultsChan := make(chan result, len(urls))
	r1, r2, timeout1, timeout2 := parallelCalls(resultsChan, &wg, urls, client)

	// in case we receive 429 from both services we will trust the default
	if timeout1 && timeout2 {
		return def
	}

	// simple case, just going with the option with the least timeout received from the backends
	shouldUseFirstUrl := r1.TimeLimit > r2.TimeLimit
	if shouldUseFirstUrl {
		log.Printf("CHOOSEN: Server choose the: %s", r1.Url)
		return r1.Url
	} else {
		log.Printf("CHOOSEN: Server choose the: %s", r2.Url)
		return r2.Url
	}
}

func parallelCalls(resultsChan chan result, wg *sync.WaitGroup, urls []string, client *http.Client) (result, result, bool, bool) {
	for idx, url := range urls {
		wg.Add(1)

		go func(idx int, url string) {
			defer wg.Done()

			r, err := client.Get(url)

			if err != nil {
				result := result{
					Err:        err,
					StatusCode: r.StatusCode,
				}
				resultsChan <- result
				return
			}
			defer r.Body.Close()

			var res response
			if err = json.NewDecoder(r.Body).Decode(&res); err != nil {
				log.Printf("ERROR: error while parsing json")

				result := result{
					Err:        err,
					StatusCode: r.StatusCode,
				}
				resultsChan <- result
				return
			}

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
	timeout1 := r1.StatusCode == http.StatusTooManyRequests
	timeout2 := r2.StatusCode == http.StatusTooManyRequests
	return r1, r2, timeout1, timeout2
}
