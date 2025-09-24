package db

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedis returns a Redis client if REDIS_URL set; otherwise nil.
func NewRedis() *redis.Client {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		return nil
	}
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil
	}
	c := redis.NewClient(opt)
	// quick ping
	_ = c.Ping(context.Background()).Err()
	return c
}

// CacheSet sets a key with TTL.
func CacheSet(r *redis.Client, key string, val string, ttl time.Duration) {
	if r == nil {
		return
	}
	_ = r.Set(context.Background(), key, val, ttl).Err()
}

// CacheGet gets a key.
func CacheGet(r *redis.Client, key string) (string, error) {
	if r == nil {
		return "", redis.Nil
	}
	return r.Get(context.Background(), key).Result()
}
