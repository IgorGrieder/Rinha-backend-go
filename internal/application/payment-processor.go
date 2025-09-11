package application

import (
	"time"

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
	payment = payment.NewPaymentWithTimeStamp()

	if err := p.q.Enqueue(queueName, payment); err != nil {
		return err
	}

	return nil
}

func (p *PaymentProcessor) GetAll(startDate, endDate time.Time) {
	// Convert time.Time to Unix timestamps for the score range
	startScore := float64(startDate.Unix())
	endScore := float64(endDate.Unix())

	p.r.GetPayments(startScore, endScore)

}
