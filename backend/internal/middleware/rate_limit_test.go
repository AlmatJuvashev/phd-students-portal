package middleware

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginRateLimiter_WithMiniredis(t *testing.T) {
	// Setup miniredis
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	// Connect go-redis to the mock
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	rl := NewLoginRateLimiter(rdb)
	ctx := context.Background()
	identifier := "user@example.com"

	t.Run("Initially Allowed", func(t *testing.T) {
		allowed, ttl, err := rl.CheckAllowed(ctx, identifier)
		require.NoError(t, err)
		assert.True(t, allowed)
		assert.Zero(t, ttl)
	})

	t.Run("Increment Failures", func(t *testing.T) {
		// record 4 times
		for i := 0; i < 4; i++ {
			err := rl.RecordFailure(ctx, identifier)
			require.NoError(t, err)
		}

		// Should still be allowed
		allowed, _, err := rl.CheckAllowed(ctx, identifier)
		require.NoError(t, err)
		assert.True(t, allowed)
	})

	t.Run("Lockout Triggered", func(t *testing.T) {
		// record 5th failure
		err := rl.RecordFailure(ctx, identifier)
		require.NoError(t, err)

		// Check Lockout
		allowed, ttl, err := rl.CheckAllowed(ctx, identifier)
		require.NoError(t, err)
		assert.False(t, allowed)
		assert.Greater(t, ttl, time.Duration(0))
		assert.LessOrEqual(t, ttl, 30*time.Minute)
	})

	t.Run("Reset Clears Lockout", func(t *testing.T) {
		rl.Reset(ctx, identifier)

		allowed, ttl, err := rl.CheckAllowed(ctx, identifier)
		require.NoError(t, err)
		assert.True(t, allowed)
		assert.Zero(t, ttl)
	})

	t.Run("Redis Connection Error Handling", func(t *testing.T) {
		// Simulate nil client
		rlNil := NewLoginRateLimiter(nil)
		
		// CheckAllowed returns true (fail open)
		allowed, _, err := rlNil.CheckAllowed(ctx, "fail-open")
		assert.NoError(t, err)
		assert.True(t, allowed)

		// RecordFailure no-op
		err = rlNil.RecordFailure(ctx, "fail-open")
		assert.NoError(t, err)

		// Reset no-op
		rlNil.Reset(ctx, "fail-open")
	})
	
	t.Run("Redis Error Simulation", func(t *testing.T) {
		// Create a separate miniredis instance then kill it to simulate connection failure
		mr2, err := miniredis.Run()
		require.NoError(t, err)
		
		rdb2 := redis.NewClient(&redis.Options{
			Addr: mr2.Addr(),
		})
		
		rl2 := NewLoginRateLimiter(rdb2)
		mr2.Close() // Kill redis
		
		// CheckAllowed -> Should fail open (return true) but return error? 
		// Code logic: } else if err != nil { return true, 0, err }
		
		allowed, _, err := rl2.CheckAllowed(ctx, "redis-down")
		assert.Error(t, err) // Expect error
		assert.True(t, allowed) // Fail open confirmed
		
		err = rl2.RecordFailure(ctx, "redis-down")
		assert.Error(t, err)
	})
}
