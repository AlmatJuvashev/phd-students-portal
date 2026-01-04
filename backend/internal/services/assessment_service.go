package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/jmoiron/sqlx/types"
)

var ErrForbidden = errors.New("forbidden")

type AttemptAlreadyInProgressError struct {
	Attempt *models.AssessmentAttempt
}

func (e *AttemptAlreadyInProgressError) Error() string {
	return "attempt already in progress"
}

type MaxAttemptsReachedError struct {
	MaxAttempts int
}

func (e *MaxAttemptsReachedError) Error() string {
	return fmt.Sprintf("max attempts reached (%d)", e.MaxAttempts)
}

type CooldownActiveError struct {
	RetryAfter time.Duration
}

func (e *CooldownActiveError) Error() string {
	return "cooldown active"
}

type AttemptAutoSubmittedError struct {
	Attempt *models.AssessmentAttempt
	Reason  string
}

func (e *AttemptAutoSubmittedError) Error() string {
	if e.Reason == "" {
		return "attempt auto-submitted"
	}
	return fmt.Sprintf("attempt auto-submitted: %s", e.Reason)
}

type AssessmentService struct {
	repo repository.AssessmentRepository
}

func NewAssessmentService(repo repository.AssessmentRepository) *AssessmentService {
	return &AssessmentService{repo: repo}
}

func (s *AssessmentService) CreateAssessment(ctx context.Context, a models.Assessment) (*models.Assessment, error) {
	if a.TenantID == "" {
		return nil, errors.New("tenant_id is required")
	}
	if a.CourseOfferingID == "" {
		return nil, errors.New("course_offering_id is required")
	}
	if a.Title == "" {
		return nil, errors.New("title is required")
	}
	if len(a.SecuritySettings) == 0 {
		a.SecuritySettings = types.JSONText([]byte(`{}`))
	}
	return s.repo.CreateAssessment(ctx, a)
}

func (s *AssessmentService) ListAssessments(ctx context.Context, tenantID string, courseOfferingID string) ([]models.Assessment, error) {
	return s.repo.ListAssessments(ctx, tenantID, courseOfferingID)
}

func (s *AssessmentService) GetAssessmentForTaking(ctx context.Context, tenantID, assessmentID string) (*models.Assessment, []models.Question, error) {
	assessment, err := s.repo.GetAssessment(ctx, assessmentID)
	if err != nil {
		return nil, nil, err
	}
	if assessment.TenantID != tenantID {
		return nil, nil, ErrForbidden
	}

	questions, err := s.repo.GetAssessmentQuestions(ctx, assessmentID)
	if err != nil {
		return nil, nil, err
	}

	// Hide correctness from student while taking.
	for qi := range questions {
		for oi := range questions[qi].Options {
			questions[qi].Options[oi].IsCorrect = false
		}
	}

	return assessment, questions, nil
}

func (s *AssessmentService) UpdateAssessment(ctx context.Context, tenantID string, a models.Assessment) (*models.Assessment, error) {
	current, err := s.repo.GetAssessment(ctx, a.ID)
	if err != nil {
		return nil, err
	}
	if current.TenantID != tenantID {
		return nil, ErrForbidden
	}

	// Preserve immutable fields.
	a.TenantID = current.TenantID
	a.CreatedBy = current.CreatedBy
	if len(a.SecuritySettings) == 0 {
		a.SecuritySettings = current.SecuritySettings
	}

	if err := s.repo.UpdateAssessment(ctx, a); err != nil {
		return nil, err
	}
	return s.repo.GetAssessment(ctx, a.ID)
}

func (s *AssessmentService) DeleteAssessment(ctx context.Context, tenantID, assessmentID string) error {
	current, err := s.repo.GetAssessment(ctx, assessmentID)
	if err != nil {
		return err
	}
	if current.TenantID != tenantID {
		return ErrForbidden
	}
	return s.repo.DeleteAssessment(ctx, assessmentID)
}

