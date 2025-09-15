package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/IgorGrieder/Rinha-backend-go/internal/domain"
	"github.com/IgorGrieder/Rinha-backend-go/internal/ports"
)

type PaymentProcessor struct {
	r ports.Repository
	q ports.Queue
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

func NewPaymentProcessor(r ports.Repository, q ports.Queue) *PaymentProcessor {
	return &PaymentProcessor{r, q}
}

func (p *PaymentProcessor) ProcessPayment(queueName string, payment *domain.InternalPayment) error {
	payment = payment.NewPaymentWithTimeStamp()

	if err := p.q.Enqueue(queueName, payment); err != nil {
		return err
	}

	return nil
}

func (p *PaymentProcessor) GetAll(startDate, endDate time.Time) ([]domain.InternalPayment, error) {
	// Convert time.Time to Unix timestamps for the score range
	startScore := float64(startDate.Unix())
	endScore := float64(endDate.Unix())

	payments, err := p.r.GetPayments(startScore, endScore)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to fetch payments based on the given dates")
	}

	return payments, nil
}

func (p *PaymentProcessor) ProcessWorker(data []byte, fallbackAddr, defaultAddr string) error {
	client := &http.Client{Timeout: 1 * time.Second}

	for {

		urlToCall := decideServer(defaultAddr, fallbackAddr, client)
		r, err := client.Post(urlToCall, "application/json", bytes.NewBuffer(data))

		if err != nil {
			log.Printf("ERROR: Failed to POST payment data: %v", err)
			continue
		}

		isDefault := urlToCall == defaultAddr

		isErrorStatus := r.StatusCode != http.StatusOK
		if isErrorStatus {
			log.Printf("WARN: Server responded with non-200 status: %s", r.Status)
		} else {
			log.Printf("INFO: Successfully processed payment job.")

			// Writting to redis the actual value, if we get an error we must queue again the job
			if err = p.r.SetValue(job.Id.String(), job, isDefault); err != nil {
				// We could send to and DLQ but in this case I will just schedule an go routine
				// for some time to provide the proper write

				return err

			}

			break
		}

		r.Body.Close()
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
