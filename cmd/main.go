package main

import (
	"context"
	"fmt"

	"github.com/IgorGrieder/Rinha-backend-go/cmd/http"
	"github.com/IgorGrieder/Rinha-backend-go/cmd/queue"
	"github.com/IgorGrieder/Rinha-backend-go/internal/adapters"
	"github.com/IgorGrieder/Rinha-backend-go/internal/application"
	"github.com/IgorGrieder/Rinha-backend-go/internal/config"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Configs
	cfg := config.NewConfig()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.REDIS_ADDR, cfg.REDIS_PORT),
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	fmt.Println("STARTING redis")
	ctx := context.Background()

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Could not connect to Redis: %v", err)
	}

	fmt.Printf("Successfully connected to Redis. PING response: %s\n", pong)

	// Composition Root
	r := adapters.NewRepository(redisClient, cfg.HASH_DEFAULT, cfg.HASH_FALLBACK)
	q := adapters.NewQueue(redisClient)
	s := application.NewPaymentProcessor(r, q)
	w := queue.NewWorker(r, s, q, cfg.QUEUE, cfg.DEFAULT_ADDR, cfg.FALLBACK_ADDR)

	// spawning workers to read from the queue
	for idx := range cfg.WORKERS {
		go w.StartPaymentQueue(idx)
	}

	http.StartServer(cfg, s)
}
