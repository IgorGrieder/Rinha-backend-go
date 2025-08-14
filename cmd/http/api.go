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
		"/POST api/payments",
		ProcessPaymentHandler(s, r, res),
	)

	svr := &http.Server{Addr: fmt.Sprintf(":%d", cfg.PORT), Handler: mux}

	if err := svr.ListenAndServe(); err != nil {
		os.Exit(1)
	}
}
