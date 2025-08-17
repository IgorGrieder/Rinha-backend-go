package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IgorGrieder/Rinha-backend-go/internal/domain"
	"github.com/IgorGrieder/Rinha-backend-go/internal/ports"
)

func ProcessPaymentHandler(s ports.PaymentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var payment domain.Payment
		err := json.NewDecoder(r.Body).Decode(&payment)

		if err != nil {
			fmt.Println("Failed body parsing")
			return
		}

		s.ProcessPayment(&payment)
		w.WriteHeader(http.StatusOK)
	}
}

func GetSummaryHandler(s ports.PaymentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.GetAll()
	}
}
