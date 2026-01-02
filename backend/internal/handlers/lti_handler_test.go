package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/handlers"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Re-using MockLTIRepository pattern as defined in services_test, but defined here to avoid import cycles or shared package complexity for now
type HMockLTIRepo struct {
	mock.Mock
}

func (m *HMockLTIRepo) CreateTool(ctx context.Context, params models.CreateToolParams) (*models.LTITool, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LTITool), args.Error(1)
}

func (m *HMockLTIRepo) GetTool(ctx context.Context, id string) (*models.LTITool, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LTITool), args.Error(1)
}

func (m *HMockLTIRepo) ListTools(ctx context.Context, tenantID string) ([]models.LTITool, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).([]models.LTITool), args.Error(1)
}

func (m *HMockLTIRepo) GetToolByClientID(ctx context.Context, clientID string) (*models.LTITool, error) {
	args := m.Called(ctx, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LTITool), args.Error(1)
}

func (m *HMockLTIRepo) CreateKey(ctx context.Context, key models.LTIKey) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *HMockLTIRepo) GetActiveKey(ctx context.Context) (*models.LTIKey, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LTIKey), args.Error(1)
}

func (m *HMockLTIRepo) ListActiveKeys(ctx context.Context) ([]models.LTIKey, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.LTIKey), args.Error(1)
}

func TestLTIHandler_RegisterTool(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(HMockLTIRepo)
	cfg := config.AppConfig{}
	svc := services.NewLTIService(mockRepo, cfg)
	h := handlers.NewLTIHandler(svc, cfg)

	t.Run("Success", func(t *testing.T) {
		reqBody := models.CreateToolParams{
			TenantID: "t1", Name: "Zoom", ClientID: "c1", InitiateLoginURL: "https://z.com/init", RedirectionURIs: []string{"https://z.com/launch"},
		}
		expected := &models.LTITool{ID: "new", Name: "Zoom"}
		mockRepo.On("CreateTool", mock.Anything, reqBody).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		jb, _ := json.Marshal(reqBody)
		c.Request = httptest.NewRequest("POST", "/lti/tools", bytes.NewBuffer(jb))

		h.RegisterTool(c)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestLTIHandler_LoginInit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(HMockLTIRepo)
	cfg := config.AppConfig{IssuerURL: "https://my-lms.com"}
	svc := services.NewLTIService(mockRepo, cfg)
	h := handlers.NewLTIHandler(svc, cfg)

	t.Run("Redirect Success", func(t *testing.T) {
		tool := &models.LTITool{ID: "t1", InitiateLoginURL: "https://ext.com/login", ClientID: "cid", DeploymentID: "did"}
		mockRepo.On("GetTool", mock.Anything, "t1").Return(tool, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/lti/login_init?tool_id=t1&target_link_uri=https://my-lms.com/res", nil)
		c.Set("userID", "u1")


		h.LoginInit(c)
		assert.Equal(t, http.StatusFound, w.Code)
		loc := w.Header().Get("Location")
		assert.Contains(t, loc, "https://ext.com/login")
		assert.Contains(t, loc, "login_hint=u1")
	})

	t.Run("Missing Params", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// Missing tool_id
		c.Request = httptest.NewRequest("GET", "/lti/login_init?target_link_uri=foo", nil)

		h.LoginInit(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		toolID := "t-missing"
		// Mock repo to return nil for this tool -> Service returns "tool not found" error
		mockRepo.On("GetTool", mock.Anything, toolID).Return((*models.LTITool)(nil), nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/lti/login_init?tool_id="+toolID+"&target_link_uri=foo", nil)
		// Ensure userID is set if middleware would set it, though optional for this specific error path check
		c.Set("userID", "u1")

		h.LoginInit(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestLTIHandler_GetJWKS(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(HMockLTIRepo)
	cfg := config.AppConfig{}
	svc := services.NewLTIService(mockRepo, cfg)
	h := handlers.NewLTIHandler(svc, cfg)

	mockRepo.On("ListActiveKeys", mock.Anything).Return([]models.LTIKey{}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/.well-known/jwks.json", nil)

	h.GetJWKS(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "keys")
}
