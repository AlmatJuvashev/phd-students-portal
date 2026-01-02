package solver

import (
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCalculateCost(t *testing.T) {
	// Standard Setup
	now := time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC)
	
	// Create minimal problem instance
	instance := ProblemInstance{
		Rooms: map[string]models.Room{
			"R1": {ID: "R1", Capacity: 50, DepartmentID: strPtr("CS")},
			"R2": {ID: "R2", Capacity: 20, DepartmentID: nil},
		},
		Sessions: map[string]SessionData{
			"S1": {ID: "S1", MaxStudents: 30, DepartmentID: "CS", InstructorID: "I1"},
			"S2": {ID: "S2", MaxStudents: 40, DepartmentID: "CS"},
		},
	}

	config := DefaultConfig()
	config.WeightHardConflict = 1000
	config.WeightCapacity = 10
	config.WeightUtilization = 1
	config.CapacityConstraint = "HARD"
	config.DepartmentConstraint = "SOFT"

	solver := NewSchedulerSolver(config)

	t.Run("Valid Solution", func(t *testing.T) {
		sol := &Solution{
			Assignments: map[string]Assignment{
				"S1": {SessionID: "S1", RoomID: "R1", StartTime: now, EndTime: now.Add(1 * time.Hour)},
			},
		}
		
		cost := solver.CalculateCost(sol, instance)
		// Capacity utilization: 50 - 30 = 20. Sqrt(20) = 4.472...
		// Cost = 0*Hard + 4.472*1 (WeightUtilization)
		assert.InDelta(t, 4.47, cost, 0.01)
		assert.True(t, sol.IsValid)
	})

	t.Run("Hard Constraint Violation (Overlap)", func(t *testing.T) {
		sol := &Solution{
			Assignments: map[string]Assignment{
				"S1": {SessionID: "S1", RoomID: "R1", StartTime: now, EndTime: now.Add(1 * time.Hour)},
				"S2": {SessionID: "S2", RoomID: "R1", StartTime: now, EndTime: now.Add(1 * time.Hour)},
			},
		}
		
		cost := solver.CalculateCost(sol, instance)
		assert.Greater(t, cost, 1000.0)
		assert.False(t, sol.IsValid)
	})

	t.Run("Capacity Violation - Hard", func(t *testing.T) {
		solver.Config.CapacityConstraint = "HARD"
		// Ensure Dept matches to avoid noise
		instance.Rooms["R2"] = models.Room{ID: "R2", Capacity: 20, DepartmentID: strPtr("CS")}
		
		sol := &Solution{
			Assignments: map[string]Assignment{
				// S2 (40) > R2 (20)
				"S2": {SessionID: "S2", RoomID: "R2", StartTime: now, EndTime: now.Add(1 * time.Hour)},
			},
		}
		
		cost := solver.CalculateCost(sol, instance)
		assert.Greater(t, cost, 1000.0)
	})

	t.Run("Capacity Violation - Soft", func(t *testing.T) {
		solver.Config.CapacityConstraint = "SOFT"
		// Ensure Dept matches to avoid noise
		instance.Rooms["R2"] = models.Room{ID: "R2", Capacity: 20, DepartmentID: strPtr("CS")}
		
		sol := &Solution{
			Assignments: map[string]Assignment{
				"S2": {SessionID: "S2", RoomID: "R2", StartTime: now, EndTime: now.Add(1 * time.Hour)},
			},
		}
		
		cost := solver.CalculateCost(sol, instance)
		// Hard = 0. Soft = WeightCapacity(10) * WeightUtilization(1) = 10? 
		// Wait, softPenalties += w. w=10.
		// TotalCost = soft * 1 = 10.
		assert.InDelta(t, 10.0, cost, 0.1)
	})
}

func strPtr(s string) *string {
	return &s
}
