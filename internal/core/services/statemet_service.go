package services

import (
	"context"
	"finanzas-mvp/internal/core/domain"
	"finanzas-mvp/internal/core/ports"
)

type StatementService struct {
	categorizer *Categorizer
	stmtRepo    ports.StatementRepository
}

func NewStatementService(
	categorizer *Categorizer,
	statementRepo ports.StatementRepository,
) *StatementService {

	return &StatementService{
		categorizer: categorizer,
		stmtRepo:    statementRepo,
	}
}

func (s *StatementService) SaveStatementWithMovements(
	ctx context.Context,
	statement domain.Statement,
	movements []domain.Movement,
) error {

	return s.stmtRepo.SaveWithMovements(ctx, &statement, movements)
}
