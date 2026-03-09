package ports

import (
	"context"
	"finanzas-mvp/internal/core/domain"
)

type MovementRepository interface {
	SaveBatch(ctx context.Context, movements []domain.Movement) error
}

type StatementRepository interface {
	// SaveWithMovements hace la operación transaccional: guarda statement y movements y devuelve error si falla.
	SaveWithMovements(ctx context.Context, statement *domain.Statement, movements []domain.Movement) error
}
