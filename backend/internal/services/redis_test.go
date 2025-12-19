package services_test

import (
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestNewRedis(t *testing.T) {
	// Test Case 1: Valid URL
	// NewRedis parses the URL and returns a client. It does not connect immediately,
	// so this safe to run without a real Redis instance for unit testing logic.
	client := services.NewRedis("redis://localhost:6379/0")
	assert.NotNil(t, client)
	
	// Test Case 2: Invalid URL
	// Should return nil
	clientInvalid := services.NewRedis("not-a-redis-url")
	assert.Nil(t, clientInvalid)
}
