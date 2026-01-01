package solver

import (
	"strconv"
	"strings"
)

// parseTimeMins converts "HH:MM" or "HH:MM:SS" to minutes from midnight
func parseTimeMins(tS string) int {
	parts := strings.Split(tS, ":")
	if len(parts) < 2 { return 0 }
	h, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	return h*60 + m
}
