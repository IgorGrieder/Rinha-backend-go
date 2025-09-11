package ports

import (
	"time"

	"github.com/IgorGrieder/Rinha-backend-go/internal/domain"
)

type PaymentService interface {
	ProcessPayment(queueName string, p *domain.InternalPayment) error
	GetAll(startDate, endDate time.Time) ([]domain.InternalPayment, error)
}
