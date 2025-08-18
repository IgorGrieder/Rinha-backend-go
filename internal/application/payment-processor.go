package application

import (
	"context"

	"github.com/IgorGrieder/Rinha-backend-go/internal/domain"
	"github.com/IgorGrieder/Rinha-backend-go/internal/ports"
)

type PaymentProcessor struct {
	r ports.Repository
	q ports.Queue
}

func NewPaymentProcessor(r ports.Repository, q ports.Queue) ports.PaymentService {
	return &PaymentProcessor{r, q}
}

func (p *PaymentProcessor) ProcessPayment(payment *domain.Payment) {
	ctx := context.Background()
	p.q.Enqueue(ctx, payment)
}

func (p *PaymentProcessor) GetAll() {

}
