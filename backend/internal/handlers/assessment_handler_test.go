package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAssessmentHandler_CreateAssessment_Success(t *testing.T) {
	mockRepo := new(mockAssessmentRepoForItemBank)
	svc := services.NewAssessmentService(mockRepo)
	h := NewAssessmentHandler(svc)

	gin.SetMode(gin.TestMode)

	mockRepo.On("CreateAssessment", mock.Anything, mock.MatchedBy(func(a models.Assessment) bool {
		return a.TenantID == "tenant-1" && a.CreatedBy == "user-1" && a.CourseOfferingID == "off-1" && a.Title == "Midterm"
	})).Return(&models.Assessment{
		ID:               "ass-1",
		TenantID:         "tenant-1",
		CourseOfferingID: "off-1",
		Title:            "Midterm",
	}, nil)

	payload, _ := json.Marshal(map[string]any{
		"course_offering_id": "off-1",
		"title":              "Midterm",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/assessments", bytes.NewBuffer(payload))
	c.Set("tenant_id", "tenant-1")
	c.Set("userID", "user-1")

	h.CreateAssessment(c)
	require.Equal(t, http.StatusCreated, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestAssessmentHandler_GetAssessment_HidesCorrectness(t *testing.T) {
	mockRepo := new(mockAssessmentRepoForItemBank)
	svc := services.NewAssessmentService(mockRepo)
	h := NewAssessmentHandler(svc)

	gin.SetMode(gin.TestMode)

	mockRepo.On("GetAssessment", mock.Anything, "ass-1").Return(&models.Assessment{
		ID:       "ass-1",
		TenantID: "tenant-1",
	}, nil)
	mockRepo.On("GetAssessmentQuestions", mock.Anything, "ass-1").Return([]models.Question{
		{
			ID:    "q1",
			Type:  models.QuestionTypeMCQ,
			Stem:  "Question 1",
			Options: []models.QuestionOption{
				{ID: "o1", Text: "A", IsCorrect: true},
				{ID: "o2", Text: "B", IsCorrect: false},
			},
		},
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/assessments/ass-1", nil)
	c.Params = gin.Params{{Key: "id", Value: "ass-1"}}
	c.Set("tenant_id", "tenant-1")
	c.Set("userID", "user-1")

	h.GetAssessment(c)
	require.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Questions []models.Question `json:"questions"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Len(t, resp.Questions, 1)
	require.Len(t, resp.Questions[0].Options, 2)
	require.False(t, resp.Questions[0].Options[0].IsCorrect)
	require.False(t, resp.Questions[0].Options[1].IsCorrect)
}

func TestAssessmentHandler_StartAttempt_ConflictWhenInProgress(t *testing.T) {
	mockRepo := new(mockAssessmentRepoForItemBank)
	svc := services.NewAssessmentService(mockRepo)
	h := NewAssessmentHandler(svc)

	gin.SetMode(gin.TestMode)

	mockRepo.On("GetAssessment", mock.Anything, "ass-1").Return(&models.Assessment{
		ID:       "ass-1",
		TenantID: "tenant-1",
	}, nil)

	mockRepo.On("ListAttemptsByAssessmentAndStudent", mock.Anything, "ass-1", "stud-1").Return([]models.AssessmentAttempt{
		{
			ID:           "att-1",
			AssessmentID: "ass-1",
			StudentID:    "stud-1",
			StartedAt:    time.Now(),
			Status:       models.AttemptStatusInProgress,
		},
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/assessments/ass-1/attempts", nil)
	c.Params = gin.Params{{Key: "id", Value: "ass-1"}}
	c.Set("tenant_id", "tenant-1")
	c.Set("userID", "stud-1")

	h.StartAttempt(c)
	require.Equal(t, http.StatusConflict, w.Code)
}

func TestAssessmentHandler_SubmitResponse_ForbiddenWhenNotOwner(t *testing.T) {
	mockRepo := new(mockAssessmentRepoForItemBank)
	svc := services.NewAssessmentService(mockRepo)
	h := NewAssessmentHandler(svc)

	gin.SetMode(gin.TestMode)

	mockRepo.On("GetAttempt", mock.Anything, "att-1").Return(&models.AssessmentAttempt{
		ID:           "att-1",
		AssessmentID: "ass-1",
		StudentID:    "someone-else",
		StartedAt:    time.Now(),
		Status:       models.AttemptStatusInProgress,
	}, nil)

	body, _ := json.Marshal(map[string]any{"question_id": "q1", "option_id": "o1"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/attempts/att-1/response", bytes.NewBuffer(body))
	c.Params = gin.Params{{Key: "id", Value: "att-1"}}
	c.Set("tenant_id", "tenant-1")
	c.Set("userID", "stud-1")

	h.SubmitResponse(c)
	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestAssessmentHandler_GetAttemptDetails_RevealsCorrectAfterSubmit(t *testing.T) {
	mockRepo := new(mockAssessmentRepoForItemBank)
	svc := services.NewAssessmentService(mockRepo)
	h := NewAssessmentHandler(svc)

	gin.SetMode(gin.TestMode)

	mockRepo.On("GetAttempt", mock.Anything, "att-1").Return(&models.AssessmentAttempt{
		ID:           "att-1",
		AssessmentID: "ass-1",
		StudentID:    "stud-1",
		StartedAt:    time.Now().Add(-30 * time.Minute),
		Status:       models.AttemptStatusSubmitted,
	}, nil)
	mockRepo.On("GetAssessment", mock.Anything, "ass-1").Return(&models.Assessment{
		ID:       "ass-1",
		TenantID: "tenant-1",
	}, nil)
	mockRepo.On("GetAssessmentQuestions", mock.Anything, "ass-1").Return([]models.Question{
		{
			ID:    "q1",
			Type:  models.QuestionTypeMCQ,
			Stem:  "Question 1",
			Options: []models.QuestionOption{
				{ID: "o1", Text: "A", IsCorrect: true},
				{ID: "o2", Text: "B", IsCorrect: false},
			},
		},
	}, nil)
	mockRepo.On("ListResponses", mock.Anything, "att-1").Return([]models.ItemResponse{
		{AttemptID: "att-1", QuestionID: "q1", SelectedOptionID: func() *string { s := "o1"; return &s }(), Score: 1, IsCorrect: true},
	}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/attempts/att-1", nil)
	c.Params = gin.Params{{Key: "id", Value: "att-1"}}
	c.Set("tenant_id", "tenant-1")
	c.Set("userID", "stud-1")

	h.GetAttemptDetails(c)
	require.Equal(t, http.StatusOK, w.Code)

	var payload struct {
		Questions []models.Question `json:"questions"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &payload))
	require.Len(t, payload.Questions, 1)
	require.True(t, payload.Questions[0].Options[0].IsCorrect)
}
