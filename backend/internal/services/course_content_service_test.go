package services

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepo for service testing
type MockCourseContentRepo struct {
	mock.Mock
}

func (m *MockCourseContentRepo) CreateModule(ctx context.Context, mod *models.CourseModule) error {
	args := m.Called(ctx, mod)
	return args.Error(0)
}
func (m *MockCourseContentRepo) GetModule(ctx context.Context, id string) (*models.CourseModule, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.CourseModule), args.Error(1)
}
func (m *MockCourseContentRepo) ListModules(ctx context.Context, courseID string) ([]models.CourseModule, error) {
	args := m.Called(ctx, courseID)
	return args.Get(0).([]models.CourseModule), args.Error(1)
}
func (m *MockCourseContentRepo) UpdateModule(ctx context.Context, mod *models.CourseModule) error {
	args := m.Called(ctx, mod)
	return args.Error(0)
}
func (m *MockCourseContentRepo) DeleteModule(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockCourseContentRepo) CreateLesson(ctx context.Context, l *models.CourseLesson) error {
	args := m.Called(ctx, l)
	return args.Error(0)
}
func (m *MockCourseContentRepo) GetLesson(ctx context.Context, id string) (*models.CourseLesson, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.CourseLesson), args.Error(1)
}
func (m *MockCourseContentRepo) ListLessons(ctx context.Context, moduleID string) ([]models.CourseLesson, error) {
	args := m.Called(ctx, moduleID)
	return args.Get(0).([]models.CourseLesson), args.Error(1)
}
func (m *MockCourseContentRepo) UpdateLesson(ctx context.Context, l *models.CourseLesson) error {
	args := m.Called(ctx, l)
	return args.Error(0)
}
func (m *MockCourseContentRepo) DeleteLesson(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockCourseContentRepo) CreateActivity(ctx context.Context, a *models.CourseActivity) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}
func (m *MockCourseContentRepo) GetActivity(ctx context.Context, id string) (*models.CourseActivity, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.CourseActivity), args.Error(1)
}
func (m *MockCourseContentRepo) ListActivities(ctx context.Context, lessonID string) ([]models.CourseActivity, error) {
	args := m.Called(ctx, lessonID)
	return args.Get(0).([]models.CourseActivity), args.Error(1)
}
func (m *MockCourseContentRepo) UpdateActivity(ctx context.Context, a *models.CourseActivity) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}
func (m *MockCourseContentRepo) DeleteActivity(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Ensure MockRepo implements interface
var _ repository.CourseContentRepository = (*MockCourseContentRepo)(nil)

func TestCourseContentService_ValidateContent(t *testing.T) {
	mockRepo := new(MockCourseContentRepo)
	svc := NewCourseContentService(mockRepo)
	ctx := context.Background()

	// 1. Valid Quiz
	validQuizJSON := `{"timeLimit": 60, "passingScore": 80, "questions": [{"id": "q1", "type": "multiple_choice", "text": "Q1", "points": 10}]}`
	a1 := &models.CourseActivity{
		LessonID: "l1",
		Title:    "Valid Quiz",
		Type:     "quiz",
		Content:  validQuizJSON,
	}
	mockRepo.On("CreateActivity", ctx, a1).Return(nil)
	err := svc.CreateActivity(ctx, a1)
	assert.NoError(t, err)

	// 2. Invalid JSON for Quiz
	a2 := &models.CourseActivity{
		LessonID: "l1",
		Title:    "Invalid JSON",
		Type:     "quiz",
		Content:  `{ "broken": `,
	}
	// Validate should fail before calling repo
	err = svc.CreateActivity(ctx, a2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid quiz config json")

	// 3. Valid Survey
	validSurveyJSON := `{"anonymous": true, "showProgressBar": true, "questions": [{"id": "s1", "type": "rating_stars", "text": "Rate", "required": true}]}`
	a3 := &models.CourseActivity{
		LessonID: "l1",
		Title:    "Valid Survey",
		Type:     "survey",
		Content:  validSurveyJSON,
	}
	mockRepo.On("CreateActivity", ctx, a3).Return(nil)
	err = svc.CreateActivity(ctx, a3)
	assert.NoError(t, err)
}
