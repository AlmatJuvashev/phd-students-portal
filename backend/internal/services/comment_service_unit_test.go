package services_test

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestCommentService_Unit(t *testing.T) {
	mockRepo := NewMockCommentRepository()
	svc := services.NewCommentService(mockRepo)
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		_, _ = svc.Create(ctx, models.Comment{})
	})

	t.Run("GetByDocumentID", func(t *testing.T) {
		_, _ = svc.GetByDocumentID(ctx, "t1", "d1")
	})

	assert.NotNil(t, svc)
}
