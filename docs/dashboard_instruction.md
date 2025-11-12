# Admin Panel Implementation Plan

## 1. Goals

Provide a role-aware administrative interface with:

- Consistent layout (sidebar + main content) using shadcn/ui.
- Separation of concerns: layout, pages, shared components.
- Secure role-based access enforcement (superadmin vs admin).
- CRUD creation flows for admins and regular users (students, advisors, chairs).
- Visibility into student progress (completion %, active node, last activity).

## 2. Roles & Permissions Matrix

| Feature / Page                          | superadmin | admin | student | advisor | chair |
| --------------------------------------- | ---------- | ----- | ------- | ------- | ----- |
| Access admin layout                     | ✅         | ✅    | ❌      | ❌      | ❌    |
| Create Admins (create admin accounts)   | ✅         | ❌    | ❌      | ❌      | ❌    |
| Create Users (students/advisors/chairs) | ✅         | ✅    | ❌      | ❌      | ❌    |
| Student Progress Overview               | ✅         | ✅    | ❌      | ❌      | ❌    |
| View Dashboard (aggregate stats)        | ✅         | ✅    | ❌      | ❌      | ❌    |

Notes:

- superadmin may create admins.
- admin cannot create superadmin.
- Both admin & superadmin see student progress.

## 3. High-Level Architecture

```
frontend/src/
  layouts/
    AdminLayout.tsx          # layout skeleton (sidebar + content)
  pages/admin/
    Dashboard.tsx
    CreateAdmins.tsx         # guarded: superadmin only
    CreateUsers.tsx          # guarded: admin + superadmin
    StudentProgress.tsx      # guarded: admin + superadmin
  components/admin/
    SidebarNav.tsx           # navigation items rendered based on role
    UserCreateForm.tsx       # base form reused (admin variant + superadmin variant)
    AdminCreateForm.tsx      # specialized for admin creation
    StudentProgressTable.tsx # table + filtering
    RoleBadge.tsx            # small visual for role
  lib/admin/
    access.ts                # role predicates / guard helpers
    api.ts                   # admin API helpers
```

## 4. Routing Strategy

- All admin pages nested under `/admin/*`.
- Add `ProtectedRoute` (already partially designed) extension: `requiredAnyRole?: Role[]`.
- Example route config:

```tsx
<Route
  path="/admin"
  element={
    <ProtectedRoute requiredAnyRole={["admin", "superadmin"]}>
      <AdminLayout />
    </ProtectedRoute>
  }
>
  <Route index element={<Dashboard />} />
  <Route
    path="create-admins"
    element={
      <ProtectedRoute requiredAnyRole={["superadmin"]}>
        <CreateAdmins />
      </ProtectedRoute>
    }
  />
  <Route
    path="create-users"
    element={
      <ProtectedRoute requiredAnyRole={["admin", "superadmin"]}>
        <CreateUsers />
      </ProtectedRoute>
    }
  />
  <Route
    path="student-progress"
    element={
      <ProtectedRoute requiredAnyRole={["admin", "superadmin"]}>
        <StudentProgress />
      </ProtectedRoute>
    }
  />
</Route>
```

## 5. Backend Endpoints (Target State)

| Method | Endpoint                    | Purpose                                   | Auth Role        |
| ------ | --------------------------- | ----------------------------------------- | ---------------- |
| POST   | /api/users                  | Create user (student/advisor/chair/admin) | admin/superadmin |
| GET    | /api/users?role=student     | List users (filterable)                   | admin/superadmin |
| GET    | /api/admin/student-progress | Aggregate journey progress                | admin/superadmin |
| GET    | /api/admin/stats            | Dashboard summary (counts)                | admin/superadmin |

### 5.1 Student Progress Payload Model

```json
[
  {
    "id": "uuid",
    "name": "Иванов И.И.",
    "email": "ivanov@example.com",
    "role": "student",
    "progress": {
      "completed_nodes": 34,
      "total_nodes": 120,
      "percent": 28.3,
      "current_node_id": "S1_publications_list",
      "last_submission_at": "2025-11-10T12:03:00Z"
    }
  }
]
```

### 5.2 Deriving Progress

- Query `node_instances` joined with `nodes` to count total per playbook.
- Aggregate states: `done|completed|approved` contribute to completion count.
- Latest activity: MAX(updated_at) from submissions or node_instances.

## 6. UI Components Breakdown

### 6.1 SidebarNav

- Collapsible on mobile using `<Sheet />`.
- Shows menu groups:
  - Core: Dashboard, Student Progress
  - Management: Create Users, (Create Admins if superadmin)
- Active route highlighting via `useLocation()`.

### 6.2 AdminLayout

Responsibilities:

- Provide sidebar + topbar.
- Inject `<Outlet />` region.
- Provide context: currentUser role for children.

### 6.3 CreateAdmins Page

- Form fields: first_name, last_name, email.
- Hidden role = `admin` (superadmin cannot create superadmin by default unless future toggle added).
- After success: toast + optional quick-create another.
- Duplicate email handling.

### 6.4 CreateUsers Page

- Role select: student | advisor | chair.
- Form component reuses base validation.
- Bulk creation (future enhancement) placeholder.

