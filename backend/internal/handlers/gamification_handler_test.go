package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGamificationRepo struct {
	mock.Mock
}

func (m *mockGamificationRepo) RecordXPEvent(ctx context.Context, event models.XPEvent) error {
	return m.Called(ctx, event).Error(0)
}
func (m *mockGamificationRepo) UpsertUserXP(ctx context.Context, tenantID, userID string, amount int) error {
	return m.Called(ctx, tenantID, userID, amount).Error(0)
}
func (m *mockGamificationRepo) GetUserStats(ctx context.Context, userID string) (*models.UserXP, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserXP), args.Error(1)
}
func (m *mockGamificationRepo) UpdateUserLevel(ctx context.Context, userID string, level int) error {
	return m.Called(ctx, userID, level).Error(0)
}
func (m *mockGamificationRepo) GetLevelByXP(ctx context.Context, totalXP int) (int, error) {
	args := m.Called(ctx, totalXP)
	return args.Int(0), args.Error(1)
}
func (m *mockGamificationRepo) ListBadges(ctx context.Context, tenantID string) ([]models.Badge, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Badge), args.Error(1)
}
func (m *mockGamificationRepo) GetUserBadges(ctx context.Context, userID string) ([]models.UserBadge, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserBadge), args.Error(1)
}
func (m *mockGamificationRepo) GetLeaderboard(ctx context.Context, tenantID string, limit int) ([]models.LeaderboardEntry, error) {
	args := m.Called(ctx, tenantID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.LeaderboardEntry), args.Error(1)
}
func (m *mockGamificationRepo) CreateBadge(ctx context.Context, b *models.Badge) error {
	return m.Called(ctx, b).Error(0)
}
func (m *mockGamificationRepo) AwardBadge(ctx context.Context, userID, badgeID string) error {
	return m.Called(ctx, userID, badgeID).Error(0)
}
func (m *mockGamificationRepo) WithTransaction(ctx context.Context, fn func(repo repository.GamificationRepository) error) error {
	return fn(m)
}

func TestGamificationHandler_GetMyStats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockGamificationRepo)
	svc := services.NewGamificationService(repo)
	h := NewGamificationHandler(svc)

	repo.On("GetUserStats", mock.Anything, "u1").Return(&models.UserXP{UserID: "u1", TotalXP: 100}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/stats", nil)
	c.Set("userID", "u1")

	h.GetMyStats(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGamificationHandler_GetLeaderboard(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockGamificationRepo)
	svc := services.NewGamificationService(repo)
	h := NewGamificationHandler(svc)

	repo.On("GetLeaderboard", mock.Anything, "t1", 10).Return([]models.LeaderboardEntry{}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/leaderboard?limit=10", nil)
	c.Set("tenant_id", "t1")

	h.GetLeaderboard(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGamificationHandler_CreateBadge(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockGamificationRepo)
	svc := services.NewGamificationService(repo)
	h := NewGamificationHandler(svc)

	badge := models.Badge{Name: "New Badge", Code: "NB1"}
	repo.On("CreateBadge", mock.Anything, mock.Anything).Return(nil)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(badge)
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/badges", bytes.NewBuffer(body))
	c.Set("tenant_id", "t1")

	h.CreateBadge(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGamificationHandler_GetMyBadges(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockGamificationRepo)
	svc := services.NewGamificationService(repo)
	h := NewGamificationHandler(svc)

	repo.On("GetUserBadges", mock.Anything, "u1").Return([]models.UserBadge{{ID: "ub1"}}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/badges/mine", nil)
	c.Set("userID", "u1")

	h.GetMyBadges(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGamificationHandler_ListAllBadges(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockGamificationRepo)
	svc := services.NewGamificationService(repo)
	h := NewGamificationHandler(svc)

	repo.On("ListBadges", mock.Anything, "t1").Return([]models.Badge{}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/badges", nil)
	c.Set("tenant_id", "t1")

	h.ListAllBadges(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
