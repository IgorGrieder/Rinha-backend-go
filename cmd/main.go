package main

import (
	"fmt"

	"github.com/IgorGrieder/Rinha-backend-go/cmd/http"
	"github.com/IgorGrieder/Rinha-backend-go/cmd/queue"
	"github.com/IgorGrieder/Rinha-backend-go/internal/adapters"
	"github.com/IgorGrieder/Rinha-backend-go/internal/application"
	"github.com/IgorGrieder/Rinha-backend-go/internal/config"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.NewConfig()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.REDIS_ADDR, cfg.REDIS_PORT),
		Password: "",
		DB:       0,
		Protocol: 2,
	})
	r := adapters.NewRepository(redisClient)
	s := application.NewPaymentProcessor(r)
	queue.StartPaymentQueue(redisClient)
	http.StartServer(cfg, redisClient, s)
}
