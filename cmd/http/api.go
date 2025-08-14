package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/IgorGrieder/Rinha-backend-go/internal/config"
)

func StartServer(cfg *config.Config) {
	mux := http.NewServeMux()

	svr := &http.Server{Addr: fmt.Sprintf(":%d", cfg.PORT), Handler: mux}

	if err := svr.ListenAndServe(); err != nil {
		os.Exit(1)
	}

}
