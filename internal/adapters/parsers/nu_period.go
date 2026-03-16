package parsers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ExtractNuPeriodMonth(text string) (string, error) {

	reCut := regexp.MustCompile(`(?i)Fecha\s+de\s+corte\s*([0-9]{1,2})\s+([A-Z]{3})\s+([0-9]{4})`)

	if m := reCut.FindStringSubmatch(text); len(m) == 4 {

		day, _ := strconv.Atoi(m[1])
		month := monthMap[strings.ToUpper(m[2])]
		year, _ := strconv.Atoi(m[3])

		t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

		return t.Format("2006-01"), nil
	}

	// ejemplo: 31 DIC - 30 ENE 2026
	reRange := regexp.MustCompile(`([0-9]{1,2})\s+([A-Z]{3})\s*-\s*([0-9]{1,2})\s+([A-Z]{3})\s+([0-9]{4})`)

	if m := reRange.FindStringSubmatch(text); len(m) == 6 {

		month := monthMap[strings.ToUpper(m[4])]
		year, _ := strconv.Atoi(m[5])

		t := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

		return t.Format("2006-01"), nil
	}

	return "", fmt.Errorf("no se encontró fecha de corte en extracto Nu")
}
