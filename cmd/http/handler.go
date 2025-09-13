package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/IgorGrieder/Rinha-backend-go/internal/domain"
	"github.com/IgorGrieder/Rinha-backend-go/internal/ports"
)

type PaymentsResponse struct {
	Payments []domain.InternalPayment `json:"payments"`
}

func ProcessPaymentHandler(s ports.PaymentService, queueName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var payment domain.Payment
		err := json.NewDecoder(r.Body).Decode(&payment)
		payment.RequestedAt = time.Now().UTC()
		defer r.Body.Close()

		if err != nil {
			fmt.Println("Failed body parsing")
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		if err = s.ProcessPayment(queueName, domain.PaymentMapper(payment)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("ERROR: failed inserting an payment to be handled: %s", err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func GetSummaryHandler(s ports.PaymentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var dateFilter domain.DateFilter
		err := json.NewDecoder(r.Body).Decode(&dateFilter)
		defer r.Body.Close()

		if err != nil {
			fmt.Println("Failed body parsing")
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		payments, err := s.GetAll(dateFilter.StartDate, dateFilter.EndDate)

		if len(payments) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		paymentsResponse := &PaymentsResponse{Payments: payments}
		jsonToSend, err := json.Marshal(paymentsResponse)

		if err != nil {
			fmt.Println("Failed parsing the return json")
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonToSend)
	}
}
