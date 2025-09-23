package domain

import (
	"time"

	"github.com/google/uuid"
)

type DateFilter struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

type Payment struct {
	Id          uuid.UUID `json:"correlationId"`
	Amount      float32   `json:"amount"`
	RequestedAt time.Time `json:"requestedAt"`
}

type InternalPayment struct {
	Id                 uuid.UUID `json:"Id"`
	Amount             float32   `json:"Amount"`
	RequestedAt        time.Time `json:"RequestedAt"`
	IsDefaultProcessor bool      `json:"IsDefaultProcessor"`
}

func (p InternalPayment) NewPaymentWithTimeStamp() *InternalPayment {
	return &InternalPayment{
		Id:                 p.Id,
		Amount:             p.Amount,
		IsDefaultProcessor: p.IsDefaultProcessor,
		RequestedAt:        time.Now().UTC(),
	}
}

func PaymentMapper(payment Payment) *InternalPayment {
	return &InternalPayment{
		Id:          payment.Id,
		Amount:      payment.Amount,
		RequestedAt: payment.RequestedAt,
	}
}
