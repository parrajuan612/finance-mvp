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

// ParseBancolombiaText convierte el texto plano extraído del PDF de Bancolombia
// en una lista de domain.Movement (campos clave: Date, Description, Amount, Type).
// Los IDs quedan como uuid.Nil porque en esta etapa aún no guardamos en BD.
func ParseBancolombiaText(text string) ([]domain.Movement, error) {
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("texto vacío")
	}

	// 1) Intentar extraer año desde encabezado "DESDE: YYYY/MM/DD   HASTA: YYYY/MM/DD"
	year := time.Now().Year()
	reYear := regexp.MustCompile(`HASTA:\s*([0-9]{4})/[0-9]{2}/[0-9]{2}`)
	if m := reYear.FindStringSubmatch(text); len(m) >= 2 {
		if y, err := strconv.Atoi(m[1]); err == nil && y > 1900 {
			year = y
		}
	}

	// 2) Preparaciones: líneas y regexes
	lines := strings.Split(text, "\n")
	// Línea que empieza con fecha d/m o dd/mm
	reDateLine := regexp.MustCompile(`^\s*([0-9]{1,2})/([0-9]{1,2})\s+(.*)$`)
	// números con miles (1,234,567.89) o sin miles (.90) ; captura -. and commas
	reNum := regexp.MustCompile(`[-+]?\d{1,3}(?:,\d{3})*(?:\.\d+)?|[-+]?\.\d+`)

	var movements []domain.Movement

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		lineTrim := strings.TrimSpace(line)
		if lineTrim == "" {
			continue
		}

		// detecta línea que inicia con fecha
		if m := reDateLine.FindStringSubmatch(line); len(m) >= 4 {
			dayStr := m[1]
			monthStr := m[2]
			rest := m[3] // aquí viene descripción + (posiblemente) cantidades

			// parsear día y mes
			day, _ := strconv.Atoi(dayStr)
			month, _ := strconv.Atoi(monthStr)

			// descripción tentativa: quitamos números comunes al final (si existen)
			// pero primero buscamos números en la misma línea
			nums := reNum.FindAllString(rest, -1)

			var amountStr string
			var balanceStr string
			description := rest

			// Heurística 1: si hay números en la línea, el último suele ser el monto (o balance).
			// Tomamos el último como amount tentativa.
			if len(nums) > 0 {
				amountStr = nums[len(nums)-1]
				// description = rest hasta la posición del amount (intento)
				idx := strings.LastIndex(rest, amountStr)
				if idx > 0 {
					description = strings.TrimSpace(rest[:idx])
				} else {
					description = strings.TrimSpace(rest)
				}
			} else {
				// Heurística 2: mirar las próximas 1-2 líneas por números (amount o balance).
				for j := i + 1; j <= i+2 && j < len(lines); j++ {
					nl := strings.TrimSpace(lines[j])
					if nl == "" {
						continue
					}
					found := reNum.FindAllString(nl, -1)
					if len(found) > 0 {
						// si no teníamos amount, tomar el primer match como amount
						if amountStr == "" {
							amountStr = found[0]
						} else if balanceStr == "" {
							balanceStr = found[len(found)-1]
						}
					}
				}
			}

			// Si amountStr está vacío: intentar buscar en la siguiente línea la primera cantidad.
			if amountStr == "" {
				if i+1 < len(lines) {
					found := reNum.FindAllString(lines[i+1], -1)
					if len(found) > 0 {
						amountStr = found[0]
					}
				}
			}

			// Limpieza y parse de amountStr
			amount := 0.0
			if amountStr != "" {
				amount = parseNumber(amountStr)
			}

			// Decidir signo según la presencia de '-' en la cadena encontrada.
			// Si el número no tenía signo, intentamos inferir: si parece grande y no tiene '-', lo dejamos positivo.
			if strings.Contains(amountStr, "-") || strings.HasPrefix(strings.TrimSpace(rest), "-") {
				amount = -abs(amount)
			}

			// formar fecha con año detectado (cuidado: si mes < startMonth, podría pertenecer a año anterior;
			// para MVP asumimos todas las operaciones dentro del rango del periodo)
			loc := time.UTC
			date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)

			// Normalizar descripción: eliminar múltiples espacios y caracteres extraños unicode
			description = normalizeSpaces(description)

			// crear movimiento
			mov := domain.Movement{
				ID:          uuid.Nil,
				UserID:      uuid.Nil,
				AccountID:   uuid.Nil,
				StatementID: nil,
				CategoryID:  uuid.Nil,
				Date:        date,
				Description: description,
				Amount:      amount,
				Type:        domain.TypeExpense,
				CreatedAt:   time.Now(),
			}
			if amount >= 0 {
				mov.Type = domain.TypeIncome
			} else {
				// monto negativo -> expense
				mov.Type = domain.TypeExpense
			}

			movements = append(movements, mov)
		}
	}

	if len(movements) == 0 {
		return nil, fmt.Errorf("no se detectaron movimientos (se aplicaron heurísticas). Revisa el texto o ajusta el parser")
	}

	return movements, nil
}

// parseNumber quita comas y convierte a float64. Maneja formatos como ".90" -> 0.90
func parseNumber(s string) float64 {
	s = strings.TrimSpace(s)
	// normalizar signos y punto decimal
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
	}
	// quitar signos +/-
	s = strings.TrimPrefix(s, "+")
	s = strings.TrimPrefix(s, "-")
	// quitar comas de miles
	s = strings.ReplaceAll(s, ",", "")
	// si empieza con '.' añadir 0
	if strings.HasPrefix(s, ".") {
		s = "0" + s
	}
	// si está vacío -> 0
	if s == "" {
		return 0.0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	if neg {
		return -f
	}
	return f
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func normalizeSpaces(s string) string {
	// reemplaza múltiples espacios por uno, trim
	// también normaliza algunos caracteres NO-ASCII comunes del PDF
	s = strings.ReplaceAll(s, "\t", " ")
	// quitar caracteres de control raros
	s = strings.ReplaceAll(s, "\u00A0", " ") // NBSP
	// compactar espacios
	reMulti := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(reMulti.ReplaceAllString(s, " "))
}
