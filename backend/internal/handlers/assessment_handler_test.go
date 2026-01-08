package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockAssessmentRepo struct {
	mock.Mock
}

func (m *mockAssessmentRepo) CreateQuestionBank(ctx context.Context, bank models.QuestionBank) (*models.QuestionBank, error) {
	args := m.Called(ctx, bank)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.QuestionBank), args.Error(1)
}
func (m *mockAssessmentRepo) GetQuestionBank(ctx context.Context, id string) (*models.QuestionBank, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.QuestionBank), args.Error(1)
}
func (m *mockAssessmentRepo) ListQuestionBanks(ctx context.Context, tenantID string) ([]models.QuestionBank, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.QuestionBank), args.Error(1)
}
func (m *mockAssessmentRepo) UpdateQuestionBank(ctx context.Context, bank models.QuestionBank) error {
	return m.Called(ctx, bank).Error(0)
}
func (m *mockAssessmentRepo) DeleteQuestionBank(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockAssessmentRepo) CreateQuestion(ctx context.Context, q models.Question) (*models.Question, error) {
	args := m.Called(ctx, q)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Question), args.Error(1)
}
func (m *mockAssessmentRepo) GetQuestion(ctx context.Context, id string) (*models.Question, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Question), args.Error(1)
}
func (m *mockAssessmentRepo) ListQuestionsByBank(ctx context.Context, bankID string) ([]models.Question, error) {
	args := m.Called(ctx, bankID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Question), args.Error(1)
}
func (m *mockAssessmentRepo) UpdateQuestion(ctx context.Context, q models.Question) error {
	return m.Called(ctx, q).Error(0)
}
func (m *mockAssessmentRepo) DeleteQuestion(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockAssessmentRepo) CreateAssessment(ctx context.Context, a models.Assessment) (*models.Assessment, error) {
	args := m.Called(ctx, a)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Assessment), args.Error(1)
}
func (m *mockAssessmentRepo) GetAssessment(ctx context.Context, id string) (*models.Assessment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Assessment), args.Error(1)
}
func (m *mockAssessmentRepo) ListAssessments(ctx context.Context, tenantID string, courseOfferingID string) ([]models.Assessment, error) {
	args := m.Called(ctx, tenantID, courseOfferingID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Assessment), args.Error(1)
}
func (m *mockAssessmentRepo) UpdateAssessment(ctx context.Context, a models.Assessment) error {
	return m.Called(ctx, a).Error(0)
}
func (m *mockAssessmentRepo) DeleteAssessment(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockAssessmentRepo) CreateAttempt(ctx context.Context, attempt models.AssessmentAttempt) (*models.AssessmentAttempt, error) {
	args := m.Called(ctx, attempt)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.AssessmentAttempt), args.Error(1)
}
func (m *mockAssessmentRepo) ListAttemptsByAssessmentAndStudent(ctx context.Context, assessmentID, studentID string) ([]models.AssessmentAttempt, error) {
	args := m.Called(ctx, assessmentID, studentID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.AssessmentAttempt), args.Error(1)
}
func (m *mockAssessmentRepo) SaveItemResponse(ctx context.Context, response models.ItemResponse) error {
	return m.Called(ctx, response).Error(0)
}
func (m *mockAssessmentRepo) CompleteAttempt(ctx context.Context, attemptID string, score float64) error {
	return m.Called(ctx, attemptID, score).Error(0)
}
func (m *mockAssessmentRepo) GetAttempt(ctx context.Context, id string) (*models.AssessmentAttempt, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.AssessmentAttempt), args.Error(1)
}
func (m *mockAssessmentRepo) ListResponses(ctx context.Context, attemptID string) ([]models.ItemResponse, error) {
	args := m.Called(ctx, attemptID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.ItemResponse), args.Error(1)
}
func (m *mockAssessmentRepo) LogProctoringEvent(ctx context.Context, log models.ProctoringLog) error {
	return m.Called(ctx, log).Error(0)
}
func (m *mockAssessmentRepo) CountProctoringEvents(ctx context.Context, attemptID string) (int, error) {
	args := m.Called(ctx, attemptID)
	return args.Int(0), args.Error(1)
}
func (m *mockAssessmentRepo) GetAssessmentQuestions(ctx context.Context, assessmentID string) ([]models.Question, error) {
	args := m.Called(ctx, assessmentID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.Question), args.Error(1)
}

