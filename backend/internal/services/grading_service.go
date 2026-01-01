package services

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/jmoiron/sqlx/types"
)

type GradingService struct {
	repo repository.GradingRepository
}

func NewGradingService(repo repository.GradingRepository) *GradingService {
	return &GradingService{repo: repo}
}

// --- Schemas ---

func (s *GradingService) CreateSchema(ctx context.Context, schema *models.GradingSchema) error {
	if schema.Name == "" {
		return errors.New("name is required")
	}
	// Validate JSON scale structure? 
	// Assuming frontend sends valid JSON, but ideally we parse it here to verify.
	return s.repo.CreateSchema(ctx, schema)
}

func (s *GradingService) ListSchemas(ctx context.Context, tenantID string) ([]models.GradingSchema, error) {
	return s.repo.ListSchemas(ctx, tenantID)
}

func (s *GradingService) GetDefaultSchema(ctx context.Context, tenantID string) (*models.GradingSchema, error) {
	return s.repo.GetDefaultSchema(ctx, tenantID)
}

// --- Grading Logic ---

type GradeRule struct {
	MinPercent float64 `json:"min"`
	Grade      string  `json:"grade"`
}

// SubmitGrade calculates the letter grade and saves the entry.
func (s *GradingService) SubmitGrade(ctx context.Context, entry *models.GradebookEntry, tenantID string) error {
	if entry.CourseOfferingID == "" || entry.ActivityID == "" || entry.StudentID == "" {
		return errors.New("offering, activity, and student IDs are required")
	}
	if entry.MaxScore <= 0 {
		return errors.New("max_score must be positive")
	}

	// 1. Fetch Schema (simplify: just use default for tenant for now, or fetch course specific later)
	schema, err := s.repo.GetDefaultSchema(ctx, tenantID)
	if err != nil {
		// Fallback or error? For MVP, error if no schema.
		// Or perform a hardcoded fallback.
		return errors.New("no default grading schema found for tenant")
	}

	// 2. Calculate Percentage
	percent := (entry.Score / entry.MaxScore) * 100

	// 3. Determine Letter Grade
	letterGrade, err := calculateLetterGrade(percent, schema.Scale)
	if err != nil {
		return err
	}
	entry.Grade = letterGrade
	
	entry.GradedAt = time.Now()

	return s.repo.CreateEntry(ctx, entry)
}


func calculateLetterGrade(percent float64, scaleJSON types.JSONText) (string, error) {
	var rules []GradeRule
	if err := json.Unmarshal(scaleJSON, &rules); err != nil {
		return "", err
	}

	// Sort rules descending by MinPercent to find highest match
	// e.g. [{90, A}, {80, B}, {70, C}]
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].MinPercent > rules[j].MinPercent
	})

	for _, rule := range rules {
		if percent >= rule.MinPercent {
			return rule.Grade, nil
		}
	}
	
	// If below all thresholds, return "F" or empty?
	// Implicitly usually 0 is F. If rules cover 0, we are good.
	// If not, return "F" as fallback?
	return "F", nil
}

func (s *GradingService) ListStudentGrades(ctx context.Context, studentID string) ([]models.GradebookEntry, error) {
	return s.repo.ListStudentEntries(ctx, studentID)
}
