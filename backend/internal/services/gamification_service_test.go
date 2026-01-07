package services

import (
	"context"
	"database/sql"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGamificationService_AwardXP(t *testing.T) {
	ctx := context.Background()
	repo := new(MockGamificationRepository)
	svc := NewGamificationService(repo)

	tenantID := "t1"
	userID := "u1"
	amount := 100

	// 1. Record Event
	repo.On("RecordXPEvent", ctx, mock.Anything).Return(nil)
	// 2. Update User XP
	repo.On("UpsertUserXP", ctx, tenantID, userID, amount).Return(nil)
	// 3. Check Level Up
	repo.On("GetUserStats", ctx, userID).Return(&models.UserXP{UserID: userID, TotalXP: 50, Level: 1}, nil)
	// 4. Get Level By XP
	repo.On("GetLevelByXP", ctx, 50).Return(1, nil)

	err := svc.AwardXP(ctx, tenantID, userID, amount, "test", "test", "test")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestGamificationService_AwardXP_LevelUp(t *testing.T) {
	ctx := context.Background()
	repo := new(MockGamificationRepository)
	svc := NewGamificationService(repo)

	tenantID := "t1"
	userID := "u1"
	amount := 500

	repo.On("RecordXPEvent", ctx, mock.Anything).Return(nil)
	repo.On("UpsertUserXP", ctx, tenantID, userID, amount).Return(nil)
	// Current stats (after upsert in real DB, but repo stats in mock here)
	repo.On("GetUserStats", ctx, userID).Return(&models.UserXP{UserID: userID, TotalXP: 600, Level: 1}, nil)
	repo.On("GetLevelByXP", ctx, 600).Return(5, nil)
	repo.On("UpdateUserLevel", ctx, userID, 5).Return(nil)

	err := svc.AwardXP(ctx, tenantID, userID, amount, "test", "test", "test")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestGamificationService_ListBadges(t *testing.T) {
	ctx := context.Background()
	repo := new(MockGamificationRepository)
	svc := NewGamificationService(repo)

	tenantID := "t1"
	repo.On("ListBadges", ctx, tenantID).Return([]models.Badge{{ID: "b1", Name: "Pioneer"}}, nil)

	badges, err := svc.ListBadges(ctx, tenantID)
	assert.NoError(t, err)
	assert.Len(t, badges, 1)
	assert.Equal(t, "Pioneer", badges[0].Name)
}

func TestGamificationService_GetLeaderboard(t *testing.T) {
	ctx := context.Background()
	repo := new(MockGamificationRepository)
	svc := NewGamificationService(repo)

	tenantID := "t1"
	repo.On("GetLeaderboard", ctx, tenantID, 10).Return([]models.LeaderboardEntry{{UserID: "u1", TotalXP: 1000}}, nil)

	board, err := svc.GetLeaderboard(ctx, tenantID, 10)
	assert.NoError(t, err)
	assert.Len(t, board, 1)
	assert.Equal(t, 1000, board[0].TotalXP)
}

func TestGamificationService_GetUserStats(t *testing.T) {
	ctx := context.Background()
	repo := new(MockGamificationRepository)
	svc := NewGamificationService(repo)

	t.Run("Found", func(t *testing.T) {
		expected := &models.UserXP{UserID: "u1", TotalXP: 100}
		repo.On("GetUserStats", ctx, "u1").Return(expected, nil).Once()
		res, err := svc.GetUserStats(ctx, "u1")
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("NotFound", func(t *testing.T) {
		repo.On("GetUserStats", ctx, "u2").Return(nil, sql.ErrNoRows).Once()
		res, err := svc.GetUserStats(ctx, "u2")
		assert.NoError(t, err)
		assert.Equal(t, 0, res.TotalXP)
	})
}

func TestGamificationService_GetUserBadges(t *testing.T) {
	ctx := context.Background()
	repo := new(MockGamificationRepository)
	svc := NewGamificationService(repo)

	repo.On("GetUserBadges", ctx, "u1").Return([]models.UserBadge{{ID: "ub1"}}, nil)
	res, err := svc.GetUserBadges(ctx, "u1")
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestGamificationService_CreateBadge(t *testing.T) {
	ctx := context.Background()
	repo := new(MockGamificationRepository)
	svc := NewGamificationService(repo)

	badge := &models.Badge{Name: "New"}
	repo.On("CreateBadge", ctx, mock.Anything).Return(nil)
	err := svc.CreateBadge(ctx, badge)
	assert.NoError(t, err)
}

func TestGamificationService_AwardBadge(t *testing.T) {
	ctx := context.Background()
	repo := new(MockGamificationRepository)
	svc := NewGamificationService(repo)

	repo.On("AwardBadge", ctx, "u1", "b1").Return(nil)
	err := svc.AwardBadge(ctx, "u1", "b1")
	assert.NoError(t, err)
}
