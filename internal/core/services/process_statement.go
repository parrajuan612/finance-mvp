package services

import (
	"finanzas-mvp/internal/adapters/parsers"
	"finanzas-mvp/internal/core/domain"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type StatementService struct {
	categorizer *Categorizer // Inyectamos el servicio de categorización
}

func NewStatementService(cat *Categorizer) *StatementService {
	return &StatementService{
		categorizer: cat,
	}
}

func (s *StatementService) ProcessStatement(file io.Reader, filename string, password string, bankID string) ([]domain.Movement, error) {
	savedPath, err := s.saveFile(file, filename)
	if err != nil {
		return nil, err
	}
	defer os.Remove(savedPath)

	text, err := parsers.ExtractText(savedPath, password)
	if err != nil {
		return nil, err
	}

	parser, err := parsers.NewParser(bankID)
	if err != nil {
		return nil, err
	}

	movements, err := parser.Parse(text)
	if err != nil {
		return nil, err
	}

	for i := range movements {

		s.categorizer.Categorize(&movements[i])

	}

	return movements, nil
}

func (s *StatementService) saveFile(file io.Reader, filename string) (string, error) {

	newFilename := fmt.Sprintf("%d_%s", time.Now().Unix(), filename)

	path := filepath.Join("storage", newFilename)

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(out, file)
	if err != nil {
		out.Close()
		return "", err
	}

	err = out.Close()
	if err != nil {
		return "", err
	}

	return path, nil
}
