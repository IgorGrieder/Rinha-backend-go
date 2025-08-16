package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/IgorGrieder/Rinha-backend-go/internal/config"
	"github.com/IgorGrieder/Rinha-backend-go/internal/ports"
)

func StartServer(cfg *config.Config, s ports.PaymentService) {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"/POST /payments",
		ProcessPaymentHandler(s),
	)

	mux.HandleFunc(
		"/GET /payments-summary",
		GetSummaryHandler(s),
	)

	svr := &http.Server{Addr: fmt.Sprintf(":%d", cfg.PORT), Handler: mux}

	if err := svr.ListenAndServe(); err != nil {
		fmt.Println("Server crashed")
		os.Exit(1)
	}
}
