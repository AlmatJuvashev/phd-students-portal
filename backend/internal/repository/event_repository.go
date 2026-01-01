package repository

import (
	"context"
	"errors"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, event *models.Event, attendees []string) error
	GetEvents(ctx context.Context, userID, tenantID string, start, end time.Time) ([]models.Event, error)
	GetEvent(ctx context.Context, eventID string) (*models.Event, error)
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, eventID, userID string) error
}

type SQLEventRepository struct {
	db *sqlx.DB
}

func NewSQLEventRepository(db *sqlx.DB) *SQLEventRepository {
	return &SQLEventRepository{db: db}
}

func (r *SQLEventRepository) CreateEvent(ctx context.Context, event *models.Event, attendees []string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert event
	query := `
		INSERT INTO events (tenant_id, title, description, start_time, end_time, event_type, location, meeting_type, meeting_url, physical_address, color, creator_id, recurrence_type, recurrence_end)
		VALUES (:tenant_id, :title, :description, :start_time, :end_time, :event_type, :location, :meeting_type, :meeting_url, :physical_address, :color, :creator_id, :recurrence_type, :recurrence_end)
		RETURNING id, created_at, updated_at`
	
	rows, err := tx.NamedQuery(query, event)
	if err != nil {
		return err
	}
	if rows.Next() {
		rows.Scan(&event.ID, &event.CreatedAt, &event.UpdatedAt)
	}
	rows.Close()

	// Insert attendees
	if len(attendees) > 0 {
		attendeeQuery := `
			INSERT INTO event_attendees (event_id, tenant_id, user_id, status)
			VALUES ($1, $2, $3, $4)`
		
		for _, userID := range attendees {
			_, err := tx.ExecContext(ctx, attendeeQuery, event.ID, event.TenantID, userID, models.EventAttendeeStatusPending)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *SQLEventRepository) GetEvents(ctx context.Context, userID, tenantID string, start, end time.Time) ([]models.Event, error) {
	query := `
		SELECT e.* FROM events e
		LEFT JOIN event_attendees ea ON e.id = ea.event_id
		WHERE (e.creator_id = $1 OR ea.user_id = $1)
		AND e.tenant_id = $2
		AND e.start_time >= $3 AND e.end_time <= $4
		ORDER BY e.start_time ASC`
	
	var events []models.Event
	err := r.db.SelectContext(ctx, &events, query, userID, tenantID, start, end)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (r *SQLEventRepository) GetEvent(ctx context.Context, eventID string) (*models.Event, error) {
	query := `SELECT * FROM events WHERE id = $1`
	var event models.Event
	err := r.db.GetContext(ctx, &event, query, eventID)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *SQLEventRepository) UpdateEvent(ctx context.Context, event *models.Event) error {
	query := `
		UPDATE events 
		SET title=:title, description=:description, start_time=:start_time, end_time=:end_time, 
			event_type=:event_type, location=:location, meeting_type=:meeting_type, meeting_url=:meeting_url,
			physical_address=:physical_address, color=:color, recurrence_type=:recurrence_type, recurrence_end=:recurrence_end, updated_at=NOW()
		WHERE id=:id AND creator_id=:creator_id`
	
	result, err := r.db.NamedExecContext(ctx, query, event)
	if err != nil {
		return err
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("event not found or unauthorized")
	}
	return nil
}

func (r *SQLEventRepository) DeleteEvent(ctx context.Context, eventID, userID string) error {
	query := `DELETE FROM events WHERE id = $1 AND creator_id = $2`
	result, err := r.db.ExecContext(ctx, query, eventID, userID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("event not found or unauthorized")
	}
	return nil
}
