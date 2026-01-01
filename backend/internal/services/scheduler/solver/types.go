package solver

import (
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
)

// TimeSlot represents a 15-minute chunk offset from the start of the week.
// E.g., Monday 08:00 might be slot 32 (if day starts at 00:00).
// For simplicity, we can assume a day-based slot system or just use time.Time.
// Using integer slots is 100x faster for array indexing.
type TimeSlot int

// ProblemInstance holds all immutable data for the scheduling problem.
// Optimized for O(1) lookups during the solver loop.
type ProblemInstance struct {
	Sessions     map[string]SessionData // SessionID -> Flattened Data
	Rooms        map[string]models.Room // RoomID -> Room Model
	Instructors  map[string]models.User // InstructorID -> User Model (optional)
	Dependencies map[string][]string    // SessionID -> List of prerequisite SessionIDs
}

// SessionData is a flattened representation of a session for the solver.
// It combines ClassSession and CourseOffering data to avoid joins during solving.
type SessionData struct {
	ID           string
	DurationMins int    // Duration in minutes
	MaxStudents  int    // Capacity required (from CourseOffering)
	InstructorID string // Preferred/Required Instructor
	DepartmentID string // Department of the Course (optional)
	FixedRoomID  string // If strictly pre-assigned (Hard Constraint)
	OriginalTime time.Time
}

// Assignment represents a decision for a single session.
type Assignment struct {
	SessionID string
	RoomID    string
	StartTime time.Time
	EndTime   time.Time
}

// Solution represents a complete schedule state.
type Solution struct {
	Assignments map[string]Assignment // SessionID -> Assignment
	Score       float64              // Cached score (Energy)
	IsValid     bool                 // True if all Hard Constraints are met
}

// SolverConfig controls the Simulated Annealing parameters.
type SolverConfig struct {
	MaxIterations int     `json:"max_iterations"`
	InitialTemp   float64 `json:"initial_temp"`
	CoolingRate   float64 `json:"cooling_rate"`
	
	// Weights for Cost Function
	WeightHardConflict float64 `json:"weight_hard_conflict"` // Overlaps
	WeightCapacity     float64 `json:"weight_capacity"`      // Room too small
	WeightUtilization  float64 `json:"weight_utilization"`   // Waste
	WeightLocality     float64 `json:"weight_locality"`      // Different buildings
	WeightInstructor   float64 `json:"weight_instructor"`    // Instructor constraints

	// Constraint Severities: "HARD", "SOFT", "OFF"
	DepartmentConstraint  string `json:"constraint_department"`   // Dept matching
	CapacityConstraint    string `json:"constraint_capacity"`     // Room size
	TimeConflictConstraint string `json:"constraint_time_conflict"`// Overlaps (usually HARD)
}

// DefaultConfig returns standard parameters for the solver.
func DefaultConfig() SolverConfig {
	return SolverConfig{
		MaxIterations:      10000,
		InitialTemp:        1000.0,
		CoolingRate:        0.995,
		
		WeightHardConflict: 1000.0, // Base weight for Hard
		WeightCapacity:     500.0,  // Base weight for Capacity (if Hard)
		WeightUtilization:  10.0,   // Soft
		WeightLocality:     20.0,
		WeightInstructor:   50.0,

		DepartmentConstraint:   "HARD",
		CapacityConstraint:     "HARD",
		TimeConflictConstraint: "HARD",
	}
}
