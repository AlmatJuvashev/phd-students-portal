# Scheduling Algorithm Strategy for PhD Portal

## Overview
The scheduling feature ("Scheduler") solves the **University Timetabling Problem**, assigning Courses to Rooms and Timeslots while satisfying constraints.

## Proposed Strategy: Two-Stage Implementation

### Stage 1: The Enforcer (Conflict Detection)
**Goal**: Prevent invalid data entry during manual scheduling.
**Mechanism**: Strict validation at the Service/Database level.
**Technique**: PostgreSQL Range Types (`TSTZRANGE`) with overlapping constraints.

**Strict Constraints (Hard Constraints)**:
1.  **Room Conflict**: A Room cannot host two events at the same time.
2.  **Instructor Conflict**: An Instructor cannot teach two courses at the same time.
3.  **Cohort Conflict**: A Student Group (Cohort) cannot appear in two places at once.

### Stage 2: The Solver (Auto-Scheduler)
**Goal**: Automatically generate valid schedules for unassigned courses.
**Algorithm Recommendation**: **Constraint Satisfaction Problem (CSP) with Backtracking**.

**Why CSP?**
-   **Exactness**: It finds a valid solution if one exists.
-   **Simplicity**: Easier to implement reliably in Go than Genetic Algorithms for this scale.
-   **Library**: Can use `gnboorse/centipede` or custom DFS solver.

**Model Definition**:
-   **Variables**: Course Sessions (e.g., "Bio101 - Lecture 1").
-   **Domains**: Cartesian product of Available Rooms × Available Timeslots.
-   **Constraints**: The Hard Constraints defined in Stage 1 plus Soft Constraints (e.g., "Minimize gaps for students").

## Implementation Roadmap
1.  **Database**: Implement `room_id`, `timeslot` (start/end), `instructor_id`, `cohort_id` on `ScheduleEvent` table.
2.  **Validator**: Implement `ValidateEventCollision(tx, event)` using SQL `&&` operator for range overlap.
3.  **UI**: Admin Ops Scheduler (Manual Drag-and-Drop) powered by Validator.
4.  **Auto-Scheduler**: Background worker using CSP solver (Future Phase).
