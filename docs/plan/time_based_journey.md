# Time-Based Journey Implementation Analysis

## Goal Description
Assess the feasibility, benefits, and drawbacks of implementing a **Time-Based / Deadline-Based Journey** alongside the current **Self-Paced (Node-Based) Journey**.

The current application uses a dependency graph (Nodes & Edges) where progress is determined by completing prerequisites. A time-based view would introduce calendar dates, hard/soft deadlines, and time-sensitive milestones.

## Pros & Cons Analysis

### Benefits (Pros)
1.  **Administrative Compliance**: PhD programs often have strict university deadlines (e.g., "Submit annual report by June 15th"). A pure self-journey map treats these as just another step, missing the urgency.
2.  **Better Planning**: Students can visualize their long-term trajectory. A "Gantt Chart" view helps them realize if they are falling behind schedule, even if they are completing daily tasks accurately.
3.  **Proactive Notifications**: Enabling deadlines allows the system to send reminders ("Your topic approval is due in 3 days"), reducing missed administrative gates.
4.  **Hybrid Utility**:
    *   *Self-Journey* handles the **"Logic"** (Step A requires Step B).
    *   *Time-Journey* handles the **"Reality"** (Step A must be done by October).

### Drawbacks (Cons)
1.  **Conflict & Anxiety**: A student might be "Locked" out of a node (missing prerequisite) but the Timeline says it's "Overdue". This creates stress. The UI needs to handle this "Blocked but Urgent" state carefully.
2.  **Maintenance Overhead**: Dates change every academic year. The system would need robust "Cohorts" or "Intake Years" logic so that 2024 students get 2024 dates, and 2025 students get 2025 dates.
3.  **False Rigidity**: Research is unpredictable. Hard-coding a timeline might make students feel they are "failing" just because experiments took longer, even if the program allows flexibility.
4.  **UI Complexity**: Maintaining two synchronized views (Map vs. Timeline) doubles the frontend complexity.

---

## Proposed Implementation Concepts

### 1. Data Structure Updates
We need to extend the `Playbook` and `NodeDef` schemas.
*   **Cohorts/Intakes**: Define start dates for different groups (e.g., "Fall 2024 Intake").
*   **Relative Deadlines**: Define deadlines relative to the start date (e.g., `T + 3 months`).
*   **Hard vs. Soft Deadlines**:
    *   *Soft*: "Recommended completion by..." (Green/Yellow status).
    *   *Hard*: "Must be submitted by..." (Red status, potentially blocking).

```typescript
// Proposed Schema Extension
type NodeDef = {
  // ... existing fields
  timing?: {
    type: "relative" | "absolute";
    offset_days?: number; // e.g., 90 days after start
    fixed_date?: string;  // e.g., "2024-12-31" (rarely used, prefer cohorts)
    duration_days?: number; // expected time to complete
    critical: boolean; // Is this a hard deadline?
  };
}
```

### 2. Visualization Options
Instead of replacing the Map, we add a toggle or a secondary view.

#### A. The "Gantt" View
A horizontal timeline showing all nodes.
*   **X-Axis**: Time (Months/Years).
*   **Bars**: Node duration.
*   **Visuals**:
    *   Completed nodes: Solid Blue.
    *   Future nodes: Faded Gray.
    *   Overdue nodes: Red outline.
    *   "Today" marker line.

#### B. The "Deadline" Overlay
Keep the existing Node Map but decorate nodes with date badges.
*   Nodes approaching deadline get a 🕒 icon and "Due: Oct 15" label.
*   Overdue nodes pulse Red.

### 3. Notification Engine
Leverage the recently implemented Notification Center.
*   **Cron Job**: Runs daily to check `Node.timing`.
*   **Triggers**:
    *   7 days before deadline: "Heads up..."
    *   1 day before deadline: "Urgent..."
    *   Overdue: "Action required..."

## Recommendation
**Implement a Hybrid "Deadline Overlay" first.**
1.  Add `timing` metadata to critical nodes (e.g., Annual Attestation, Thesis Defense).
2.  Display these dates on the existing Journey Map nodes.
3.  Send notifications based on these dates.

This avoids the complexity of building a full Gantt chart immediately while capturing 80% of the value (reminding students of critical dates).
