package domain

import (
	"time"

	"github.com/google/uuid"
)

type StatementStatus string

const (
	StatusPending   StatementStatus = "pending"
	StatusProcessed StatementStatus = "processed"
	StatusError     StatementStatus = "error"
)

type Statement struct {
	ID          uuid.UUID       `json:"id"`
	UserID      uuid.UUID       `json:"user_id"`
	AccountID   uuid.UUID       `json:"account_id"`
	BankID      uuid.UUID       `json:"bank_id"`
	FileName    string          `json:"file_name"`
	PeriodMonth string          `json:"period_month"`
	UploadDate  time.Time       `json:"upload_date"`
	Status      StatementStatus `json:"status"`
}
