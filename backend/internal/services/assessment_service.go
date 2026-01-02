package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/jmoiron/sqlx/types"
)

type AssessmentService struct {
	repo repository.AssessmentRepository
}

func NewAssessmentService(repo repository.AssessmentRepository) *AssessmentService {
	return &AssessmentService{repo: repo}
}

// CreateAttempt initializes a new assessment attempt for a student
func (s *AssessmentService) CreateAttempt(ctx context.Context, assessmentID, studentID string) (*models.AssessmentAttempt, error) {
	// 1. Fetch Assessment
	assessment, err := s.repo.GetAssessment(ctx, assessmentID)
	if err != nil {
		return nil, err
	}

	// 2. Check Availability
	now := time.Now()
	if assessment.AvailableFrom != nil && now.Before(*assessment.AvailableFrom) {
		return nil, errors.New("assessment is not yet available")
	}
	if assessment.AvailableUntil != nil && now.After(*assessment.AvailableUntil) {
		return nil, errors.New("assessment is closed")
	}

	// 3. Create Attempt
	// TODO: Check for existing attempts if retake policy exists
	return s.repo.CreateAttempt(ctx, models.AssessmentAttempt{
		AssessmentID: assessmentID,
		StudentID:    studentID,
	})
}

// SubmitResponse saves a student's answer to a specific question
func (s *AssessmentService) SubmitResponse(ctx context.Context, attemptID, questionID string, optionID *string, text *string) error {
	// Validate attempt exists and is in progress
	attempt, err := s.repo.GetAttempt(ctx, attemptID)
	if err != nil {
		return err
	}
	if attempt.Status != models.AttemptStatusInProgress {
		return errors.New("attempt is not in progress")
	}
	// TODO: Validate Time Limit

	response := models.ItemResponse{
		AttemptID:        attemptID,
		QuestionID:       questionID,
		SelectedOptionID: optionID,
		TextResponse:     text,
		IsCorrect:        false, // Will be calculated on completion or here
		Score:            0,
	}

	return s.repo.SaveItemResponse(ctx, response)
}

// CompleteAttempt finishes the exam and runs auto-grading
func (s *AssessmentService) CompleteAttempt(ctx context.Context, attemptID string) (*models.AssessmentAttempt, error) {
	attempt, err := s.repo.GetAttempt(ctx, attemptID)
	if err != nil {
		return nil, err
	}
	if attempt.Status != models.AttemptStatusInProgress {
		return attempt, nil // Already completed
	}

	// 1. Fetch Questions and Student Responses
	questions, err := s.repo.GetAssessmentQuestions(ctx, attempt.AssessmentID)
	if err != nil {
		return nil, err
	}
	responses, err := s.repo.ListResponses(ctx, attemptID)
	if err != nil {
		return nil, err
	}

	// 2. Map Responses for O(1) loop up
	responseMap := make(map[string]models.ItemResponse)
	for _, r := range responses {
		responseMap[r.QuestionID] = r
	}

	// 3. Calculate Score
	totalScore := 0.0
	
	// Iterate through all questions to grade
	for _, q := range questions {
		if resp, exists := responseMap[q.ID]; exists {
			score, isCorrect := s.calculateScore(q, resp)
			
			// Update the individual response with correctness and score
			resp.Score = score
			resp.IsCorrect = isCorrect
			// Persist grading result for this item
			// Note: This calls DB in loop. For high perf, use batch update. For MVP, loop is acceptable.
			_ = s.repo.SaveItemResponse(ctx, resp) 
			
			totalScore += score
		}
	}

	// 4. Finalize Attempt
	err = s.repo.CompleteAttempt(ctx, attemptID, totalScore)
	if err != nil {
		return nil, err
	}

	return s.repo.GetAttempt(ctx, attemptID)
}

// AutoGrade calculates score for a specific response against the question
// This should be called during CompleteAttempt
func (s *AssessmentService) calculateScore(question models.Question, response models.ItemResponse) (float64, bool) {
	if question.Type == models.QuestionTypeMCQ || question.Type == models.QuestionTypeTrueFalse {
		for _, opt := range question.Options {
			if opt.IsCorrect && response.SelectedOptionID != nil && opt.ID == *response.SelectedOptionID {
				return question.PointsDefault, true
			}
		}
	}
	// Add other types logic
	return 0, false
}

// ReportProctoringEvent logs an incident and checks against security policy
func (s *AssessmentService) ReportProctoringEvent(ctx context.Context, attemptID string, event models.ProctoringEventType, meta map[string]interface{}) error {
	// 1. Log Event
	metadataJSON, _ := json.Marshal(meta)
	log := models.ProctoringLog{
		AttemptID: attemptID,
		EventType: event,
		Metadata:  types.JSONText(metadataJSON),
	}
	if err := s.repo.LogProctoringEvent(ctx, log); err != nil {
		return err
	}

	// 2. Check Policy (Optional: If strict mode is on)
	attempt, err := s.repo.GetAttempt(ctx, attemptID)
	if err != nil {
		return err
	}
	assessment, err := s.repo.GetAssessment(ctx, attempt.AssessmentID)
	if err != nil {
		return err
	}

	var settings models.SecuritySettings
	if len(assessment.SecuritySettings) > 0 {
		_ = json.Unmarshal(assessment.SecuritySettings, &settings)
	}

	// 3. Enforce Limit
	if settings.MaxViolations > 0 && settings.AutoSubmitOnLimit {
		count, err := s.repo.CountProctoringEvents(ctx, attemptID)
		if err != nil {
			return err
		}
		if count >= settings.MaxViolations {
			// Auto Terminate
			return s.terminateAttempt(ctx, attemptID, "Security Violation Limit Reached")
		}
	}

	return nil
}

func (s *AssessmentService) terminateAttempt(ctx context.Context, attemptID, reason string) error {
	// Logic to force complete with 0 score or current score
	// For now, we reuse CompleteAttempt but typically we might mark it as "Flagged"
	// Reusing CompleteAttempt for MVP
	_, err := s.CompleteAttempt(ctx, attemptID)
	return err
}
