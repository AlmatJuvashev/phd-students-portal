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
	// Setup: 2 Sessions, 2 Rooms. Overlap in time.
	// Both sessions initially assigned to Room 101.
	// Solver should move one to Room 102 to minimize conflict penalty (assuming hard constraint weight).

	// Sessions at same time
	start := time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC)
	s1 := SessionData{ID: "s1", DurationMins: 60, MaxStudents: 10, OriginalTime: start}
	s2 := SessionData{ID: "s2", DurationMins: 60, MaxStudents: 10, OriginalTime: start}

	// Rooms
	r1 := models.Room{ID: "r1", Capacity: 20}
	r2 := models.Room{ID: "r2", Capacity: 20}

	instance := ProblemInstance{
		Sessions: map[string]SessionData{"s1": s1, "s2": s2},
		Rooms:    map[string]models.Room{"r1": r1, "r2": r2},
	}

	cfg := DefaultConfig()
	cfg.MaxIterations = 1000
	// Set high penalty for conflicts
	cfg.WeightHardConflict = 10000.0

	solver := NewSchedulerSolver(cfg)
	// Seed/Initial solution might put them in different rooms by chance (random greedy).
	// But we want to ensure *optimization* works.
	// Since GenerateInitialSolution is random, we can just run Solve.
	// If it fails to separate them, score will be high.

	sol, err := solver.Solve(context.Background(), instance)
	assert.NoError(t, err)

	a1 := sol.Assignments["s1"]
	a2 := sol.Assignments["s2"]

	// They must be in different rooms since time matches
	assert.NotEqual(t, a1.RoomID, a2.RoomID, "Sessions at same time must be in different rooms")
}

func TestSolver_Mutate(t *testing.T) {
	s1 := SessionData{ID: "s1", DurationMins: 60, OriginalTime: time.Now()}
	s2 := SessionData{ID: "s2", DurationMins: 60, OriginalTime: time.Now()}
	r1 := models.Room{ID: "r1"}
	r2 := models.Room{ID: "r2"}

	instance := ProblemInstance{
		Sessions: map[string]SessionData{"s1": s1, "s2": s2},
		Rooms:    map[string]models.Room{"r1": r1, "r2": r2},
	}
	
	sol := &Solution{
		Assignments: map[string]Assignment{
			"s1": {SessionID: "s1", RoomID: "r1"},
			"s2": {SessionID: "s2", RoomID: "r1"},
		},
	}
	
	solver := NewSchedulerSolver(DefaultConfig())
	
	// Test correctness of mutation: It should produce a valid assignment for existing keys
	// Since mutation is random (Strategy 0 or 1), we run it multiple times to catch potential panics or corruptions
	for i := 0; i < 50; i++ {
		mutated := solver.Mutate(sol, instance)
		assert.NotNil(t, mutated)
		assert.Equal(t, 2, len(mutated.Assignments))
		assert.Contains(t, []string{"r1", "r2"}, mutated.Assignments["s1"].RoomID)
		assert.Contains(t, []string{"r1", "r2"}, mutated.Assignments["s2"].RoomID)
	}
}

func TestSolver_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cfg := DefaultConfig()
	cfg.MaxIterations = 1000000 // Long run
	solver := NewSchedulerSolver(cfg)

	instance := ProblemInstance{
		Sessions: map[string]SessionData{"s1": {ID: "s1"}},
		Rooms:    map[string]models.Room{"r1": {ID: "r1"}},
	}
    
	// Cancel immediately
	cancel()
	
	start := time.Now()
	_, err := solver.Solve(ctx, instance)
	elapsed := time.Since(start)

	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	assert.Less(t, elapsed, 100*time.Millisecond, "Solver should return immediately on cancelled context")
}
