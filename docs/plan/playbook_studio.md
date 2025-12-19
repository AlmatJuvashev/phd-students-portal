# Playbook Studio Specification (v1)
## Online Course + Journey Map Constructor (DB-First, JSON-Optional)

### Status
- Version: **v1 (MVP, production-ready)**
- Goal: Replace `playbook.json` as the primary source of truth.

---

## 1. Product Goals

### 1.1 Primary goals
1) **Make `playbook.json` obsolete** as runtime configuration by moving authoring to an online editor with publish/versioning.
2) Provide an **excellent authoring UX**:
- fast creation and editing
- minimal friction for repetitive work (templates, duplication, command palette)
- safe publishing (draft vs published with preview + validation)
3) Support **universal education** scenarios:
- university admissions + research workflows
- K-12 learning flows
- prep-center exam preparation journeys

### 1.2 Non-goals (v1)
- Real-time multi-user co-editing with conflict resolution (can be added later)
- LTI/OneRoster/QTI/SCORM authoring or integration (separate milestones)
- Payment/billing and CRM

---

## 2. Core Concepts & Terminology

### 2.1 Entities
- **Course**: container for learning content and workflows.
- **Module**: section inside a course.
- **Lesson**: a unit of learning inside a module.
- **Activity**: an actionable unit within a lesson (content, assignment, checklist, etc.).
- **Journey**: a graph or ordered path of steps that represent a workflow; steps can link to activities or be independent “workflow nodes”.

### 2.2 Draft vs Published
- **Draft**: editable version (autosaved).
- **Published**: immutable snapshot (used by runtime).
- Publishing creates a **new published version** and preserves history.

### 2.3 “Obsolete JSON” principle
- DB is the source of truth.
- JSON exists only as:
  - import/export format
  - backup
  - template distribution

---

## 3. Personas & Permissions (v1)

### 3.1 Roles
- **Org Admin**: full control across org.
- **Teacher/Builder**: can create/edit courses and journeys if granted.
- **Student**: read-only published content.

### 3.2 Permission rules (must implement)
- Only **Org Admin** or **Builder** can edit drafts or publish.
- Students only view **published** versions.
- A Builder can edit only courses they own or are assigned to.

### 3.3 Auditing (must implement)
Log actions:
- publish/unpublish
- version restore
- archive/unarchive nodes
- role changes
- import/export
- document requirement changes (security-sensitive)

---

## 4. UX Requirements (Excellent UX, Minimum Feature Set)

### 4.1 Editing ergonomics (must have)
- Autosave (debounced, e.g. 500–1500 ms) with a visible status:
  - “Saving…” → “Saved ✓ at 13:41”
- Undo/Redo:
  - per-builder session (client-side stack)
- Keyboard shortcuts:
  - `Cmd/Ctrl + K` command palette
  - `Cmd/Ctrl + Z` undo, `Shift + Cmd/Ctrl + Z` redo
  - `/` focus search
  - `N` add new step (context-aware)

### 4.2 Builder views (must have)
- **Outline view** (primary): handles 100–500+ steps reliably
- **Flow view** (secondary): visual graph, collapsible groups
- “Insert between steps” affordance (`+`) with type picker

### 4.3 Safety & clarity (must have)
- Draft vs Published badge always visible
- Publish is explicit:
  - publish dialog shows warnings/errors and “what changed”
- Validation panel:
  - errors block publish
  - warnings allow publish (optional config)
- Preview as Student:
  - simulate lock/unlock
  - show due dates
  - show required docs
  - show approvals

### 4.4 Quality requirements
- Fast: outline operations remain responsive with 500 nodes
- Accessible: keyboard navigation for key actions, readable contrast
- Resilient: builder can recover from refresh (draft persisted)

---

## 5. Data Model (Recommended Hybrid: JSONB Authoring + Version Snapshots)

### 5.1 Why Hybrid
- Fast to ship and iterate on UX (JSON editing is flexible)
- Versioning is simple (snapshot per publish)
- Later analytics can be added by denormalizing/ETL if needed

