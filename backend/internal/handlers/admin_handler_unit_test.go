package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MinimalMockAdminRepo struct {
	repository.AdminRepository
}

func (m *MinimalMockAdminRepo) ListStudentsForMonitor(ctx context.Context, f models.FilterParams) ([]models.StudentMonitorRow, error) {
	return []models.StudentMonitorRow{{ID: "s1", Name: "Student One"}}, nil
}
func (m *MinimalMockAdminRepo) GetAdvisorsForStudents(ctx context.Context, ids []string) (map[string][]models.AdvisorSummary, error) {
	return nil, nil
}
func (m *MinimalMockAdminRepo) GetDoneCountsForStudents(ctx context.Context, ids []string) (map[string]int, error) {
	return nil, nil
}
func (m *MinimalMockAdminRepo) GetLastUpdatesForStudents(ctx context.Context, ids []string) (map[string]time.Time, error) {
	return nil, nil
}
func (m *MinimalMockAdminRepo) GetRPRequiredForStudents(ctx context.Context, ids []string) (map[string]bool, error) {
	return nil, nil
}

func TestAdminHandler_MonitorStudents_Unit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	mockRepo := &MinimalMockAdminRepo{}
	pbm := &playbook.Manager{}
	svc := services.NewAdminService(mockRepo, pbm, config.AppConfig{}, nil)
	jSvc := services.NewJourneyService(nil, pbm, config.AppConfig{}, nil, nil, nil)
	h := handlers.NewAdminHandler(config.AppConfig{}, pbm, svc, jSvc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenant_id", "t1")
		c.Set("role", "admin")
		c.Next()
	})
	r.GET("/admin/monitor", h.MonitorStudents)

	req, _ := http.NewRequest("GET", "/admin/monitor", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp, 1)
	assert.Equal(t, "Student One", resp[0]["name"])
}
