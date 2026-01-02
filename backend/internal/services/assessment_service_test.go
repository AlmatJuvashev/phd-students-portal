package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAssessmentRepo
type MockAssessmentRepo struct {
	mock.Mock
}

func (m *MockAssessmentRepo) CreateQuestionBank(ctx context.Context, bank models.QuestionBank) (*models.QuestionBank, error) {
	args := m.Called(ctx, bank)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.QuestionBank), args.Error(1)
}
func (m *MockAssessmentRepo) GetQuestionBank(ctx context.Context, id string) (*models.QuestionBank, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.QuestionBank), args.Error(1)
}
func (m *MockAssessmentRepo) ListQuestionBanks(ctx context.Context, tenantID string) ([]models.QuestionBank, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).([]models.QuestionBank), args.Error(1)
}
func (m *MockAssessmentRepo) CreateQuestion(ctx context.Context, q models.Question) (*models.Question, error) {
	args := m.Called(ctx, q)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Question), args.Error(1)
}
func (m *MockAssessmentRepo) GetQuestion(ctx context.Context, id string) (*models.Question, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Question), args.Error(1)
}
func (m *MockAssessmentRepo) ListQuestionsByBank(ctx context.Context, bankID string) ([]models.Question, error) {
	args := m.Called(ctx, bankID)
	return args.Get(0).([]models.Question), args.Error(1)
}
func (m *MockAssessmentRepo) CreateAssessment(ctx context.Context, a models.Assessment) (*models.Assessment, error) {
	args := m.Called(ctx, a)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Assessment), args.Error(1)
}
func (m *MockAssessmentRepo) GetAssessment(ctx context.Context, id string) (*models.Assessment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Assessment), args.Error(1)
}
func (m *MockAssessmentRepo) CreateAttempt(ctx context.Context, attempt models.AssessmentAttempt) (*models.AssessmentAttempt, error) {
	args := m.Called(ctx, attempt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AssessmentAttempt), args.Error(1)
}
func (m *MockAssessmentRepo) SaveItemResponse(ctx context.Context, response models.ItemResponse) error {
	args := m.Called(ctx, response)
	return args.Error(0)
}
func (m *MockAssessmentRepo) CompleteAttempt(ctx context.Context, attemptID string, score float64) error {
	args := m.Called(ctx, attemptID, score)
	return args.Error(0)
}
func (m *MockAssessmentRepo) GetAttempt(ctx context.Context, id string) (*models.AssessmentAttempt, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AssessmentAttempt), args.Error(1)
}
func (m *MockAssessmentRepo) ListResponses(ctx context.Context, attemptID string) ([]models.ItemResponse, error) {
	args := m.Called(ctx, attemptID)
	return args.Get(0).([]models.ItemResponse), args.Error(1)
}
func (m *MockAssessmentRepo) GetAssessmentQuestions(ctx context.Context, assessmentID string) ([]models.Question, error) {
	args := m.Called(ctx, assessmentID)
	return args.Get(0).([]models.Question), args.Error(1)
}
func (m *MockAssessmentRepo) UpdateQuestion(ctx context.Context, q models.Question) error {
	args := m.Called(ctx, q)
	return args.Error(0)
}
func (m *MockAssessmentRepo) DeleteQuestion(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockAssessmentRepo) LogProctoringEvent(ctx context.Context, log models.ProctoringLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}
func (m *MockAssessmentRepo) CountProctoringEvents(ctx context.Context, attemptID string) (int, error) {
	args := m.Called(ctx, attemptID)
	return args.Get(0).(int), args.Error(1)
}

