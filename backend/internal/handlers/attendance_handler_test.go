package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// -- Mocks for Attendance --

type MockAttendanceService struct {
	mock.Mock
}

func (m *MockAttendanceService) BatchRecordAttendance(ctx context.Context, sessionID string, updates []models.ClassAttendance, teacherID string) error {
	args := m.Called(ctx, sessionID, updates, teacherID)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}

func (m *MockAttendanceService) GetSessionAttendance(ctx context.Context, sessionID string) ([]models.ClassAttendance, error) {
	return nil, nil // Not used in this handler test
}

func TestAttendanceHandler_BatchRecordAttendance(t *testing.T) {
	mockSvc := new(MockAttendanceService)
	handler := handlers.NewAttendanceHandler(mockSvc)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "teacher-1")
	})
	router.POST("/attendance/:session_id", handler.BatchRecordAttendance)

	sessionID := "sess-1"
	updates := []handlers.AttendanceUpdate{
		{StudentID: "s1", Status: "PRESENT"},
	}
	payload := handlers.BatchAttendanceRequest{Updates: updates}
	body, _ := json.Marshal(payload)

	// Expectations
	mockSvc.On("BatchRecordAttendance", mock.Anything, sessionID, mock.AnythingOfType("[]models.ClassAttendance"), "teacher-1").Return(nil)

	// Requests
	req, _ := http.NewRequest(http.MethodPost, "/attendance/"+sessionID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Error Case: Service returns error
	mockSvc.On("BatchRecordAttendance", mock.Anything, "sess-error", mock.Anything, "teacher-1").Return(assert.AnError)
	
	reqErr, _ := http.NewRequest(http.MethodPost, "/attendance/sess-error", bytes.NewBuffer(body))
	wErr := httptest.NewRecorder()
	router.ServeHTTP(wErr, reqErr)
	assert.Equal(t, http.StatusInternalServerError, wErr.Code)
}
