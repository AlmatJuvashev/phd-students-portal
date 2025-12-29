package db_test

import (
	"os"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/db"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedis_Nil(t *testing.T) {
	os.Unsetenv("REDIS_URL")
	client := db.NewRedis()
	assert.Nil(t, client)

	// Test nil safety
	db.CacheSet(nil, "key", "val", 0)
	val, err := db.CacheGet(nil, "key")
	assert.Equal(t, redis.Nil, err)
	assert.Empty(t, val)
}

func TestRedis_InvalidURL(t *testing.T) {
	os.Setenv("REDIS_URL", "invalid url")
	defer os.Unsetenv("REDIS_URL")
	
	client := db.NewRedis()
	assert.Nil(t, client)
}
