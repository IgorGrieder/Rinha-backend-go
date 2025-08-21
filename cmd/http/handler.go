package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/IgorGrieder/Rinha-backend-go/internal/domain"
	"github.com/IgorGrieder/Rinha-backend-go/internal/ports"
)

func ProcessPaymentHandler(s ports.PaymentService, queueName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var payment domain.Payment
		err := json.NewDecoder(r.Body).Decode(&payment)
		payment.RequestedAt = time.Now().UTC()
		defer r.Body.Close()

		if err != nil {
			fmt.Println("Failed body parsing")
			return
		}

		s.ProcessPayment(queueName, &payment)
		w.WriteHeader(http.StatusOK)
	}
}

func GetSummaryHandler(s ports.PaymentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.GetAll()
	}
}
