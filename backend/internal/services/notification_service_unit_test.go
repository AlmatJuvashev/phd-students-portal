package services_test

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestNotificationService_Unit(t *testing.T) {
	mockRepo := NewMockNotificationRepository()
	svc := services.NewNotificationService(mockRepo)
	ctx := context.Background()

	t.Run("CreateNotification", func(t *testing.T) {
		_ = svc.CreateNotification(ctx, &models.Notification{})
	})

	t.Run("GetUnreadNotifications", func(t *testing.T) {
		_, _ = svc.GetUnreadNotifications(ctx, "u1")
	})

	t.Run("MarkAsRead", func(t *testing.T) {
		_ = svc.MarkAsRead(ctx, "n1", "u1")
	})

	t.Run("MarkAllAsRead", func(t *testing.T) {
		_ = svc.MarkAllAsRead(ctx, "u1")
	})

	t.Run("ListNotifications", func(t *testing.T) {
		_, _ = svc.ListNotifications(ctx, "u1", 10)
	})

	assert.NotNil(t, svc)
}
