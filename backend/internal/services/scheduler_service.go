package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type SchedulerService struct {
	repo         repository.SchedulerRepository
	resourceRepo repository.ResourceRepository
}

func NewSchedulerService(repo repository.SchedulerRepository, resourceRepo repository.ResourceRepository) *SchedulerService {
	return &SchedulerService{
		repo:         repo,
		resourceRepo: resourceRepo,
	}
}

// --- Terms ---
func (s *SchedulerService) CreateTerm(ctx context.Context, t *models.AcademicTerm) error {
	if t.Code == "" || t.Name == "" {
		return errors.New("code and name are required")
	}
	if t.EndDate.Before(t.StartDate) {
		return errors.New("end date must be after start date")
	}
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return s.repo.CreateTerm(ctx, t)
}

func (s *SchedulerService) ListTerms(ctx context.Context, tenantID string) ([]models.AcademicTerm, error) {
	return s.repo.ListTerms(ctx, tenantID)
}

// --- Offerings & Staff ---
func (s *SchedulerService) CreateOffering(ctx context.Context, o *models.CourseOffering) error {
	if o.CourseID == "" || o.TermID == "" {
		return errors.New("course_id and term_id are required")
	}
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	if o.Status == "" {
		o.Status = "DRAFT"
	}
	return s.repo.CreateOffering(ctx, o)
}

func (s *SchedulerService) AddStaff(ctx context.Context, staff *models.CourseStaff) error {
	staff.CreatedAt = time.Now()
	return s.repo.AddStaff(ctx, staff)
}

// --- Scheduling & Conflict Detection ---

// ConflictError represents a scheduling conflict
type ConflictError struct {
	Reason string
}

func (e *ConflictError) Error() string {
	return e.Reason
}

// ScheduleSession creates a session ONLY if no conflicts exist
func (s *SchedulerService) ScheduleSession(ctx context.Context, session *models.ClassSession) error {
	// Basic Validation
	if session.CourseOfferingID == "" || session.Date.IsZero() {
		return errors.New("offering_id and date are required")
	}
	
	// Conflict Checks
	if err := s.CheckConflicts(ctx, session); err != nil {
		return err
	}

	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()
	return s.repo.CreateSession(ctx, session)
}

// CheckConflicts checks Room and Instructor availability + Capacity
func (s *SchedulerService) CheckConflicts(ctx context.Context, session *models.ClassSession) error {
	// Fetch offering for capacity check
	offering, err := s.repo.GetOffering(ctx, session.CourseOfferingID)
	if err != nil {
		return fmt.Errorf("failed to get offering: %w", err)
	}

	// 1. Room Conflict & Capacity Check
	if session.RoomID != nil {
		// Capacity Check
		room, err := s.resourceRepo.GetRoom(ctx, *session.RoomID)
		if err != nil {
			return fmt.Errorf("failed to get room: %w", err)
		}
		
		// Logic: If Offering needs more seats than room has -> Error
		// Note: We use MaxCapacity as a safe bet, or CurrentEnrolled if stricter
		if offering.MaxCapacity > room.Capacity {
			return &ConflictError{Reason: fmt.Sprintf("Room capacity (%d) is less than offering max capacity (%d)", room.Capacity, offering.MaxCapacity)}
		}

		// Overlap Check
		sessions, err := s.repo.ListSessionsByRoom(ctx, *session.RoomID, session.Date, session.Date)
		if err != nil {
			return err
		}
		if hasOverlap(session, sessions) {
			return &ConflictError{Reason: fmt.Sprintf("Room %s is already booked at this time", *session.RoomID)}
		}
	}

	// 2. Instructor Conflict
	if session.InstructorID != nil {
		sessions, err := s.repo.ListSessionsByInstructor(ctx, *session.InstructorID, session.Date, session.Date)
		if err != nil {
			return err
		}
		if hasOverlap(session, sessions) {
			return &ConflictError{Reason: "Instructor is already teaching another class at this time"}
		}
	}

	return nil
}

// hasOverlap Helper: Assumes sessions are on the same day. 
// Compares Time Strings "HH:MM".
func hasOverlap(target *models.ClassSession, existing []models.ClassSession) bool {
	tStart := parseTime(target.StartTime)
	tEnd := parseTime(target.EndTime)

	for _, e := range existing {
		if e.ID == target.ID { continue } // Skip self if updating
		eStart := parseTime(e.StartTime)
		eEnd := parseTime(e.EndTime)

		// Overlap formula: (StartA < EndB) and (EndA > StartB)
		if tStart < eEnd && tEnd > eStart {
			return true
		}
	}
	return false
}

// Other CRUDs...
func (s *SchedulerService) ListSessions(ctx context.Context, offeringID string, start, end time.Time) ([]models.ClassSession, error) {
	return s.repo.ListSessions(ctx, offeringID, start, end)
}

// parseTime helper converts "14:30" to minutes from midnight for easy comparison
func parseTime(hm string) int {
	var h, m int
	fmt.Sscanf(hm, "%d:%d", &h, &m)
	return h*60 + m
}

