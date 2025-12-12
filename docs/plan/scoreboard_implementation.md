# Scoreboard Feature Implementation Plan

## Goal Description
Gamify the student journey by adding a **Scoreboard** (Leaderboard) that tracks progress via XP points. This fosters healthy competition and provides a clear metric of success.

## Core Rules
1.  **XP Logic**:
    *   **Standard Node**: 100 XP per 'done' state.
    *   **Conditional (Level 3) Node**: 0 XP (excluded from score).
2.  **Visualization**:
    *   **Top 5 Students**: Avatar, Name, Total XP.
    *   **Average Score**: Benchmark for the cohort.
    *   **Current User**: Their specific Rank and Score.

## Technical Implementation

### 1. Backend (Go)
**New Endpoint**: `GET /api/v1/journey/scoreboard`

**Logic**:
1.  Fetch all `journey_states` where `state = 'done'`.
2.  Fetch active Playbook from `PlaybookManager`.
3.  **Filter**: Iterate through states.
    *   Get `NodeDef` from Manager.
    *   Identify World ID. If `World == "W3"`, skip.
    *   Else, `Score += 100`.
4.  **Aggregate**: Group by `user_id`.
5.  **Enrich**: Fetch User details (Name, Avatar) for the ranked list.
6.  **Calculate Stats**:
    *   `Average`: Sum(All Scores) / Count(Active Students).
    *   `Rank`: Sort list, find index of requesting user.

**Response Structure**:
```json
{
  "top_5": [
    { "user_id": "...", "name": "Almat J.", "avatar": "...", "score": 2500 }
  ],
  "average_score": 1800,
  "me": {
    "score": 2200,
    "rank": 3
  }
}
```

### 2. Frontend (React)
**New Component**: `ScoreboardWidget`
*   **Location**: Inside `JourneyMap.tsx` HUD (e.g., a "Trophy" icon that opens a Popover/Sheet).
*   **Design**:
    *   **Header**: "Leaderboard".
    *   **List**: 1. 🥇 User A (3000 XP) ... 5. User E.
    *   **Divider**: "..."
    *   **Footer**: "You are #12 (1500 XP)".
    *   **Visuals**: Use `lucide-react` icons (Trophy, Medal).

## Database Changes
None required. We leverage existing `journey_states` and `users` tables.

## Risks & Considerations
*   **Privacy**: Some universities allow public leaderboards, others strictly forbid showing student names.
    *   *Mitigation*: Add a feature flag `ENABLE_PUBLIC_LEADERBOARD`. If false, only show "You vs Average". For now, we build the full version as requested.
*   **Performance**: Aggregating *all* journey states on every request might be slow if we hit 10k students.
    *   *Mitigation*: Cache the result in Redis for 1-5 minutes (`SCOREBOARD_CACHE`).
