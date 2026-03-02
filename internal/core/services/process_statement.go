package services

import (
	"finanzas-mvp/internal/adapters/parsers"
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

func (s *StatementService) ProcessStatement(file io.Reader, filename string, password string, bankID string) error {
	savedPath, err := s.saveFile(file, filename)
	if err != nil {
		return err
	}
	defer os.Remove(savedPath) // Limpieza: borrar archivo tras procesar

	text, err := parsers.ExtractText(savedPath, password)
	if err != nil {
		return err
	}

	parser, err := parsers.NewParser(bankID)
	if err != nil {
		return err
	}

	movements, err := parser.Parse(text)
	if err != nil {
		return err
	}

	for i := range movements {

		s.categorizer.Categorize(&movements[i])

		fmt.Printf("%s | %-40s | %10.2f | %-10s | CAT: %d (%s)\n",
			movements[i].Date.Format("2006-01-02"),
			movements[i].Description,
			movements[i].Amount,
			movements[i].Type,
			movements[i].CategoryID,   // El ID numérico (1, 2, 3...)
			movements[i].CategoryName, // El nombre legible
		)
	}

	return nil
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
