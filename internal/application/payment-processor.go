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

func (p *PaymentProcessor) ProcessPayment(queueName string, payment *domain.InternalPayment) error {
	ctx := context.Background()
	payment = payment.NewPaymentWithTimeStamp()

	if err := p.q.Enqueue(ctx, queueName, payment); err != nil {
		return err
	}

	return nil
}

func (p *PaymentProcessor) GetAll() {

}
