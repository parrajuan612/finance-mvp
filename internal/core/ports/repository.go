package ports

import (
	"context"
	"finanzas-mvp/internal/core/domain"
)

type MovementRepository interface {
	SaveBatch(ctx context.Context, movements []domain.Movement) error
}
