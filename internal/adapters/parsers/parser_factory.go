package parsers

import (
	"finanzas-mvp/internal/core/ports"
	"fmt"
)

const (
	BankBancolombia = "1"
	BankNu          = "2"
)

func NewParser(bankID string) (ports.MovementParser, error) {

	switch bankID {

	case BankBancolombia:
		return NewBancolombiaParser(), nil

	case BankNu:
		return NewNuParser(), nil

	default:
		return nil, fmt.Errorf("unsupported bank")
	}
}
