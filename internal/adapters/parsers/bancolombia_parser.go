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

type BancolombiaParser struct{}

func NewBancolombiaParser() *BancolombiaParser {
	return &BancolombiaParser{}
}

func (p *BancolombiaParser) Parse(text string) ([]domain.Movement, error) {
	return ParseBancolombiaText(text)
}

// ParseBancolombiaText_v3: agrupa por bloques (línea con fecha + siguientes líneas),
// extrae números y decide cuál es VALOR y cuál SALDO con heurísticas basadas en prevBalance.
func ParseBancolombiaText(text string) ([]domain.Movement, error) {
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("texto vacío")
	}

	year := time.Now().Year()
	reYear := regexp.MustCompile(`HASTA:\s*([0-9]{4})/[0-9]{2}/[0-9]{2}`)
	if m := reYear.FindStringSubmatch(text); len(m) >= 2 {
		if y, err := strconv.Atoi(m[1]); err == nil && y > 1900 {
			year = y
		}
	}

	lines := strings.Split(text, "\n")

	// detecta línea que empieza con fecha d/m o dd/mm
	reDateLine := regexp.MustCompile(`^\s*([0-9]{1,2})/([0-9]{1,2})\s+(.*)$`)
	// regex para extraer números (miles con comas y decimales)
	reNum := regexp.MustCompile(`[-+]?\d{1,3}(?:,\d{3})*(?:\.\d+)?|[-+]?\.\d+`)
	// regex para detectar si una línea tiene sólo números / espacios (nos ayuda a anexar balances en líneas separadas)
	reNumericLine := regexp.MustCompile(`^[\s\d,.\-+]+$`)
	// regex para recortar números finales de una descripción
	reTrailingNums := regexp.MustCompile(`(?:\s+[-+]?\d{1,3}(?:,\d{3})*(?:\.\d+)?\s*)+$`)

	type block struct {
		day   int
		month int
		lines []string
	}

	// 1) construir bloques por transacción (cada bloque comienza con una línea que contiene la fecha)
	var blocks []block
	var curr *block
	for i := 0; i < len(lines); i++ {
		ln := lines[i]
		if strings.TrimSpace(ln) == "" {
			continue
		}
		if m := reDateLine.FindStringSubmatch(ln); len(m) >= 4 {
			// nueva transacción
			day, _ := strconv.Atoi(m[1])
			month, _ := strconv.Atoi(m[2])
			curr = &block{day: day, month: month, lines: []string{m[3]}}
			blocks = append(blocks, *curr)
			// anexar siguientes 1-2 líneas si son "numéricas" o cortas y con números (pueden ser saldo/valor en otra línea)
			for k := 1; k <= 2 && i+k < len(lines); k++ {
				next := strings.TrimSpace(lines[i+k])
				if next == "" {
					continue
				}
				// si la línea siguiente es mayormente numérica o contiene números y no es muy larga, la anexamos
				if reNumericLine.MatchString(next) || (len(reNum.FindAllString(next, -1)) >= 1 && len(next) < 80) {
					curr.lines = append(curr.lines, next)
					// marcar como consumida
					i += 1
				} else {
					break
				}
			}
		} else {
			// línea que no inicia con fecha: si hay un bloque activo, anexar si parece pertenecer
			if curr != nil {
				trim := strings.TrimSpace(ln)
				// anexar sólo si la línea no es un encabezado / footer evidente
				if trim != "" {
					curr.lines = append(curr.lines, trim)
				}
			}
		}
	}

	var movements []domain.Movement
	var prevBalance *float64

	// 2) procesar cada bloque
	for _, b := range blocks {
		joined := strings.Join(b.lines, " ")
		// extraer todos los números del bloque
		numStrs := reNum.FindAllString(joined, -1)
		nums := make([]float64, 0, len(numStrs))
		for _, s := range numStrs {
			nums = append(nums, parseNumber(s))
		}

		var value float64
		var foundValue bool
		var balance float64
		var foundBalance bool

		// heurística para decidir SALDO/VALOR
		if len(nums) >= 2 {
			// si existe prevBalance, elegir como SALDO el número más cercano a prevBalance
			if prevBalance != nil {
				// encontrar índice con mínima diferencia respecto a prevBalance
				minIdx := 0
				minDiff := abs(nums[0] - *prevBalance)
				for idx := 1; idx < len(nums); idx++ {
					d := abs(nums[idx] - *prevBalance)
					if d < minDiff {
						minDiff = d
						minIdx = idx
					}
				}
				// asumimos que el valor es el número inmediatamente anterior al saldo si existe
				balance = nums[minIdx]
				foundBalance = true
				if minIdx-1 >= 0 {
					value = nums[minIdx-1]
					foundValue = true
				} else {
					// si no hay anterior, intentamos inferir value = balance - prevBalance
					valCand := balance - *prevBalance
					if !isAbsHuge(valCand) {
						value = valCand
						foundValue = true
					}
				}
			} else {
				// no hay prevBalance: asumimos el último número es SALDO y el anterior es VALOR
				balance = nums[len(nums)-1]
				foundBalance = true
				value = nums[len(nums)-2]
				foundValue = true
			}
		} else if len(nums) == 1 {
			// sólo 1 número: preferimos interpretarlo como SALDO si prevBalance existe (y calcular value)
			one := nums[0]
			if prevBalance != nil {
				balance = one
				foundBalance = true
				amt := balance - *prevBalance
				if !isAbsHuge(amt) {
					value = amt
					foundValue = true
				}
			} else {
				// sin prevBalance lo tratamos como valor directo
				value = one
				foundValue = true
			}
		} else {
			// no hay números: saltamos
			continue
		}

		// Si no encontramos un valor válido, saltar
		if !foundValue {
			// actualizar prevBalance si encontramos balance
			if foundBalance {
				prevBalance = &balance
			}
			continue
		}

		// Normalizar descripción: eliminar números finales y compactar espacios
		desc := joined
		// eliminar trailing numeric tokens
		desc = reTrailingNums.ReplaceAllString(desc, "")
		desc = normalizeSpaces(desc)

		// fecha
		date := time.Date(year, time.Month(b.month), b.day, 0, 0, 0, 0, time.UTC)

		typ := domain.TypeIncome
		if value < 0 {
			typ = domain.TypeExpense
		}

		mov := domain.Movement{
			ID:          uuid.Nil,
			UserID:      uuid.Nil,
			AccountID:   uuid.Nil,
			StatementID: nil,
			CategoryID:  uuid.Nil,
			Date:        date,
			Description: desc,
			Amount:      value,
			Type:        typ,
			CreatedAt:   time.Now(),
		}
		movements = append(movements, mov)

		// actualizar prevBalance
		if foundBalance {
			prevBalance = &balance
		} else {
			// si no había balance y pudimos inferirlo (por ejemplo si tuvimos prevBalance y value calculado),
			// podemos estimar nuevo balance = prevBalance + value
			if prevBalance != nil {
				est := *prevBalance + value
				prevBalance = &est
			}
		}
	}

	if len(movements) == 0 {
		return nil, fmt.Errorf("no se detectaron movimientos")
	}
	return movements, nil
}

// parseNumber quita comas y convierte a float64. Maneja formatos como ".90" -> 0.90
func parseNumber(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
	}
	s = strings.TrimPrefix(s, "+")
	s = strings.TrimPrefix(s, "-")
	s = strings.ReplaceAll(s, ",", "")
	if strings.HasPrefix(s, ".") {
		s = "0" + s
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
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
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.ReplaceAll(s, "\u00A0", " ")
	reMulti := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(reMulti.ReplaceAllString(s, " "))
}

// isAbsHuge evita aceptar diferencias absurdas (umbral arbitrario grande)
func isAbsHuge(v float64) bool {
	if v < 0 {
		v = -v
	}
	return v > 1e9
}
