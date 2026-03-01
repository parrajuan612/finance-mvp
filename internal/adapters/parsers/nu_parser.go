package parsers

import "finanzas-mvp/internal/core/domain"

type NuParser struct{}

func NewNuParser() *NuParser {
	return &NuParser{}
}

func (p *NuParser) Parse(text string) ([]domain.Movement, error) {
	return ParseNuText(text)
}

func ParseNuText(text string) ([]domain.Movement, error) {
	return nil, nil
}
