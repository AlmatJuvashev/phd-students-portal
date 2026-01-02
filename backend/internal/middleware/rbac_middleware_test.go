package middleware_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/middleware"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRBACRepo needed for AuthzService
type MockRBACRepo struct {
	mock.Mock
}

func (m *MockRBACRepo) GetUserRolesInContext(ctx context.Context, userID uuid.UUID, contextType string, contextID uuid.UUID) ([]models.RoleDef, error) {
	args := m.Called(ctx, userID, contextType, contextID)
	return args.Get(0).([]models.RoleDef), args.Error(1)
}

func (m *MockRBACRepo) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]string, error) {
	args := m.Called(ctx, roleID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRBACRepo) CreateRole(ctx context.Context, role models.RoleDef) error { return nil }
func (m *MockRBACRepo) AssignRoleToUser(ctx context.Context, assignment models.UserContextRole) error { return nil }

func TestRequirePermission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Allowed Access", func(t *testing.T) {
		mockRepo := new(MockRBACRepo)
		
		// Setup expectations: User has role in Global context that gives permission
		userID := uuid.New()
		roleID := uuid.New()
		mockRepo.On("GetUserRolesInContext", mock.Anything, userID, models.ContextGlobal, uuid.Nil).Return([]models.RoleDef{{ID: roleID}}, nil)
		mockRepo.On("GetRolePermissions", mock.Anything, roleID).Return([]string{"view_dashboard"}, nil)

		authzSvc := services.NewAuthzService(mockRepo)
		mw := middleware.NewRBACMiddleware(authzSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Set("userID", userID.String())

		// Test Global Permission Check
		mw.RequirePermission("view_dashboard", models.ContextGlobal, "")(c)

		assert.Equal(t, 200, w.Code) // Should call Next()
	})

	t.Run("Denied Access", func(t *testing.T) {
		mockRepo := new(MockRBACRepo)
		userID := uuid.New()
		
		// Setup expectations: User has NO roles in Global context
		mockRepo.On("GetUserRolesInContext", mock.Anything, userID, models.ContextGlobal, uuid.Nil).Return([]models.RoleDef{}, nil)

		authzSvc := services.NewAuthzService(mockRepo)
		mw := middleware.NewRBACMiddleware(authzSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Set("userID", userID.String())

		mw.RequirePermission("dangerous_action", models.ContextGlobal, "")(c)

		assert.Equal(t, 403, w.Code)
	})

	t.Run("Context Extraction and Check", func(t *testing.T) {
		mockRepo := new(MockRBACRepo)
		userID := uuid.New()
		courseID := uuid.New()
		roleID := uuid.New()

		// 1. Global check fails
		mockRepo.On("GetUserRolesInContext", mock.Anything, userID, models.ContextGlobal, uuid.Nil).Return([]models.RoleDef{}, nil)
		
		// 2. Course context check succeeds
		mockRepo.On("GetUserRolesInContext", mock.Anything, userID, models.ContextCourse, courseID).Return([]models.RoleDef{{ID: roleID}}, nil)
		mockRepo.On("GetRolePermissions", mock.Anything, roleID).Return([]string{"course.view"}, nil)

		authzSvc := services.NewAuthzService(mockRepo)
		mw := middleware.NewRBACMiddleware(authzSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/courses/"+courseID.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: courseID.String()}} // Simulate URL param
		c.Set("userID", userID.String())

		mw.RequirePermission("course.view", models.ContextCourse, "id")(c)

		assert.Equal(t, 200, w.Code)
	})
}
