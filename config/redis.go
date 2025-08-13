package config

import (
    redisadapter "github.com/IgorGrieder/Rinha-backend-go/adapters/redis"
    "github.com/IgorGrieder/Rinha-backend-go/ports"
)

// RedisConfig holds the minimal options to configure a Redis adapter.
// Use this in your composition root and pass the returned CachePort
// into your application/services for easy testing and substitution.
type RedisConfig struct {
    Addr     string
    Password string
    DB       int
    Protocol int
}

// NewRedis constructs a new Redis adapter instance implementing ports.CachePort.
// No globals, no side effects, easy to mock in tests.
func NewRedis(cfg RedisConfig) ports.CachePort {
    return redisadapter.NewAdapter(redisadapter.Options{
        Addr:     cfg.Addr,
        Password: cfg.Password,
        DB:       cfg.DB,
        Protocol: cfg.Protocol,
    })
}