### 6.5 StudentProgress Page

Features MVP:

- Table columns: Name, Email, % Progress (Progress bar), Current Node (tooltip title), Last Activity (relative time), Actions (View Journey).
- Filters: search by name/email, role=student implicitly.
  Future enhancements: export CSV, column sort, pagination.

## 7. State Management

- Use local component state + SWR/React Query (optional future) for caching.
- API wrapper `lib/admin/api.ts` for fetch with auth token.

## 8. Access Control Utilities (lib/admin/access.ts)

```ts
export const isSuperAdmin = (u?: User) => u?.role === "superadmin";
export const isAdmin = (u?: User) =>
  u?.role === "admin" || u?.role === "superadmin";
export const requireAny = (u: User | null, roles: string[]) =>
  !!u && roles.includes(u.role);
```

## 9. Milestones & Suggested Commits

| Milestone | Scope                          | Commit Message Example                                                 |
| --------- | ------------------------------ | ---------------------------------------------------------------------- |
| M1        | Layout shell + routing guards  | feat(admin): add AdminLayout with sidebar and protected nested routes  |
| M2        | CreateAdmins page (superadmin) | feat(admin): implement superadmin Create Admins form + API integration |
| M3        | CreateUsers page               | feat(admin): add user creation page for students/advisors/chairs       |
| M4        | Backend progress endpoint      | feat(api): add student progress aggregation endpoint                   |
| M5        | StudentProgress UI             | feat(admin): implement student progress table with filters             |
| M6        | Polish & role badges           | feat(admin): add role badges, loading states, toasts                   |

## 10. Detailed AI Implementation Prompts

### Prompt A (Milestone M1)

"""
Implement AdminLayout with sidebar using shadcn/ui.
Requirements:

- Create `layouts/AdminLayout.tsx`.
- Sidebar: fixed on desktop (w-64), hidden on mobile with hamburger (Sheet).
- Topbar: page title placeholder + user avatar + logout.
- Navigation items (Dashboard, Student Progress, Create Users, Create Admins (superadmin only)).
- Add routes under /admin with ProtectedRoute wrapper.
- Highlight active link.
- Add skeleton loading state.
  """

### Prompt B (Milestone M2)

"""
Create `pages/admin/CreateAdmins.tsx`:

- Use Card + Form (shadcn form components if available or simple form controls).
- Fields: first_name, last_name, email.
- POST /api/users with role=admin.
- On success: show toast, reset form.
- Guard: only superadmin. If not authorized -> redirect /admin.
  """

### Prompt C (Milestone M3)

"""
Create `pages/admin/CreateUsers.tsx`:

- Fields: first_name, last_name, email, role (select: student|advisor|chair).
- Validate required fields + email format.
- Show error from backend if duplicate.
- After submit: toast + keep values cleared except role.
- Show recent created users (optional future placeholder section).
  """

### Prompt D (Milestone M4 Backend)

"""
Add endpoint GET /api/admin/student-progress.
Steps:

- Query: select users where role='student'.
- Join node_instances to count completed.
- Total nodes: from playbook active cache (pbManager) or static query against nodes table.
- Return JSON array as described in plan section 5.1.
- Protect endpoint: admin or superadmin only (middleware check).
  """

### Prompt E (Milestone M5 UI)

"""
Create `pages/admin/StudentProgress.tsx`:

- Fetch /api/admin/student-progress.
- Show table with columns: Name, Email, Progress (bar + %), Current Node (id with tooltip from playbook title), Last Activity (format distance), Actions (View Journey -> navigate /node/<currentNodeId>).
- Add search input filtering client-side by name/email.
- Add empty state.
- Add loading + error states.
  """

### Prompt F (Milestone M6 Polish)

"""
Enhance admin UI:

- Add RoleBadge component.
- Add error boundary around admin routes.
- Add toasts for success/error on all forms.
- Improve responsive layout (sidebar collapse, mobile spacing).
- Add aria-labels and test keyboard navigation.
  """

## 11. Risks & Mitigations

| Risk                      | Impact                    | Mitigation                                                   |
| ------------------------- | ------------------------- | ------------------------------------------------------------ |
| Role checks duplicated    | Inconsistent behavior     | Centralize in `access.ts` and ProtectedRoute extension       |
| Large progress query slow | Slow admin UI             | Add indexes on (users.role), (node_instances.user_id, state) |
| Frontend state staleness  | Outdated progress display | Add manual refresh button, later adopt SWR/React Query       |
| Unauthorized access race  | Flicker on load           | Gate rendering until auth state resolved                     |

## 12. Success Criteria

- All routes enforce correct access.
- Creating users/admins works end-to-end with visible toasts.
- Student progress loads < 1s with sample dataset (opt: instrumentation).
- Layout responsive: mobile (sheet), desktop (persistent sidebar).
- Clear commit history aligned to milestones.

## 13. Future Enhancements (Backlog)

- Bulk CSV import of students.
- Export student progress as CSV.
- Impersonation mode (superadmin -> student view).
- Audit log page.
- Pagination + server-side filtering for large cohorts.

---

Generated implementation blueprint ready for stepwise execution.
