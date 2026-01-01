package services

import (
	"context"
	"errors"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type CurriculumService struct {
	repo repository.CurriculumRepository
}

func NewCurriculumService(repo repository.CurriculumRepository) *CurriculumService {
	return &CurriculumService{repo: repo}
}

// --- Programs ---

func (s *CurriculumService) CreateProgram(ctx context.Context, p *models.Program) error {
	if p.TenantID == "" {
		return errors.New("tenant_id is required")
	}
	if p.Title == "" && p.Name == "" {
		return errors.New("title is required")
	}
	// Default legacy name if missing
	if p.Name == "" {
		p.Name = p.Code // Fallback
	}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return s.repo.CreateProgram(ctx, p)
}

func (s *CurriculumService) GetProgram(ctx context.Context, id string) (*models.Program, error) {
	return s.repo.GetProgram(ctx, id)
}

func (s *CurriculumService) ListPrograms(ctx context.Context, tenantID string) ([]models.Program, error) {
	return s.repo.ListPrograms(ctx, tenantID)
}

func (s *CurriculumService) UpdateProgram(ctx context.Context, p *models.Program) error {
	p.UpdatedAt = time.Now()
	return s.repo.UpdateProgram(ctx, p)
}

func (s *CurriculumService) DeleteProgram(ctx context.Context, id string) error {
	return s.repo.DeleteProgram(ctx, id)
}

// --- Courses ---

func (s *CurriculumService) CreateCourse(ctx context.Context, c *models.Course) error {
	if c.TenantID == "" {
		return errors.New("tenant_id is required")
	}
	if c.Title == "" {
		return errors.New("title is required")
	}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return s.repo.CreateCourse(ctx, c)
}

func (s *CurriculumService) ListCourses(ctx context.Context, tenantID string, programID *string) ([]models.Course, error) {
	return s.repo.ListCourses(ctx, tenantID, programID)
}

func (s *CurriculumService) UpdateCourse(ctx context.Context, c *models.Course) error {
	c.UpdatedAt = time.Now()
	return s.repo.UpdateCourse(ctx, c)
}

func (s *CurriculumService) DeleteCourse(ctx context.Context, id string) error {
	return s.repo.DeleteCourse(ctx, id)
}

// --- Journey Map ---

func (s *CurriculumService) CreateJourneyMap(ctx context.Context, jm *models.JourneyMap) error {
	jm.CreatedAt = time.Now()
	return s.repo.CreateJourneyMap(ctx, jm)
}

func (s *CurriculumService) GetJourneyMap(ctx context.Context, programID string) (*models.JourneyMap, error) {
	return s.repo.GetJourneyMapByProgram(ctx, programID)
}

func (s *CurriculumService) CreateNodeDefinition(ctx context.Context, nd *models.JourneyNodeDefinition) error {
	nd.CreatedAt = time.Now()
	return s.repo.CreateNodeDefinition(ctx, nd)
}

func (s *CurriculumService) GetNodeDefinitions(ctx context.Context, journeyMapID string) ([]models.JourneyNodeDefinition, error) {
	return s.repo.GetNodeDefinitions(ctx, journeyMapID)
}

// --- Cohorts ---

func (s *CurriculumService) CreateCohort(ctx context.Context, c *models.Cohort) error {
	c.CreatedAt = time.Now()
	return s.repo.CreateCohort(ctx, c)
}

func (s *CurriculumService) ListCohorts(ctx context.Context, programID string) ([]models.Cohort, error) {
	return s.repo.ListCohorts(ctx, programID)
}
