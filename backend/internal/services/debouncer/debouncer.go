package debouncer

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type NotificationDebouncer struct {
	client   *redis.Client
	window   time.Duration
	disabled bool
}

func NewNotificationDebouncer() *NotificationDebouncer {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6381"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	windowMinutes := 10 // Default: 10 minutes
	if env := os.Getenv("NOTIFICATION_DEBOUNCE_MINUTES"); env != "" {
		if parsed, err := time.ParseDuration(env + "m"); err == nil {
			windowMinutes = int(parsed.Minutes())
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	disabled := false
	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Printf("Redis connection failed (debouncing disabled): %v\n", err)
		disabled = true
	}

	return &NotificationDebouncer{
		client:   client,
		window:   time.Duration(windowMinutes) * time.Minute,
		disabled: disabled,
	}
}

// ShouldNotify checks if a notification should be sent based on debouncing rules
// Returns true if notification should be sent, false if it should be skipped
func (d *NotificationDebouncer) ShouldNotify(ctx context.Context, userID string, nodeID, eventType string) bool {
	if d.disabled {
		return true // Always notify if Redis is unavailable
	}

	key := d.buildKey(userID, nodeID, eventType)

	// Check if key exists (notification was sent recently)
	exists, err := d.client.Exists(ctx, key).Result()
	if err != nil {
		fmt.Printf("Redis error checking debounce key: %v\n", err)
		return true // On error, allow notification
	}

	if exists > 0 {
		// Key exists, skip notification
		return false
	}

	// Key doesn't exist, set it and allow notification
	if err := d.client.Set(ctx, key, time.Now().Unix(), d.window).Err(); err != nil {
		fmt.Printf("Redis error setting debounce key: %v\n", err)
	}

	return true
}

func (d *NotificationDebouncer) buildKey(userID string, nodeID, eventType string) string {
	return fmt.Sprintf("notif:debounce:%s:%s:%s", userID, nodeID, eventType)
}

func (d *NotificationDebouncer) Close() error {
	if d.client != nil {
		return d.client.Close()
	}
	return nil
}
