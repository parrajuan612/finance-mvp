package postgres

import (
	"time"

	"github.com/google/uuid"
)

type StatementModel struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;index"`
	AccountID   uuid.UUID `gorm:"type:uuid;index"`
	BankID      uuid.UUID `gorm:"type:uuid;index"`
	FileName    string
	FilePath    string
	FileHash    string `gorm:"index"`
	FileSize    int64
	MimeType    string
	PeriodMonth string
	UploadDate  time.Time
	Status      string
	ProcessedAt *time.Time
	Attempts    int
	ErrorDetail string
	CreatedAt   time.Time
}

type MovementModel struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID       uuid.UUID  `gorm:"type:uuid;index"`
	AccountID    uuid.UUID  `gorm:"type:uuid;index"`
	StatementID  *uuid.UUID `gorm:"type:uuid;index"`
	CategoryID   uuid.UUID  `gorm:"type:uuid;index"`
	Date         time.Time  `gorm:"type:date;index"`
	Description  string
	Amount       float64
	Type         string
	MovementHash string `gorm:"index"`
	CreatedAt    time.Time
}
