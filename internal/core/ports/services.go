package ports

import (
	"context"
	"finanzas-mvp/internal/core/domain"
)

type ParserService interface {
	Parse(ctx context.Context, filePath string) ([]domain.Movement, error)
}
