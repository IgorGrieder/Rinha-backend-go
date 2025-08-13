package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

type helloResponse struct {
    Message string `json:"message"`
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
            return
        }
        name := r.URL.Query().Get("name")
        if name == "" {
            name = "World"
        }
        w.Header().Set("Content-Type", "application/json")
        _ = json.NewEncoder(w).Encode(helloResponse{Message: fmt.Sprintf("Hello, %s!", name)})
    })
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
            return
        }
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok"))
    })

    port := getEnv("PORT", "8080")
    srv := &http.Server{
        Addr:              ":" + port,
        Handler:           mux,
        ReadHeaderTimeout: 5 * time.Second,
        ReadTimeout:       10 * time.Second,
        WriteTimeout:      10 * time.Second,
        IdleTimeout:       60 * time.Second,
    }

    // graceful shutdown
    shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    go func() {
        log.Printf("HTTP server listening on :%s", port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("server error: %v", err)
        }
    }()

    <-shutdownCtx.Done()
    log.Println("shutting down...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("graceful shutdown failed: %v", err)
        _ = srv.Close()
    }
    log.Println("server stopped")
}

func getEnv(key, fallback string) string {
    v := os.Getenv(key)
    if v == "" {
        return fallback
    }
    return v
}
