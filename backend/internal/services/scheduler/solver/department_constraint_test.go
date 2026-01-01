package solver

import (
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestSolver_DepartmentConstraints(t *testing.T) {
	// Setup
	anatomyDept := "dept-anatomy"
	mathDept := "dept-math"

	// Rooms
	roomAnatomy := models.Room{ID: "r-anatomy", Capacity: 50, DepartmentID: &anatomyDept}
	roomMath := models.Room{ID: "r-math", Capacity: 50, DepartmentID: &mathDept}
	roomUniversal := models.Room{ID: "r-univ", Capacity: 50, DepartmentID: nil}

	// Sessions
	sessAnatomy := SessionData{ID: "s-anatomy", MaxStudents: 30, DepartmentID: anatomyDept, DurationMins: 60}
	sessMath := SessionData{ID: "s-math", MaxStudents: 30, DepartmentID: mathDept, DurationMins: 60}
	sessGeneral := SessionData{ID: "s-general", MaxStudents: 30, DepartmentID: "", DurationMins: 60}

	instance := ProblemInstance{
		Sessions: map[string]SessionData{
			"s-anatomy": sessAnatomy,
			"s-math":    sessMath,
			"s-general": sessGeneral,
		},
		Rooms: map[string]models.Room{
			"r-anatomy": roomAnatomy,
			"r-math":    roomMath,
			"r-univ":    roomUniversal,
		},
		Instructors: make(map[string]models.User),
	}

	solver := NewSchedulerSolver(DefaultConfig())

	// Test Case 1: Anatomy Session in Anatomy Room -> Valid
	sol1 := &Solution{
		Assignments: map[string]Assignment{
			"s-anatomy": {SessionID: "s-anatomy", RoomID: "r-anatomy", StartTime: time.Now(), EndTime: time.Now().Add(time.Hour)},
		},
	}
	cost1 := solver.CalculateCost(sol1, instance)
	// Expect 0 hard penalties (assuming capacity ok)
	// Note: CalculateCost returns weighted sum. Hard penalty weight is very high (say 1000).
	// If cost is small (soft constraints only), then valid.
	assert.Less(t, cost1, 100.0, "Anatomy in Anatomy Room should be low cost")

	// Test Case 2: Anatomy Session in Math Room -> Invalid
	sol2 := &Solution{
		Assignments: map[string]Assignment{
			"s-anatomy": {SessionID: "s-anatomy", RoomID: "r-math", StartTime: time.Now(), EndTime: time.Now().Add(time.Hour)},
		},
	}
	cost2 := solver.CalculateCost(sol2, instance)
	assert.Greater(t, cost2, 100.0, "Anatomy in Math Room should be high cost")

	// Test Case 3: Anatomy Session in Universal Room -> Invalid (Strict Rule)
	sol3 := &Solution{
		Assignments: map[string]Assignment{
			"s-anatomy": {SessionID: "s-anatomy", RoomID: "r-univ", StartTime: time.Now(), EndTime: time.Now().Add(time.Hour)},
		},
	}
	cost3 := solver.CalculateCost(sol3, instance)
	assert.Greater(t, cost3, 100.0, "Anatomy in Universal Room should be high cost (Strict)")

	// Test Case 4: General Session in Universal Room -> Valid
	sol4 := &Solution{
		Assignments: map[string]Assignment{
			"s-general": {SessionID: "s-general", RoomID: "r-univ", StartTime: time.Now(), EndTime: time.Now().Add(time.Hour)},
		},
	}
	cost4 := solver.CalculateCost(sol4, instance)
	assert.Less(t, cost4, 100.0, "General in Universal Room should be low cost")

	// Test Case 5: General Session in Anatomy Room -> Invalid
	sol5 := &Solution{
		Assignments: map[string]Assignment{
			"s-general": {SessionID: "s-general", RoomID: "r-anatomy", StartTime: time.Now(), EndTime: time.Now().Add(time.Hour)},
		},
	}
	cost5 := solver.CalculateCost(sol5, instance)
	assert.Greater(t, cost5, 100.0, "General in Anatomy Room should be high cost")

}
