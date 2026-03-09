package postgres

import (
	"time"

	"github.com/google/uuid"
)

type StatementModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	AccountID   uuid.UUID `gorm:"type:uuid;not null"`
	BankID      int       `gorm:"type:int4;not null"`
	FileName    string    `gorm:"type:varchar(255);not null"`
	PeriodMonth string    `gorm:"type:varchar(7);not null"`
	UploadDate  time.Time `gorm:"type:timestamp;default:now();not null"`
	Status      string    `gorm:"type:movements.statement_status;default:pending;not null"`
}

type MovementModel struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null"`
	AccountID   uuid.UUID  `gorm:"type:uuid;not null"`
	StatementID *uuid.UUID `gorm:"type:uuid"`
	CategoryID  int        `gorm:"type:int4;not null"`
	Date        time.Time  `gorm:"type:date;not null"`
	Description string     `gorm:"type:text;not null"`
	Amount      float64    `gorm:"type:numeric(15,2);not null"`
	Type        string     `gorm:"type:movements.transaction_type;not null"`
	IsActive    bool       `gorm:"column:is_active;default:true;not null"`
	CreatedAt   time.Time  `gorm:"column:created_at;default:now();not null"`
}

func (StatementModel) TableName() string {
	return "movements.statements"
}

func (MovementModel) TableName() string {
	return "movements.movements"
}
