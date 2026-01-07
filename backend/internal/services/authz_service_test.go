package services_test

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func (m *MockRBACRepo) CreateRole(ctx context.Context, role models.RoleDef) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}
func (m *MockRBACRepo) AssignRoleToUser(ctx context.Context, assignment models.UserContextRole) error {
	args := m.Called(ctx, assignment)
	return args.Error(0)
}

func TestAuthzService_HasPermission(t *testing.T) {
	mockRepo := new(MockRBACRepo)
	svc := services.NewAuthzService(mockRepo)
	userID := uuid.New()
	courseID := uuid.New()
	
	// Case 1: Global Admin overrides everything
	t.Run("Global Admin Permission", func(t *testing.T) {
		globalRole := models.RoleDef{ID: uuid.New(), Name: "GlobalAdmin"}
		mockRepo.On("GetUserRolesInContext", mock.Anything, userID, models.ContextGlobal, uuid.Nil).Return([]models.RoleDef{globalRole}, nil).Once()
		mockRepo.On("GetRolePermissions", mock.Anything, globalRole.ID).Return([]string{"*"}, nil).Once()

		allowed, err := svc.HasPermission(context.Background(), userID, "course.edit", models.ContextCourse, courseID)
		assert.NoError(t, err)
		assert.True(t, allowed)
	})

	// Case 2: Specific Course Permission
	t.Run("Specific Course Permission", func(t *testing.T) {
		// Global check returns no roles
		mockRepo.On("GetUserRolesInContext", mock.Anything, userID, models.ContextGlobal, uuid.Nil).Return([]models.RoleDef{}, nil).Once()
		
		// Course check returns "Instructor"
		instructorRole := models.RoleDef{ID: uuid.New(), Name: "Instructor"}
		mockRepo.On("GetUserRolesInContext", mock.Anything, userID, models.ContextCourse, courseID).Return([]models.RoleDef{instructorRole}, nil).Once()
		mockRepo.On("GetRolePermissions", mock.Anything, instructorRole.ID).Return([]string{"course.edit", "grade.edit"}, nil).Once()

		allowed, err := svc.HasPermission(context.Background(), userID, "grade.edit", models.ContextCourse, courseID)
		assert.NoError(t, err)
		assert.True(t, allowed)
	})

	// Case 3: No Permission
	t.Run("No Permission", func(t *testing.T) {
		mockRepo.On("GetUserRolesInContext", mock.Anything, userID, models.ContextGlobal, uuid.Nil).Return([]models.RoleDef{}, nil).Once()
		mockRepo.On("GetUserRolesInContext", mock.Anything, userID, models.ContextCourse, courseID).Return([]models.RoleDef{}, nil).Once()

		allowed, err := svc.HasPermission(context.Background(), userID, "grade.edit", models.ContextCourse, courseID)
		assert.NoError(t, err)
		assert.False(t, allowed)
	})
}
