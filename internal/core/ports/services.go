package ports

import (
	"context"
	"finanzas-mvp/internal/core/domain"
	"io"
)

type StatementService interface {
	ProcessStatement(
		file io.Reader,
		fileName string,
		password string,
		bankID string,
	) ([]domain.Movement, string, error)

	SaveStatementWithMovements(
		ctx context.Context,
		statement domain.Statement,
		movements []domain.Movement,
	) error
}
