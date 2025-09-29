package ports

import (
	"github.com/IgorGrieder/Rinha-backend-go/internal/domain"
)

type Repository interface {
	SetValue(key string, value domain.InternalPayment, isDefault bool) error
	GetPayments(startScore, endScore float64) ([]domain.InternalPayment, error)
}

type Queue interface {
	Enqueue(queueName string, payment *domain.InternalPayment) error
	Dequeue(queueName string) ([]string, error)
}
