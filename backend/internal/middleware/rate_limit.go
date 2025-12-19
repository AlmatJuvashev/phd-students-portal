package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// LoginRateLimiter handles logic for tracking failed login attempts
type LoginRateLimiter struct {
	rds *redis.Client
}

func NewLoginRateLimiter(rds *redis.Client) *LoginRateLimiter {
	return &LoginRateLimiter{rds: rds}
}

// CheckAllowed checks if the user/IP is currently locked out.
// Returns active=true if allowed, false if locked.
// Also returns time until unlock (0 if allowed).
func (rl *LoginRateLimiter) CheckAllowed(ctx context.Context, identifier string) (bool, time.Duration, error) {
	key := fmt.Sprintf("rate_limit:login:%s", identifier)
	
	val, err := rl.rds.Get(ctx, key).Int()
	if err == redis.Nil {
		return true, 0, nil
	} else if err != nil {
		return true, 0, err // Fail open on redis error? Or fail closed? standardized on fail open for availability
	}

	if val >= 5 {
		// Get TTL
		ttl, err := rl.rds.TTL(ctx, key).Result()
		if err != nil {
			return false, 30 * time.Minute, nil
		}
		return false, ttl, nil
	}

	return true, 0, nil
}

// RecordFailure increments the failure count for this identifier.
// On 1st failure, sets expiry to 30 mins.
func (rl *LoginRateLimiter) RecordFailure(ctx context.Context, identifier string) error {
	key := fmt.Sprintf("rate_limit:login:%s", identifier)
	
	// Increment
	val, err := rl.rds.Incr(ctx, key).Result()
	if err != nil {
		return err
	}

	// If this is the first failure (val=1), set expiration window
	if val == 1 {
		rl.rds.Expire(ctx, key, 30*time.Minute)
	}

	return nil
}

// Reset clears the failure count (e.g. on successful login)
func (rl *LoginRateLimiter) Reset(ctx context.Context, identifier string) {
	key := fmt.Sprintf("rate_limit:login:%s", identifier)
	rl.rds.Del(ctx, key)
}