func TestAssessmentHandler_CreateAssessment_Success(t *testing.T) {
	mockRepo := new(mockAssessmentRepo)
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
}

func TestAssessmentHandler_ListAssessments(t *testing.T) {
	mockRepo := new(mockAssessmentRepo)
	svc := services.NewAssessmentService(mockRepo)
	h := NewAssessmentHandler(svc)

	gin.SetMode(gin.TestMode)

	mockRepo.On("ListAssessments", mock.Anything, "t1", "co1").Return([]models.Assessment{{ID: "a1"}}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/assessments?course_offering_id=co1", nil)
	c.Set("tenant_id", "t1")

	h.ListAssessments(c)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestAssessmentHandler_UpdateAssessment(t *testing.T) {
	mockRepo := new(mockAssessmentRepo)
	svc := services.NewAssessmentService(mockRepo)
	h := NewAssessmentHandler(svc)

	gin.SetMode(gin.TestMode)

	mockRepo.On("GetAssessment", mock.Anything, "a1").Return(&models.Assessment{ID: "a1", TenantID: "t1"}, nil)
	mockRepo.On("UpdateAssessment", mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("GetAssessment", mock.Anything, "a1").Return(&models.Assessment{ID: "a1", Title: "Updated"}, nil)

	body, _ := json.Marshal(map[string]any{"title": "Updated"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/assessments/a1", bytes.NewBuffer(body))
	c.Params = gin.Params{{Key: "id", Value: "a1"}}
	c.Set("tenant_id", "t1")

	h.UpdateAssessment(c)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestAssessmentHandler_DeleteAssessment(t *testing.T) {
	mockRepo := new(mockAssessmentRepo)
	svc := services.NewAssessmentService(mockRepo)
	h := NewAssessmentHandler(svc)

	gin.SetMode(gin.TestMode)

	mockRepo.On("GetAssessment", mock.Anything, "a1").Return(&models.Assessment{ID: "a1", TenantID: "t1"}, nil)
	mockRepo.On("DeleteAssessment", mock.Anything, "a1").Return(nil)

	w := httptest.NewRecorder()
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", "t1")
		c.Next()
	})
	r.DELETE("/assessments/:id", h.DeleteAssessment)
	
	req, _ := http.NewRequest("DELETE", "/assessments/a1", nil)
	r.ServeHTTP(w, req)
	
	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestAssessmentHandler_CompleteAttempt(t *testing.T) {
	mockRepo := new(mockAssessmentRepo)
	svc := services.NewAssessmentService(mockRepo)
	h := NewAssessmentHandler(svc)

	gin.SetMode(gin.TestMode)

	// Calls in CompleteAttempt
	mockRepo.On("GetAttempt", mock.Anything, "at1").Return(&models.AssessmentAttempt{ID: "at1", StudentID: "s1", AssessmentID: "a1", Status: models.AttemptStatusInProgress}, nil).Once()
	mockRepo.On("GetAssessment", mock.Anything, "a1").Return(&models.Assessment{ID: "a1", TenantID: "t1"}, nil).Once()
	
	// Calls in completeAttempt (internal)
	mockRepo.On("GetAttempt", mock.Anything, "at1").Return(&models.AssessmentAttempt{ID: "at1", StudentID: "s1", AssessmentID: "a1", Status: models.AttemptStatusInProgress}, nil).Once()
	mockRepo.On("GetAssessmentQuestions", mock.Anything, "a1").Return([]models.Question{}, nil)
	mockRepo.On("ListResponses", mock.Anything, "at1").Return([]models.ItemResponse{}, nil)
	mockRepo.On("CompleteAttempt", mock.Anything, "at1", 0.0).Return(nil)
	mockRepo.On("GetAttempt", mock.Anything, "at1").Return(&models.AssessmentAttempt{ID: "at1", Score: 85}, nil).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/attempts/at1/complete", nil)
	c.Params = gin.Params{{Key: "id", Value: "at1"}}
	c.Set("tenant_id", "t1")
	c.Set("userID", "s1")

	h.CompleteAttempt(c)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestAssessmentHandler_LogProctoringEvent(t *testing.T) {
	mockRepo := new(mockAssessmentRepo)
	svc := services.NewAssessmentService(mockRepo)
	h := NewAssessmentHandler(svc)

	gin.SetMode(gin.TestMode)

	mockRepo.On("GetAttempt", mock.Anything, "at1").Return(&models.AssessmentAttempt{ID: "at1", StudentID: "s1", AssessmentID: "a1"}, nil)
	mockRepo.On("GetAssessment", mock.Anything, "a1").Return(&models.Assessment{ID: "a1", TenantID: "t1"}, nil)
	mockRepo.On("LogProctoringEvent", mock.Anything, mock.MatchedBy(func(l models.ProctoringLog) bool {
		return l.EventType == models.ProctoringEventTabSwitch
	})).Return(nil)

	body, _ := json.Marshal(map[string]any{"event_type": "TAB_SWITCH", "metadata": map[string]any{"count": 1}})
	w := httptest.NewRecorder()
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", "t1")
		c.Set("userID", "s1")
		c.Next()
	})
	r.POST("/attempts/:id/log", h.LogProctoringEvent)
	
	req, _ := http.NewRequest("POST", "/attempts/at1/log", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestAssessmentHandler_GetAssessment_HidesCorrectness(t *testing.T) {
	mockRepo := new(mockAssessmentRepo)
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

	h.GetAssessment(c)
	require.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Questions []models.Question `json:"questions"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Len(t, resp.Questions, 1)
	require.False(t, resp.Questions[0].Options[0].IsCorrect)
}

func TestAssessmentHandler_StartAttempt(t *testing.T) {
	mockRepo := new(mockAssessmentRepo)
	svc := services.NewAssessmentService(mockRepo)
	h := NewAssessmentHandler(svc)

	gin.SetMode(gin.TestMode)

	mockRepo.On("GetAssessment", mock.Anything, "ass-1").Return(&models.Assessment{
		ID:       "ass-1",
		TenantID: "t1",
	}, nil)
	mockRepo.On("ListAttemptsByAssessmentAndStudent", mock.Anything, "ass-1", "s1").Return([]models.AssessmentAttempt{}, nil)
	mockRepo.On("CreateAttempt", mock.Anything, mock.Anything).Return(&models.AssessmentAttempt{ID: "at1"}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/assessments/ass-1/attempts", nil)
	c.Params = gin.Params{{Key: "id", Value: "ass-1"}}
	c.Set("tenant_id", "t1")
	c.Set("userID", "s1")

	h.StartAttempt(c)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestAssessmentHandler_SubmitResponse(t *testing.T) {
	mockRepo := new(mockAssessmentRepo)
	svc := services.NewAssessmentService(mockRepo)
	h := NewAssessmentHandler(svc)

	gin.SetMode(gin.TestMode)

	mockRepo.On("GetAttempt", mock.Anything, "at1").Return(&models.AssessmentAttempt{ID: "at1", StudentID: "s1", AssessmentID: "a1", Status: models.AttemptStatusInProgress}, nil)
	mockRepo.On("GetAssessment", mock.Anything, "a1").Return(&models.Assessment{ID: "a1", TenantID: "t1"}, nil)
	mockRepo.On("SaveItemResponse", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(map[string]any{"question_id": "q1", "option_id": "o1"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/attempts/at1/response", bytes.NewBuffer(body))
	c.Params = gin.Params{{Key: "id", Value: "at1"}}
	c.Set("tenant_id", "t1")
	c.Set("userID", "s1")

	h.SubmitResponse(c)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestAssessmentHandler_GetAttemptDetails(t *testing.T) {
	mockRepo := new(mockAssessmentRepo)
	svc := services.NewAssessmentService(mockRepo)
	h := NewAssessmentHandler(svc)

	gin.SetMode(gin.TestMode)

	mockRepo.On("GetAttempt", mock.Anything, "at1").Return(&models.AssessmentAttempt{ID: "at1", StudentID: "s1", AssessmentID: "a1", Status: models.AttemptStatusSubmitted}, nil)
	mockRepo.On("GetAssessment", mock.Anything, "a1").Return(&models.Assessment{ID: "a1", TenantID: "t1"}, nil)
	mockRepo.On("GetAssessmentQuestions", mock.Anything, "a1").Return([]models.Question{{ID: "q1"}}, nil)
	mockRepo.On("ListResponses", mock.Anything, "at1").Return([]models.ItemResponse{{QuestionID: "q1"}}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/attempts/at1", nil)
	c.Params = gin.Params{{Key: "id", Value: "at1"}}
	c.Set("tenant_id", "t1")
	c.Set("userID", "s1")

	h.GetAttemptDetails(c)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestAssessmentHandler_ErrorPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockAssessmentRepo)
	svc := services.NewAssessmentService(repo)
	h := NewAssessmentHandler(svc)

	t.Run("CreateAssessment_InvalidJSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/assessments", bytes.NewBufferString("invalid json"))
		h.CreateAssessment(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetAssessment_NotFound", func(t *testing.T) {
		repo.On("GetAssessment", mock.Anything, "a1").Return(nil, nil).Once()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "a1"}}
		c.Request, _ = http.NewRequest("GET", "/assessments/a1", nil)
		c.Set("tenant_id", "t1")
		h.GetAssessment(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("DeleteAssessment_ServiceError", func(t *testing.T) {
		repo.On("GetAssessment", mock.Anything, "a1").Return(&models.Assessment{ID: "a1"}, nil).Once()
		repo.On("DeleteAssessment", mock.Anything, "a1").Return(assert.AnError).Once()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "a1"}}
		c.Request, _ = http.NewRequest("DELETE", "/assessments/a1", nil)
		h.DeleteAssessment(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAssessmentHandler_ListMyAttempts(t *testing.T) {
	mockRepo := new(mockAssessmentRepo)
	svc := services.NewAssessmentService(mockRepo)
	h := NewAssessmentHandler(svc)

	gin.SetMode(gin.TestMode)

	mockRepo.On("GetAssessment", mock.Anything, "a1").Return(&models.Assessment{ID: "a1", TenantID: "t1"}, nil)
	mockRepo.On("ListAttemptsByAssessmentAndStudent", mock.Anything, "a1", "s1").Return([]models.AssessmentAttempt{{ID: "at1"}}, nil)

	w := httptest.NewRecorder()
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", "t1")
		c.Set("userID", "s1")
		c.Next()
	})
	r.GET("/assessments/:id/my-attempts", h.ListMyAttempts)

	req, _ := http.NewRequest("GET", "/assessments/a1/my-attempts", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestAssessmentHandler_DeleteAssessment_Forbidden(t *testing.T) {
	mockRepo := new(mockAssessmentRepo)
	svc := services.NewAssessmentService(mockRepo)
	h := NewAssessmentHandler(svc)

	gin.SetMode(gin.TestMode)

	mockRepo.On("GetAssessment", mock.Anything, "a1").Return(&models.Assessment{ID: "a1", TenantID: "t2"}, nil) // Different tenant

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("DELETE", "/assessments/a1", nil)
	c.Params = gin.Params{{Key: "id", Value: "a1"}}
	c.Set("tenant_id", "t1")

	h.DeleteAssessment(c)
	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestAssessmentHandler_StartAttempt_Errors(t *testing.T) {
	mockRepo := new(mockAssessmentRepo)
	svc := services.NewAssessmentService(mockRepo)
	h := NewAssessmentHandler(svc)

	gin.SetMode(gin.TestMode)

	t.Run("MaxAttempts", func(t *testing.T) {
		settings, _ := json.Marshal(models.SecuritySettings{MaxAttempts: 1})
		mockRepo.On("GetAssessment", mock.Anything, "a1").Return(&models.Assessment{
			ID:               "a1",
			TenantID:         "t1",
			SecuritySettings: types.JSONText(settings),
		}, nil).Once()
		mockRepo.On("ListAttemptsByAssessmentAndStudent", mock.Anything, "a1", "s1").Return([]models.AssessmentAttempt{{ID: "prev"}}, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/assessments/a1/attempts", nil)
		c.Params = gin.Params{{Key: "id", Value: "a1"}}
		c.Set("tenant_id", "t1")
		c.Set("userID", "s1")

		h.StartAttempt(c)
		require.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "MAX_ATTEMPTS_REACHED")
	})

	t.Run("AttemptInProgress", func(t *testing.T) {
		mockRepo.On("GetAssessment", mock.Anything, "a1").Return(&models.Assessment{ID: "a1", TenantID: "t1"}, nil).Once()
		mockRepo.On("ListAttemptsByAssessmentAndStudent", mock.Anything, "a1", "s1").Return([]models.AssessmentAttempt{{ID: "prev", Status: models.AttemptStatusInProgress}}, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/assessments/a1/attempts", nil)
		c.Params = gin.Params{{Key: "id", Value: "a1"}}
		c.Set("tenant_id", "t1")
		c.Set("userID", "s1")

		h.StartAttempt(c)
		require.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "ATTEMPT_IN_PROGRESS")
	})

	t.Run("CooldownActive", func(t *testing.T) {
		settings, _ := json.Marshal(models.SecuritySettings{CooldownMinutes: 60})
		mockRepo.On("GetAssessment", mock.Anything, "a1").Return(&models.Assessment{
			ID:               "a1",
			TenantID:         "t1",
			SecuritySettings: types.JSONText(settings),
		}, nil).Once()
		
		lastAttemptTime := time.Now().Add(-30 * time.Minute)
		mockRepo.On("ListAttemptsByAssessmentAndStudent", mock.Anything, "a1", "s1").Return([]models.AssessmentAttempt{{ID: "prev", Status: models.AttemptStatusSubmitted, StartedAt: lastAttemptTime, FinishedAt: &lastAttemptTime}}, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/assessments/a1/attempts", nil)
		c.Params = gin.Params{{Key: "id", Value: "a1"}}
		c.Set("tenant_id", "t1")
		c.Set("userID", "s1")

		h.StartAttempt(c)
		require.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.Contains(t, w.Body.String(), "COOLDOWN_ACTIVE")
	})
}

func TestAssessmentHandler_SubmitResponse_Errors(t *testing.T) {
	mockRepo := new(mockAssessmentRepo)
	svc := services.NewAssessmentService(mockRepo)
	h := NewAssessmentHandler(svc)
	gin.SetMode(gin.TestMode)

	t.Run("AutoSubmitted", func(t *testing.T) {
		timeLimit := 10
		startTime := time.Now().Add(-15 * time.Minute)
		mockRepo.On("GetAttempt", mock.Anything, "at1").Return(&models.AssessmentAttempt{
			ID: "at1", StudentID: "s1", AssessmentID: "a1", Status: models.AttemptStatusInProgress, StartedAt: startTime,
		}, nil).Once()
		mockRepo.On("GetAssessment", mock.Anything, "a1").Return(&models.Assessment{
			ID: "a1", TenantID: "t1", TimeLimitMinutes: &timeLimit,
		}, nil).Once()
		
		// Internal auto-submit calls
		mockRepo.On("GetAttempt", mock.Anything, "at1").Return(&models.AssessmentAttempt{ID: "at1", Status: models.AttemptStatusInProgress, AssessmentID: "a1"}, nil).Once()
		mockRepo.On("GetAssessmentQuestions", mock.Anything, "a1").Return([]models.Question{}, nil)
		mockRepo.On("ListResponses", mock.Anything, "at1").Return([]models.ItemResponse{}, nil)
		mockRepo.On("CompleteAttempt", mock.Anything, "at1", 0.0).Return(nil)
		mockRepo.On("GetAttempt", mock.Anything, "at1").Return(&models.AssessmentAttempt{ID: "at1", Status: models.AttemptStatusSubmitted}, nil).Once()

		body, _ := json.Marshal(map[string]any{"question_id": "q1", "option_id": "o1"})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/attempts/at1/response", bytes.NewBuffer(body))
		c.Params = gin.Params{{Key: "id", Value: "at1"}}
		c.Set("tenant_id", "t1")
		c.Set("userID", "s1")

		h.SubmitResponse(c)
		require.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "ATTEMPT_AUTO_SUBMITTED")
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/attempts/at1/response", bytes.NewBufferString("invalid"))
		h.SubmitResponse(c)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Forbidden", func(t *testing.T) {
		mockRepo.On("GetAttempt", mock.Anything, "at1").Return(&models.AssessmentAttempt{
			ID: "at1", StudentID: "other-student", AssessmentID: "a1",
		}, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/attempts/at1/response", bytes.NewBufferString("{}"))
		c.Params = gin.Params{{Key: "id", Value: "at1"}}
		c.Set("tenant_id", "t1")
		c.Set("userID", "s1")

		h.SubmitResponse(c)
		require.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockRepo.On("GetAttempt", mock.Anything, "at1").Return(&models.AssessmentAttempt{
			ID: "at1", StudentID: "s1", AssessmentID: "a1", Status: models.AttemptStatusInProgress,
		}, nil).Once()
		mockRepo.On("GetAssessment", mock.Anything, "a1").Return(&models.Assessment{
			ID: "a1", TenantID: "t1",
		}, nil).Once()
		mockRepo.On("SaveItemResponse", mock.Anything, mock.Anything).Return(assert.AnError).Once()

		body, _ := json.Marshal(map[string]any{"question_id": "q1", "option_id": "o1"})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/attempts/at1/response", bytes.NewBuffer(body))
		c.Params = gin.Params{{Key: "id", Value: "at1"}}
		c.Set("tenant_id", "t1")
		c.Set("userID", "s1")

		h.SubmitResponse(c)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAssessmentHandler_OtherForbidden(t *testing.T) {
	mockRepo := new(mockAssessmentRepo)
	svc := services.NewAssessmentService(mockRepo)
	h := NewAssessmentHandler(svc)
	gin.SetMode(gin.TestMode)

	t.Run("GetAssessment_Forbidden", func(t *testing.T) {
		mockRepo.On("GetAssessment", mock.Anything, "a1").Return(&models.Assessment{ID: "a1", TenantID: "t2"}, nil).Once()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "a1"}}
		c.Request, _ = http.NewRequest("GET", "/assessments/a1", nil)
		c.Set("tenant_id", "t1")
		h.GetAssessment(c)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("UpdateAssessment_Forbidden", func(t *testing.T) {
		mockRepo.On("GetAssessment", mock.Anything, "a1").Return(&models.Assessment{ID: "a1", TenantID: "t2"}, nil).Once()
		body, _ := json.Marshal(map[string]any{"title": "Updated"})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "a1"}}
		c.Request, _ = http.NewRequest("PUT", "/assessments/a1", bytes.NewBuffer(body))
		c.Set("tenant_id", "t1")
		h.UpdateAssessment(c)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}


func TestAssessmentHandler_ExtendedErrorPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("LogProctoringEvent_BindError", func(t *testing.T) {
		repo := new(mockAssessmentRepo)
		svc := services.NewAssessmentService(repo)
		h := NewAssessmentHandler(svc)
		
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/attempts/at1/log", bytes.NewBufferString("invalid json"))
		h.LogProctoringEvent(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("LogProctoringEvent_Forbidden", func(t *testing.T) {
		repo := new(mockAssessmentRepo)
		svc := services.NewAssessmentService(repo)
		h := NewAssessmentHandler(svc)
		
		repo.On("GetAttempt", mock.Anything, "at1").Return(&models.AssessmentAttempt{ID: "at1", AssessmentID: "a1"}, nil)
		repo.On("GetAssessment", mock.Anything, "a1").Return(&models.Assessment{ID: "a1", TenantID: "other-tenant"}, nil)
		
		body, _ := json.Marshal(map[string]any{"event_type": "TAB_SWITCH", "metadata": map[string]any{}})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/attempts/at1/log", bytes.NewBuffer(body))
		c.Set("tenant_id", "t1")
		c.Set("userID", "u1")
		c.Params = gin.Params{{Key: "id", Value: "at1"}}

		h.LogProctoringEvent(c)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("LogProctoringEvent_ServiceError", func(t *testing.T) {
		repo := new(mockAssessmentRepo)
		svc := services.NewAssessmentService(repo)
		h := NewAssessmentHandler(svc)

		repo.On("GetAttempt", mock.Anything, "at1").Return(nil, assert.AnError) // Repo failure
		
		body, _ := json.Marshal(map[string]any{"event_type": "TAB_SWITCH", "metadata": map[string]any{}})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/attempts/at1/log", bytes.NewBuffer(body))
		c.Params = gin.Params{{Key: "id", Value: "at1"}}
		c.Set("tenant_id", "t1")

		h.LogProctoringEvent(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("CompleteAttempt_Forbidden", func(t *testing.T) {
		repo := new(mockAssessmentRepo)
		svc := services.NewAssessmentService(repo)
		h := NewAssessmentHandler(svc)
		
		repo.On("GetAttempt", mock.Anything, "at1").Return(&models.AssessmentAttempt{ID: "at1", AssessmentID: "a1"}, nil)
		repo.On("GetAssessment", mock.Anything, "a1").Return(&models.Assessment{ID: "a1", TenantID: "other"}, nil) // Tenant mismatch

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/attempts/at1/complete", nil)
		c.Params = gin.Params{{Key: "id", Value: "at1"}}
		c.Set("tenant_id", "t1")
		c.Set("userID", "u1")

		h.CompleteAttempt(c)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("CompleteAttempt_ServiceError", func(t *testing.T) {
		repo := new(mockAssessmentRepo)
		svc := services.NewAssessmentService(repo)
		h := NewAssessmentHandler(svc)
		
		repo.On("GetAttempt", mock.Anything, "at1").Return(nil, assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/attempts/at1/complete", nil)
		c.Params = gin.Params{{Key: "id", Value: "at1"}}
		c.Set("tenant_id", "t1")

		h.CompleteAttempt(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("ListMyAttempts_Forbidden", func(t *testing.T) {
		repo := new(mockAssessmentRepo)
		svc := services.NewAssessmentService(repo)
		h := NewAssessmentHandler(svc)
		
		repo.On("GetAssessment", mock.Anything, "a1").Return(&models.Assessment{ID: "a1", TenantID: "other"}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/assessments/a1/my-attempts", nil)
		c.Params = gin.Params{{Key: "id", Value: "a1"}}
		c.Set("tenant_id", "t1")

		h.ListMyAttempts(c)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("ListMyAttempts_ServiceError", func(t *testing.T) {
		repo := new(mockAssessmentRepo)
		svc := services.NewAssessmentService(repo)
		h := NewAssessmentHandler(svc)
		
		repo.On("GetAssessment", mock.Anything, "a1").Return(nil, assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/assessments/a1/my-attempts", nil)
		c.Params = gin.Params{{Key: "id", Value: "a1"}}
		c.Set("tenant_id", "t1")

		h.ListMyAttempts(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
	
	t.Run("GetAttemptDetails_Forbidden", func(t *testing.T) {
		repo := new(mockAssessmentRepo)
		svc := services.NewAssessmentService(repo)
		h := NewAssessmentHandler(svc)

		repo.On("GetAttempt", mock.Anything, "at1").Return(&models.AssessmentAttempt{ID: "at1", AssessmentID: "a1"}, nil)
		repo.On("GetAssessment", mock.Anything, "a1").Return(&models.Assessment{ID: "a1", TenantID: "other"}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/attempts/at1", nil)
		c.Params = gin.Params{{Key: "id", Value: "at1"}}
		c.Set("tenant_id", "t1")
		c.Set("userID", "u1")

		h.GetAttemptDetails(c)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("GetAttemptDetails_ServiceError", func(t *testing.T) {
		repo := new(mockAssessmentRepo)
		svc := services.NewAssessmentService(repo)
		h := NewAssessmentHandler(svc)

		repo.On("GetAttempt", mock.Anything, "at1").Return(nil, assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/attempts/at1", nil)
		c.Params = gin.Params{{Key: "id", Value: "at1"}}
		c.Set("tenant_id", "t1")

		h.GetAttemptDetails(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
