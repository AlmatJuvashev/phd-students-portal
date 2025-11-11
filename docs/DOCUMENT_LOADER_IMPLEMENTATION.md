# Admin Panel Implementation Plan

## 1. Goals
Provide a role-aware administrative interface with:
- Consistent layout (sidebar + main content) using shadcn/ui.
- Separation of concerns: layout, pages, shared components.
- Secure role-based access enforcement (superadmin vs admin).
- CRUD creation flows for admins and regular users (students, advisors, chairs).
- Visibility into student progress (completion %, active node, last activity).

## 2. Roles & Permissions Matrix
| Feature / Page | superadmin | admin | student | advisor | chair |
|----------------|------------|-------|---------|---------|-------|
| Access admin layout | ✅ | ✅ | ❌ | ❌ | ❌ |
| Create Admins (create admin accounts) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Create Users (students/advisors/chairs) | ✅ | ✅ | ❌ | ❌ | ❌ |
| Student Progress Overview | ✅ | ✅ | ❌ | ❌ | ❌ |
| View Dashboard (aggregate stats) | ✅ | ✅ | ❌ | ❌ | ❌ |

Notes:
- superadmin may create admins.
- admin cannot create superadmin.
- Both admin & superadmin see student progress.

## 3. High-Level Architecture
```
frontend/src/
  layouts/
    AdminLayout.tsx
  pages/admin/
    Dashboard.tsx
    CreateAdmins.tsx
    CreateUsers.tsx
    StudentProgress.tsx
  components/admin/
    SidebarNav.tsx
    UserCreateForm.tsx
    AdminCreateForm.tsx
    StudentProgressTable.tsx
    RoleBadge.tsx
  lib/admin/
    access.ts
    api.ts
```

## 4. Routing Strategy
```tsx
<Route path="/admin" element={
  <ProtectedRoute requiredAnyRole={["admin","superadmin"]}>
    <AdminLayout />
  </ProtectedRoute>}>
  <Route index element={<Dashboard />} />
  <Route path="create-admins" element={
    <ProtectedRoute requiredAnyRole={["superadmin"]}>
      <CreateAdmins />
    </ProtectedRoute>} />
  <Route path="create-users" element={
    <ProtectedRoute requiredAnyRole={["admin","superadmin"]}>
      <CreateUsers />
    </ProtectedRoute>} />
  <Route path="student-progress" element={
    <ProtectedRoute requiredAnyRole={["admin","superadmin"]}>
      <StudentProgress />
    </ProtectedRoute>} />
</Route>
```

## 5. Backend Endpoints
| Method | Endpoint | Purpose | Auth |
|--------|----------|---------|------|
| POST | /api/users | Create user (student/advisor/chair/admin) | admin/superadmin |
| GET | /api/users?role=student | List students | admin/superadmin |
| GET | /api/admin/student-progress | Progress aggregation | admin/superadmin |
| GET | /api/admin/stats | Dashboard KPIs | admin/superadmin |

### 5.1 Student Progress Payload
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
- Count total nodes from playbook.
- Completed states: done|completed|approved.
- Latest activity: MAX(updated_at).

## 6. UI Components
- SidebarNav (Sheet mobile)
- AdminLayout (sidebar + topbar + Outlet)
- CreateAdmins form
- CreateUsers form
- StudentProgressTable (filters, progress bar)
- RoleBadge

## 7. State Management
Local state; later React Query optional.

## 8. Access Utilities (access.ts)
```ts
export const isSuperAdmin = (u?: User) => u?.role === 'superadmin';
export const isAdmin = (u?: User) => u && (u.role === 'admin' || u.role === 'superadmin');
export const requireAny = (u: User | null, roles: string[]) => !!u && roles.includes(u.role);
```

## 9. Milestones & Commits
| M | Scope | Commit |
|---|-------|--------|
| M1 | Layout + routing | feat(admin): add AdminLayout with sidebar and protected routes |
| M2 | CreateAdmins | feat(admin): implement superadmin admin creation form |
| M3 | CreateUsers | feat(admin): add user creation page (students/advisors/chairs) |
| M4 | Progress endpoint | feat(api): add student progress aggregation endpoint |
| M5 | Progress UI | feat(admin): student progress table with filters |
| M6 | Polish | feat(admin): role badges, toasts, loading states |

## 10. AI Prompts
Prompt A (M1):
“Implement AdminLayout with sidebar (shadcn). Desktop fixed w-64, mobile Sheet. Add nav items role-aware. Integrate ProtectedRoute. Add skeleton loading.”

Prompt B (M2):
“Create CreateAdmins.tsx: fields first_name, last_name, email. POST /api/users role=admin. Guard superadmin only. Toast on success.”

Prompt C (M3):
“Create CreateUsers.tsx: fields first_name, last_name, email, role select (student|advisor|chair). Validation + duplicate handling. Toast + clear form.”

Prompt D (M4):
“Add GET /api/admin/student-progress. Aggregate completion stats per student. Return payload spec. Guard admin/superadmin.”

Prompt E (M5):
“Build StudentProgress.tsx: fetch progress, table columns (Name, Email, Progress %, Current Node tooltip, Last Activity, Action). Client search filter.”

Prompt F (M6):
“Polish admin UI: RoleBadge, toasts, responsive sidebar, error boundary, aria labels.”

## 11. Risks
| Risk | Mitigation |
|------|------------|
| Dup role logic | Centralize in access.ts |
| Slow progress query | Indexes + limit pagination |
| Unauthorized flash | Delay render until auth ready |
| Stale data | Manual refresh button |

## 12. Success Criteria
- Correct RBAC on all routes.
- User/admin creation functional.
- Progress page loads <1s with sample data.
- Responsive layout mobile/desktop.
- Clear commit history by milestone.

## 13. Backlog
- CSV bulk import
- Progress export
- Impersonation
- Audit log
- Pagination & server filtering

---
Blueprint ready for execution.