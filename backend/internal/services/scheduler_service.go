package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/scheduler/solver"
)

type SchedulerService struct {
	repo           repository.SchedulerRepository
	resourceRepo   repository.ResourceRepository
	curriculumRepo repository.CurriculumRepository
}

func NewSchedulerService(repo repository.SchedulerRepository, resourceRepo repository.ResourceRepository, curriculumRepo repository.CurriculumRepository) *SchedulerService {
	return &SchedulerService{
		repo:           repo,
		resourceRepo:   resourceRepo,
		curriculumRepo: curriculumRepo,
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

// ScheduleSession creates a session ONLY if no conflicts exist (or only soft warnings)
func (s *SchedulerService) ScheduleSession(ctx context.Context, session *models.ClassSession) ([]string, error) {
	// Basic Validation
	if session.CourseOfferingID == "" || session.Date.IsZero() {
		return nil, errors.New("offering_id and date are required")
	}
	
	// Conflict Checks
	warnings, err := s.CheckConflicts(ctx, session)
	if err != nil {
		return nil, err
	}

	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()
	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}
	
	return warnings, nil
}

// CheckConflicts checks Room and Instructor availability + Capacity + Dept Constraints
// Returns warnings ([]string) and error (if Critical Conflict)
func (s *SchedulerService) CheckConflicts(ctx context.Context, session *models.ClassSession) ([]string, error) {
	var warnings []string
	cfg := solver.DefaultConfig() // Use defaults for Manual Entry

	// Fetch offering for capacity check
	offering, err := s.repo.GetOffering(ctx, session.CourseOfferingID)
	if err != nil {
		return nil, fmt.Errorf("failed to get offering: %w", err)
	}

	// 1. Room Conflict & Capacity Check
	if session.RoomID != nil {
		room, err := s.resourceRepo.GetRoom(ctx, *session.RoomID)
		if err != nil {
			return nil, fmt.Errorf("failed to get room: %w", err)
		}
		
		// A. Capacity Check
		if offering.MaxCapacity > room.Capacity {
			msg := fmt.Sprintf("Room capacity (%d) is less than offering max capacity (%d)", room.Capacity, offering.MaxCapacity)
			if cfg.CapacityConstraint == "HARD" {
				return nil, &ConflictError{Reason: msg}
			} else if cfg.CapacityConstraint == "SOFT" {
				warnings = append(warnings, "Warning: "+msg)
			}
		}

		// B. Overlap Check
		sessions, err := s.repo.ListSessionsByRoom(ctx, *session.RoomID, session.Date, session.Date)
		if err != nil {
			return nil, err
		}
		if hasOverlap(session, sessions) {
			msg := fmt.Sprintf("Room %s is already booked at this time", *session.RoomID)
			// Overlaps are arguably always Hard for physical rooms, but let's respect config
			if cfg.TimeConflictConstraint == "HARD" {
				return nil, &ConflictError{Reason: msg}
			} else {
				warnings = append(warnings, "Warning: "+msg)
			}
		}

		// C. Department Check
		// Try to fetch course metadata
		course, err := s.curriculumRepo.GetCourse(ctx, offering.CourseID)
		if err == nil && course != nil && course.DepartmentID != nil {
			roomDept := ""
			if room.DepartmentID != nil {
				roomDept = *room.DepartmentID
			}
			if *course.DepartmentID != roomDept {
				msg := fmt.Sprintf("Department Mismatch: Course is '%s', Room is '%s'", *course.DepartmentID, roomDept)
				if cfg.DepartmentConstraint == "HARD" {
					return nil, &ConflictError{Reason: msg}
				} else if cfg.DepartmentConstraint == "SOFT" {
					warnings = append(warnings, "Warning: "+msg)
				}
			}
		}
	}

	// 2. Instructor Conflict
	if session.InstructorID != nil {
		sessions, err := s.repo.ListSessionsByInstructor(ctx, *session.InstructorID, session.Date, session.Date)
		if err != nil {
			return nil, err
		}
		if hasOverlap(session, sessions) {
			return nil, &ConflictError{Reason: "Instructor is already teaching another class at this time"}
		}
	}

	return warnings, nil
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


// AutoSchedule runs the optimizer for a given Term
func (s *SchedulerService) AutoSchedule(ctx context.Context, tenantID, termID string, config *solver.SolverConfig) (*solver.Solution, error) {
	// 1. Fetch Data
	// A. Sessions
	sessions, err := s.repo.ListSessionsForTerm(ctx, termID)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}
	if len(sessions) == 0 {
		return nil, errors.New("no sessions found for this term")
	}

	// B. Rooms (Assume ResourceRepo has ListRooms)
	rooms, err := s.resourceRepo.ListRooms(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list rooms: %w", err)
	}

	// C. Offerings (for Capacity info)
	offerings, err := s.repo.ListOfferings(ctx, tenantID, termID)
	if err != nil {
		return nil, fmt.Errorf("failed to list offerings: %w", err)
	}
	offeringMap := make(map[string]models.CourseOffering)
	courseIDs := make([]string, 0, len(offerings))
	for _, o := range offerings {
		offeringMap[o.ID] = o
		courseIDs = append(courseIDs, o.CourseID)
	}

	// D. Courses (for Department Info)
	// We fetch all courses for the tenant to map CourseID -> DepartmentID
	allCourses, err := s.curriculumRepo.ListCourses(ctx, tenantID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list courses: %w", err)
	}
	courseDeptMap := make(map[string]string)
	for _, c := range allCourses {
		if c.DepartmentID != nil {
			courseDeptMap[c.ID] = *c.DepartmentID
		}
	}

	// 2. Build Problem Instance
	instance := solver.ProblemInstance{
		Sessions:     make(map[string]solver.SessionData),
		Rooms:        make(map[string]models.Room),
		Instructors:  make(map[string]models.User),
		Dependencies: make(map[string][]string),
	}

	for _, room := range rooms {
		instance.Rooms[room.ID] = room
	}

	for _, sess := range sessions {
		offering, ok := offeringMap[sess.CourseOfferingID]
		if !ok { continue }

		// Parse duration
		sTime := parseTime(sess.StartTime)
		eTime := parseTime(sess.EndTime)
		duration := eTime - sTime

		// Instructor
		instrID := ""
		if sess.InstructorID != nil {
			instrID = *sess.InstructorID
		}

		// Department
		deptID := courseDeptMap[offering.CourseID]

		instance.Sessions[sess.ID] = solver.SessionData{
			ID:           sess.ID,
			DurationMins: duration,
			MaxStudents:  offering.MaxCapacity,
			InstructorID: instrID,
			DepartmentID: deptID,
			FixedRoomID:  "",
			OriginalTime: sess.Date,
		}
	}

	// 3. Run Solver
	cfg := solver.DefaultConfig()
	if config != nil {
		cfg = *config
	}
	// Limit iterations for synchronous HTTP handling safety, unless overridden to higher
	if config == nil || config.MaxIterations == 0 {
		cfg.MaxIterations = 5000 
	}
	
	slv := solver.NewSchedulerSolver(cfg)
	solution, err := slv.Solve(ctx, instance)
	if err != nil {
		return nil, err
	}

	return solution, nil
}
