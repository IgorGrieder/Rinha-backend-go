package domain

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	Id          uuid.UUID `json:"correlationId"`
	Amount      float32   `json:"amount"`
	RequestedAt time.Time `json:"requestedAt"`
}

type InternalPayment struct {
	Id                 uuid.UUID `json:"IdId"`
	Amount             float32   `json:"Amount"`
	RequestedAt        time.Time `json:"RequestedAt"`
	IsDefaultProcessor bool      `json:"IsDefaultProcessor"`
}

func PaymentMapper(payment Payment) *InternalPayment {
	return &InternalPayment{
		Id:          payment.Id,
		Amount:      payment.Amount,
		RequestedAt: payment.RequestedAt,
	}
}
