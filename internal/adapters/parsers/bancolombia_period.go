// file: internal/adapters/parsers/bancolombia_period.go
package parsers

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func ExtractPeriodMonth(text string) (string, error) {

	reYMD := regexp.MustCompile(`(?i)HASTA:\s*([0-9]{4})[/-]([0-9]{1,2})[/-]([0-9]{1,2})`)
	if m := reYMD.FindStringSubmatch(text); len(m) == 4 {
		year, _ := strconv.Atoi(m[1])
		month, _ := strconv.Atoi(m[2])
		day, _ := strconv.Atoi(m[3])
		t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		return t.Format("2006-01"), nil
	}
	reDMY := regexp.MustCompile(`(?i)HASTA:\s*([0-9]{1,2})[/-]([0-9]{1,2})[/-]([0-9]{4})`)
	if m := reDMY.FindStringSubmatch(text); len(m) == 4 {
		day, _ := strconv.Atoi(m[1])
		month, _ := strconv.Atoi(m[2])
		year, _ := strconv.Atoi(m[3])
		t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		return t.Format("2006-01"), nil
	}
	reRange := regexp.MustCompile(`(?i)DESDE:\s*([0-9]{1,4}[\/\-][0-9]{1,2}[\/\-][0-9]{1,4}).{0,40}?HASTA:\s*([0-9]{1,4}[\/\-][0-9]{1,2}[\/\-][0-9]{1,4})`)
	if m := reRange.FindStringSubmatch(text); len(m) >= 3 {
		// intentamos parsear la segunda fecha robustamente usando los dos regexes anteriores
		if pm, err := ExtractPeriodMonth("HASTA: " + m[2]); err == nil {
			return pm, nil
		}
	}

	return "", fmt.Errorf("no se encontró fecha 'HASTA' en el texto")
}
