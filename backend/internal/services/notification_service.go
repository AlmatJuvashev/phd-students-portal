package services

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type NotificationService struct {
	repo repository.NotificationRepository
}

func NewNotificationService(repo repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) CreateNotification(ctx context.Context, notif *models.Notification) error {
	return s.repo.Create(ctx, notif)
}

func (s *NotificationService) GetUnreadNotifications(ctx context.Context, userID string) ([]models.Notification, error) {
	return s.repo.GetUnread(ctx, userID)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID, userID string) error {
	return s.repo.MarkAsRead(ctx, notificationID, userID)
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}

func (s *NotificationService) ListNotifications(ctx context.Context, userID string, limit int) ([]models.Notification, error) {
	return s.repo.ListByRecipient(ctx, userID, limit)
}
