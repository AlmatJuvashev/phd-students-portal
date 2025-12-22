package services

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type DictionaryService struct {
	repo repository.DictionaryRepository
}

func NewDictionaryService(repo repository.DictionaryRepository) *DictionaryService {
	return &DictionaryService{repo: repo}
}

// --- Programs ---

func (s *DictionaryService) ListPrograms(ctx context.Context, tenantID string, activeOnly bool) ([]models.Program, error) {
	return s.repo.ListPrograms(ctx, tenantID, activeOnly)
}

func (s *DictionaryService) CreateProgram(ctx context.Context, tenantID, name, code string) (string, error) {
	return s.repo.CreateProgram(ctx, tenantID, name, code)
}

func (s *DictionaryService) UpdateProgram(ctx context.Context, tenantID, id string, name, code string, isActive *bool) error {
	return s.repo.UpdateProgram(ctx, tenantID, id, name, code, isActive)
}

func (s *DictionaryService) DeleteProgram(ctx context.Context, tenantID, id string) error {
	return s.repo.DeleteProgram(ctx, tenantID, id)
}

// --- Specialties ---

func (s *DictionaryService) ListSpecialties(ctx context.Context, tenantID string, activeOnly bool, programID string) ([]models.Specialty, error) {
	return s.repo.ListSpecialties(ctx, tenantID, activeOnly, programID)
}

func (s *DictionaryService) CreateSpecialty(ctx context.Context, tenantID, name, code string, programIDs []string) (string, error) {
	return s.repo.CreateSpecialty(ctx, tenantID, name, code, programIDs)
}

func (s *DictionaryService) UpdateSpecialty(ctx context.Context, tenantID, id, name, code string, isActive *bool, programIDs []string) error {
	return s.repo.UpdateSpecialty(ctx, tenantID, id, name, code, isActive, programIDs)
}

func (s *DictionaryService) DeleteSpecialty(ctx context.Context, tenantID, id string) error {
	return s.repo.DeleteSpecialty(ctx, tenantID, id)
}

// --- Cohorts ---

func (s *DictionaryService) ListCohorts(ctx context.Context, tenantID string, activeOnly bool) ([]models.Cohort, error) {
	return s.repo.ListCohorts(ctx, tenantID, activeOnly)
}

func (s *DictionaryService) CreateCohort(ctx context.Context, tenantID, name, startDate, endDate string) (string, error) {
	return s.repo.CreateCohort(ctx, tenantID, name, startDate, endDate)
}

func (s *DictionaryService) UpdateCohort(ctx context.Context, tenantID, id, name, startDate, endDate string, isActive *bool) error {
	return s.repo.UpdateCohort(ctx, tenantID, id, name, startDate, endDate, isActive)
}

func (s *DictionaryService) DeleteCohort(ctx context.Context, tenantID, id string) error {
	return s.repo.DeleteCohort(ctx, tenantID, id)
}

// --- Departments ---

func (s *DictionaryService) ListDepartments(ctx context.Context, tenantID string, activeOnly bool) ([]models.Department, error) {
	return s.repo.ListDepartments(ctx, tenantID, activeOnly)
}

func (s *DictionaryService) CreateDepartment(ctx context.Context, tenantID, name, code string) (string, error) {
	return s.repo.CreateDepartment(ctx, tenantID, name, code)
}

func (s *DictionaryService) UpdateDepartment(ctx context.Context, tenantID, id, name, code string, isActive *bool) error {
	return s.repo.UpdateDepartment(ctx, tenantID, id, name, code, isActive)
}

func (s *DictionaryService) DeleteDepartment(ctx context.Context, tenantID, id string) error {
	return s.repo.DeleteDepartment(ctx, tenantID, id)
}
