package services_test

import (
	"errors"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestIsDuplicateKeyError(t *testing.T) {
	// 1. PQ Error
	pqErr := &pq.Error{Code: "23505"}
	assert.True(t, services.IsDuplicateKeyError(pqErr))

	pqErrOther := &pq.Error{Code: "12345"}
	assert.False(t, services.IsDuplicateKeyError(pqErrOther))

	// 2. String check
	errStr := errors.New("duplicate key value violates unique constraint")
	assert.True(t, services.IsDuplicateKeyError(errStr))

	errNormal := errors.New("something else")
	assert.False(t, services.IsDuplicateKeyError(errNormal))
}
