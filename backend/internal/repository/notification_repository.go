package repository

import (
	"context"
	"fmt"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type NotificationRepository interface {
	Create(ctx context.Context, notif *models.Notification) error
	GetUnread(ctx context.Context, userID string) ([]models.Notification, error)
	MarkAsRead(ctx context.Context, notificationID, userID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
}

type SQLNotificationRepository struct {
	db *sqlx.DB
}

func NewSQLNotificationRepository(db *sqlx.DB) *SQLNotificationRepository {
	return &SQLNotificationRepository{db: db}
}

func (r *SQLNotificationRepository) Create(ctx context.Context, notif *models.Notification) error {
	query := `
		INSERT INTO notifications (recipient_id, actor_id, title, message, link, type, tenant_id)
		VALUES (:recipient_id, :actor_id, :title, :message, :link, :type, :tenant_id)
		RETURNING id, created_at, updated_at`

	rows, err := r.db.NamedQueryContext(ctx, query, notif)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&notif.ID, &notif.CreatedAt, &notif.UpdatedAt)
	}
	return nil
}

func (r *SQLNotificationRepository) GetUnread(ctx context.Context, userID string) ([]models.Notification, error) {
	query := `
		SELECT * FROM notifications 
		WHERE recipient_id = $1 AND is_read = FALSE 
		ORDER BY created_at DESC`

	var notifs []models.Notification
	err := r.db.SelectContext(ctx, &notifs, query, userID)
	if err != nil {
		return nil, err
	}
	return notifs, nil
}

func (r *SQLNotificationRepository) MarkAsRead(ctx context.Context, notificationID, userID string) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE id = $1 AND recipient_id = $2`
	res, err := r.db.ExecContext(ctx, query, notificationID, userID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("notification not found or not owned by user")
	}
	return nil
}

func (r *SQLNotificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE recipient_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
