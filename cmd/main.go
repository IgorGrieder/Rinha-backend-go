package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IgorGrieder/Rinha-backend-go/config"
)

func main() {
	ctx := context.Background()
	if err := config.InitRedis(ctx); err != nil {
		log.Fatalf("failed to init redis: %v", err)
	}
	defer func() {
		if err := config.CloseRedis(); err != nil {
			log.Printf("error closing redis: %v", err)
		}
	}()

	cache := config.Redis()
	key := "hello"
	value := "Ola"
	if err := cache.Set(ctx, key, value, 10*time.Second); err != nil {
		log.Fatalf("set failed: %v", err)
	}
	got, err := cache.Get(ctx, key)
	if err != nil {
		log.Fatalf("get failed: %v", err)
	}
	fmt.Println(got)
}
