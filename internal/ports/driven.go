package ports

import (
	"github.com/IgorGrieder/Rinha-backend-go/internal/domain"
)

type Repository interface {
	SetValue(key string, value domain.InternalPayment, isDefault bool) error
	GetValue(key string) (string, error)
}

type Queue interface {
	Enqueue(queueName string, payment *domain.InternalPayment) error
	Dequeue(queueName string) []string
}
