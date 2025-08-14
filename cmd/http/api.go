package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/IgorGrieder/Rinha-backend-go/internal/config"
	"github.com/IgorGrieder/Rinha-backend-go/internal/ports"
	"github.com/redis/go-redis/v9"
)

func StartServer(cfg *config.Config, redis *redis.Client, s ports.PaymentService) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/POST api/payments",
		func(res http.ResponseWriter, req *http.Request) {
			s.ProcessPayment()
		},
	)

	svr := &http.Server{Addr: fmt.Sprintf(":%d", cfg.PORT), Handler: mux}

	if err := svr.ListenAndServe(); err != nil {
		os.Exit(1)
	}
}
