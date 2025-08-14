package application

import (
	"fmt"

	"github.com/IgorGrieder/Rinha-backend-go/internal/ports"
)

type PaymentProcessor struct {
	r ports.Repository
}

func NewPaymentProcessor(r ports.Repository) *PaymentProcessor {
	return &PaymentProcessor{r}
}

func (p *PaymentProcessor) ProcessPayment() {
	fmt.Print("Hellooo")
}
