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

type Adapter struct {
	client *goredis.Client
}

func NewAdapter(opts Options) *Adapter {
	client := goredis.NewClient(&goredis.Options{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
		Protocol: opts.Protocol,
	})
	return &Adapter{client: client}
}

func (a *Adapter) Ping(ctx context.Context) error {
	return a.client.Ping(ctx).Err()
}

func (a *Adapter) Close() error {
	return a.client.Close()
}

func (a *Adapter) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return a.client.Set(ctx, key, value, ttl).Err()
}

func (a *Adapter) Get(ctx context.Context, key string) (string, error) {
	return a.client.Get(ctx, key).Result()
}
