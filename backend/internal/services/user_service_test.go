package services_test

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, filter repository.UserFilter, pagination repository.Pagination) ([]models.User, int, error) {
	args := m.Called(ctx, filter, pagination)
	return args.Get(0).([]models.User), args.Int(1), args.Error(2)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, id string, hash string) error {
	args := m.Called(ctx, id, hash)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateAvatar(ctx context.Context, id string, avatarURL string) error {
	args := m.Called(ctx, id, avatarURL)
	return args.Error(0)
}

func (m *MockUserRepository) SetActive(ctx context.Context, id string, active bool) error {
	args := m.Called(ctx, id, active)
	return args.Error(0)
}

func (m *MockUserRepository) Exists(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) EmailExists(ctx context.Context, email string, excludeUserID string) (bool, error) {
	args := m.Called(ctx, email, excludeUserID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) LinkAdvisor(ctx context.Context, studentID, advisorID string) error {
	args := m.Called(ctx, studentID, advisorID)
	return args.Error(0)
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := services.NewUserService(mockRepo, nil) // Mock Redis as nil

	ctx := context.Background()
	req := services.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Role:      "student",
	}

	// Expect username generation check
	// Note: UserService generates username. It calls GetByUsername to check existence.
	// We'll simulate it doesn't exist on first try.
	mockRepo.On("GetByUsername", ctx, mock.AnythingOfType("string")).Return((*models.User)(nil), nil).Once()

	// Expect Create call
	mockRepo.On("Create", ctx, mock.MatchedBy(func(u *models.User) bool {
		return u.FirstName == "John" && 
		       u.LastName == "Doe" && 
			   u.Email == "john@example.com" &&
			   len(u.PasswordHash) > 0 // Password should be hashed
	})).Return("new-uuid", nil)

	user, tempPass, err := svc.CreateUser(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "new-uuid", user.ID)
	assert.NotEmpty(t, user.Username)
	assert.NotEmpty(t, tempPass)

	mockRepo.AssertExpectations(t)
}

func TestUserService_ListUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := services.NewUserService(mockRepo, nil)

	ctx := context.Background()
	filter := repository.UserFilter{Role: "student"}
	pagination := repository.Pagination{Limit: 10, Offset: 0}

	expectedUsers := []models.User{
		{ID: "1", FirstName: "Alice"},
		{ID: "2", FirstName: "Bob"},
	}

	mockRepo.On("List", ctx, filter, pagination).Return(expectedUsers, 2, nil)

	users, total, err := svc.ListUsers(ctx, filter, pagination)
	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, users, 2)
	assert.Equal(t, "Alice", users[0].FirstName)
}
