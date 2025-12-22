package services

import (
	"strings"

	"github.com/lib/pq"
)

// IsDuplicateKeyError checks if an error is a postgres unique constraint violation
func IsDuplicateKeyError(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505" // unique_violation
	}
	// Fallback string check if driver changes or wrapping occurs
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}
