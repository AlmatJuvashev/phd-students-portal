package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

// AuditService handles audit-related business logic
type AuditService struct {
	repo         repository.AuditRepository
	curriculumRepo repository.CurriculumRepository
}

func NewAuditService(repo repository.AuditRepository, curriculumRepo repository.CurriculumRepository) *AuditService {
	return &AuditService{repo: repo, curriculumRepo: curriculumRepo}
}

// --- Learning Outcomes ---

func (s *AuditService) ListLearningOutcomes(ctx context.Context, tenantID string, programID, courseID *string) ([]models.LearningOutcome, error) {
	return s.repo.ListLearningOutcomes(ctx, tenantID, programID, courseID)
}

func (s *AuditService) GetLearningOutcome(ctx context.Context, id string) (*models.LearningOutcome, error) {
	return s.repo.GetLearningOutcome(ctx, id)
}

func (s *AuditService) CreateLearningOutcome(ctx context.Context, outcome *models.LearningOutcome, changedByID string) error {
	if err := s.repo.CreateLearningOutcome(ctx, outcome); err != nil {
		return err
	}
	
	// Log the change
	newValueJSON, _ := json.Marshal(outcome)
	return s.repo.LogCurriculumChange(ctx, &models.CurriculumChangeLog{
		TenantID:   outcome.TenantID,
		EntityType: "outcome",
		EntityID:   outcome.ID,
		Action:     "created",
		NewValue:   string(newValueJSON),
		ChangedBy:  changedByID,
	})
}

func (s *AuditService) UpdateLearningOutcome(ctx context.Context, outcome *models.LearningOutcome, changedByID string) error {
	// Get old value for audit
	old, _ := s.repo.GetLearningOutcome(ctx, outcome.ID)
	
	if err := s.repo.UpdateLearningOutcome(ctx, outcome); err != nil {
		return err
	}
	
	// Log the change
	oldValueJSON, _ := json.Marshal(old)
	newValueJSON, _ := json.Marshal(outcome)
	return s.repo.LogCurriculumChange(ctx, &models.CurriculumChangeLog{
		TenantID:   outcome.TenantID,
		EntityType: "outcome",
		EntityID:   outcome.ID,
		Action:     "updated",
		OldValue:   string(oldValueJSON),
		NewValue:   string(newValueJSON),
		ChangedBy:  changedByID,
	})
}

func (s *AuditService) DeleteLearningOutcome(ctx context.Context, tenantID, id, changedByID string) error {
	// Get old value for audit
	old, _ := s.repo.GetLearningOutcome(ctx, id)
	
	if err := s.repo.DeleteLearningOutcome(ctx, id); err != nil {
		return err
	}
	
	// Log the change
	oldValueJSON, _ := json.Marshal(old)
	return s.repo.LogCurriculumChange(ctx, &models.CurriculumChangeLog{
		TenantID:   tenantID,
		EntityType: "outcome",
		EntityID:   id,
		Action:     "deleted",
		OldValue:   string(oldValueJSON),
		ChangedBy:  changedByID,
	})
}

// --- Outcome Assessments ---

func (s *AuditService) LinkOutcomeToAssessment(ctx context.Context, outcomeID, nodeDefID string, weight float64) error {
	return s.repo.LinkOutcomeToAssessment(ctx, outcomeID, nodeDefID, weight)
}

func (s *AuditService) GetOutcomeAssessments(ctx context.Context, outcomeID string) ([]models.OutcomeAssessment, error) {
	return s.repo.GetOutcomeAssessments(ctx, outcomeID)
}

// --- Curriculum Change Log ---

func (s *AuditService) ListCurriculumChanges(ctx context.Context, filter models.AuditReportFilter) ([]models.CurriculumChangeLog, error) {
	return s.repo.ListCurriculumChanges(ctx, filter)
}

// --- Report Helpers ---

type ProgramSummaryReport struct {
	Program       models.Program              `json:"program"`
	Courses       []models.Course             `json:"courses"`
	Outcomes      []models.LearningOutcome    `json:"outcomes"`
	TotalCredits  int                         `json:"total_credits"`
	TotalCourses  int                         `json:"total_courses"`
	TotalOutcomes int                         `json:"total_outcomes"`
	GeneratedAt   time.Time                   `json:"generated_at"`
}

func (s *AuditService) GenerateProgramSummary(ctx context.Context, tenantID, programID string) (*ProgramSummaryReport, error) {
	program, err := s.curriculumRepo.GetProgram(ctx, programID)
	if err != nil {
		return nil, err
	}
	
	courses, err := s.curriculumRepo.ListCourses(ctx, tenantID, &programID)
	if err != nil {
		return nil, err
	}
	
	outcomes, err := s.repo.ListLearningOutcomes(ctx, tenantID, &programID, nil)
	if err != nil {
		return nil, err
	}
	
	totalCredits := 0
	for _, c := range courses {
		totalCredits += c.Credits
	}
	
	return &ProgramSummaryReport{
		Program:       *program,
		Courses:       courses,
		Outcomes:      outcomes,
		TotalCredits:  totalCredits,
		TotalCourses:  len(courses),
		TotalOutcomes: len(outcomes),
		GeneratedAt:   time.Now(),
	}, nil
}
