package ports

import "finanzas-mvp/internal/core/domain"

type MovementParser interface {
	Parse(text string) ([]domain.Movement, string, error)
}
