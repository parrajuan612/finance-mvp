package parsers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"finanzas-mvp/internal/core/domain"

	"github.com/google/uuid"
)

var (
	// Busca el año en la cabecera del documento (ej: HASTA: 2025/09/30)
	reYear = regexp.MustCompile(`HASTA:\s*([0-9]{4})/[0-9]{2}/[0-9]{2}`)
	// Detecta el inicio de una fila de transacción (espacios opcionales + DD/MM + al menos un espacio)
	reRow = regexp.MustCompile(`^\s*([0-9]{1,2})/([0-9]{1,2})\s+`)
	// Usa 2 o más espacios consecutivos para separar las columnas
	reColumns = regexp.MustCompile(`\s{2,}`)
)

type BancolombiaParser struct{}

func NewBancolombiaParser() *BancolombiaParser {
	return &BancolombiaParser{}
}

func (p *BancolombiaParser) Parse(text string) ([]domain.Movement, error) {
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("texto vacío")
	}

	year := p.extractYear(text)
	lines := strings.Split(text, "\n")
	var movements []domain.Movement

	for _, line := range lines {
		movement, ok := p.parseLine(line, year)
		if !ok {
			continue
		}
		movements = append(movements, movement)
	}

	if len(movements) == 0 {
		return nil, fmt.Errorf("no se detectaron movimientos válidos")
	}

	return movements, nil
}

// parseLine se encarga de extraer la data de un string
func (p *BancolombiaParser) parseLine(line string, year int) (domain.Movement, bool) {
	matches := reRow.FindStringSubmatch(line)
	if len(matches) < 3 {
		return domain.Movement{}, false
	}

	parts := reColumns.Split(strings.TrimSpace(line), -1)
	if len(parts) < 4 {
		return domain.Movement{}, false
	}

	// matches[1] es el día, matches[2] es el mes extraído
	day, _ := strconv.Atoi(matches[1])
	month, _ := strconv.Atoi(matches[2])

	description := strings.TrimSpace(parts[1])
	amount := parseNumber(parts[len(parts)-2])

	if amount == 0 {
		return domain.Movement{}, false
	}

	return p.buildMovement(day, month, year, description, amount), true
}

// buildMovement crea la entidad del dominio
func (p *BancolombiaParser) buildMovement(d, m, y int, desc string, amt float64) domain.Movement {
	movementType := domain.TypeIncome
	if amt < 0 {
		movementType = domain.TypeExpense
	}

	return domain.Movement{
		ID:          uuid.Nil,
		Date:        time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC),
		Description: desc,
		Amount:      amt,
		Type:        movementType,
		CreatedAt:   time.Now(),
	}
}

// extractYear obtener la fecha del documento
func (p *BancolombiaParser) extractYear(text string) int {
	year := time.Now().Year()
	if m := reYear.FindStringSubmatch(text); len(m) >= 2 {
		if y, err := strconv.Atoi(m[1]); err == nil && y > 1900 {
			year = y
		}
	}
	return year
}

// parseNumber asegura que los montos bancarios se conviertan correctamente a float
func parseNumber(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	isNegative := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "+")
	s = strings.TrimPrefix(s, "-")
	s = strings.ReplaceAll(s, ",", "")

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}

	if isNegative {
		return -f
	}
	return f
}
