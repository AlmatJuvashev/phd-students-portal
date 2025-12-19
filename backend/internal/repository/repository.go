package repository

import (
	"errors"
)

var (
	ErrNotFound = errors.New("record not found")
)

// Pagination defines standard pagination parameters
type Pagination struct {
	Limit  int
	Offset int
}

// UserFilter defines criteria for listing users
type UserFilter struct {
	Role       string
	Program    string
	Specialty  string
	Department string
	Cohort     string
	Active     *bool // nil = all, true = active, false = inactive
	Search     string
}