// CreateAttempt initializes a new assessment attempt for a student
func (s *AssessmentService) CreateAttempt(ctx context.Context, tenantID, assessmentID, studentID string) (*models.AssessmentAttempt, error) {
	// 1. Fetch Assessment
	assessment, err := s.repo.GetAssessment(ctx, assessmentID)
	if err != nil {
		return nil, err
	}
	if assessment.TenantID != tenantID {
		return nil, ErrForbidden
	}

	// 2. Check Availability
	now := time.Now()
	if assessment.AvailableFrom != nil && now.Before(*assessment.AvailableFrom) {
		return nil, errors.New("assessment is not yet available")
	}
	if assessment.AvailableUntil != nil && now.After(*assessment.AvailableUntil) {
		return nil, errors.New("assessment is closed")
	}

	// 3. Retake / in-progress policy (configured via security_settings for now)
	var settings models.SecuritySettings
	if len(assessment.SecuritySettings) > 0 {
		_ = json.Unmarshal(assessment.SecuritySettings, &settings)
	}

	attempts, err := s.repo.ListAttemptsByAssessmentAndStudent(ctx, assessmentID, studentID)
	if err != nil {
		return nil, err
	}

	// If there's an in-progress attempt, either reuse it or auto-submit if time is up.
	for i := range attempts {
		if attempts[i].Status != models.AttemptStatusInProgress {
			continue
		}

		if s.isAttemptExpired(assessment, &attempts[i], now) {
			_, _ = s.completeAttempt(ctx, tenantID, studentID, attempts[i].ID)
			break
		}

		return nil, &AttemptAlreadyInProgressError{Attempt: &attempts[i]}
	}

	if settings.MaxAttempts > 0 || settings.CooldownMinutes > 0 {
		// Enforce max attempts based on completed attempts (submitted/graded).
		completed := 0
		var lastFinished *time.Time
		for _, a := range attempts {
			if a.Status == models.AttemptStatusInProgress {
				continue
			}
			completed++
			if a.FinishedAt != nil && (lastFinished == nil || a.FinishedAt.After(*lastFinished)) {
				t := *a.FinishedAt
				lastFinished = &t
			}
		}

		if settings.MaxAttempts > 0 && completed >= settings.MaxAttempts {
			return nil, &MaxAttemptsReachedError{MaxAttempts: settings.MaxAttempts}
		}

		if settings.CooldownMinutes > 0 && lastFinished != nil {
			retryAt := lastFinished.Add(time.Duration(settings.CooldownMinutes) * time.Minute)
			if now.Before(retryAt) {
				return nil, &CooldownActiveError{RetryAfter: time.Until(retryAt)}
			}
		}
	}

	// 3. Create Attempt
	return s.repo.CreateAttempt(ctx, models.AssessmentAttempt{
		AssessmentID: assessmentID,
		StudentID:    studentID,
	})
}

