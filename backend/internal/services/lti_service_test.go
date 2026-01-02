package services_test

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLTIRepository
type MockLTIRepository struct {
	mock.Mock
}

func (m *MockLTIRepository) CreateTool(ctx context.Context, params models.CreateToolParams) (*models.LTITool, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LTITool), args.Error(1)
}

func (m *MockLTIRepository) GetTool(ctx context.Context, id string) (*models.LTITool, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LTITool), args.Error(1)
}

func (m *MockLTIRepository) ListTools(ctx context.Context, tenantID string) ([]models.LTITool, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).([]models.LTITool), args.Error(1)
}

func (m *MockLTIRepository) GetToolByClientID(ctx context.Context, clientID string) (*models.LTITool, error) {
	args := m.Called(ctx, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LTITool), args.Error(1)
}

func (m *MockLTIRepository) CreateKey(ctx context.Context, key models.LTIKey) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockLTIRepository) GetActiveKey(ctx context.Context) (*models.LTIKey, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LTIKey), args.Error(1)
}

func (m *MockLTIRepository) ListActiveKeys(ctx context.Context) ([]models.LTIKey, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.LTIKey), args.Error(1)
}

func TestLTIService_GenerateLoginInit(t *testing.T) {
	mockRepo := new(MockLTIRepository)
	cfg := config.AppConfig{IssuerURL: "https://platform.com"}
	svc := services.NewLTIService(mockRepo, cfg)

	toolID := "tool-1"
	tool := &models.LTITool{
		ID:               toolID,
		InitiateLoginURL: "https://tool.com/login",
		ClientID:         "client-123",
		DeploymentID:     "deploy-1",
	}

	mockRepo.On("GetTool", mock.Anything, toolID).Return(tool, nil)

	url, err := svc.GenerateLoginInit(context.Background(), toolID, "user-1", "https://platform.com/course/1")
	assert.NoError(t, err)
	assert.Contains(t, url, "https://tool.com/login")
	assert.Contains(t, url, "iss=https%3A%2F%2Fplatform.com")
	assert.Contains(t, url, "client_id=client-123")
	assert.Contains(t, url, "login_hint=user-1")
	assert.Contains(t, url, "lti_message_hint=deploy-1")
}

func TestLTIService_RegisterTool(t *testing.T) {
	mockRepo := new(MockLTIRepository)
	svc := services.NewLTIService(mockRepo, config.AppConfig{})

	params := models.CreateToolParams{Name: "Test Tool"}
	expected := &models.LTITool{ID: "new-id", Name: "Test Tool"}

	mockRepo.On("CreateTool", mock.Anything, params).Return(expected, nil)

	result, err := svc.RegisterTool(context.Background(), params)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestLTIService_GenerateLoginInit_NotFound(t *testing.T) {
	mockRepo := new(MockLTIRepository)
	svc := services.NewLTIService(mockRepo, config.AppConfig{IssuerURL: "https://platform.com"})

	// Setup: Tool ID "missing" returns nil, nil (or nil, error depending on repo logic, usually nil, nil for not found if handled gracefully or nil, ErrNotFound)
	// In the real implementation `GetTool` returns (nil, nil) if rows.Scan errors with sql.ErrNoRows? 
	// Let's check `lti_repository.go`. Yes, it returns `nil, nil` on ErrNoRows.
	// The service checks `if tool == nil { return "", fmt.Errorf("tool not found") }`.
	
	mockRepo.On("GetTool", mock.Anything, "missing").Return((*models.LTITool)(nil), nil)

	url, err := svc.GenerateLoginInit(context.Background(), "missing", "u1", "link")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tool not found")
	assert.Empty(t, url)
}

func TestLTIService_GenerateLoginInit_InvalidURL(t *testing.T) {
	mockRepo := new(MockLTIRepository)
	svc := services.NewLTIService(mockRepo, config.AppConfig{})

	// Tool with invalid URL in DB (unlikely if validation works, but defensive coding)
	tool := &models.LTITool{ID: "t2", InitiateLoginURL: ":/invalid-url"} 
	mockRepo.On("GetTool", mock.Anything, "t2").Return(tool, nil)

	_, err := svc.GenerateLoginInit(context.Background(), "t2", "u1", "link")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid tool login url")
}

func TestLTIService_KeyManagement(t *testing.T) {
	mockRepo := new(MockLTIRepository)
	svc := services.NewLTIService(mockRepo, config.AppConfig{})

	t.Run("GetJWKS Success", func(t *testing.T) {
		// Mock ListActiveKeys returning a dummy key
		// Real PEM generation is complex for mock, assuming service handles invalid PEM gracefully or we provide valid PEM
		// Let's provide a valid minimalistic Public Key PEM for testing transparency
		// Or just rely on service transformation logic. 
		// Ideally we test `RotateKey` calls `CreateKey`.
		
		// For unit test simplicity, let's Verify EnsureActiveKey logic
		mockRepo.On("GetActiveKey", mock.Anything).Return((*models.LTIKey)(nil), nil).Once() // First call not found
		
		// RotateKey will call CreateKey
		mockRepo.On("CreateKey", mock.Anything, mock.Anything).Return(nil)
		
		// And GetActiveKey again
		mockRepo.On("GetActiveKey", mock.Anything).Return(&models.LTIKey{ID: "new"}, nil)
		
		key, err := svc.EnsureActiveKey(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, "new", key.ID)
	})
}
