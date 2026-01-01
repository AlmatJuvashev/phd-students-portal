package solver

import (
	"context"
	"math"
	"math/rand"
	"time"
)

// SchedulerSolver is the main entry point for the optimization engine.
type SchedulerSolver struct {
	Config SolverConfig
}

func NewSchedulerSolver(cfg SolverConfig) *SchedulerSolver {
	return &SchedulerSolver{Config: cfg}
}

// Solve runs the Simulated Annealing algorithm to find the best schedule.
func (s *SchedulerSolver) Solve(ctx context.Context, instance ProblemInstance) (*Solution, error) {
	rand.Seed(time.Now().UnixNano())

	// 1. Generate Initial Solution (Random or Greedy)
	currentSolution := s.GenerateInitialSolution(instance)
	currentSolution.Score = s.CalculateCost(currentSolution, instance)
	
	bestSolution := currentSolution.Clone()
	
	temp := s.Config.InitialTemp

	// 2. Annealing Loop
	for i := 0; i < s.Config.MaxIterations; i++ {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return bestSolution, ctx.Err()
		default:
		}

		// 3. Mutate: Generate neighbor
		neighbor := s.Mutate(currentSolution, instance)
		neighbor.Score = s.CalculateCost(neighbor, instance)

		// 4. Acceptance Probability
		delta := neighbor.Score - currentSolution.Score // We want to MINIMIZE score (Energy)
		// If delta < 0 (Neighbor is better), we accept.
		// If delta > 0 (Neighbor is worse), we might accept.

		if delta < 0 || s.shouldAccept(delta, temp) {
			currentSolution = neighbor
			
			// Update global best
			if currentSolution.Score < bestSolution.Score {
				bestSolution = currentSolution.Clone()
			}
		}

		// 5. Cool down
		temp *= s.Config.CoolingRate
		
		// Termination Optimization: If temp is essentially zero, break early
		if temp < 0.001 {
			break
		}
	}

	return bestSolution, nil
}

// shouldAccept employs the Metropolis criterion
func (s *SchedulerSolver) shouldAccept(delta, temp float64) bool {
	probability := math.Exp(-delta / temp)
	return rand.Float64() < probability
}

// GenerateInitialSolution creates a starting point using a fast greedy fit.
func (s *SchedulerSolver) GenerateInitialSolution(instance ProblemInstance) *Solution {
	sol := &Solution{
		Assignments: make(map[string]Assignment),
	}

	// Convert map keys to slice for random access if needed, or just iterate
	// For initial solution, we just iterate.
	roomIDs := make([]string, 0, len(instance.Rooms))
	for rID := range instance.Rooms {
		roomIDs = append(roomIDs, rID)
	}

	for _, session := range instance.Sessions {
		// MVP: Random assignment
		// Optimization: Try to find a room with enough capacity first
		var bestRoomID string
		
		// 1. Try to find a fitting room
		for _, rID := range roomIDs {
			if instance.Rooms[rID].Capacity >= session.MaxStudents {
				bestRoomID = rID
				break
			}
		}
		// 2. Fallback to random if no perfect fit
		if bestRoomID == "" && len(roomIDs) > 0 {
			bestRoomID = roomIDs[rand.Intn(len(roomIDs))]
		}
		
		// Time: Keep original time if set, or random?
		// Assuming we are RESCHEDULING, we might want to stick to original times initially
		// or random shifts. Let's stick to original time +/- 0 hours for now.
		startTime := session.OriginalTime
		endTime := startTime.Add(time.Duration(session.DurationMins) * time.Minute)

		sol.Assignments[session.ID] = Assignment{
			SessionID: session.ID,
			RoomID:    bestRoomID,
			StartTime: startTime,
			EndTime:   endTime,
		}
	}
	return sol
}

// Mutate modifies the solution slightly.
// Strategies:
// 1. Change Room (Move to different room, same time)
// 2. Change Time (Move to different time, same room) - NOT IMPLEMENTED for MVP, assuming fixed slots for now?
//    Actually, we should implement simple time shifts.
// 3. Swap (Swap rooms/times between two sessions)
func (s *SchedulerSolver) Mutate(sol *Solution, instance ProblemInstance) *Solution {
	newSol := sol.Clone()
	
	if len(instance.Sessions) == 0 {
		return newSol
	}

	// Pick a random session to mutate
	// We need a slice of IDs to pick randomly
	sessionIDs := make([]string, 0, len(instance.Sessions))
	for id := range instance.Sessions {
		sessionIDs = append(sessionIDs, id)
	}
	targetID := sessionIDs[rand.Intn(len(sessionIDs))]
	targetAssign := newSol.Assignments[targetID]

	// Mutation Strategy Selector
	strategy := rand.Intn(2) // 0 or 1

	if strategy == 0 {
		// Strategy 0: Change Room
		// Pick random room
		roomIDs := make([]string, 0, len(instance.Rooms))
		for rID := range instance.Rooms {
			roomIDs = append(roomIDs, rID)
		}
		if len(roomIDs) > 0 {
			newRoomID := roomIDs[rand.Intn(len(roomIDs))]
			targetAssign.RoomID = newRoomID
			newSol.Assignments[targetID] = targetAssign
		}
	} else {
		// Strategy 1: Swap Rooms with another random session
		// (This preserves time slots and number of assigned sessions)
		otherID := sessionIDs[rand.Intn(len(sessionIDs))]
		if otherID != targetID {
			otherAssign := newSol.Assignments[otherID]
			
			// Swap Rooms
			tempRoom := targetAssign.RoomID
			targetAssign.RoomID = otherAssign.RoomID
			otherAssign.RoomID = tempRoom
			
			newSol.Assignments[targetID] = targetAssign
			newSol.Assignments[otherID] = otherAssign
		}
	}
	
	// Future: Strategy 2 - Change Time (requires strict boundaries and slot logic)

	return newSol
}



// Clone creates a deep copy of the solution
func (old *Solution) Clone() *Solution {
	newSol := &Solution{
		Assignments: make(map[string]Assignment, len(old.Assignments)),
		Score:       old.Score,
		IsValid:     old.IsValid,
	}
	for k, v := range old.Assignments {
		newSol.Assignments[k] = v
	}
	return newSol
}
