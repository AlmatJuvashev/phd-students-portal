package services_test

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestContactService_Unit(t *testing.T) {
	mockRepo := NewMockContactRepository()
	svc := services.NewContactService(mockRepo)
	ctx := context.Background()

	t.Run("ListPublic", func(t *testing.T) {
		_, _ = svc.ListPublic(ctx, "t1")
	})

	t.Run("ListAdmin", func(t *testing.T) {
		_, _ = svc.ListAdmin(ctx, "t1", true)
	})

	t.Run("Create", func(t *testing.T) {
		_, _ = svc.Create(ctx, "t1", models.Contact{})
	})

	t.Run("Update", func(t *testing.T) {
		// Empty updates should return nil immediately
		err := svc.Update(ctx, "t1", "c1", nil)
		assert.NoError(t, err)

		// Non-empty updates
		_ = svc.Update(ctx, "t1", "c1", map[string]interface{}{"name": "New Name"})
	})

	t.Run("Delete", func(t *testing.T) {
		_ = svc.Delete(ctx, "t1", "c1")
	})

	assert.NotNil(t, svc)
}
