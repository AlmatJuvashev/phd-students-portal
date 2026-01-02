package solver

import (
	"math"
)

// CalculateCost computes total energy (penalty) of the solution.
// Higher = Worse. Goal is 0 hard constraints and minimal soft constraints.
func (s *SchedulerSolver) CalculateCost(sol *Solution, instance ProblemInstance) float64 {
	hardPenalties := 0.0
	softPenalties := 0.0

	// Pre-process: Group assignments by Room and Instructor for O(N) overlap checks
	// This is expensive inside the loop. 
	// Optimization: The SA "Move" function should effectively update these costs incrementally
	// but for MVP we will recalculate from scratch to ensure correctness first.
	
	byRoom := make(map[string][]Assignment)
	byInstructor := make(map[string][]Assignment)

	// Single pass to build maps and check simple constraints
	for _, assign := range sol.Assignments {
		// 1. Capacity Check (Hard)
		room, ok := instance.Rooms[assign.RoomID]
		if !ok { continue } // Should not happen in valid encoding
		
		session, ok := instance.Sessions[assign.SessionID]
		if !ok { continue }


		// Capacity Check
		if room.Capacity < session.MaxStudents { 
			w := s.getWeight(s.Config.CapacityConstraint, s.Config.WeightCapacity)
			if s.Config.CapacityConstraint == "HARD" {
				hardPenalties += w
			} else {
				softPenalties += w
			}
		}
		

		// Department Constraint
		roomDept := ""
		if room.DepartmentID != nil {
			roomDept = *room.DepartmentID
		}
		
		if roomDept != session.DepartmentID {
			w := s.getWeight(s.Config.DepartmentConstraint, s.Config.WeightHardConflict)
			if s.Config.DepartmentConstraint == "HARD" {
				hardPenalties += w
			} else {
				softPenalties += w
			}
		}
		
		// Attribute Check
		if reqs := session.Requirements; len(reqs) > 0 {
			if attrs, ok := instance.RoomAttributes[assign.RoomID]; ok {
				for _, r := range reqs {
					if val, has := attrs[r.Key]; !has || val != r.Value {
						// Missing Requirement
						// Treat as Soft for now unless critical config added
						softPenalties += s.Config.WeightUtilization * 5 // High soft penalty
					}
				}
			} else {
				// Room has no attributes but requirements exist -> Penalty
				softPenalties += s.Config.WeightUtilization * 5
			}
		}

		byRoom[assign.RoomID] = append(byRoom[assign.RoomID], assign)
		if session.InstructorID != "" {
			byInstructor[session.InstructorID] = append(byInstructor[session.InstructorID], assign)
			
			// Check Unavailability
			if unavailList, exists := instance.Unavailability[session.InstructorID]; exists {
				for _, slot := range unavailList {
					// Check Day of Week
					// assignment.StartTime is full Time.Time. Need to extract Day/WebDay (0=Sun, 1=Mon)
					// Go's time.Weekday(): Sunday=0
					if int(assign.StartTime.Weekday()) == slot.DayOfWeek {
						// Check Time Overlap
						// We need to parse slot.StartTime "HH:MM:SS" -> time.Time for comparison on that day
						// Or compare hours/minutes
						// Assuming standard day, let's compare HH*60+MM
						
						sH, sM, _ := assign.StartTime.Clock()
						eH, eM, _ := assign.EndTime.Clock()
						sMins := sH*60 + sM
						eMins := eH*60 + eM
						
						// Parse Slot time (naive string parse "HH:MM")
						uStartMins := parseTimeMins(slot.StartTime)
						uEndMins := parseTimeMins(slot.EndTime)
						
						if sMins < uEndMins && eMins > uStartMins {
							// Overlap with unavailable slot
							if slot.IsUnavailable {
								softPenalties += s.Config.WeightInstructor
							}
						}
					}
				}
			}
		}


		
		// Soft: Utilization
		waste := float64(room.Capacity - session.MaxStudents)
		if waste > 0 {
			softPenalties += math.Sqrt(waste) // Sub-linear penalty
		}
	}

	// 2. Overlap Checks (Hard) using the groups
	hardPenalties += countOverlaps(byRoom)
	hardPenalties += countOverlaps(byInstructor)
	
	// Cohort Overlaps (Need to build map first or do it in loop above)
	// Let's rebuild map here for clarity/correctness despite O(N) cost
	byCohort := make(map[string][]Assignment)
	for _, assign := range sol.Assignments {
		if sess, ok := instance.Sessions[assign.SessionID]; ok {
			for _, cid := range sess.Cohorts {
				byCohort[cid] = append(byCohort[cid], assign)
			}
		}
	}
	hardPenalties += countOverlaps(byCohort)

	// Total Energy
	// WeightHardConflict should be large enough to dominate any soft constraints
	totalCost := (hardPenalties * s.Config.WeightHardConflict) + 
				 (softPenalties * s.Config.WeightUtilization)
	
	// Check feasibility
	sol.IsValid = (hardPenalties == 0)
	
	return totalCost
}

// countOverlaps checks for time intersections within a group of assignments
func countOverlaps(groups map[string][]Assignment) float64 {
	conflicts := 0.0
	for _, assignments := range groups {
		// O(N^2) for assignments in same resource. Usually N is small (classes in Room 101).
		for i := 0; i < len(assignments); i++ {
			for j := i + 1; j < len(assignments); j++ {
				if overlap(assignments[i], assignments[j]) {
					conflicts += 1.0
				}
			}
		}
	}
	return conflicts
}

func overlap(a, b Assignment) bool {
	// (StartA < EndB) and (EndA > StartB)
	return a.StartTime.Before(b.EndTime) && a.EndTime.After(b.StartTime)
}

func (s *SchedulerSolver) getWeight(severity string, defaultWeight float64) float64 {
	switch severity {
	case "HARD":
		return s.Config.WeightHardConflict // Always use high penalty
	case "SOFT":
		return defaultWeight // Use specific soft weight
	default:
		return 0 // OFF
	}
}
