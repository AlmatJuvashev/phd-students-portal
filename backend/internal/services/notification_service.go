package services

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type NotificationService struct {
	db *sqlx.DB
}

func NewNotificationService(db *sqlx.DB) *NotificationService {
	return &NotificationService{db: db}
}

func (s *NotificationService) CreateNotification(ctx context.Context, notif *models.Notification) error {
	query := `
		INSERT INTO notifications (recipient_id, actor_id, title, message, link, type)
		VALUES (:recipient_id, :actor_id, :title, :message, :link, :type)
		RETURNING id, created_at, updated_at`
	
	rows, err := s.db.NamedQueryContext(ctx, query, notif)
	if err != nil {
		return err
	}
	if rows.Next() {
		rows.Scan(&notif.ID, &notif.CreatedAt, &notif.UpdatedAt)
	}
	rows.Close()
	return nil
}

func (s *NotificationService) GetUnreadNotifications(ctx context.Context, userID string) ([]models.Notification, error) {
	query := `
		SELECT * FROM notifications 
		WHERE recipient_id = $1 AND is_read = FALSE 
		ORDER BY created_at DESC`
	
	var notifs []models.Notification
	err := s.db.SelectContext(ctx, &notifs, query, userID)
	if err != nil {
		return nil, err
	}
	return notifs, nil
}

func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID, userID string) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE id = $1 AND recipient_id = $2`
	_, err := s.db.ExecContext(ctx, query, notificationID, userID)
	return err
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE recipient_id = $1`
	_, err := s.db.ExecContext(ctx, query, userID)
	return err
}
