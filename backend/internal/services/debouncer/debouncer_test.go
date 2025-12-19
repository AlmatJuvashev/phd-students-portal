package debouncer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNotificationDebouncer(t *testing.T) {
	// Test that creating a debouncer doesn't panic, even without Redis
	d := NewNotificationDebouncer()
	assert.NotNil(t, d)
	// When Redis is not available, disabled should be true
	assert.True(t, d.disabled, "Expected debouncer to be disabled when Redis is unavailable")
	d.Close()
}

func TestBuildKey(t *testing.T) {
	d := &NotificationDebouncer{}
	key := d.buildKey("user-123", "node-abc", "document_submitted")
	expected := "notif:debounce:user-123:node-abc:document_submitted"
	assert.Equal(t, expected, key)
}

func TestShouldNotify_DisabledDebouncer(t *testing.T) {
	// When debouncer is disabled (no Redis), should always return true
	d := &NotificationDebouncer{disabled: true}
	ctx := context.Background()

	result := d.ShouldNotify(ctx, "user-1", "node-1", "event-1")
	assert.True(t, result, "Disabled debouncer should always allow notifications")

	// Multiple calls should still return true
	result = d.ShouldNotify(ctx, "user-1", "node-1", "event-1")
	assert.True(t, result, "Disabled debouncer should always allow notifications")
}

func TestClose_NilClient(t *testing.T) {
	d := &NotificationDebouncer{client: nil}
	err := d.Close()
	assert.NoError(t, err, "Close should not error with nil client")
}
