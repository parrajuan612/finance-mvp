package postgres

import (
	"context"
	"time"

	"finanzas-mvp/internal/core/domain"
	"finanzas-mvp/internal/core/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type statementRepository struct {
	db *gorm.DB
}

func NewStatementRepository(db *gorm.DB) ports.StatementRepository {
	return &statementRepository{db: db}
}

func (r *statementRepository) SaveWithMovements(ctx context.Context, statement *domain.Statement, movements []domain.Movement) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1) Crear StatementModel
		stmtModel := StatementModel{
			ID:          uuid.Nil, // GORM generará UUID
			UserID:      statement.UserID,
			AccountID:   statement.AccountID,
			BankID:      statement.BankID,
			FileName:    statement.FileName,
			PeriodMonth: statement.PeriodMonth,
			UploadDate:  time.Now(),
			Status:      string(domain.StatusPending),
		}

		if err := tx.Create(&stmtModel).Error; err != nil {
			return err
		}

		// Push back created ID into domain.Statement
		statement.ID = stmtModel.ID

		// 2) Mapear movimientos a MovementModel y asignar statement_id
		movModels := make([]MovementModel, 0, len(movements))
		for _, mv := range movements {
			stID := stmtModel.ID
			movModels = append(movModels, MovementModel{
				ID:          uuid.Nil,
				UserID:      mv.UserID,
				AccountID:   mv.AccountID,
				StatementID: &stID,
				CategoryID:  mv.CategoryID,
				Date:        mv.Date,
				Description: mv.Description,
				Amount:      mv.Amount,
				Type:        string(mv.Type),
				CreatedAt:   time.Now(),
			})
		}

		if len(movModels) > 0 {
			if err := tx.Create(&movModels).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
