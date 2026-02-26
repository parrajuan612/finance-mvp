package ports

import (
	"context"
	"finanzas-mvp/internal/core/domain"
)

// StatementParser es el contrato que implementarán los parsers por banco/formato.
type StatementParser interface {
	// Parse debe devolver los movimientos extraídos y el periodo normalizado (ej "2026-01")
	Parse(ctx context.Context, filePath string, password string) ([]domain.Movement, string, error)
}
