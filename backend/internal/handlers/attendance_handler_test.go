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

type mockAttendanceRepo struct {
	mock.Mock
}

func (m *mockAttendanceRepo) BatchUpsertAttendance(ctx context.Context, sessionID string, records []models.ClassAttendance, recordedBy string) error {
	args := m.Called(ctx, sessionID, records, recordedBy)
	return args.Error(0)
}
func (m *mockAttendanceRepo) GetSessionAttendance(ctx context.Context, sessionID string) ([]models.ClassAttendance, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ClassAttendance), args.Error(1)
}
func (m *mockAttendanceRepo) GetStudentAttendance(ctx context.Context, studentID string) ([]models.ClassAttendance, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ClassAttendance), args.Error(1)
}
func (m *mockAttendanceRepo) RecordAttendance(ctx context.Context, sessionID string, record models.ClassAttendance) error {
	return m.Called(ctx, sessionID, record).Error(0)
}

func TestAttendanceHandler_BatchRecordAttendance(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockAttendanceRepo)
	svc := services.NewAttendanceService(repo)
	h := NewAttendanceHandler(svc)

	updates := BatchAttendanceRequest{
		Updates: []AttendanceUpdate{
			{StudentID: "s1", Status: "present"},
		},
	}
	repo.On("BatchUpsertAttendance", mock.Anything, "s1", mock.Anything, "u1").Return(nil)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(updates)
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/sessions/s1/attendance", bytes.NewBuffer(body))
	c.Params = gin.Params{{Key: "session_id", Value: "s1"}}
	c.Set("userID", "u1")

	h.BatchRecordAttendance(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAttendanceHandler_GetSessionAttendance(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockAttendanceRepo)
	svc := services.NewAttendanceService(repo)
	h := NewAttendanceHandler(svc)

	repo.On("GetSessionAttendance", mock.Anything, "s1").Return([]models.ClassAttendance{}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/sessions/s1/attendance", nil)
	c.Params = gin.Params{{Key: "session_id", Value: "s1"}}

	h.GetSessionAttendance(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
