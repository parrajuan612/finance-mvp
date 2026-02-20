package domain

import (
	"time"

	"github.com/google/uuid"
)

type Bank struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Account struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	BankID   uuid.UUID `json:"bank_id"`
	Name     string    `json:"name"`
	CreateAt time.Time `json:"created_at"`
}
