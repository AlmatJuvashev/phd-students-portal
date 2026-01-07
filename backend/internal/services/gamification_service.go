package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/google/uuid"
)

type GamificationService struct {
	repo repository.GamificationRepository
}

func NewGamificationService(repo repository.GamificationRepository) *GamificationService {
	return &GamificationService{repo: repo}
}

// AwardXP adds XP to a user and handles leveling up
func (s *GamificationService) AwardXP(ctx context.Context, tenantID, userID string, amount int, eventType, sourceType, sourceID string) error {
	return s.repo.WithTransaction(ctx, func(repo repository.GamificationRepository) error {
		// 1. Record Event
		xpEvent := models.XPEvent{
			ID:         uuid.NewString(),
			TenantID:   tenantID,
			UserID:     userID,
			EventType:  eventType,
			XPAmount:   amount,
			SourceType: sourceType,
			SourceID:   sourceID,
			CreatedAt:  time.Now(),
		}
		if err := repo.RecordXPEvent(ctx, xpEvent); err != nil {
			return err
		}

		// 2. Update User XP
		if err := repo.UpsertUserXP(ctx, tenantID, userID, amount); err != nil {
			return err
		}

		// 3. Check Level Up
		currentStats, err := repo.GetUserStats(ctx, userID)
		if err != nil {
			return err
		}

		nextLevel, err := repo.GetLevelByXP(ctx, currentStats.TotalXP)
		if err == nil && nextLevel > currentStats.Level {
			if err := repo.UpdateUserLevel(ctx, userID, nextLevel); err != nil {
				return err
			}
			// Could emit "LevelUp" event here for notifications
		}

		return nil
	})
}

// GetUserStats returns XP stats for a user
func (s *GamificationService) GetUserStats(ctx context.Context, userID string) (*models.UserXP, error) {
	stats, err := s.repo.GetUserStats(ctx, userID)
	if err == sql.ErrNoRows {
		// Return default stats if not found
		return &models.UserXP{UserID: userID, TotalXP: 0, Level: 1}, nil
	}
	return stats, err
}

// ListBadges returns all badges for a tenant
func (s *GamificationService) ListBadges(ctx context.Context, tenantID string) ([]models.Badge, error) {
	return s.repo.ListBadges(ctx, tenantID)
}

// GetUserBadges returns badges earned by a user
func (s *GamificationService) GetUserBadges(ctx context.Context, userID string) ([]models.UserBadge, error) {
	return s.repo.GetUserBadges(ctx, userID)
}

// GetLeaderboard returns top users by XP
func (s *GamificationService) GetLeaderboard(ctx context.Context, tenantID string, limit int) ([]models.LeaderboardEntry, error) {
	if limit > 50 {
		limit = 50
	}
	return s.repo.GetLeaderboard(ctx, tenantID, limit)
}

// Admin: Create Badge
func (s *GamificationService) CreateBadge(ctx context.Context, b *models.Badge) error {
	b.ID = uuid.NewString()
	b.CreatedAt = time.Now()
	return s.repo.CreateBadge(ctx, b)
}

// Helper: Award Badge
func (s *GamificationService) AwardBadge(ctx context.Context, userID, badgeID string) error {
	return s.repo.AwardBadge(ctx, userID, badgeID)
}
