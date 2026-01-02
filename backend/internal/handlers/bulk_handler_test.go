package handlers_test

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// -- Mocks --

type MockBulkService struct {
	mock.Mock
}

func (m *MockBulkService) ImportStudents(ctx context.Context, r io.Reader, tenantID string) (int, []error) {
	args := m.Called(ctx, mock.Anything, tenantID) // Cannot match Reader easily, assume it flows
	return args.Int(0), args.Get(1).([]error)
}

func TestBulkHandler_BulkEnrollStudents(t *testing.T) {
	mockSvc := new(MockBulkService)
	handler := handlers.NewBulkHandler(mockSvc)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("tenantID", "curr-tenant") // Mock Middleware
	})
	router.POST("/bulk", handler.BulkEnrollStudents)

	// Create Multipart Form
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "students.csv")
	part.Write([]byte("name,email\na,b"))
	writer.Close()

	// Expectations
	mockSvc.On("ImportStudents", mock.Anything, mock.Anything, "curr-tenant").Return(5, []error{})

	// Request
	req, _ := http.NewRequest(http.MethodPost, "/bulk", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
	
	// Error Case: No File
	reqNoFile, _ := http.NewRequest(http.MethodPost, "/bulk", nil)
	wNoFile := httptest.NewRecorder()
	router.ServeHTTP(wNoFile, reqNoFile)
	assert.Equal(t, http.StatusBadRequest, wNoFile.Code)
}
