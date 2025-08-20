package ports

import "github.com/IgorGrieder/Rinha-backend-go/internal/domain"

type PaymentService interface {
	ProcessPayment(queueName string, p *domain.Payment)
	GetAll()
}
