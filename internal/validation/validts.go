package validation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	unitDuration = regexp.MustCompile(`(?i)^(\d+\s*(?:h|hr|hrs|hour|hours))?\s*(\d+\s*(?:m|min|mins|minute|minutes))?\s*(\d+\s*(?:s|sec|secs|second|seconds))?$`)
	bareNumber   = regexp.MustCompile(`^\d+$`)
	mmss         = regexp.MustCompile(`^\d{1,3}:\d{2}$`)
	hhmmss       = regexp.MustCompile(`^\d+:\d{2}:\d{2}$`)
)

func SanitiseTimestamp(input string) (string, bool) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", false
	}

	// Bare number means seconds.
	if bareNumber.MatchString(input) {
		n, err := strconv.Atoi(input)
		if err != nil {
			return "", false
		}
		return formatYTDLP(n), true
	}

	// Accept mm:ss.
	if mmss.MatchString(input) {
		parts := strings.Split(input, ":")
		minutes, _ := strconv.Atoi(parts[0])
		seconds, _ := strconv.Atoi(parts[1])
		if seconds > 59 {
			return "", false
		}
		return formatYTDLP(minutes*60 + seconds), true
	}

	// Accept hh:mm:ss.
	if hhmmss.MatchString(input) {
		parts := strings.Split(input, ":")
		hours, _ := strconv.Atoi(parts[0])
		minutes, _ := strconv.Atoi(parts[1])
		seconds, _ := strconv.Atoi(parts[2])
		if minutes > 59 || seconds > 59 {
			return "", false
		}
		return formatYTDLP(hours*3600 + minutes*60 + seconds), true
	}

	// Parse unit-based durations.
	if !unitDuration.MatchString(input) {
		return "", false
	}

	// Normalize whitespace and case before extracting components.
	normalized := strings.ToLower(strings.ReplaceAll(input, " ", ""))
	var total int
	var num strings.Builder
	for _, r := range normalized {
		if r >= '0' && r <= '9' {
			num.WriteRune(r)
			continue
		}
		n, err := strconv.Atoi(num.String())
		if err != nil {
			return "", false
		}
		switch r {
		case 'h':
			total += n * 3600
		case 'm':
			total += n * 60
		case 's':
			total += n
		default:
			// Skip letters in units like "hours", "minutes", "seconds".
		}
		num.Reset()
	}
	if num.Len() > 0 {
		return "", false
	}

	return formatYTDLP(total), true
}

func formatYTDLP(seconds int) string {
	h := seconds / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60
	return fmt.Sprintf("*%02d:%02d:%02d-inf", h, m, s)
}

func IsValidTimestamp(input string) bool {
	var regex = regexp.MustCompile(`^\*\d{2}:\d{2}:\d{2}-inf$`)
	return regex.MatchString(input)
}
