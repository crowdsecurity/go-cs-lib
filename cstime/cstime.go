package cstime

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseDuration parses a string representing a duration, and supports
// days as a unit (e.g., "2d", "2d3h", "24h", "2h45m").
func ParseDurationWithDays(input string) (time.Duration, error) {
	var total time.Duration

	s := strings.TrimSpace(input)

	if s == "" {
		return 0, errors.New("empty duration string")
	}

	negative := false
	if strings.HasPrefix(s, "-") {
		negative = true
		s = s[1:]
	}

	// Check if 'd' is present and only at the start
	if i := strings.IndexByte(s, 'd'); i != -1 {
		dayPart := s[:i]
		rest := s[i+1:]

		// 'd' must be at the start, so rest must not contain another 'd'
		if strings.Contains(rest, "d") {
			return 0, fmt.Errorf("invalid duration %q: 'd' can appear only as the first unit", input)
		}

		// Day part must be a valid integer
		dayVal, err := strconv.Atoi(dayPart)
		if err != nil {
			return 0, fmt.Errorf("invalid day value in duration %q", input)
		}

		total += time.Duration(dayVal) * 24 * time.Hour
		s = rest
	}

	if strings.Contains(s, "-") {
		return 0, fmt.Errorf("misplaced negative sign in duration %q", input)
	}

	if s != "" {
		// Delegate remaining part to time.ParseDuration
		dur, err := time.ParseDuration(s)
		if err != nil {
			return 0, err
		}
		total += dur
	}

	if negative {
		total = -total
	}

	return total, nil
}
