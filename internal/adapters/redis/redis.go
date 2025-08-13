package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Options struct {
	Addr     string
	Password string
	DB       int
	Protocol int
}

type Redis struct {
	client *goredis.Client
}

func NewAdapter(opts Options) *Redis {
	client := goredis.NewClient(&goredis.Options{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
		Protocol: opts.Protocol,
	})
	return &Redis{client: client}
}

func (a *Redis) Ping(ctx context.Context) error {
	return a.client.Ping(ctx).Err()
}

func (a *Redis) Close() error {
	return a.client.Close()
}

func (a *Redis) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return a.client.Set(ctx, key, value, ttl).Err()
}

func (a *Redis) Get(ctx context.Context, key string) (string, error) {
	return a.client.Get(ctx, key).Result()
}
