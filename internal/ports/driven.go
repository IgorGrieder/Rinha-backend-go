package ports

import (
	"context"
	"time"

	"github.com/IgorGrieder/Rinha-backend-go/internal/domain"
)

type Repository interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type Queue interface {
	Enqueue(ctx context.Context, payment *domain.Payment)
}
