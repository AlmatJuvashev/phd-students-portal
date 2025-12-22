package services

import (
	"context"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type CalendarService struct {
	repo repository.EventRepository
}

func NewCalendarService(repo repository.EventRepository) *CalendarService {
	return &CalendarService{repo: repo}
}

func (s *CalendarService) CreateEvent(ctx context.Context, event *models.Event, attendees []string) error {
	return s.repo.CreateEvent(ctx, event, attendees)
}

func (s *CalendarService) GetEvents(ctx context.Context, userID, tenantID string, start, end time.Time) ([]models.Event, error) {
	return s.repo.GetEvents(ctx, userID, tenantID, start, end)
}

func (s *CalendarService) GetEvent(ctx context.Context, eventID string) (*models.Event, error) {
	return s.repo.GetEvent(ctx, eventID)
}

func (s *CalendarService) UpdateEvent(ctx context.Context, event *models.Event) error {
	return s.repo.UpdateEvent(ctx, event)
}

func (s *CalendarService) DeleteEvent(ctx context.Context, eventID, userID string) error {
	return s.repo.DeleteEvent(ctx, eventID, userID)
}