// SubmitResponse saves a student's answer to a specific question
func (s *AssessmentService) SubmitResponse(ctx context.Context, tenantID, attemptID, studentID, questionID string, optionID *string, text *string) error {
	// Validate attempt exists and is in progress
	attempt, err := s.repo.GetAttempt(ctx, attemptID)
	if err != nil {
		return err
	}
	if attempt.StudentID != studentID {
		return ErrForbidden
	}
	if attempt.Status != models.AttemptStatusInProgress {
		return errors.New("attempt is not in progress")
	}

	assessment, err := s.repo.GetAssessment(ctx, attempt.AssessmentID)
	if err != nil {
		return err
	}
	if assessment.TenantID != tenantID {
		return ErrForbidden
	}

	// Enforce time limit by auto-submitting on access.
	if s.isAttemptExpired(assessment, attempt, time.Now()) {
		completed, _ := s.completeAttempt(ctx, tenantID, studentID, attemptID)
		return &AttemptAutoSubmittedError{Attempt: completed, Reason: "TIME_LIMIT_EXCEEDED"}
	}

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
func (s *AssessmentService) CompleteAttempt(ctx context.Context, tenantID, attemptID, studentID string) (*models.AssessmentAttempt, error) {
	attempt, err := s.repo.GetAttempt(ctx, attemptID)
	if err != nil {
		return nil, err
	}
	if attempt.StudentID != studentID {
		return nil, ErrForbidden
	}
	assessment, err := s.repo.GetAssessment(ctx, attempt.AssessmentID)
	if err != nil {
		return nil, err
	}
	if assessment.TenantID != tenantID {
		return nil, ErrForbidden
	}

	return s.completeAttempt(ctx, tenantID, studentID, attemptID)
}

func (s *AssessmentService) completeAttempt(ctx context.Context, tenantID, studentID, attemptID string) (*models.AssessmentAttempt, error) {
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
	totalPossible := 0.0
	
	// Iterate through all questions to grade
	for _, q := range questions {
		totalPossible += q.PointsDefault
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

	percentageScore := 0.0
	if totalPossible > 0 {
		percentageScore = (totalScore / totalPossible) * 100
	}

	// 4. Finalize Attempt
	err = s.repo.CompleteAttempt(ctx, attemptID, percentageScore)
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
func (s *AssessmentService) ReportProctoringEvent(ctx context.Context, tenantID, attemptID, studentID string, event models.ProctoringEventType, meta map[string]interface{}) error {
	attempt, err := s.repo.GetAttempt(ctx, attemptID)
	if err != nil {
		return err
	}
	if attempt.StudentID != studentID {
		return ErrForbidden
	}

	assessment, err := s.repo.GetAssessment(ctx, attempt.AssessmentID)
	if err != nil {
		return err
	}
	if assessment.TenantID != tenantID {
		return ErrForbidden
	}

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
	_, err := s.completeAttempt(ctx, "", "", attemptID)
	return err
}

func (s *AssessmentService) GetAttemptDetails(ctx context.Context, tenantID, attemptID, studentID string) (*models.AssessmentAttempt, *models.Assessment, []models.Question, []models.ItemResponse, error) {
	attempt, err := s.repo.GetAttempt(ctx, attemptID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if attempt.StudentID != studentID {
		return nil, nil, nil, nil, ErrForbidden
	}

	assessment, err := s.repo.GetAssessment(ctx, attempt.AssessmentID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if assessment.TenantID != tenantID {
		return nil, nil, nil, nil, ErrForbidden
	}

	// Auto-submit on timeout.
	if attempt.Status == models.AttemptStatusInProgress && s.isAttemptExpired(assessment, attempt, time.Now()) {
		updated, _ := s.completeAttempt(ctx, tenantID, studentID, attemptID)
		if updated != nil {
			attempt = updated
		}
	}

	questions, err := s.repo.GetAssessmentQuestions(ctx, attempt.AssessmentID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	responses, err := s.repo.ListResponses(ctx, attemptID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Hide correctness while in progress; reveal on submitted/graded.
	if attempt.Status == models.AttemptStatusInProgress {
		for qi := range questions {
			for oi := range questions[qi].Options {
				questions[qi].Options[oi].IsCorrect = false
			}
		}
	}

	return attempt, assessment, questions, responses, nil
}

func (s *AssessmentService) ListMyAttempts(ctx context.Context, tenantID, assessmentID, studentID string) ([]models.AssessmentAttempt, error) {
	assessment, err := s.repo.GetAssessment(ctx, assessmentID)
	if err != nil {
		return nil, err
	}
	if assessment.TenantID != tenantID {
		return nil, ErrForbidden
	}
	return s.repo.ListAttemptsByAssessmentAndStudent(ctx, assessmentID, studentID)
}

func (s *AssessmentService) isAttemptExpired(assessment *models.Assessment, attempt *models.AssessmentAttempt, now time.Time) bool {
	if assessment == nil || attempt == nil {
		return false
	}
	if assessment.TimeLimitMinutes == nil || *assessment.TimeLimitMinutes <= 0 {
		return false
	}
	expiresAt := attempt.StartedAt.Add(time.Duration(*assessment.TimeLimitMinutes) * time.Minute)
	return now.After(expiresAt)
}
