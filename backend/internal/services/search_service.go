package services

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type SearchService struct {
	repo repository.SearchRepository
}

func NewSearchService(repo repository.SearchRepository) *SearchService {
	return &SearchService{repo: repo}
}

func (s *SearchService) GlobalSearch(ctx context.Context, query string, role string, userID string) ([]models.SearchResult, error) {
	// 5 items per category
	limit := 5
	
	// 1. Search Users
	users, err := s.repo.SearchUsers(ctx, query, role, userID, limit)
	if err != nil {
		// Log error but stick to partial results? Or return error?
		// Handler usually handles it. We'll return partial + error or just error.
		return nil, err
	}
	if users == nil {
		users = []models.SearchResult{}
	}

	// 2. Search Documents
	docs, err := s.repo.SearchDocuments(ctx, query, role, userID, limit)
	if err != nil {
		return nil, err
	}
	if docs == nil {
		docs = []models.SearchResult{}
	}

	// 3. Search Messages
	// Not implemented in repo yet, logic skipped.

	// Combine
	results := append(users, docs...)
	return results, nil
}
