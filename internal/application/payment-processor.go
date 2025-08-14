package application

import (
	"fmt"

	"github.com/IgorGrieder/Rinha-backend-go/internal/ports"
)

type PaymentProcessor struct {
	r ports.Repository
}

func NewPaymentProcessor(r ports.Repository) ports.PaymentService {
	return &PaymentProcessor{r}
}

func (p *PaymentProcessor) ProcessPayment() {
	fmt.Print("Hellooo")
}

func (p *PaymentProcessor) GetAll() {

}