### 5.2 Tables (minimal)
**Organizations**
- `organizations(id, slug, name, type, branding_json, settings_json, created_at)`

**Courses (v1)**
- `courses(id, organization_id, title, description, status, created_at, created_by)`
  - `status`: `draft | published | archived`

**Course versions**
- `course_versions(id, course_id, version, state, snapshot_json, created_at, created_by, release_notes)`
  - `state`: `draft | published`
  - `snapshot_json`: full course structure (modules/lessons/activities metadata)
  - `release_notes`: optional

**Journeys**
- `journeys(id, organization_id, course_id, status, created_at, created_by)`

**Journey versions**
- `journey_versions(id, journey_id, version, state, graph_json, created_at, created_by, release_notes, diff_summary_json)`
  - `graph_json`: nodes + edges + metadata (see schema below)
  - `diff_summary_json`: computed at publish time (optional)

**Optional (recommended for search)**
- `journey_node_index(id, journey_version_id, node_id, title, type, tags_json, is_archived)`
  - built at publish time for fast filtering/search

> Note: If you already have normalized course/module/lesson/activity tables, keep them — but treat Studio as producing a **published snapshot** that runtime reads consistently.

---

## 6. Journey Graph JSON Schema (v1)

### 6.1 Top-level structure
```json
{
  "meta": {
    "title": "PhD Onboarding Journey",
    "description": "Steps from enrollment to first-year milestones",
    "startNodeId": "node_001",
    "layout": { "type": "flow" },
    "version": 12
  },
  "nodes": [
    { "...node..." }
  ],
  "edges": [
    { "...edge..." }
  ]
}

6.3 Edge schema (required)

{
  "id": "edge_001",
  "from": "node_001",
  "to": "node_002",
  "condition": {
    "type": "completion",
    "rule": "completed"
  }
}

6.4 Node types (v1)

info — read-only informational step

content — links to lesson/activity content

assignment — student response required (text/files)

checklist — list of sub-steps with completion checkboxes

milestone — gated step requiring approval

quiz — assessment step

upload_only — documents required but no narrative response

external_link — points outside platform

7. Validation Rules (Block Publish on Errors)
7.1 Errors (block publish)

Missing meta.startNodeId or invalid node reference

Duplicate node ids

Edges referencing missing nodes

Cycles if settings.allowCycles = false

Unreachable nodes from start (unless explicitly allowed)

Invalid due rules (negative days, missing anchors)

Node has approval.required = true but no approval roles

Document requirements invalid (minCount > maxCount, empty mime list)

Node links invalid (activityId missing for link.kind=activity)

7.2 Warnings (allow publish with notice)

Orphan nodes intentionally hidden by visibility rules

Very large documents list (e.g., > 20 required docs) — usability warning

Nodes with long titles/description beyond suggested limits

Multiple start-like nodes (nodes with no incoming edges) besides start node

7.3 “Fix-it” UX requirements

Validation panel links each issue to the offending node

One-click actions where possible (e.g., “set as start node”)

8. Runtime Semantics (How the Published Journey Behaves)
8.1 Status computation (v1)

Each node for a student is one of:

locked — prereqs not satisfied

available — unlocked, not completed

in_review — submitted, awaiting approval/grading

completed — done/approved

overdue — available but past due date

8.2 Completion criteria (v1)

info/content: student can mark as completed (optional org setting)

assignment/quiz/upload_only: completed when submission exists and accepted/graded (or auto-graded)

checklist: completed when all sub-steps are checked

milestone: completed only when approved by authorized role

9. Publishing & Versioning
9.1 Draft lifecycle

Draft is autosaved frequently.

Draft changes do not affect students until published.

9.2 Publishing process (v1)

Validate draft.

Generate diff summary (optional, recommended).

Create new journey_versions row with state=published.

Mark journey status published.

Runtime reads latest published version for that course/org.

9.3 Rollback (v1)

Admin can restore any previous published version as “current” (creates a new published version pointing to that snapshot or sets pointer).

9.4 Release notes (v1)

Publish dialog includes optional release notes.

10. Import/Export (Retire playbook.json Safely)
10.1 Import (v1)

playbook.json → draft journey version

map node types and fields; store unknown fields in metadata_json to avoid data loss

show import summary:

nodes imported

edges imported

fields mapped with warnings

10.2 Export (v1)

Export the published journey version to a stable JSON format

Include metadata: org, course, version, exported_at

10.3 Migration strategy

Runtime fallback:

If no published DB journey exists, optionally use legacy JSON (temporary)

Convert orgs progressively.

When stable, remove JSON fallback.

11. API Specification (v1)
11.1 Studio: Courses

GET /studio/courses (list courses)

POST /studio/courses (create)

GET /studio/courses/:courseId (details)

PUT /studio/courses/:courseId (update metadata)

POST /studio/courses/:courseId/publish (publish course snapshot)

GET /studio/courses/:courseId/versions (course versions)

POST /studio/courses/:courseId/clone (duplicate)

11.2 Studio: Journeys

GET /studio/courses/:courseId/journey/draft

PUT /studio/courses/:courseId/journey/draft (autosave)

POST /studio/courses/:courseId/journey/validate

POST /studio/courses/:courseId/journey/publish

GET /studio/courses/:courseId/journey/versions

POST /studio/courses/:courseId/journey/restore (restore a version)

POST /studio/courses/:courseId/journey/import

GET /studio/courses/:courseId/journey/export

11.3 Runtime (student-facing)

GET /courses/:courseId/journey → latest published graph + computed status per node

POST /journey/nodes/:nodeId/complete

POST /journey/nodes/:nodeId/submit (for assignments/upload-only)

POST /journey/nodes/:nodeId/approve (authorized roles)

POST /journey/nodes/:nodeId/return (authorized roles)

Note: You can keep runtime endpoints aligned with your existing API; Studio endpoints are admin-only.

12. UI Specifications (v1)
12.1 Information architecture

/admin/studio — landing

/admin/studio/courses — course library

/admin/studio/courses/:id/builder — course structure builder

/admin/studio/courses/:id/journey — journey builder (outline + flow)

/admin/studio/courses/:id/preview — preview as student

/admin/studio/templates — templates gallery (optional v1)

/admin/studio/settings — org-level builder settings (optional v1)

12.2 Course Library UI requirements

Search (title, tags)

Filters:

status (draft/published/archived)

owner

Actions:

create new

duplicate

archive

export

12.3 Journey Builder UI requirements

Top toolbar:

version badge

autosave indicator

validate

preview

publish

Left panel:

outline tree (chapters/groups)

search nodes

Main area:

Outline list (sortable) OR Flow canvas

Right inspector drawer (tabs):

Basics

Requirements

Gating

Due dates

Approval

Visibility

Link (activity/external)

Validation panel:

errors/warnings list

click → focus node

12.4 Preview UI requirements

Student-like rendering:

locked/available/overdue/in_review/completed

Simulated user role selection (student/teacher/admin)

Simulated time offset (optional):

“today + 10 days” to test overdue logic

13. Performance & Reliability Requirements

Autosave must not spam server:

debounce; avoid saving if no changes

Large graphs:

Outline view must handle 500 nodes

Flow view can be limited/collapsed in v1

All operations must remain tenant-scoped

Publish must be atomic:

transaction wrapping version creation + status updates

14. Acceptance Criteria (Core v1 for Playbook Studio)

A) Obsolete JSON

Runtime can load a published journey from DB and function without any JSON file.

B) Authoring

Admin creates a course, creates a journey, publishes it, and students see it.

C) Safety

Draft edits do not affect students until publish.

Validation blocks publish on errors.

D) UX

Autosave works; undo/redo works; “insert between nodes” works; preview works.

E) Migration

Existing playbook.json can be imported with a clear summary.

15. Implementation Notes (Practical Recommendations)

Build Outline-first. Flow canvas is “nice to have” but outline must be perfect.

Store icons as simple strings; map to a curated icon set.

Keep node configs forward-compatible:

unknown fields should be preserved in JSON for future.

Add a “Template” concept later; in v1, duplication is enough.