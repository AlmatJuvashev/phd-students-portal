package handlers

import (
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

type mockAnalyticsRepo struct {
	mock.Mock
}

func (m *mockAnalyticsRepo) GetStudentsByStage(ctx context.Context) ([]models.StudentStageStats, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.StudentStageStats), args.Error(1)
}
func (m *mockAnalyticsRepo) GetAdvisorLoad(ctx context.Context) ([]models.AdvisorLoadStats, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.AdvisorLoadStats), args.Error(1)
}
func (m *mockAnalyticsRepo) GetOverdueTasks(ctx context.Context) ([]models.OverdueTaskStats, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.OverdueTaskStats), args.Error(1)
}
func (m *mockAnalyticsRepo) GetTotalStudents(ctx context.Context, filter models.FilterParams) (int, error) {
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}
func (m *mockAnalyticsRepo) GetNodeCompletionCount(ctx context.Context, nodeID string, filter models.FilterParams) (int, error) {
	args := m.Called(ctx, nodeID, filter)
	return args.Int(0), args.Error(1)
}
func (m *mockAnalyticsRepo) GetDurationForNodes(ctx context.Context, nodeIDs []string, filter models.FilterParams) ([]float64, error) {
	args := m.Called(ctx, nodeIDs, filter)
	return args.Get(0).([]float64), args.Error(1)
}
func (m *mockAnalyticsRepo) GetBottleneck(ctx context.Context, filter models.FilterParams) (string, int, error) {
	args := m.Called(ctx, filter)
	return args.String(0), args.Int(1), args.Error(2)
}
func (m *mockAnalyticsRepo) GetProfileFlagCount(ctx context.Context, key string, minVal float64, filter models.FilterParams) (int, error) {
	args := m.Called(ctx, key, minVal, filter)
	return args.Int(0), args.Error(1)
}
func (m *mockAnalyticsRepo) SaveRiskSnapshot(ctx context.Context, s *models.RiskSnapshot) error {
	return m.Called(ctx, s).Error(0)
}
func (m *mockAnalyticsRepo) GetStudentRiskHistory(ctx context.Context, studentID string) ([]models.RiskSnapshot, error) {
	args := m.Called(ctx, studentID)
	return args.Get(0).([]models.RiskSnapshot), args.Error(1)
}
func (m *mockAnalyticsRepo) GetHighRiskStudents(ctx context.Context, threshold float64) ([]models.RiskSnapshot, error) {
	args := m.Called(ctx, threshold)
	return args.Get(0).([]models.RiskSnapshot), args.Error(1)
}
func (m *mockAnalyticsRepo) GetSessionAttendance(ctx context.Context, sessionID string) ([]models.ClassAttendance, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).([]models.ClassAttendance), args.Error(1)
}

func TestAnalyticsHandler_GetMonitorMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mockAnalyticsRepo)
	svc := services.NewAnalyticsService(mockRepo, nil, nil, nil)
	h := NewAnalyticsHandler(svc)

	filter := models.FilterParams{TenantID: "t1"}

	mockRepo.On("GetTotalStudents", mock.Anything, filter).Return(10, nil)
	mockRepo.On("GetNodeCompletionCount", mock.Anything, mock.Anything, filter).Return(5, nil)
	mockRepo.On("GetDurationForNodes", mock.Anything, mock.Anything, filter).Return([]float64{1.0}, nil)
	mockRepo.On("GetBottleneck", mock.Anything, filter).Return("node1", 1, nil)
	mockRepo.On("GetProfileFlagCount", mock.Anything, mock.Anything, mock.Anything, filter).Return(2, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/monitor", nil)
	c.Set("tenant_id", "t1")

	h.GetMonitorMetrics(c)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(10), resp["total_students_count"])
}

func TestAnalyticsHandler_GetHighRiskStudents(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mockAnalyticsRepo)
	svc := services.NewAnalyticsService(mockRepo, nil, nil, nil)
	h := NewAnalyticsHandler(svc)

	mockRepo.On("GetHighRiskStudents", mock.Anything, 50.0).Return([]models.RiskSnapshot{}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/high-risk?threshold=50.0", nil)

	h.GetHighRiskStudents(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
