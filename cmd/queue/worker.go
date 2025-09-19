package queue

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IgorGrieder/Rinha-backend-go/internal/domain"
	"github.com/IgorGrieder/Rinha-backend-go/internal/ports"
)

type Worker struct {
	repository   ports.Repository
	service      ports.PaymentService
	queue        ports.Queue
	queueName    string
	fallbackAddr string
	defaultAddr  string
}

func NewWorker(r ports.Repository, s ports.PaymentService, q ports.Queue, queue string, fallbackAddr string, defaultAddr string) *Worker {
	return &Worker{
		repository:   r,
		service:      s,
		queue:        q,
		queueName:    queue,
		fallbackAddr: fallbackAddr,
		defaultAddr:  defaultAddr,
	}
}

func (w *Worker) StartPaymentQueue(workerId int) {
	log.Printf("Starting payment queue worker %d...", workerId)

	// The worker job should be just reading from the queue and using the service for the rest
	// I will alter on change the responsabilty
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

		if err := w.service.ProcessWorker(data, w.fallbackAddr, w.defaultAddr); err != nil {
			// This is the correct way to check the error type
			if jsonErr, ok := err.(*domain.JSONParsingError); !ok {
				// If we get an error we will use a goroutine to handle it
				go func() {
					// We won't care about the error or not in this situation
					time.Sleep(5 * time.Minute)
					w.queue.Enqueue(w.queueName, &job)
				}()
			} else {
				log.Printf("JSON Parsing Error occurred: %v", jsonErr)
			}
		}

	}
}
