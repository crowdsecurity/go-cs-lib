package cstime

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseDuration parses a string representing a duration, and supports
// days as a unit (e.g., "2d", "2d3h", "24h", "2h45m")
func ParseDuration(s string) (time.Duration, error) {
	daysDuration := time.Duration(0)

	if strings.Contains(s, "d") {
		parts := strings.Split(s, "d")
		if len(parts) != 2 {
			return 0, fmt.Errorf("invalid duration format '%s'", s)
		}
		
		days, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}

		daysDuration = time.Hour * time.Duration(24 * days)

		s = parts[1]

		if s == "" {
			return daysDuration, nil
		}
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, err
	}
		
	if (daysDuration < 0) {
		return daysDuration - d, nil
	}

	return daysDuration + d, nil
}
