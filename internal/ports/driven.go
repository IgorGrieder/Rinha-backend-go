package ports

import (
	"context"

	"github.com/IgorGrieder/Rinha-backend-go/internal/domain"
)

type Repository interface {
	SetValue(ctx context.Context, key string, value domain.InternalPayment, isDefault bool) error
	GetValue(ctx context.Context, key string) (string, error)
}

type Queue interface {
	Enqueue(ctx context.Context, queueName string, payment *domain.InternalPayment) error
	Dequeue(ctx context.Context, queueName string) []string
}
