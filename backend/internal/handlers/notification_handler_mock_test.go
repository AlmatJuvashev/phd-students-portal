package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// -- Mocks --

type MockNotificationRepo struct {
	mock.Mock
}

func (m *MockNotificationRepo) Create(ctx context.Context, n *models.Notification) error {
	return m.Called(ctx, n).Error(0)
}
func (m *MockNotificationRepo) GetUnread(ctx context.Context, userID string) ([]models.Notification, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Notification), args.Error(1)
}
func (m *MockNotificationRepo) MarkAsRead(ctx context.Context, id, userID string) error {
	return m.Called(ctx, id, userID).Error(0)
}
func (m *MockNotificationRepo) MarkAllAsRead(ctx context.Context, userID string) error {
	return m.Called(ctx, userID).Error(0)
}
func (m *MockNotificationRepo) ListByRecipient(ctx context.Context, userID string, limit int) ([]models.Notification, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Notification), args.Error(1)
}
func (m *MockNotificationRepo) CountUnread(ctx context.Context, userID string) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func TestNotificationHandler_ErrorPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("GetUnread_ServiceError", func(t *testing.T) {
		mockRepo := new(MockNotificationRepo)
		svc := services.NewNotificationService(mockRepo)
		h := handlers.NewNotificationHandler(svc)

		mockRepo.On("GetUnread", mock.Anything, "u1").Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/notifications/unread", nil)
		c.Set("userID", "u1")

		h.GetUnread(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("MarkAsRead_ServiceError", func(t *testing.T) {
		mockRepo := new(MockNotificationRepo)
		svc := services.NewNotificationService(mockRepo)
		h := handlers.NewNotificationHandler(svc)

		mockRepo.On("MarkAsRead", mock.Anything, "n1", "u1").Return(errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/notifications/n1/read", nil)
		c.Params = gin.Params{{Key: "id", Value: "n1"}}
		c.Set("userID", "u1")

		h.MarkAsRead(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("MarkAllAsRead_ServiceError", func(t *testing.T) {
		mockRepo := new(MockNotificationRepo)
		svc := services.NewNotificationService(mockRepo)
		h := handlers.NewNotificationHandler(svc)

		mockRepo.On("MarkAllAsRead", mock.Anything, "u1").Return(errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/notifications/read-all", nil)
		c.Set("userID", "u1")

		h.MarkAllAsRead(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("GetNotifications_ServiceError", func(t *testing.T) {
		mockRepo := new(MockNotificationRepo)
		svc := services.NewNotificationService(mockRepo)
		h := handlers.NewNotificationHandler(svc)

		mockRepo.On("ListByRecipient", mock.Anything, "u1", 50).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/notifications", nil)
		c.Set("userID", "u1")

		h.GetNotifications(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
