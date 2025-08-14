package http

import (
	"net/http"

	"github.com/IgorGrieder/Rinha-backend-go/internal/ports"
)

func ProcessPaymentHandler(s ports.PaymentService, r http.ResponseWriter, res *http.Request) http.HandlerFunc {
	return func(r http.ResponseWriter, res *http.Request) {
		s.ProcessPayment()
	}
}

func GetSummaryHandler(s ports.PaymentService, r http.ResponseWriter, res *http.Request) http.HandlerFunc {
	return func(r http.ResponseWriter, res *http.Request) {
		s.GetAll()
	}
}
