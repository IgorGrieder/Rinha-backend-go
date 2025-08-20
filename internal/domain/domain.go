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
