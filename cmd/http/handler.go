package http

import (
	"net/http"

	"github.com/IgorGrieder/Rinha-backend-go/internal/ports"
)

func ProcessPaymentHandler(s ports.PaymentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.ProcessPayment()
	}
}

func GetSummaryHandler(s ports.PaymentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.GetAll()
	}
}
