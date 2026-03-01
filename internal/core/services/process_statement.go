package services

import (
	"finanzas-mvp/internal/adapters/parsers"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type StatementService struct{}

func NewStatementService() *StatementService {
	return &StatementService{}
}

func (s *StatementService) ProcessStatement(file io.Reader, filename string, password string, bankID string) error {

	savedPath, err := s.saveFile(file, filename)
	if err != nil {
		return err
	}
	text, err := parsers.ExtractText(savedPath, password)
	if err != nil {
		return err
	}
	fmt.Println("====== RAW TEXT START ======")
	fmt.Println(text)
	fmt.Println("====== RAW TEXT END ======")
	parser, err := parsers.NewParser(bankID)
	if err != nil {
		return err
	}
	movements, err := parser.Parse(text)
	if err != nil {
		return err
	}
	for _, m := range movements {
		fmt.Printf("%s | %-60s | %12.2f | %s\n",
			m.Date.Format("2006-01-02"),
			m.Description,
			m.Amount,
			m.Type,
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
