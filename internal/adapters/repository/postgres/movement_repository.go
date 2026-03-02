package postgres

import (
	"context"
	"finanzas-mvp/internal/core/domain"
	"finanzas-mvp/internal/core/ports"

	"gorm.io/gorm"
)

type movementRepository struct {
	db *gorm.DB
}

func NewMovementRepository(db *gorm.DB) ports.MovementRepository {
	return &movementRepository{db: db}
}

func (r *movementRepository) SaveBatch(ctx context.Context, movements []domain.Movement) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&movements).Error; err != nil {
			return err
		}
		return nil
	})
}
