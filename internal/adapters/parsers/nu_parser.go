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
	// Detecta el inicio de transacción: ej. "29 ENE" o "31 DIC 2025"
	// Grupo 1: Día, Grupo 2: Mes, Grupo 3: Año (Opcional), Grupo 4: Resto de la línea
	reNuRow = regexp.MustCompile(`^\s*(\d{2})\s+([A-Z]{3})\b(?:\s+(\d{4}))?\s+(.+)`)

	// Detecta la segunda línea cuando el año baja de renglón (ej: "2026         Colombi")
	reNuYearLine = regexp.MustCompile(`^(\d{4})\s*(.*)`)

	// Mapeo de meses de texto a número
	monthMap = map[string]int{
		"ENE": 1, "FEB": 2, "MAR": 3, "ABR": 4, "MAY": 5, "JUN": 6,
		"JUL": 7, "AGO": 8, "SEP": 9, "OCT": 10, "NOV": 11, "DIC": 12,
	}
)

type NuParser struct{}

func NewNuParser() *NuParser {
	return &NuParser{}
}

func (p *NuParser) Parse(text string) ([]domain.Movement, error) {
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("texto vacío")
	}

	lines := strings.Split(text, "\n")
	var movements []domain.Movement
	defaultYear := time.Now().Year() // Por si no encontramos el año

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		matches := reNuRow.FindStringSubmatch(line)
		if len(matches) == 0 {
			continue // No es una fila de transacción
		}

		day, _ := strconv.Atoi(matches[1])
		month := monthMap[matches[2]]
		yearStr := matches[3]
		restOfLine := matches[4]

		year := defaultYear
		if yearStr != "" {
			year, _ = strconv.Atoi(yearStr)
		}

		// Separar el resto de la línea por espacios grandes
		parts := reColumns.Split(strings.TrimSpace(restOfLine), -1)
		if len(parts) < 2 {
			continue
		}

		description := strings.TrimSpace(parts[0])
		amountStr := parts[1]

		// --- Magia para Nu: Revisar la siguiente línea si el año no estaba en la primera ---
		if yearStr == "" && i+1 < len(lines) {
			nextLine := strings.TrimSpace(lines[i+1])
			if yMatch := reNuYearLine.FindStringSubmatch(nextLine); len(yMatch) > 0 {
				year, _ = strconv.Atoi(yMatch[1])

				// Si hay texto adicional (ej: "Colombi"), lo agregamos a la descripción
				if len(yMatch) > 2 && yMatch[2] != "" {
					extraDesc := strings.TrimSpace(reColumns.Split(yMatch[2], -1)[0])
					// Evitar concatenar cosas como "↪ A capital"
					if !strings.HasPrefix(extraDesc, "↪") {
						description += " " + extraDesc
					}
				}
				i++ // Saltamos la siguiente línea porque ya la procesamos
			}
		}

		// Convertir el valor monetario
		amount := parseNuNumber(amountStr)
		if amount == 0 {
			continue
		}

		// Lógica de Nu: Diferenciar gastos de pagos
		// En el extracto de Nu todo sale positivo. Los pagos o devoluciones deben ser "Ingresos" para ti.
		lowerDesc := strings.ToLower(description)

		// 1. Detectar si es un reembolso (ej: "Devolución - Mercado Pago...")
		isRefund := strings.HasPrefix(lowerDesc, "devolución") || strings.HasPrefix(lowerDesc, "devolucion")

		// 2. Detectar si es el pago de la tarjeta (ej: "Gracias por tu pago")
		isCardPayment := strings.Contains(lowerDesc, "gracias por tu pago")

		isIncome := isRefund || isCardPayment

		finalAmount := amount
		movementType := domain.TypeIncome

		if !isIncome {
			finalAmount = -amount // Es una compra, lo volvemos negativo
			movementType = domain.TypeExpense
		}

		mov := domain.Movement{
			ID:          uuid.Nil,
			Date:        time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC),
			Description: description,
			Amount:      finalAmount,
			Type:        movementType,
			CreatedAt:   time.Now(),
		}

		movements = append(movements, mov)
	}

	if len(movements) == 0 {
		return nil, fmt.Errorf("no se detectaron movimientos en Nu")
	}

	return movements, nil
}

// parseNuNumber convierte formato de Nu ($47.000,00 -> 47000.00)
func parseNuNumber(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "$")
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}

	isNegative := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")

	// Nu usa punto para miles y coma para decimales.
	// 1. Quitamos los puntos de miles
	s = strings.ReplaceAll(s, ".", "")
	// 2. Cambiamos la coma decimal por un punto (formato estándar de Go)
	s = strings.ReplaceAll(s, ",", ".")

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}

	if isNegative {
		return -f
	}
	return f
}
