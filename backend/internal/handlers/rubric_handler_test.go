package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRubricRepository
type MockRubricRepository struct {
	mock.Mock
}

func (m *MockRubricRepository) CreateRubric(ctx context.Context, r *models.Rubric) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}
func (m *MockRubricRepository) GetRubric(ctx context.Context, id string) (*models.Rubric, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Rubric), args.Error(1)
}
func (m *MockRubricRepository) ListRubrics(ctx context.Context, courseID string) ([]models.Rubric, error) {
	args := m.Called(ctx, courseID)
	return args.Get(0).([]models.Rubric), args.Error(1)
}
func (m *MockRubricRepository) SubmitGrade(ctx context.Context, g *models.RubricGrade) error {
	args := m.Called(ctx, g)
	return args.Error(0)
}
func (m *MockRubricRepository) GetGrade(ctx context.Context, submissionID string) (*models.RubricGrade, error) {
	args := m.Called(ctx, submissionID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.RubricGrade), args.Error(1)
}

func setupRubricHandler() (*RubricHandler, *MockRubricRepository) {
	mockRepo := new(MockRubricRepository)
	svc := services.NewRubricService(mockRepo)
	handler := NewRubricHandler(svc)
	return handler, mockRepo
}

func TestRubricHandler_CreateRubric(t *testing.T) {
	handler, mockRepo := setupRubricHandler()
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		rubric := models.Rubric{
			Title:       "Test Rubric",
			Description: "Desc",
			Criteria: []models.RubricCriterion{
				{Title: "C1", Weight: 100, Levels: []models.RubricLevel{{Title: "L1", Points: 10}}},
			},
		}
		jsonBytes, _ := json.Marshal(rubric)
		c.Request, _ = http.NewRequest("POST", "/courses/c1/rubrics", bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: "c1"}}

		mockRepo.On("CreateRubric", mock.Anything, mock.AnythingOfType("*models.Rubric")).Return(nil)

		handler.CreateRubric(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestRubricHandler_ListRubrics(t *testing.T) {
	handler, mockRepo := setupRubricHandler()
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/courses/c1/rubrics", nil)
		c.Params = gin.Params{{Key: "id", Value: "c1"}}

		mockRepo.On("ListRubrics", mock.Anything, "c1").Return([]models.Rubric{{ID: "r1"}}, nil)

		handler.ListRubrics(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var list []models.Rubric
		json.Unmarshal(w.Body.Bytes(), &list)
		assert.Len(t, list, 1)
	})
}

func TestRubricHandler_GetRubric(t *testing.T) {
	handler, mockRepo := setupRubricHandler()
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/rubrics/r1", nil)
		c.Params = gin.Params{{Key: "id", Value: "r1"}}

		mockRepo.On("GetRubric", mock.Anything, "r1").Return(&models.Rubric{ID: "r1"}, nil)

		handler.GetRubric(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/rubrics/r99", nil)
		c.Params = gin.Params{{Key: "id", Value: "r99"}}

		mockRepo.On("GetRubric", mock.Anything, "r99").Return(nil, assert.AnError)

		handler.GetRubric(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestRubricHandler_SubmitGrade(t *testing.T) {
	handler, mockRepo := setupRubricHandler()
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		input := services.GradeInput{
			RubricID: "r1",
			Selections: []struct{
				CriterionID string `json:"criterion_id"`
				LevelID     string `json:"level_id"`
			}{
				{CriterionID: "c1", LevelID: "l1"},
			},
		}
		jsonBytes, _ := json.Marshal(input)
		c.Request, _ = http.NewRequest("POST", "/submissions/s1/rubric_grade", bytes.NewBuffer(jsonBytes))
		c.Params = gin.Params{{Key: "id", Value: "s1"}}
		c.Set("userID", "grader-1") // Mock middleware setting user

		// Mock GetRubric for validation
		rubric := &models.Rubric{
			ID: "r1",
			Criteria: []models.RubricCriterion{
				{
					ID: "c1", 
					Weight: 1.0, // Assuming 1.0 multiplier
					Levels: []models.RubricLevel{
						{ID: "l1", CriterionID: "c1", Points: 10.0},
					},
				},
			},
		}
		mockRepo.On("GetRubric", mock.Anything, "r1").Return(rubric, nil)

		// Mock SubmitGrade
		mockRepo.On("SubmitGrade", mock.Anything, mock.AnythingOfType("*models.RubricGrade")).Return(nil)

		handler.SubmitGrade(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}


