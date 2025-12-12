# Journey Map Constructor Implementation Plan

## Goal Description
Empower university staff and teachers to create, edit, and manage custom Journey Maps (Playbooks) through a visual "No-Code" interface, moving away from hardcoded TypeScript definitions.

## Core Components

### 1. Visual Builder UI
A drag-and-drop interface for constructing the graph.
*   **Tech Stack**: Recommend **React Flow** or **TanStack Query** + **SVG** for the canvas. A custom implementation is also viable given our specific "World/Node" structure.
*   **Features**:
    *   **Palette**: Drag "Action Nodes", "Form Nodes", "Gateway Nodes" onto the canvas.
    *   **Connectivity**: Draw lines to define `prerequisites` (dependencies) and `outcomes` (branches).
    *   **Node Editor**: A side panel to configure internal node logic (Fields, Uploads, Deadlines, Roles).
    *   **World Management**: Group nodes into "Worlds" (Year 1, Year 2, etc.) dynamically.

### 2. Backend & Persistence
Transition `Playbook` from a static TS object to a stored entity.

#### Database Schema (PostgreSQL)
```sql
CREATE TABLE playbooks (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title JSONB NOT NULL, -- { "en": "...", "ru": "..." }
  version INT NOT NULL DEFAULT 1,
  is_active BOOLEAN DEFAULT false,
  structure JSONB NOT NULL, -- The full Node/World definition
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Store active journey instances for students
CREATE TABLE student_journeys (
  id UUID PRIMARY KEY,
  student_id UUID NOT NULL,
  playbook_id UUID REFERENCES playbooks(id),
  current_state JSONB, -- Node states
  started_at TIMESTAMPTZ DEFAULT NOW()
);
```

### 3. Versioning Strategy (Critical)
University curriculums change. We must handle "In-Flight" journeys carefully.
*   **Immutable Versions**: Once a Playbook is "Published", it is immutable.
*   **Migration Logic**:
    *   **Option A (Legacy Support)**: Existing students stay on their started Playbook version (v1). New students get v2.
    *   **Option B (Hot Migration)**: Harder. Requires mapping old Node IDs to new Node IDs to migrate active student states.
    *   **Recommendation**: Start with **Option A**. It's safer and matches academic "Catalogs" (Student follows the requirements of their entry year).

## Implementation Phases

### Phase 1: The "Code-Behind" API
1.  Create `playbooks` table.
2.  Create CRUD endpoints (`GET /playbooks`, `POST /playbooks`, `PUT /playbooks/:id`).
3.  Refactor Frontend to fetch `playbook` from API instead of importing from `lib/playbook.ts`.
4.  *Result*: We can update the JSON in DB to change the app, even without a UI yet.

### Phase 2: The Visual Editor
1.  Implement the Canvas (Nodes & Edges).
2.  Implement the Property Inspector (Form builder for node fields).
3.  Implement "Layout Engine" to auto-organize the graph.
4.  *Result*: Staff can visually design flows.

### Phase 3: Assignment & Role Based Access
1.  Allow assigning specific Playbooks to specific Departments or Specialties.
2.  "My Playbooks" dashboard for Professors.

## Risks & Considerations
*   **Complexity of Form Builder**: Recreating a robust form builder (for the "Form" nodes) is a project in itself. Consider using a simplified schema first (just "Upload" or "Text Input") before full dynamic forms.
*   **Validation**: Infinite loops in prerequisites must be detected by the builder to prevent crashing the journey engine.
