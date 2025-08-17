package domain

import "github.com/google/uuid"

type Payment struct {
	Id     uuid.UUID `json:"correlationId"`
	Amount float32   `json:"amount"`
}
