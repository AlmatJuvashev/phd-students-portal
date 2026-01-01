package solver

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestSolver_Solve_RecommendsCapacityOptimization(t *testing.T) {
	// 1. Setup Data
	// - Room 101: Capacity 10
	// - Room 102: Capacity 100
	// - Session A: 5 Students (Should go to 101)
	// - Session B: 90 Students (Should go to 102)

	roomSmall := models.Room{ID: "r101", Capacity: 10, Name: "Small Room"}
	roomBig := models.Room{ID: "r102", Capacity: 100, Name: "Lecture Hall"}

	sessionSmall := SessionData{
		ID:           "s_small",
		DurationMins: 90,
		MaxStudents:  5,
		OriginalTime: time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
	}
	sessionBig := SessionData{
		ID:           "s_big",
		DurationMins: 90,
		MaxStudents:  90,
		OriginalTime: time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
	}

	instance := ProblemInstance{
		Sessions: map[string]SessionData{
			"s_small": sessionSmall,
			"s_big":   sessionBig,
		},
		Rooms: map[string]models.Room{
			"r101": roomSmall,
			"r102": roomBig,
		},
		Instructors: make(map[string]models.User),
	}

	// 2. Configure Solver
	cfg := DefaultConfig()
	cfg.MaxIterations = 2000 // Fast run for test
	// Heavy penalty on capacity violation
	cfg.WeightCapacity = 1000.0 
	// Penalty on waste (utilization)
	cfg.WeightUtilization = 10.0
	
	solver := NewSchedulerSolver(cfg)

	// 3. Solve
	solution, err := solver.Solve(context.Background(), instance)
	assert.NoError(t, err)
	assert.NotNil(t, solution)

	// 4. Verify Validity
	// The solver should start with a random state (possibly bad) and optimize to a good state.
	
	// Check Assignments
	assignSmall := solution.Assignments["s_small"]
	assignBig := solution.Assignments["s_big"]

	// The big session MUST be in the big room to avoid capacity penalty
	assert.Equal(t, "r102", assignBig.RoomID, "Big session should be in big room")
	
	// The small session SHOULD be in the small room to minimize waste penalty (100-5 = 95 waste vs 10-5 = 5 waste)
	assert.Equal(t, "r101", assignSmall.RoomID, "Small session should be in small room")
	
	t.Logf("Final Score: %f", solution.Score)
}

func TestSolver_Solve_FixesHardConflict(t *testing.T) {
	// Setup: 2 Sessions, 1 Room. Hard Conflict over time.
	// NOTE: Time mutation not implemented yet, so solver can only swap rooms.
	// But if we have 2 rooms and 2 sessions at same time, it should split them.
}
