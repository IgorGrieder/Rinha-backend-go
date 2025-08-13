package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	redisadapter "github.com/IgorGrieder/Rinha-backend-go/internal/adapters/redis"
	"github.com/IgorGrieder/Rinha-backend-go/internal/ports"
)

var (
	redisOnce      sync.Once
	redisInstance  ports.Repository
	redisInitError error
)

func InitRedis(ctx context.Context) error {
	redisOnce.Do(func() {
		addr := getEnv("REDIS_ADDR", "localhost:6379")
		password := os.Getenv("REDIS_PASSWORD")
		db := parseInt(getEnv("REDIS_DB", "0"))
		protocol := parseInt(getEnv("REDIS_PROTOCOL", "2"))

		adapter := redisadapter.NewAdapter(redisadapter.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
			Protocol: protocol,
		})

		pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		if err := adapter.Ping(pingCtx); err != nil {
			redisInitError = fmt.Errorf("redis ping failed: %w", err)
			return
		}
		redisInstance = adapter
	})
	return redisInitError
}

func Redis() ports.Repository {
	return redisInstance
}

func CloseRedis() error {
	if redisInstance != nil {
		return redisInstance.Close()
	}
	return nil
}

func getEnv(key string, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
