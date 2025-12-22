package services

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type CommentService struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) Create(ctx context.Context, comment models.Comment) (string, error) {
	return s.repo.Create(ctx, comment)
}

func (s *CommentService) GetByDocumentID(ctx context.Context, tenantID string, docID string) ([]models.Comment, error) {
	return s.repo.GetByDocumentID(ctx, tenantID, docID)
}
