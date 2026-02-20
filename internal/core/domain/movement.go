package domain

import (
	"time"

	"github.com/google/uuid"
)

type Movement struct {
	ID          uuid.UUID    `json:"id"`
	UserID      uuid.UUID    `json:"user_id"`
	AccountID   uuid.UUID    `json:"account_id"`
	StatementID *uuid.UUID   `json:"statement_id,omitempty"`
	CategoryID  uuid.UUID    `json:"category_id"`
	Date        time.Time    `json:"date"`
	Description string       `json:"description"`
	Amount      float64      `json:"amount"`
	Type        MovementType `json:"type"`
	CreatedAt   time.Time    `json:"created_at"`
}
type Category struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	IsActive bool      `json:"is_active"`
}
type MovementType string

const (
	TypeIncome  MovementType = "income"
	TypeExpense MovementType = "expense"
)
