package services_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestDocumentService_Unit(t *testing.T) {
	mockRepo := NewMockDocumentRepository()
	mockStorage := &services.MockStorageClient{}
	svc := services.NewDocumentService(mockRepo, config.AppConfig{}, mockStorage)
	ctx := context.Background()

	t.Run("CreateMetadata", func(t *testing.T) {
		_, _ = svc.CreateMetadata(ctx, services.CreateDocumentRequest{Title: "T"})
	})

	t.Run("CreateVersion", func(t *testing.T) {
		_, _ = svc.CreateVersion(ctx, "d1", "t1", "u1", models.DocumentVersion{})
	})

	t.Run("Details", func(t *testing.T) {
		_, _, _ = svc.GetDocumentDetails(ctx, "d1")
		_, _ = svc.ListUserDocuments(ctx, "u1")
		_ = svc.DeleteDocument(ctx, "d1")
	})

	t.Run("Presign Errors", func(t *testing.T) {
		svcNoStorage := services.NewDocumentService(mockRepo, config.AppConfig{}, nil)
		_, _, err := svcNoStorage.PresignUpload(ctx, "d1", "f", "c")
		assert.Error(t, err)
		assert.Equal(t, "storage not configured", err.Error())

		_, err = svcNoStorage.PresignDownload(ctx, "v1")
		assert.Error(t, err)

		_, err = svcNoStorage.PresignLatestDownload(ctx, "d1")
		assert.Error(t, err)

		mockRepo.GetVersionFunc = func(ctx context.Context, id string) (*models.DocumentVersion, error) {
			return nil, assert.AnError
		}
		_, err = svc.PresignDownload(ctx, "v1")
		assert.Error(t, err)
	})

	t.Run("Presign", func(t *testing.T) {
		mockRepo.GetVersionFunc = func(ctx context.Context, id string) (*models.DocumentVersion, error) {
			return &models.DocumentVersion{ObjectKey: sql.NullString{String: "key", Valid: true}}, nil
		}
		mockRepo.GetLatestVersionFunc = func(ctx context.Context, id string) (*models.DocumentVersion, error) {
			return &models.DocumentVersion{ObjectKey: sql.NullString{String: "key", Valid: true}}, nil
		}
		_, _, _ = svc.PresignUpload(ctx, "d1", "f", "application/pdf")
		_, _ = svc.PresignDownload(ctx, "v1")
		_, _ = svc.PresignLatestDownload(ctx, "d1")
	})

	t.Run("Helpers", func(t *testing.T) {
		_, _ = svc.GetVersionFile(ctx, "v1")
		assert.True(t, svc.IsS3Configured())
		assert.NotEmpty(t, svc.GetS3Bucket())
		
		svcEmpty := services.NewDocumentService(mockRepo, config.AppConfig{}, nil)
		assert.False(t, svcEmpty.IsS3Configured())
		assert.Empty(t, svcEmpty.GetS3Bucket())
	})

	assert.NotNil(t, svc)
}