func TestAssessmentService_ReportProctoringEvent_AutoSubmit(t *testing.T) {
	mockRepo := new(MockAssessmentRepo)
	svc := services.NewAssessmentService(mockRepo)
	ctx := context.Background()

	attemptID := "att_spy"
	
	// Settings: Max 3 violations
	settingsJSON := []byte(`{"max_violations": 3, "auto_submit_on_limit": true}`)
	
	// Mock Fetch Attempt & Assessment
	mockRepo.On("GetAttempt", ctx, attemptID).Return(&models.AssessmentAttempt{
		ID: attemptID, AssessmentID: "exam_spy", Status: models.AttemptStatusInProgress,
	}, nil)
	
	mockRepo.On("GetAssessment", ctx, "exam_spy").Return(&models.Assessment{
		ID: "exam_spy", SecuritySettings: types.JSONText(settingsJSON),
	}, nil)

	// Mock Log Event
	mockRepo.On("LogProctoringEvent", ctx, mock.AnythingOfType("models.ProctoringLog")).Return(nil)

	// Mock Count -> Return 3 (Hitting limit)
	mockRepo.On("CountProctoringEvents", ctx, attemptID).Return(3, nil)

	// Expect Termination (CompleteAttempt with score calculation)
	// We need to support GetAssessmentQuestions etc for CompleteAttempt to work, 
	// OR we assume terminateAttempt calls CompleteAttempt which calls those methods.
	// For this test, let's mock the dependencies of CompleteAttempt too.
	// 1. GetAssessmentQuestions
	mockRepo.On("GetAssessmentQuestions", ctx, "exam_spy").Return([]models.Question{}, nil)
	// 2. ListResponses
	mockRepo.On("ListResponses", ctx, attemptID).Return([]models.ItemResponse{}, nil)
	// 3. CompleteAttempt
	mockRepo.On("CompleteAttempt", ctx, attemptID, 0.0).Return(nil)
	// 4. GetAttempt (fetch result)
	mockRepo.On("GetAttempt", ctx, attemptID).Return(&models.AssessmentAttempt{ID: attemptID, Status: models.AttemptStatusSubmitted}, nil)

	// Action: Report 3rd violation
	err := svc.ReportProctoringEvent(ctx, attemptID, models.ProctoringEventTabSwitch, nil)
	
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAssessmentService_CreateAttempt(t *testing.T) {
	mockRepo := new(MockAssessmentRepo)
	svc := services.NewAssessmentService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		assessment := &models.Assessment{ID: "exam1", AvailableFrom: nil, AvailableUntil: nil}
		mockRepo.On("GetAssessment", ctx, "exam1").Return(assessment, nil)
		mockRepo.On("CreateAttempt", ctx, mock.AnythingOfType("models.AssessmentAttempt")).Return(&models.AssessmentAttempt{ID: "attempt1"}, nil)

		attempt, err := svc.CreateAttempt(ctx, "exam1", "student1")
		assert.NoError(t, err)
		assert.NotNil(t, attempt)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotStarted", func(t *testing.T) {
		future := time.Now().Add(1 * time.Hour)
		assessment := &models.Assessment{ID: "exam2", AvailableFrom: &future}
		mockRepo.On("GetAssessment", ctx, "exam2").Return(assessment, nil)

		_, err := svc.CreateAttempt(ctx, "exam2", "student1")
		assert.Error(t, err)
		assert.Equal(t, "assessment is not yet available", err.Error())
	})
}

func TestAssessmentService_CompleteAttempt_AutoGrading(t *testing.T) {
	mockRepo := new(MockAssessmentRepo)
	svc := services.NewAssessmentService(mockRepo)
	ctx := context.Background()

	// Setup Assessment Question
	q1 := models.Question{
		ID: "q1", Type: models.QuestionTypeMCQ, PointsDefault: 5.0,
		Options: []models.QuestionOption{
			{ID: "optA", IsCorrect: true},
			{ID: "optB", IsCorrect: false},
		},
	}
	q2 := models.Question{
		ID: "q2", Type: models.QuestionTypeTrueFalse, PointsDefault: 2.0,
		Options: []models.QuestionOption{
			{ID: "optTrue", IsCorrect: true},
			{ID: "optFalse", IsCorrect: false},
		},
	}

	attemptID := "attempt_grading"
	mockRepo.On("GetAttempt", ctx, attemptID).Return(&models.AssessmentAttempt{
		ID: attemptID, AssessmentID: "exam_grading", Status: models.AttemptStatusInProgress,
	}, nil).Once()

	mockRepo.On("GetAssessmentQuestions", ctx, "exam_grading").Return([]models.Question{q1, q2}, nil)

	// Mock Student Responses
	// Student answered q1 correctly (optA)
	// Student answered q2 incorrectly (optFalse)
	optA := "optA"
	optFalse := "optFalse"
	
	resp1 := models.ItemResponse{QuestionID: "q1", SelectedOptionID: &optA}
	resp2 := models.ItemResponse{QuestionID: "q2", SelectedOptionID: &optFalse}
	
	mockRepo.On("ListResponses", ctx, attemptID).Return([]models.ItemResponse{resp1, resp2}, nil)

	// Expect SaveItemResponse calls with graded result
	// q1 -> Score 5, IsCorrect true
	mockRepo.On("SaveItemResponse", ctx, mock.MatchedBy(func(r models.ItemResponse) bool {
		return r.QuestionID == "q1" && r.Score == 5.0 && r.IsCorrect == true
	})).Return(nil)
	
	// q2 -> Score 0, IsCorrect false
	mockRepo.On("SaveItemResponse", ctx, mock.MatchedBy(func(r models.ItemResponse) bool {
		return r.QuestionID == "q2" && r.Score == 0.0 && r.IsCorrect == false
	})).Return(nil)

	// Expect CompleteAttempt with Total Score = 5.0 (5 + 0)
	mockRepo.On("CompleteAttempt", ctx, attemptID, 5.0).Return(nil)

	// Re-fetch attempt at end
	mockRepo.On("GetAttempt", ctx, attemptID).Return(&models.AssessmentAttempt{
		ID: attemptID, Status: models.AttemptStatusSubmitted, Score: 5.0,
	}, nil).Once()

	result, err := svc.CompleteAttempt(ctx, attemptID)
	assert.NoError(t, err)
	assert.Equal(t, 5.0, result.Score)
	mockRepo.AssertExpectations(t)
}
