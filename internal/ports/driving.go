package ports

import "github.com/IgorGrieder/Rinha-backend-go/internal/domain"

type PaymentService interface {
	ProcessPayment(p *domain.Payment)
	GetAll()
}
