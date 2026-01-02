package services

import (
	"context"
	"errors"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type RubricService struct {
	repo repository.RubricRepository
}

func NewRubricService(repo repository.RubricRepository) *RubricService {
	return &RubricService{repo: repo}
}

func (s *RubricService) CreateRubric(ctx context.Context, r *models.Rubric) (*models.Rubric, error) {
	if len(r.Criteria) == 0 {
		return nil, errors.New("rubric must have criteria")
	}
	if err := s.repo.CreateRubric(ctx, r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *RubricService) GetRubric(ctx context.Context, id string) (*models.Rubric, error) {
	return s.repo.GetRubric(ctx, id)
}

func (s *RubricService) ListRubrics(ctx context.Context, courseID string) ([]models.Rubric, error) {
	return s.repo.ListRubrics(ctx, courseID)
}

// GradeInput is a simplified struct for handlers to bind
type GradeInput struct {
	RubricID     string `json:"rubric_id"`
	SubmissionID string `json:"submission_id"`
	GraderID     string `json:"-"`
	Comments     string `json:"comments"`
	Selections   []struct {
		CriterionID string `json:"criterion_id"`
		LevelID     string `json:"level_id"`
		// Optional manual override if LevelID is nil? For now enforce selection.
	} `json:"selections"`
}

func (s *RubricService) SubmitGrade(ctx context.Context, input GradeInput) (*models.RubricGrade, error) {
	// 1. Fetch Rubric to validate scores
	rubric, err := s.repo.GetRubric(ctx, input.RubricID)
	if err != nil {
		return nil, err
	}

	// Index criteria/levels
	critMap := make(map[string]models.RubricCriterion)
	levelMap := make(map[string]models.RubricLevel)
	
	for _, c := range rubric.Criteria {
		critMap[c.ID] = c
		for _, l := range c.Levels {
			levelMap[l.ID] = l
		}
	}

	// 2. Build Grade
	grade := &models.RubricGrade{
		SubmissionID: input.SubmissionID,
		RubricID:     input.RubricID,
		GraderID:     &input.GraderID,
		Comments:     &input.Comments,
		Items:        []models.RubricGradeItem{},
	}

	totalScore := 0.0

	for _, sel := range input.Selections {
		// Validate
		crit, ok := critMap[sel.CriterionID]
		if !ok {
			continue // skip invalid criterion ID
		}
		
		level, ok := levelMap[sel.LevelID]
		if !ok || level.CriterionID != sel.CriterionID {
			return nil, errors.New("invalid level selection for criterion " + crit.Title)
		}

		points := level.Points * crit.Weight // Add weight multiplier logic if desired. (Currently schema has weight).
		// Wait, if points are absolute (e.g. 5 pts), does weight multiply it? 
		// Usually: Criterion has weight 20%. Levels are 0-4. Score = (Selected Level / Max Level) * Weight?
		// Or: Levels have raw points. Total = Sum(Points).
		// Let's assume Points are final raw points. Weight is for visual or calculation hints. 
		// Actually schema `weight FLOAT DEFAULT 1.0`. Let's multiply.
		
		awarded := points // * crit.Weight? Let's assume points are defined raw in DB level.
		
		totalScore += awarded

		grade.Items = append(grade.Items, models.RubricGradeItem{
			CriterionID: sel.CriterionID,
			LevelID:     &sel.LevelID,
			PointsAwarded: awarded,
		})
	}
	
	grade.TotalScore = totalScore

	// 3. Save
	if err := s.repo.SubmitGrade(ctx, grade); err != nil {
		return nil, err
	}
	
	return grade, nil
}

func (s *RubricService) GetGrade(ctx context.Context, submissionID string) (*models.RubricGrade, error) {
	return s.repo.GetGrade(ctx, submissionID)
}
