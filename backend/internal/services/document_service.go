package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type DocumentService struct {
	repo    repository.DocumentRepository
	storage StorageClient
	cfg     config.AppConfig
}

func NewDocumentService(repo repository.DocumentRepository, cfg config.AppConfig, storage StorageClient) *DocumentService {
	return &DocumentService{
		repo:    repo,
		storage: storage,
		cfg:     cfg,
	}
}

type CreateDocumentRequest struct {
	Title    string
	Kind     string
	TenantID string
	UserID   string
}

func (s *DocumentService) CreateMetadata(ctx context.Context, req CreateDocumentRequest) (string, error) {
	doc := &models.Document{
		Title:    req.Title,
		Kind:     req.Kind,
		TenantID: req.TenantID,
		UserID:   req.UserID,
	}
	return s.repo.Create(ctx, doc)
}

func (s *DocumentService) CreateVersion(ctx context.Context, docID string, tenantID string, uploaderID string, fileMeta models.DocumentVersion) (string, error) {
	// 1. Verify document exists/access? (skipped for now, relying on handler/middleware)
	
	// 2. Insert Version in DB
	fileMeta.DocumentID = docID
	fileMeta.TenantID = tenantID
	fileMeta.UploadedBy = uploaderID
	
	verID, err := s.repo.CreateVersion(ctx, &fileMeta)
	if err != nil {
		return "", err
	}
	
	// 3. Update current version
	if err := s.repo.SetCurrentVersion(ctx, docID, verID); err != nil {
		return "", err
	}
	
	return verID, nil
}

func (s *DocumentService) GetDocumentDetails(ctx context.Context, id string) (*models.Document, []models.DocumentVersion, error) {
	doc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	
	versions, err := s.repo.GetVersionsByDocumentID(ctx, id)
	if err != nil && err != repository.ErrNotFound {
		return nil, nil, err
	}
	
	return doc, versions, nil
}

func (s *DocumentService) ListUserDocuments(ctx context.Context, userID string) ([]models.Document, error) {
	return s.repo.ListByUserID(ctx, userID)
}

func (s *DocumentService) DeleteDocument(ctx context.Context, id string) error {
	// TODO: Cleanup S3 files? For now just soft/hard delete from DB.
	return s.repo.Delete(ctx, id)
}

// S3 Operations

func (s *DocumentService) PresignUpload(ctx context.Context, docID string, filename string, contentType string) (string, string, error) {
	if s.storage == nil {
		return "", "", errors.New("storage not configured")
	}
	
	if err := ValidateContentType(contentType); err != nil {
		return "", "", err
	}
	
	// Key structure: {docID}/{filename}
	key := fmt.Sprintf("%s/%s", docID, filename)
	expires := GetPresignExpires()
	
	url, err := s.storage.PresignPut(ctx, key, contentType, expires)
	if err != nil {
		return "", "", err
	}
	
	return url, key, nil
}

func (s *DocumentService) PresignDownload(ctx context.Context, verID string) (string, error) {
	if s.storage == nil {
		return "", errors.New("storage not configured")
	}
	
	ver, err := s.repo.GetVersion(ctx, verID)
	if err != nil {
		return "", err
	}
	
	if !ver.ObjectKey.Valid || ver.ObjectKey.String == "" {
		return "", errors.New("version is not stored in storage")
	}
	
	expires := GetPresignExpires()
	return s.storage.PresignGet(ctx, ver.ObjectKey.String, expires)
}

func (s *DocumentService) PresignLatestDownload(ctx context.Context, docID string) (string, error) {
	if s.storage == nil {
		return "", errors.New("storage not configured")
	}
	
	ver, err := s.repo.GetLatestVersion(ctx, docID)
	if err != nil {
		return "", err
	}
	
	if !ver.ObjectKey.Valid || ver.ObjectKey.String == "" {
		return "", errors.New("latest version is not stored in storage")
	}
	
	expires := GetPresignExpires()
	return s.storage.PresignGet(ctx, ver.ObjectKey.String, expires)
}

// GetStoragePath returns the local path or S3 info for a version
// Returns: storagePath, bucket, key, mimetype, size, error
func (s *DocumentService) GetVersionFile(ctx context.Context, verID string) (*models.DocumentVersion, error) {
	return s.repo.GetVersion(ctx, verID)
}

// Helper to check storage client availability
func (s *DocumentService) IsS3Configured() bool {
	return s.storage != nil
}

// Helper to get bucket name
func (s *DocumentService) GetS3Bucket() string {
	if s.storage != nil {
		return s.storage.Bucket()
	}
	return ""
}
