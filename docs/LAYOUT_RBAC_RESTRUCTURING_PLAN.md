# Layout and RBAC Restructuring Plan

> **Status**: Draft  
> **Created**: January 8, 2026  
> **Last Updated**: January 8, 2026  
> **Version**: 1.0

## Executive Summary

This document outlines the comprehensive plan to restructure the application's layout architecture and Role-Based Access Control (RBAC) system to achieve best-in-class standards for universal educational applications. The restructuring addresses the current overloaded AdminLayout issue and introduces specialized interfaces for different user roles.

### Key Objectives

1. **Separation of Concerns** — Split monolithic AdminLayout into role-specific layouts
2. **Enhanced RBAC** — Implement industry-standard IMS Global LIS role vocabulary
3. **Multi-Role Support** — Enable users to hold multiple roles with seamless switching
4. **Configurable Workflows** — Create an approval workflow engine for multi-party processes
5. **Standards Compliance** — Complete LTI 1.3 integration and IMS Global compliance

---

## Table of Contents

1. [Current State Analysis](#current-state-analysis)
2. [Target Architecture](#target-architecture)
3. [Phase 1: RBAC Architecture](#phase-1-rbac-architecture)
4. [Phase 2: Workflow Engine](#phase-2-workflow-engine)
5. [Phase 3: Frontend Layouts](#phase-3-frontend-layouts)
6. [Phase 4: Standards Compliance](#phase-4-standards-compliance)
7. [Migration Strategy](#migration-strategy)
8. [Implementation Timeline](#implementation-timeline)

---

## Current State Analysis

### Problems Identified

| Issue                          | Impact                                  | Affected Files                                            |
| ------------------------------ | --------------------------------------- | --------------------------------------------------------- |
| Monolithic AdminLayout         | 30+ routes in single layout, poor UX    | `frontend/src/layouts/AdminLayout.tsx`                    |
| Teacher routes in AdminLayout  | Role confusion, navigation bloat        | `frontend/src/routes/index.tsx`                           |
| No separate TeacherLayout      | v11 reference has dedicated layout      | Missing vs `ui-examples/phd-journey-tracker_v11/teacher/` |
| advisor = teacher + supervisor | Role overloading                        | `backend/internal/models/user.go`                         |
| No role switching              | Multi-role users stuck in one interface | `frontend/src/contexts/AuthContext.tsx`                   |
| Incomplete LTI 1.3             | `StatusNotImplemented` in launch        | `backend/internal/handlers/lti.go`                        |

### Current Role Definitions

```go
// backend/internal/models/user.go
const (
    RoleSuperAdmin     Role = "superadmin"
    RoleAdmin          Role = "admin"
    RoleITAdmin        Role = "it_admin"
    RoleRegistrar      Role = "registrar"
    RoleContentMgr     Role = "content_manager"
    RoleStudentAffairs Role = "student_affairs"
    RoleStudent        Role = "student"
    RoleAdvisor        Role = "advisor"
    RoleSecretary      Role = "secretary"
    RoleChair          Role = "chair"
    RoleExternal       Role = "external"
)
```

### v11 Reference Architecture

The reference implementation in `ui-examples/phd-journey-tracker_v11/` demonstrates proper separation:

```
v11/
├── admin/layouts/
│   ├── StudioLayout.tsx      # Content creation
│   ├── OpsLayout.tsx         # Operations management
│   ├── ItemBankLayout.tsx    # Question banks
│   └── SchedulerLayout.tsx   # Scheduling
├── teacher/layouts/
│   └── TeacherLayout.tsx     # Instructor interface
└── student/layouts/
    └── StudentLayout.tsx     # Student portal
```

---

## Target Architecture

### Role Hierarchy

```
Platform Level:
├── superadmin              # Platform-wide administration

Tenant Level (Administrative):
├── admin                   # Tenant settings, integrations (narrowed scope)
├── hr_admin                # User lifecycle management
├── facility_manager        # Rooms, equipment, locations
├── scheduler_admin         # Schedule building and coordination

Academic Level:
├── dean                    # Program/course approval, final decisions
├── chair                   # Department oversight, analytics
├── committee_member        # Thesis defense committees

Teaching Level:
├── instructor              # Course teaching, grading, attendance
├── advisor                 # PhD supervision, thesis review, stage approval

Student Level:
├── student                 # Learning activities
├── external                # Read-only audit access
```

### Layout-to-Role Mapping

| Layout             | Route Prefix    | Roles            | Primary Functions                           |
| ------------------ | --------------- | ---------------- | ------------------------------------------- |
| `InstructorLayout` | `/teach/*`      | instructor       | Courses, Grading, Attendance                |
| `AdvisorLayout`    | `/advise/*`     | advisor          | PhD Students, Thesis Review, Stage Approval |
| `HRLayout`         | `/hr/*`         | hr_admin         | Users CRUD, Bulk Import, Access Requests    |
| `FacilityLayout`   | `/facilities/*` | facility_manager | Rooms, Equipment, Bookings                  |
| `SchedulerLayout`  | `/scheduling/*` | scheduler_admin  | Schedule Builder, Conflict Resolution       |
| `DeanLayout`       | `/academic/*`   | dean, chair      | Approval Queues, Department Analytics       |
| `AdminLayout`      | `/admin/*`      | admin            | Settings, Integrations, System Config       |
| `StudentLayout`    | `/student/*`    | student          | Learning Portal                             |
| `SuperadminLayout` | `/superadmin/*` | superadmin       | Tenant Management                           |

---

## Phase 1: RBAC Architecture

### 1.1 New Role Definitions

**File**: `backend/internal/models/user.go`

```go
// Add new roles
const (
    // ... existing roles ...

    // New roles for separation of concerns
    RoleInstructor       Role = "instructor"        // Course teaching
    RoleHRAdmin          Role = "hr_admin"          // User management
    RoleFacilityManager  Role = "facility_manager"  // Rooms/equipment
    RoleSchedulerAdmin   Role = "scheduler_admin"   // Scheduling
    RoleDean             Role = "dean"              // Academic approval
    RoleCommitteeMember  Role = "committee_member"  // Defense committees
)

// Role categories for UI grouping
var RoleCategories = map[string][]Role{
    "platform":     {RoleSuperAdmin},
    "admin":        {RoleAdmin, RoleITAdmin, RoleHRAdmin, RoleFacilityManager, RoleSchedulerAdmin},
    "academic":     {RoleDean, RoleChair, RoleCommitteeMember, RoleRegistrar},
    "teaching":     {RoleInstructor, RoleAdvisor},
    "student":      {RoleStudent, RoleExternal},
}
```

### 1.2 IMS Global LIS Role Mapping

**File**: `backend/internal/models/lis_roles.go` (new file)

```go
package models

// IMS Global LIS Role Vocabulary
// Reference: https://www.imsglobal.org/spec/lti/v1p3/#role-vocabularies

const (
    // System roles
    LISAdministrator = "http://purl.imsglobal.org/vocab/lis/v2/system/person#Administrator"
    LISSysAdmin      = "http://purl.imsglobal.org/vocab/lis/v2/system/person#SysAdmin"

    // Institution roles
    LISFaculty       = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Faculty"
    LISStaff         = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Staff"
    LISStudent       = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Student"
    LISAdvisor       = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Advisor"

    // Context (course) roles
    LISInstructor    = "http://purl.imsglobal.org/vocab/lis/v2/membership#Instructor"
    LISLearner       = "http://purl.imsglobal.org/vocab/lis/v2/membership#Learner"
    LISMentor        = "http://purl.imsglobal.org/vocab/lis/v2/membership#Mentor"
    LISContentDev    = "http://purl.imsglobal.org/vocab/lis/v2/membership#ContentDeveloper"
    LISManager       = "http://purl.imsglobal.org/vocab/lis/v2/membership#Manager"
    LISMember        = "http://purl.imsglobal.org/vocab/lis/v2/membership#Member"
    LISOfficer       = "http://purl.imsglobal.org/vocab/lis/v2/membership#Officer"
)

// InternalToLIS maps internal roles to IMS LIS URIs
var InternalToLIS = map[Role][]string{
    RoleSuperAdmin:      {LISAdministrator, LISSysAdmin},
    RoleAdmin:           {LISAdministrator},
    RoleHRAdmin:         {LISAdministrator, LISStaff},
    RoleDean:            {LISAdministrator, LISManager},
    RoleChair:           {LISManager, LISOfficer},
    RoleInstructor:      {LISInstructor, LISFaculty},
    RoleAdvisor:         {LISAdvisor, LISMentor, LISFaculty},
    RoleStudent:         {LISLearner, LISStudent},
    RoleExternal:        {LISMember},
    RoleRegistrar:       {LISStaff, LISManager},
    RoleContentMgr:      {LISContentDev},
    RoleFacilityManager: {LISStaff},
    RoleSchedulerAdmin:  {LISStaff, LISManager},
}

// LISToInternal maps LIS roles to internal roles (for LTI launches)
var LISToInternal = map[string]Role{
    LISAdministrator: RoleAdmin,
    LISInstructor:    RoleInstructor,
    LISLearner:       RoleStudent,
    LISMentor:        RoleAdvisor,
    LISContentDev:    RoleContentMgr,
    LISFaculty:       RoleInstructor,
    LISStudent:       RoleStudent,
    LISAdvisor:       RoleAdvisor,
}

// MapLISRolesToInternal converts LTI role claims to internal roles
func MapLISRolesToInternal(lisRoles []string) []Role {
    roleSet := make(map[Role]bool)
    for _, lis := range lisRoles {
        if internal, ok := LISToInternal[lis]; ok {
            roleSet[internal] = true
        }
    }

    roles := make([]Role, 0, len(roleSet))
    for role := range roleSet {
        roles = append(roles, role)
    }
    return roles
}
```

### 1.3 Multi-Role User Model

**File**: `backend/internal/models/user.go` (modifications)

```go
// User model update
type User struct {
    ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
    Email           string     `json:"email" gorm:"uniqueIndex:idx_user_email_tenant"`
    TenantID        *uuid.UUID `json:"tenant_id" gorm:"type:uuid;index"`

    // Legacy single role (for backward compatibility)
    Role            Role       `json:"role" gorm:"type:varchar(50)"`

    // Multi-role support
    Roles           []Role     `json:"roles" gorm:"-"` // Populated from user_roles table
    ActiveRole      Role       `json:"active_role" gorm:"-"` // Current session role

    // ... other fields
}

// UserRole represents a role assignment with context
type UserRole struct {
    ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
    UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;index"`
    Role        Role       `json:"role" gorm:"type:varchar(50)"`
    ContextType string     `json:"context_type"` // global, tenant, department, course
    ContextID   *uuid.UUID `json:"context_id" gorm:"type:uuid"`
    GrantedBy   uuid.UUID  `json:"granted_by" gorm:"type:uuid"`
    GrantedAt   time.Time  `json:"granted_at"`
    ExpiresAt   *time.Time `json:"expires_at"` // For temporary roles
    IsActive    bool       `json:"is_active" gorm:"default:true"`
}
```

### 1.4 Role Switching API

**File**: `backend/internal/handlers/auth.go` (additions)

```go
// SwitchRoleRequest represents the role switch payload
type SwitchRoleRequest struct {
    TargetRole  string `json:"target_role" binding:"required"`
    ContextType string `json:"context_type"` // optional: tenant, course, etc.
    ContextID   string `json:"context_id"`   // optional: specific context
}

// SwitchRoleResponse returns new token with switched role
type SwitchRoleResponse struct {
    Token       string   `json:"token"`
    ActiveRole  string   `json:"active_role"`
    AvailableRoles []string `json:"available_roles"`
    ExpiresAt   int64    `json:"expires_at"`
}

// SwitchRole handles POST /api/auth/switch-role
func (h *AuthHandler) SwitchRole(c *gin.Context) {
    userID := c.GetString("user_id")

    var req SwitchRoleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Verify user has the target role
    hasRole, err := h.authService.UserHasRole(c.Request.Context(), userID, req.TargetRole, req.ContextType, req.ContextID)
    if err != nil || !hasRole {
        c.JSON(http.StatusForbidden, gin.H{"error": "User does not have the requested role"})
        return
    }

    // Generate new token with switched role
    token, expiresAt, err := h.authService.GenerateTokenWithRole(c.Request.Context(), userID, req.TargetRole)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    // Get all available roles for response
    availableRoles, _ := h.authService.GetUserRoles(c.Request.Context(), userID)

    c.JSON(http.StatusOK, SwitchRoleResponse{
        Token:          token,
        ActiveRole:     req.TargetRole,
        AvailableRoles: availableRoles,
        ExpiresAt:      expiresAt,
    })
}
```

### 1.5 Database Migration

**File**: `backend/db/migrations/XXXX_add_new_roles_and_multi_role.up.sql`

```sql
-- Add new role values to enum (if using enum) or just document valid values
-- PostgreSQL approach for flexible roles:

-- 1. Create user_roles table for multi-role support
CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL,
    context_type VARCHAR(50) DEFAULT 'global', -- global, tenant, department, course
    context_id UUID, -- NULL for global context
    granted_by UUID REFERENCES users(id),
    granted_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP, -- NULL for permanent
    is_active BOOLEAN DEFAULT true,

    -- Prevent duplicate role assignments in same context
    UNIQUE(user_id, role, context_type, COALESCE(context_id, '00000000-0000-0000-0000-000000000000'))
);

CREATE INDEX idx_user_roles_user ON user_roles(user_id);
CREATE INDEX idx_user_roles_role ON user_roles(role);
CREATE INDEX idx_user_roles_context ON user_roles(context_type, context_id);

-- 2. Create permission_change_log for audit
CREATE TABLE IF NOT EXISTS permission_change_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    changed_by UUID REFERENCES users(id),
    change_type VARCHAR(30) NOT NULL, -- role_assigned, role_revoked, permission_granted, permission_revoked
    role VARCHAR(50),
    context_type VARCHAR(50),
    context_id UUID,
    previous_value JSONB, -- For tracking what changed
    new_value JSONB,
    reason TEXT, -- Optional reason for change
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_permission_log_user ON permission_change_log(user_id);
CREATE INDEX idx_permission_log_date ON permission_change_log(created_at);

-- 3. Migrate existing roles to user_roles table
INSERT INTO user_roles (user_id, role, context_type, granted_at)
SELECT id, role, 'global', created_at
FROM users
WHERE role IS NOT NULL
ON CONFLICT DO NOTHING;

-- 4. Add role_metadata table for role configuration
CREATE TABLE IF NOT EXISTS role_metadata (
    role VARCHAR(50) PRIMARY KEY,
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    category VARCHAR(50), -- platform, admin, academic, teaching, student
    default_landing_route VARCHAR(200),
    icon VARCHAR(50),
    is_system_role BOOLEAN DEFAULT false,
    can_be_self_assigned BOOLEAN DEFAULT false,
    requires_approval BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 5. Insert role metadata
INSERT INTO role_metadata (role, display_name, description, category, default_landing_route, is_system_role) VALUES
    ('superadmin', 'Super Administrator', 'Platform-wide administration', 'platform', '/superadmin', true),
    ('admin', 'Administrator', 'Tenant administration and settings', 'admin', '/admin', true),
    ('hr_admin', 'HR Administrator', 'User lifecycle management', 'admin', '/hr', false),
    ('facility_manager', 'Facility Manager', 'Rooms and equipment management', 'admin', '/facilities', false),
    ('scheduler_admin', 'Scheduler', 'Schedule building and coordination', 'admin', '/scheduling', false),
    ('dean', 'Dean', 'Academic program approval and oversight', 'academic', '/academic', false),
    ('chair', 'Department Chair', 'Department management and analytics', 'academic', '/academic', false),
    ('committee_member', 'Committee Member', 'Thesis defense committee participation', 'academic', '/academic', false),
    ('instructor', 'Instructor', 'Course teaching and grading', 'teaching', '/teach', false),
    ('advisor', 'Scientific Advisor', 'PhD supervision and thesis review', 'teaching', '/advise', false),
    ('student', 'Student', 'Learning activities', 'student', '/student', true),
    ('external', 'External Reviewer', 'Read-only audit access', 'student', '/external', false),
    ('registrar', 'Registrar', 'Academic records management', 'academic', '/admin/programs', false),
    ('content_manager', 'Content Manager', 'Learning content management', 'admin', '/admin/studio', false)
ON CONFLICT (role) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    description = EXCLUDED.description,
    category = EXCLUDED.category,
    default_landing_route = EXCLUDED.default_landing_route;

-- 6. Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
```

### 1.6 Updated JWT Token Structure

**File**: `backend/internal/services/auth.go` (modifications)

```go
// JWTClaims extended for multi-role support
type JWTClaims struct {
    jwt.RegisteredClaims
    UserID         string   `json:"user_id"`
    Email          string   `json:"email"`
    TenantID       string   `json:"tenant_id,omitempty"`

    // Role claims
    ActiveRole     string   `json:"active_role"`           // Currently active role
    AvailableRoles []string `json:"available_roles"`       // All assigned roles

    // Context claims (for context-specific tokens)
    ContextType    string   `json:"context_type,omitempty"` // course, department, etc.
    ContextID      string   `json:"context_id,omitempty"`

    // LIS claims for LTI compatibility
    LISRoles       []string `json:"lis_roles,omitempty"`    // IMS LIS role URIs

    // Legacy (backward compatibility)
    Role           string   `json:"role"`                   // Deprecated, use active_role
    IsSuperAdmin   bool     `json:"is_superadmin"`
}
```

### 1.7 Permission Matrix Update

**File**: `backend/internal/models/permissions.go` (new/updated)

```go
package models

// Permission categories
const (
    PermCategoryUser      = "user"
    PermCategoryCourse    = "course"
    PermCategoryProgram   = "program"
    PermCategoryGrade     = "grade"
    PermCategorySchedule  = "schedule"
    PermCategoryFacility  = "facility"
    PermCategoryWorkflow  = "workflow"
    PermCategoryAnalytics = "analytics"
    PermCategorySystem    = "system"
)

// Standard permissions
var StandardPermissions = []Permission{
    // User management
    {Slug: "user.view", Description: "View user profiles"},
    {Slug: "user.create", Description: "Create new users"},
    {Slug: "user.edit", Description: "Edit user profiles"},
    {Slug: "user.delete", Description: "Delete users"},
    {Slug: "user.assign_role", Description: "Assign roles to users"},
    {Slug: "user.bulk_import", Description: "Bulk import users"},

    // Course management
    {Slug: "course.view", Description: "View courses"},
    {Slug: "course.create", Description: "Create courses"},
    {Slug: "course.edit", Description: "Edit courses"},
    {Slug: "course.delete", Description: "Delete courses"},
    {Slug: "course.approve", Description: "Approve course proposals"},
    {Slug: "course.teach", Description: "Teach assigned courses"},

    // Grading
    {Slug: "grade.view", Description: "View grades"},
    {Slug: "grade.edit", Description: "Edit grades"},
    {Slug: "grade.approve", Description: "Approve grade changes"},
    {Slug: "grade.export", Description: "Export grades"},

    // Schedule management
    {Slug: "schedule.view", Description: "View schedules"},
    {Slug: "schedule.create", Description: "Create schedules"},
    {Slug: "schedule.edit", Description: "Edit schedules"},
    {Slug: "schedule.approve", Description: "Approve schedules"},

    // Facility management
    {Slug: "facility.view", Description: "View facilities"},
    {Slug: "facility.create", Description: "Create facilities"},
    {Slug: "facility.edit", Description: "Edit facilities"},
    {Slug: "facility.book", Description: "Book facilities"},
    {Slug: "facility.approve_booking", Description: "Approve facility bookings"},

    // Workflow management
    {Slug: "workflow.view", Description: "View workflows"},
    {Slug: "workflow.create", Description: "Create workflows"},
    {Slug: "workflow.approve", Description: "Approve workflow steps"},

    // Analytics
    {Slug: "analytics.view", Description: "View analytics"},
    {Slug: "analytics.export", Description: "Export analytics data"},

    // System
    {Slug: "system.settings", Description: "Manage system settings"},
    {Slug: "system.integrations", Description: "Manage integrations"},
    {Slug: "system.audit", Description: "View audit logs"},
}

// DefaultRolePermissions maps roles to their default permissions
var DefaultRolePermissions = map[Role][]string{
    RoleSuperAdmin: {"*"}, // All permissions

    RoleAdmin: {
        "system.settings", "system.integrations", "system.audit",
        "user.view", "analytics.view", "analytics.export",
    },

    RoleHRAdmin: {
        "user.view", "user.create", "user.edit", "user.delete",
        "user.assign_role", "user.bulk_import",
    },

    RoleFacilityManager: {
        "facility.view", "facility.create", "facility.edit",
        "facility.approve_booking",
    },

    RoleSchedulerAdmin: {
        "schedule.view", "schedule.create", "schedule.edit",
        "facility.view", "facility.book",
    },

    RoleDean: {
        "course.view", "course.approve", "program.view", "program.approve",
        "schedule.approve", "grade.approve", "workflow.approve",
        "analytics.view", "analytics.export",
    },

    RoleChair: {
        "course.view", "course.approve", "user.view",
        "analytics.view", "workflow.view",
    },

    RoleInstructor: {
        "course.view", "course.teach",
        "grade.view", "grade.edit",
        "schedule.view",
    },

    RoleAdvisor: {
        "course.view", "course.teach",
        "grade.view", "grade.edit",
        "workflow.view", "workflow.approve", // For thesis stages
    },

    RoleStudent: {
        "course.view", "grade.view", "schedule.view",
    },
}
```

---

## Phase 1 Implementation Checklist

### Backend Tasks

- [ ] **1.1** Add new role constants to `backend/internal/models/user.go`
- [ ] **1.2** Create `backend/internal/models/lis_roles.go` with IMS LIS mapping
- [ ] **1.3** Update User model with multi-role support
- [ ] **1.4** Create database migration for `user_roles`, `permission_change_log`, `role_metadata`
- [ ] **1.5** Run migration and verify data integrity
- [ ] **1.6** Update JWT claims structure in `backend/internal/services/auth.go`
- [ ] **1.7** Implement `SwitchRole` handler in `backend/internal/handlers/auth.go`
- [ ] **1.8** Add route `POST /api/auth/switch-role`
- [ ] **1.9** Update `AuthzService` to check combined permissions from all roles
- [ ] **1.10** Create `backend/internal/models/permissions.go` with permission matrix
- [ ] **1.11** Update login handler to return `available_roles` in response
- [ ] **1.12** Add permission change audit logging

### Frontend Tasks

- [ ] **1.13** Update `AuthContext` to store `activeRole` and `availableRoles`
- [ ] **1.14** Create `RoleSwitcher` component for header
- [ ] **1.15** Update `ProtectedRoute` to check against `activeRole`
- [ ] **1.16** Add API client method for `switchRole()`
- [ ] **1.17** Update login flow to handle multi-role response
- [ ] **1.18** Add role display in user menu

### Testing Tasks

- [ ] **1.19** Unit tests for role switching logic
- [ ] **1.20** Integration tests for multi-role authorization
- [ ] **1.21** Test backward compatibility with single-role users
- [ ] **1.22** Test LIS role mapping

---

## Phase 2: Configurable Workflow Engine

### 2.1 Overview

The workflow engine enables multi-party approval processes with configurable steps, timeouts, and escalation paths. It supports parallel and sequential approvals, automatic reminders, and delegation.

### 2.2 Core Concepts

```
Workflow Template
├── Steps[] (ordered)
│   ├── Step 1: Initiator submits
│   ├── Step 2: First approver
│   ├── Step 3: Second approver (parallel with Step 2?)
│   └── Step N: Final approver
├── Entity Type (course, schedule, thesis_stage, user_access)
└── Tenant-specific overrides

Workflow Instance
├── Template reference
├── Entity reference (the thing being approved)
├── Current step
├── Status (pending, approved, rejected, expired)
└── Approvals[] (decisions made)
```

### 2.3 Database Schema

**File**: `backend/db/migrations/XXXX_create_workflow_engine.up.sql`

```sql
-- Workflow Templates
CREATE TABLE workflow_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    entity_type VARCHAR(50) NOT NULL, -- course_approval, schedule_approval, thesis_stage, user_access
    is_active BOOLEAN DEFAULT true,
    is_system_template BOOLEAN DEFAULT false, -- Cannot be deleted
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    UNIQUE(tenant_id, name)
);

-- Workflow Steps (ordered within template)
CREATE TABLE workflow_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL REFERENCES workflow_templates(id) ON DELETE CASCADE,
    step_order INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,

    -- Who can approve this step
    required_role VARCHAR(50), -- Role required to approve
    required_permission VARCHAR(100), -- Or specific permission
    specific_user_id UUID REFERENCES users(id), -- Or specific user

    -- Step configuration
    is_optional BOOLEAN DEFAULT false,
    allow_delegation BOOLEAN DEFAULT true,
    parallel_with_previous BOOLEAN DEFAULT false, -- Can run in parallel with previous step

    -- Timeout handling
    timeout_days INT DEFAULT 7,
    auto_approve_on_timeout BOOLEAN DEFAULT false,
    auto_reject_on_timeout BOOLEAN DEFAULT false,
    escalation_role VARCHAR(50), -- Role to escalate to on timeout

    -- Notifications
    notify_on_pending BOOLEAN DEFAULT true,
    reminder_days INT DEFAULT 3, -- Days before timeout to send reminder

    created_at TIMESTAMP DEFAULT NOW(),

    UNIQUE(template_id, step_order)
);

-- Workflow Instances (actual running workflows)
CREATE TABLE workflow_instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL REFERENCES workflow_templates(id),
    tenant_id UUID REFERENCES tenants(id),

    -- What is being approved
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    entity_name VARCHAR(200), -- Cached for display

    -- Initiator
    initiated_by UUID NOT NULL REFERENCES users(id),
    initiated_at TIMESTAMP DEFAULT NOW(),

    -- Current state
    current_step_id UUID REFERENCES workflow_steps(id),
    current_step_order INT DEFAULT 1,
    status VARCHAR(30) DEFAULT 'pending', -- pending, approved, rejected, cancelled, expired

    -- Completion
    completed_at TIMESTAMP,
    final_decision VARCHAR(30), -- approved, rejected
    final_comment TEXT,

    -- Metadata
    metadata JSONB, -- Additional context data

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_workflow_instances_entity ON workflow_instances(entity_type, entity_id);
CREATE INDEX idx_workflow_instances_status ON workflow_instances(status);
CREATE INDEX idx_workflow_instances_pending ON workflow_instances(status, current_step_id) WHERE status = 'pending';

-- Workflow Approvals (decisions on each step)
CREATE TABLE workflow_approvals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instance_id UUID NOT NULL REFERENCES workflow_instances(id) ON DELETE CASCADE,
    step_id UUID NOT NULL REFERENCES workflow_steps(id),

    -- Who made the decision
    approver_id UUID REFERENCES users(id),
    approver_role VARCHAR(50),
    delegated_from UUID REFERENCES users(id), -- If delegated

    -- Decision
    decision VARCHAR(30) NOT NULL, -- approved, rejected, returned, delegated
    comment TEXT,

    -- Timing
    assigned_at TIMESTAMP DEFAULT NOW(),
    decided_at TIMESTAMP,
    due_at TIMESTAMP, -- Based on timeout_days

    -- Notifications sent
    notification_sent_at TIMESTAMP,
    reminder_sent_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_workflow_approvals_instance ON workflow_approvals(instance_id);
CREATE INDEX idx_workflow_approvals_pending ON workflow_approvals(approver_id, decided_at) WHERE decided_at IS NULL;

-- Workflow Delegation
CREATE TABLE workflow_delegations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    delegator_id UUID NOT NULL REFERENCES users(id),
    delegate_id UUID NOT NULL REFERENCES users(id),

    -- Scope of delegation
    workflow_type VARCHAR(50), -- NULL = all workflows
    role VARCHAR(50), -- Delegate acts as this role

    -- Validity period
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    reason TEXT,

    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),

    CHECK(end_date >= start_date)
);

CREATE INDEX idx_workflow_delegations_active ON workflow_delegations(delegate_id, is_active) WHERE is_active = true;

-- Insert default workflow templates
INSERT INTO workflow_templates (name, description, entity_type, is_system_template) VALUES
    ('Course Approval', 'Standard course approval workflow', 'course_approval', true),
    ('Schedule Approval', 'Multi-party schedule approval', 'schedule_approval', true),
    ('Thesis Stage Approval', 'PhD thesis stage progression', 'thesis_stage', true),
    ('User Access Request', 'Request for elevated access', 'user_access', true);

-- Insert steps for Course Approval template
WITH course_template AS (
    SELECT id FROM workflow_templates WHERE name = 'Course Approval' AND is_system_template = true
)
INSERT INTO workflow_steps (template_id, step_order, name, required_role, timeout_days, reminder_days) VALUES
    ((SELECT id FROM course_template), 1, 'Department Chair Review', 'chair', 5, 2),
    ((SELECT id FROM course_template), 2, 'Dean Approval', 'dean', 7, 3);

-- Insert steps for Schedule Approval template
WITH schedule_template AS (
    SELECT id FROM workflow_templates WHERE name = 'Schedule Approval' AND is_system_template = true
)
INSERT INTO workflow_steps (template_id, step_order, name, required_role, timeout_days, parallel_with_previous, reminder_days) VALUES
    ((SELECT id FROM schedule_template), 1, 'Instructor Confirmation', 'instructor', 3, false, 1),
    ((SELECT id FROM schedule_template), 2, 'Facility Availability', 'facility_manager', 3, true, 1),
    ((SELECT id FROM schedule_template), 3, 'Final Approval', 'scheduler_admin', 2, false, 1);

-- Insert steps for Thesis Stage Approval template
WITH thesis_template AS (
    SELECT id FROM workflow_templates WHERE name = 'Thesis Stage Approval' AND is_system_template = true
)
INSERT INTO workflow_steps (template_id, step_order, name, required_role, timeout_days, is_optional, reminder_days) VALUES
    ((SELECT id FROM thesis_template), 1, 'Advisor Review', 'advisor', 14, false, 5),
    ((SELECT id FROM thesis_template), 2, 'Committee Review', 'committee_member', 14, true, 5),
    ((SELECT id FROM thesis_template), 3, 'Chair Approval', 'chair', 7, false, 3);

-- Insert steps for User Access Request template
WITH access_template AS (
    SELECT id FROM workflow_templates WHERE name = 'User Access Request' AND is_system_template = true
)
INSERT INTO workflow_steps (template_id, step_order, name, required_role, timeout_days, reminder_days) VALUES
    ((SELECT id FROM access_template), 1, 'HR Review', 'hr_admin', 3, 1),
    ((SELECT id FROM access_template), 2, 'IT Approval', 'it_admin', 2, 1);
```

### 2.4 Workflow Models

**File**: `backend/internal/models/workflow.go` (new file)

```go
package models

import (
    "time"
    "github.com/google/uuid"
)

// WorkflowEntityType defines what can be approved
type WorkflowEntityType string

const (
    WorkflowEntityCourse      WorkflowEntityType = "course_approval"
    WorkflowEntitySchedule    WorkflowEntityType = "schedule_approval"
    WorkflowEntityThesisStage WorkflowEntityType = "thesis_stage"
    WorkflowEntityUserAccess  WorkflowEntityType = "user_access"
    WorkflowEntityGradeChange WorkflowEntityType = "grade_change"
    WorkflowEntityProgram     WorkflowEntityType = "program_approval"
)

// WorkflowStatus represents instance status
type WorkflowStatus string

const (
    WorkflowStatusPending   WorkflowStatus = "pending"
    WorkflowStatusApproved  WorkflowStatus = "approved"
    WorkflowStatusRejected  WorkflowStatus = "rejected"
    WorkflowStatusCancelled WorkflowStatus = "cancelled"
    WorkflowStatusExpired   WorkflowStatus = "expired"
)

// ApprovalDecision represents a step decision
type ApprovalDecision string

const (
    DecisionApproved  ApprovalDecision = "approved"
    DecisionRejected  ApprovalDecision = "rejected"
    DecisionReturned  ApprovalDecision = "returned"  // Sent back for changes
    DecisionDelegated ApprovalDecision = "delegated"
)

// WorkflowTemplate defines a reusable approval workflow
type WorkflowTemplate struct {
    ID               uuid.UUID          `json:"id" gorm:"type:uuid;primary_key"`
    TenantID         *uuid.UUID         `json:"tenant_id" gorm:"type:uuid"`
    Name             string             `json:"name"`
    Description      string             `json:"description"`
    EntityType       WorkflowEntityType `json:"entity_type"`
    IsActive         bool               `json:"is_active"`
    IsSystemTemplate bool               `json:"is_system_template"`
    CreatedBy        *uuid.UUID         `json:"created_by" gorm:"type:uuid"`
    Steps            []WorkflowStep     `json:"steps" gorm:"foreignKey:TemplateID"`
    CreatedAt        time.Time          `json:"created_at"`
    UpdatedAt        time.Time          `json:"updated_at"`
}

// WorkflowStep defines a single approval step
type WorkflowStep struct {
    ID                   uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
    TemplateID           uuid.UUID  `json:"template_id" gorm:"type:uuid"`
    StepOrder            int        `json:"step_order"`
    Name                 string     `json:"name"`
    Description          string     `json:"description"`
    RequiredRole         string     `json:"required_role"`
    RequiredPermission   string     `json:"required_permission"`
    SpecificUserID       *uuid.UUID `json:"specific_user_id" gorm:"type:uuid"`
    IsOptional           bool       `json:"is_optional"`
    AllowDelegation      bool       `json:"allow_delegation"`
    ParallelWithPrevious bool       `json:"parallel_with_previous"`
    TimeoutDays          int        `json:"timeout_days"`
    AutoApproveOnTimeout bool       `json:"auto_approve_on_timeout"`
    AutoRejectOnTimeout  bool       `json:"auto_reject_on_timeout"`
    EscalationRole       string     `json:"escalation_role"`
    NotifyOnPending      bool       `json:"notify_on_pending"`
    ReminderDays         int        `json:"reminder_days"`
    CreatedAt            time.Time  `json:"created_at"`
}

// WorkflowInstance represents a running workflow
type WorkflowInstance struct {
    ID               uuid.UUID         `json:"id" gorm:"type:uuid;primary_key"`
    TemplateID       uuid.UUID         `json:"template_id" gorm:"type:uuid"`
    Template         *WorkflowTemplate `json:"template" gorm:"foreignKey:TemplateID"`
    TenantID         *uuid.UUID        `json:"tenant_id" gorm:"type:uuid"`
    EntityType       string            `json:"entity_type"`
    EntityID         uuid.UUID         `json:"entity_id" gorm:"type:uuid"`
    EntityName       string            `json:"entity_name"`
    InitiatedBy      uuid.UUID         `json:"initiated_by" gorm:"type:uuid"`
    Initiator        *User             `json:"initiator" gorm:"foreignKey:InitiatedBy"`
    InitiatedAt      time.Time         `json:"initiated_at"`
    CurrentStepID    *uuid.UUID        `json:"current_step_id" gorm:"type:uuid"`
    CurrentStep      *WorkflowStep     `json:"current_step" gorm:"foreignKey:CurrentStepID"`
    CurrentStepOrder int               `json:"current_step_order"`
    Status           WorkflowStatus    `json:"status"`
    CompletedAt      *time.Time        `json:"completed_at"`
    FinalDecision    string            `json:"final_decision"`
    FinalComment     string            `json:"final_comment"`
    Metadata         JSON              `json:"metadata" gorm:"type:jsonb"`
    Approvals        []WorkflowApproval `json:"approvals" gorm:"foreignKey:InstanceID"`
    CreatedAt        time.Time         `json:"created_at"`
    UpdatedAt        time.Time         `json:"updated_at"`
}

// WorkflowApproval represents a decision on a step
type WorkflowApproval struct {
    ID                 uuid.UUID        `json:"id" gorm:"type:uuid;primary_key"`
    InstanceID         uuid.UUID        `json:"instance_id" gorm:"type:uuid"`
    StepID             uuid.UUID        `json:"step_id" gorm:"type:uuid"`
    Step               *WorkflowStep    `json:"step" gorm:"foreignKey:StepID"`
    ApproverID         *uuid.UUID       `json:"approver_id" gorm:"type:uuid"`
    Approver           *User            `json:"approver" gorm:"foreignKey:ApproverID"`
    ApproverRole       string           `json:"approver_role"`
    DelegatedFrom      *uuid.UUID       `json:"delegated_from" gorm:"type:uuid"`
    Decision           ApprovalDecision `json:"decision"`
    Comment            string           `json:"comment"`
    AssignedAt         time.Time        `json:"assigned_at"`
    DecidedAt          *time.Time       `json:"decided_at"`
    DueAt              *time.Time       `json:"due_at"`
    NotificationSentAt *time.Time       `json:"notification_sent_at"`
    ReminderSentAt     *time.Time       `json:"reminder_sent_at"`
    CreatedAt          time.Time        `json:"created_at"`
}

// WorkflowDelegation represents temporary delegation of approval authority
type WorkflowDelegation struct {
    ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
    DelegatorID  uuid.UUID  `json:"delegator_id" gorm:"type:uuid"`
    Delegator    *User      `json:"delegator" gorm:"foreignKey:DelegatorID"`
    DelegateID   uuid.UUID  `json:"delegate_id" gorm:"type:uuid"`
    Delegate     *User      `json:"delegate" gorm:"foreignKey:DelegateID"`
    WorkflowType string     `json:"workflow_type"` // Empty = all
    Role         string     `json:"role"`
    StartDate    time.Time  `json:"start_date"`
    EndDate      time.Time  `json:"end_date"`
    Reason       string     `json:"reason"`
    IsActive     bool       `json:"is_active"`
    CreatedAt    time.Time  `json:"created_at"`
}
```

### 2.5 Workflow Service Interface

**File**: `backend/internal/services/workflow.go` (new file)

```go
package services

import (
    "context"
    "github.com/google/uuid"
)

// WorkflowService defines workflow operations
type WorkflowService interface {
    // Template management
    GetTemplates(ctx context.Context, tenantID *uuid.UUID, entityType string) ([]models.WorkflowTemplate, error)
    GetTemplateByID(ctx context.Context, id uuid.UUID) (*models.WorkflowTemplate, error)
    CreateTemplate(ctx context.Context, template *models.WorkflowTemplate) error
    UpdateTemplate(ctx context.Context, template *models.WorkflowTemplate) error

    // Instance management
    StartWorkflow(ctx context.Context, req StartWorkflowRequest) (*models.WorkflowInstance, error)
    GetInstance(ctx context.Context, id uuid.UUID) (*models.WorkflowInstance, error)
    GetInstanceByEntity(ctx context.Context, entityType string, entityID uuid.UUID) (*models.WorkflowInstance, error)
    CancelWorkflow(ctx context.Context, instanceID uuid.UUID, reason string) error

    // Approval operations
    GetPendingApprovals(ctx context.Context, userID uuid.UUID, role string) ([]PendingApproval, error)
    GetPendingApprovalsCount(ctx context.Context, userID uuid.UUID, role string) (int, error)
    SubmitApproval(ctx context.Context, req ApprovalRequest) (*models.WorkflowApproval, error)
    DelegateApproval(ctx context.Context, approvalID uuid.UUID, delegateID uuid.UUID) error

    // Delegation management
    CreateDelegation(ctx context.Context, delegation *models.WorkflowDelegation) error
    GetActiveDelegations(ctx context.Context, userID uuid.UUID) ([]models.WorkflowDelegation, error)
    RevokeDelegation(ctx context.Context, delegationID uuid.UUID) error

    // Background jobs
    ProcessTimeouts(ctx context.Context) error
    SendReminders(ctx context.Context) error
}

// StartWorkflowRequest initiates a new workflow
type StartWorkflowRequest struct {
    TemplateID   uuid.UUID              `json:"template_id"`
    EntityType   string                 `json:"entity_type"`
    EntityID     uuid.UUID              `json:"entity_id"`
    EntityName   string                 `json:"entity_name"`
    InitiatedBy  uuid.UUID              `json:"initiated_by"`
    TenantID     *uuid.UUID             `json:"tenant_id"`
    Metadata     map[string]interface{} `json:"metadata"`
}

// ApprovalRequest submits a decision
type ApprovalRequest struct {
    InstanceID uuid.UUID `json:"instance_id"`
    StepID     uuid.UUID `json:"step_id"`
    ApproverID uuid.UUID `json:"approver_id"`
    Decision   string    `json:"decision" binding:"required,oneof=approved rejected returned"`
    Comment    string    `json:"comment"`
}

// PendingApproval represents an item awaiting approval
type PendingApproval struct {
    Instance     models.WorkflowInstance `json:"instance"`
    Step         models.WorkflowStep     `json:"step"`
    Approval     models.WorkflowApproval `json:"approval"`
    DaysUntilDue int                     `json:"days_until_due"`
    IsOverdue    bool                    `json:"is_overdue"`
}
```

### 2.6 Workflow API Handlers

**File**: `backend/internal/handlers/workflow.go` (new file)

```go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

// WorkflowHandler handles workflow API requests
type WorkflowHandler struct {
    workflowService services.WorkflowService
}

// RegisterWorkflowRoutes registers workflow API endpoints
func RegisterWorkflowRoutes(r *gin.RouterGroup, h *WorkflowHandler, authMiddleware, rbacMiddleware gin.HandlerFunc) {
    workflows := r.Group("/workflows")
    workflows.Use(authMiddleware)
    {
        // Templates (admin only)
        workflows.GET("/templates", h.ListTemplates)
        workflows.GET("/templates/:id", h.GetTemplate)
        workflows.POST("/templates", rbacMiddleware("workflow.create"), h.CreateTemplate)
        workflows.PUT("/templates/:id", rbacMiddleware("workflow.create"), h.UpdateTemplate)

        // Instances
        workflows.POST("/start", h.StartWorkflow)
        workflows.GET("/instances/:id", h.GetInstance)
        workflows.GET("/entity/:entityType/:entityId", h.GetInstanceByEntity)
        workflows.POST("/instances/:id/cancel", h.CancelWorkflow)

        // Approvals
        workflows.GET("/pending", h.GetPendingApprovals)
        workflows.GET("/pending/count", h.GetPendingApprovalsCount)
        workflows.POST("/approve", h.SubmitApproval)
        workflows.POST("/approvals/:id/delegate", h.DelegateApproval)

        // Delegations
        workflows.GET("/delegations", h.GetDelegations)
        workflows.POST("/delegations", h.CreateDelegation)
        workflows.DELETE("/delegations/:id", h.RevokeDelegation)
    }
}

// GetPendingApprovals returns items awaiting user's approval
// GET /api/workflows/pending?role=dean
func (h *WorkflowHandler) GetPendingApprovals(c *gin.Context) {
    userID := c.GetString("user_id")
    role := c.Query("role") // Optional: filter by role

    uid, err := uuid.Parse(userID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    approvals, err := h.workflowService.GetPendingApprovals(c.Request.Context(), uid, role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "pending_approvals": approvals,
        "count":             len(approvals),
    })
}

// SubmitApproval processes an approval decision
// POST /api/workflows/approve
func (h *WorkflowHandler) SubmitApproval(c *gin.Context) {
    userID := c.GetString("user_id")

    var req services.ApprovalRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    uid, _ := uuid.Parse(userID)
    req.ApproverID = uid

    approval, err := h.workflowService.SubmitApproval(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "approval": approval,
        "message":  "Approval submitted successfully",
    })
}

// StartWorkflow initiates a new workflow instance
// POST /api/workflows/start
func (h *WorkflowHandler) StartWorkflow(c *gin.Context) {
    userID := c.GetString("user_id")
    tenantID := c.GetString("tenant_id")

    var req services.StartWorkflowRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    uid, _ := uuid.Parse(userID)
    req.InitiatedBy = uid

    if tenantID != "" {
        tid, _ := uuid.Parse(tenantID)
        req.TenantID = &tid
    }

    instance, err := h.workflowService.StartWorkflow(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "instance": instance,
        "message":  "Workflow started successfully",
    })
}
```

### 2.7 Workflow Notification Service

**File**: `backend/internal/services/workflow_notifications.go`

```go
package services

// WorkflowNotificationService handles workflow-related notifications
type WorkflowNotificationService struct {
    emailService    EmailService
    pushService     PushNotificationService
    workflowService WorkflowService
}

// NotifyPendingApproval sends notification for new pending approval
func (s *WorkflowNotificationService) NotifyPendingApproval(ctx context.Context, approval *models.WorkflowApproval, instance *models.WorkflowInstance) error {
    // Get approver details
    // Send email notification
    // Send push notification if enabled
    // Update notification_sent_at
    return nil
}

// SendReminder sends reminder for approaching deadline
func (s *WorkflowNotificationService) SendReminder(ctx context.Context, approval *models.WorkflowApproval) error {
    // Check if reminder already sent
    // Send reminder email
    // Update reminder_sent_at
    return nil
}

// NotifyWorkflowComplete sends notification when workflow is complete
func (s *WorkflowNotificationService) NotifyWorkflowComplete(ctx context.Context, instance *models.WorkflowInstance) error {
    // Notify initiator
    // Notify all participants
    return nil
}
```

---

## Phase 2 Implementation Checklist

### Database Tasks

- [ ] **2.1** Create migration for workflow tables (templates, steps, instances, approvals, delegations)
- [ ] **2.2** Insert default workflow templates
- [ ] **2.3** Add indexes for performance
- [ ] **2.4** Run migration and verify schema

### Backend Tasks

- [ ] **2.5** Create `backend/internal/models/workflow.go`
- [ ] **2.6** Create `backend/internal/services/workflow.go` interface
- [ ] **2.7** Implement `WorkflowServiceImpl` with all methods
- [ ] **2.8** Create `backend/internal/handlers/workflow.go`
- [ ] **2.9** Register workflow routes in main router
- [ ] **2.10** Create `workflow_notifications.go` for email/push
- [ ] **2.11** Add background job for timeout processing
- [ ] **2.12** Add background job for reminder sending
- [ ] **2.13** Integrate workflow with course creation (auto-start workflow)
- [ ] **2.14** Integrate workflow with schedule creation

### Frontend Tasks

- [ ] **2.15** Create `useWorkflow` hook for API calls
- [ ] **2.16** Create `ApprovalQueue` component
- [ ] **2.17** Create `ApprovalDetail` modal/page
- [ ] **2.18** Create `WorkflowTimeline` component (shows approval history)
- [ ] **2.19** Add approval count badge to relevant layouts
- [ ] **2.20** Create `DelegationManager` component
- [ ] **2.21** Add workflow status indicator on entity cards/lists

### Testing Tasks

- [ ] **2.22** Unit tests for workflow service
- [ ] **2.23** Integration tests for approval flow
- [ ] **2.24** Test parallel approval steps
- [ ] **2.25** Test timeout and escalation
- [ ] **2.26** Test delegation scenarios

---

## Phase 3: Frontend Layouts

### 3.1 Layout Architecture Overview

Each role gets a dedicated layout with role-specific navigation, widgets, and workflows. Layouts share common components (header, role switcher, notifications) but have distinct sidebar navigation and landing pages.

```
frontend/src/layouts/
├── AdminLayout.tsx          # Settings, Integrations (narrowed)
├── StudentLayout.tsx        # Learning portal (existing)
├── SuperadminLayout.tsx     # Platform admin (existing)
├── InstructorLayout.tsx     # NEW: Course teaching
├── AdvisorLayout.tsx        # NEW: PhD supervision
├── HRLayout.tsx             # NEW: User management
├── FacilityLayout.tsx       # NEW: Room/equipment management
├── SchedulerLayout.tsx      # NEW: Schedule management
├── DeanLayout.tsx           # NEW: Academic approvals
└── shared/
    ├── LayoutShell.tsx      # Common layout wrapper
    ├── RoleSwitcher.tsx     # Role switching component
    ├── ApprovalBadge.tsx    # Pending approvals indicator
    └── NavigationItem.tsx   # Reusable nav item
```

### 3.2 Shared Layout Shell

**File**: `frontend/src/layouts/shared/LayoutShell.tsx`

```tsx
import React from "react";
import { Outlet } from "react-router-dom";
import { RoleSwitcher } from "./RoleSwitcher";
import { ApprovalBadge } from "./ApprovalBadge";
import { GlobalSearch } from "@/components/GlobalSearch";
import { UserMenu } from "@/components/layout/UserMenu";
import { LanguageSwitcher } from "@/components/layout/LanguageSwitcher";
import { cn } from "@/lib/utils";

interface LayoutShellProps {
  title: string;
  icon: React.ReactNode;
  sidebar: React.ReactNode;
  className?: string;
}

export function LayoutShell({
  title,
  icon,
  sidebar,
  className,
}: LayoutShellProps) {
  const [collapsed, setCollapsed] = React.useState(false);

  return (
    <div className="flex min-h-screen">
      {/* Sidebar */}
      <aside
        className={cn(
          "hidden md:flex flex-col border-r bg-background transition-all duration-200",
          collapsed ? "w-16" : "w-64"
        )}
      >
        {/* Logo/Title */}
        <div
          className={cn(
            "h-14 flex items-center border-b px-4",
            collapsed && "justify-center px-2"
          )}
        >
          <div className="flex items-center gap-2">
            {icon}
            {!collapsed && <span className="font-semibold">{title}</span>}
          </div>
        </div>

        {/* Navigation */}
        <nav className="flex-1 overflow-y-auto py-4">{sidebar}</nav>

        {/* Collapse toggle */}
        <div className="p-2 border-t">
          <button
            onClick={() => setCollapsed(!collapsed)}
            className="w-full p-2 rounded hover:bg-muted text-muted-foreground"
          >
            {collapsed ? "»" : "«"}
          </button>
        </div>
      </aside>

      {/* Main content */}
      <div className="flex-1 flex flex-col">
        {/* Top bar */}
        <header className="h-14 flex items-center justify-between border-b px-6 bg-background/95 backdrop-blur sticky top-0 z-50">
          <div className="flex items-center gap-4">
            <GlobalSearch />
          </div>
          <div className="flex items-center gap-3">
            <ApprovalBadge />
            <RoleSwitcher />
            <LanguageSwitcher />
            <UserMenu />
          </div>
        </header>

        {/* Page content */}
        <main className={cn("flex-1 p-6", className)}>
          <Outlet />
        </main>
      </div>
    </div>
  );
}
```

### 3.3 Role Switcher Component

**File**: `frontend/src/layouts/shared/RoleSwitcher.tsx`

```tsx
import React from "react";
import { useAuth } from "@/contexts/AuthContext";
import { useTranslation } from "react-i18next";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ChevronDown, Check } from "lucide-react";
import { api } from "@/api/client";

const roleIcons: Record<string, string> = {
  instructor: "👨‍🏫",
  advisor: "🎓",
  hr_admin: "👥",
  facility_manager: "🏢",
  scheduler_admin: "📅",
  dean: "🏛️",
  chair: "📊",
  admin: "⚙️",
  student: "📚",
};

const roleLabels: Record<string, string> = {
  instructor: "Instructor",
  advisor: "Scientific Advisor",
  hr_admin: "HR Admin",
  facility_manager: "Facility Manager",
  scheduler_admin: "Scheduler",
  dean: "Dean",
  chair: "Department Chair",
  admin: "Administrator",
  student: "Student",
};

export function RoleSwitcher() {
  const { t } = useTranslation("common");
  const { user, activeRole, availableRoles, setActiveRole } = useAuth();
  const [switching, setSwitching] = React.useState(false);

  if (!availableRoles || availableRoles.length <= 1) {
    // Single role user - show badge only
    return (
      <Badge variant="secondary" className="px-3 py-1">
        {roleIcons[activeRole]} {roleLabels[activeRole] || activeRole}
      </Badge>
    );
  }

  const handleSwitchRole = async (targetRole: string) => {
    if (targetRole === activeRole) return;

    setSwitching(true);
    try {
      const response = await api("/auth/switch-role", {
        method: "POST",
        body: JSON.stringify({ target_role: targetRole }),
      });

      // Update auth context with new token
      setActiveRole(targetRole, response.token);

      // Redirect to role's default landing page
      const landingRoutes: Record<string, string> = {
        instructor: "/teach",
        advisor: "/advise",
        hr_admin: "/hr",
        facility_manager: "/facilities",
        scheduler_admin: "/scheduling",
        dean: "/academic",
        chair: "/academic",
        admin: "/admin",
        student: "/student",
      };

      window.location.href = landingRoutes[targetRole] || "/";
    } catch (error) {
      console.error("Failed to switch role:", error);
    } finally {
      setSwitching(false);
    }
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="sm" disabled={switching}>
          {roleIcons[activeRole]} {roleLabels[activeRole] || activeRole}
          <ChevronDown className="ml-2 h-4 w-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-56">
        <DropdownMenuLabel>{t("Switch Role")}</DropdownMenuLabel>
        <DropdownMenuSeparator />
        {availableRoles.map((role) => (
          <DropdownMenuItem
            key={role}
            onClick={() => handleSwitchRole(role)}
            className="flex items-center justify-between"
          >
            <span>
              {roleIcons[role]} {roleLabels[role] || role}
            </span>
            {role === activeRole && <Check className="h-4 w-4" />}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
```

### 3.4 Instructor Layout

**File**: `frontend/src/layouts/InstructorLayout.tsx`

```tsx
import React from "react";
import { NavLink, useLocation } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { LayoutShell } from "./shared/LayoutShell";
import { cn } from "@/lib/utils";
import {
  LayoutDashboard,
  BookOpen,
  ClipboardCheck,
  Users,
  Calendar,
  MessageSquare,
  GraduationCap,
} from "lucide-react";

const navItems = [
  { icon: LayoutDashboard, label: "Overview", path: "/teach" },
  { icon: BookOpen, label: "My Courses", path: "/teach/courses" },
  {
    icon: ClipboardCheck,
    label: "Grading Hub",
    path: "/teach/grading",
    badge: true,
  },
  { icon: Users, label: "Students", path: "/teach/students" },
  { icon: Calendar, label: "Schedule", path: "/teach/schedule" },
  { icon: MessageSquare, label: "Messages", path: "/teach/messages" },
];

function InstructorNav() {
  const { t } = useTranslation("common");
  const location = useLocation();

  return (
    <div className="space-y-1 px-3">
      <div className="px-3 py-2 text-xs font-semibold text-muted-foreground uppercase tracking-wider">
        Workspace
      </div>
      {navItems.map((item) => (
        <NavLink
          key={item.path}
          to={item.path}
          className={cn(
            "flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors",
            location.pathname === item.path
              ? "bg-primary/10 text-primary"
              : "text-muted-foreground hover:bg-muted hover:text-foreground"
          )}
        >
          <item.icon className="h-4 w-4" />
          <span>
            {t(`instructor.nav.${item.label.toLowerCase()}`, item.label)}
          </span>
          {item.badge && (
            <span className="ml-auto bg-red-100 text-red-600 text-xs px-2 py-0.5 rounded-full">
              5
            </span>
          )}
        </NavLink>
      ))}
    </div>
  );
}

export function InstructorLayout() {
  return (
    <LayoutShell
      title="Instructor"
      icon={<GraduationCap className="h-5 w-5 text-primary" />}
      sidebar={<InstructorNav />}
    />
  );
}

export default InstructorLayout;
```

### 3.5 Advisor Layout

**File**: `frontend/src/layouts/AdvisorLayout.tsx`

```tsx
import React from "react";
import { NavLink, useLocation } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { LayoutShell } from "./shared/LayoutShell";
import { cn } from "@/lib/utils";
import {
  LayoutDashboard,
  Users,
  FileText,
  CheckSquare,
  Route,
  MessageSquare,
  Calendar,
  Lightbulb,
} from "lucide-react";

const navItems = [
  { icon: LayoutDashboard, label: "Overview", path: "/advise" },
  { icon: Users, label: "My Students", path: "/advise/students" },
  {
    icon: FileText,
    label: "Document Review",
    path: "/advise/documents",
    badge: true,
  },
  {
    icon: CheckSquare,
    label: "Stage Approvals",
    path: "/advise/approvals",
    badge: true,
  },
  { icon: Route, label: "Thesis Progress", path: "/advise/thesis-tracker" },
  { icon: Calendar, label: "Committee Schedule", path: "/advise/committees" },
  { icon: MessageSquare, label: "Messages", path: "/advise/messages" },
];

function AdvisorNav() {
  const { t } = useTranslation("common");
  const location = useLocation();

  return (
    <div className="space-y-1 px-3">
      <div className="px-3 py-2 text-xs font-semibold text-muted-foreground uppercase tracking-wider">
        Supervision
      </div>
      {navItems.map((item) => (
        <NavLink
          key={item.path}
          to={item.path}
          className={cn(
            "flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors",
            location.pathname === item.path
              ? "bg-primary/10 text-primary"
              : "text-muted-foreground hover:bg-muted hover:text-foreground"
          )}
        >
          <item.icon className="h-4 w-4" />
          <span>
            {t(
              `advisor.nav.${item.label.toLowerCase().replace(/ /g, "_")}`,
              item.label
            )}
          </span>
        </NavLink>
      ))}
    </div>
  );
}

export function AdvisorLayout() {
  return (
    <LayoutShell
      title="Scientific Advisor"
      icon={<Lightbulb className="h-5 w-5 text-amber-500" />}
      sidebar={<AdvisorNav />}
    />
  );
}

export default AdvisorLayout;
```

### 3.6 HR Layout

**File**: `frontend/src/layouts/HRLayout.tsx`

```tsx
import React from "react";
import { NavLink, useLocation } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { LayoutShell } from "./shared/LayoutShell";
import { cn } from "@/lib/utils";
import {
  LayoutDashboard,
  Users,
  UserPlus,
  Upload,
  Shield,
  FileCheck,
  History,
  UserCog,
} from "lucide-react";

const navItems = [
  { icon: LayoutDashboard, label: "Overview", path: "/hr" },
  { icon: Users, label: "All Users", path: "/hr/users" },
  { icon: UserPlus, label: "Add User", path: "/hr/users/new" },
  { icon: Upload, label: "Bulk Import", path: "/hr/import" },
  { icon: Shield, label: "Role Management", path: "/hr/roles" },
  {
    icon: FileCheck,
    label: "Access Requests",
    path: "/hr/access-requests",
    badge: true,
  },
  { icon: History, label: "Audit Log", path: "/hr/audit" },
];

function HRNav() {
  const { t } = useTranslation("common");
  const location = useLocation();

  return (
    <div className="space-y-1 px-3">
      <div className="px-3 py-2 text-xs font-semibold text-muted-foreground uppercase tracking-wider">
        User Management
      </div>
      {navItems.map((item) => (
        <NavLink
          key={item.path}
          to={item.path}
          className={cn(
            "flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors",
            location.pathname === item.path
              ? "bg-primary/10 text-primary"
              : "text-muted-foreground hover:bg-muted hover:text-foreground"
          )}
        >
          <item.icon className="h-4 w-4" />
          <span>
            {t(
              `hr.nav.${item.label.toLowerCase().replace(/ /g, "_")}`,
              item.label
            )}
          </span>
        </NavLink>
      ))}
    </div>
  );
}

export function HRLayout() {
  return (
    <LayoutShell
      title="HR Administration"
      icon={<UserCog className="h-5 w-5 text-blue-500" />}
      sidebar={<HRNav />}
    />
  );
}

export default HRLayout;
```

### 3.7 Facility Layout

**File**: `frontend/src/layouts/FacilityLayout.tsx`

```tsx
import React from "react";
import { NavLink, useLocation } from "react-router-dom";
import { LayoutShell } from "./shared/LayoutShell";
import { cn } from "@/lib/utils";
import {
  LayoutDashboard,
  Building,
  DoorOpen,
  Laptop,
  Calendar,
  Wrench,
  BarChart3,
  Building2,
} from "lucide-react";

const navItems = [
  { icon: LayoutDashboard, label: "Overview", path: "/facilities" },
  { icon: Building, label: "Buildings", path: "/facilities/buildings" },
  { icon: DoorOpen, label: "Rooms", path: "/facilities/rooms" },
  { icon: Laptop, label: "Equipment", path: "/facilities/equipment" },
  {
    icon: Calendar,
    label: "Bookings",
    path: "/facilities/bookings",
    badge: true,
  },
  { icon: Wrench, label: "Maintenance", path: "/facilities/maintenance" },
  { icon: BarChart3, label: "Utilization", path: "/facilities/reports" },
];

function FacilityNav() {
  const location = useLocation();

  return (
    <div className="space-y-1 px-3">
      <div className="px-3 py-2 text-xs font-semibold text-muted-foreground uppercase tracking-wider">
        Facility Management
      </div>
      {navItems.map((item) => (
        <NavLink
          key={item.path}
          to={item.path}
          className={cn(
            "flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors",
            location.pathname === item.path
              ? "bg-primary/10 text-primary"
              : "text-muted-foreground hover:bg-muted hover:text-foreground"
          )}
        >
          <item.icon className="h-4 w-4" />
          <span>{item.label}</span>
        </NavLink>
      ))}
    </div>
  );
}

export function FacilityLayout() {
  return (
    <LayoutShell
      title="Facilities"
      icon={<Building2 className="h-5 w-5 text-emerald-500" />}
      sidebar={<FacilityNav />}
    />
  );
}

export default FacilityLayout;
```

### 3.8 Scheduler Layout

**File**: `frontend/src/layouts/SchedulerLayout.tsx`

```tsx
import React from "react";
import { NavLink, useLocation } from "react-router-dom";
import { LayoutShell } from "./shared/LayoutShell";
import { cn } from "@/lib/utils";
import {
  LayoutDashboard,
  Calendar,
  CalendarPlus,
  AlertTriangle,
  CheckCircle,
  Users,
  Clock,
  CalendarClock,
} from "lucide-react";

const navItems = [
  { icon: LayoutDashboard, label: "Overview", path: "/scheduling" },
  { icon: Calendar, label: "Master Schedule", path: "/scheduling/calendar" },
  { icon: CalendarPlus, label: "Create Schedule", path: "/scheduling/create" },
  {
    icon: AlertTriangle,
    label: "Conflicts",
    path: "/scheduling/conflicts",
    badge: true,
  },
  {
    icon: CheckCircle,
    label: "Pending Approvals",
    path: "/scheduling/approvals",
    badge: true,
  },
  {
    icon: Users,
    label: "Instructor Load",
    path: "/scheduling/instructor-load",
  },
  { icon: Clock, label: "Time Slots", path: "/scheduling/time-slots" },
];

function SchedulerNav() {
  const location = useLocation();

  return (
    <div className="space-y-1 px-3">
      <div className="px-3 py-2 text-xs font-semibold text-muted-foreground uppercase tracking-wider">
        Scheduling
      </div>
      {navItems.map((item) => (
        <NavLink
          key={item.path}
          to={item.path}
          className={cn(
            "flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors",
            location.pathname === item.path
              ? "bg-primary/10 text-primary"
              : "text-muted-foreground hover:bg-muted hover:text-foreground"
          )}
        >
          <item.icon className="h-4 w-4" />
          <span>{item.label}</span>
        </NavLink>
      ))}
    </div>
  );
}

export function SchedulerLayout() {
  return (
    <LayoutShell
      title="Scheduler"
      icon={<CalendarClock className="h-5 w-5 text-purple-500" />}
      sidebar={<SchedulerNav />}
    />
  );
}

export default SchedulerLayout;
```

### 3.9 Dean/Academic Layout

**File**: `frontend/src/layouts/DeanLayout.tsx`

```tsx
import React from "react";
import { NavLink, useLocation } from "react-router-dom";
import { LayoutShell } from "./shared/LayoutShell";
import { cn } from "@/lib/utils";
import {
  LayoutDashboard,
  FileCheck,
  BookOpen,
  GraduationCap,
  Users,
  BarChart3,
  Settings,
  Landmark,
} from "lucide-react";

const navItems = [
  { icon: LayoutDashboard, label: "Overview", path: "/academic" },
  {
    icon: FileCheck,
    label: "Approval Queue",
    path: "/academic/approvals",
    badge: true,
  },
  { icon: BookOpen, label: "Courses", path: "/academic/courses" },
  { icon: GraduationCap, label: "Programs", path: "/academic/programs" },
  { icon: Users, label: "Committees", path: "/academic/committees" },
  { icon: BarChart3, label: "Analytics", path: "/academic/analytics" },
  { icon: Settings, label: "Policies", path: "/academic/policies" },
];

function DeanNav() {
  const location = useLocation();

  return (
    <div className="space-y-1 px-3">
      <div className="px-3 py-2 text-xs font-semibold text-muted-foreground uppercase tracking-wider">
        Academic Affairs
      </div>
      {navItems.map((item) => (
        <NavLink
          key={item.path}
          to={item.path}
          className={cn(
            "flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors",
            location.pathname === item.path
              ? "bg-primary/10 text-primary"
              : "text-muted-foreground hover:bg-muted hover:text-foreground"
          )}
        >
          <item.icon className="h-4 w-4" />
          <span>{item.label}</span>
        </NavLink>
      ))}
    </div>
  );
}

export function DeanLayout() {
  return (
    <LayoutShell
      title="Academic Affairs"
      icon={<Landmark className="h-5 w-5 text-indigo-500" />}
      sidebar={<DeanNav />}
    />
  );
}

export default DeanLayout;
```

### 3.10 Updated Routes Configuration

**File**: `frontend/src/routes/index.tsx` (modifications)

```tsx
// Add new layout imports
const InstructorLayout = lazy(() =>
  import("@/layouts/InstructorLayout").then((m) => ({
    default: m.InstructorLayout,
  }))
);
const AdvisorLayout = lazy(() =>
  import("@/layouts/AdvisorLayout").then((m) => ({ default: m.AdvisorLayout }))
);
const HRLayout = lazy(() =>
  import("@/layouts/HRLayout").then((m) => ({ default: m.HRLayout }))
);
const FacilityLayout = lazy(() =>
  import("@/layouts/FacilityLayout").then((m) => ({
    default: m.FacilityLayout,
  }))
);
const SchedulerLayout = lazy(() =>
  import("@/layouts/SchedulerLayout").then((m) => ({
    default: m.SchedulerLayout,
  }))
);
const DeanLayout = lazy(() =>
  import("@/layouts/DeanLayout").then((m) => ({ default: m.DeanLayout }))
);

// Add new route groups
export const router = createBrowserRouter([
  // ... existing routes ...

  // Instructor routes
  {
    path: "/teach",
    element: (
      <ProtectedRoute requiredAnyRole={["instructor", "advisor"]}>
        {WithSuspense(<InstructorLayout />)}
      </ProtectedRoute>
    ),
    errorElement: <RouteErrorBoundary />,
    children: [
      { index: true, element: WithSuspense(<TeacherDashboard />) },
      { path: "courses", element: WithSuspense(<TeacherCoursesPage />) },
      {
        path: "courses/:courseId",
        element: WithSuspense(<TeacherCourseDetail />),
      },
      { path: "grading", element: WithSuspense(<TeacherGradingPage />) },
      { path: "students", element: WithSuspense(<StudentTracker />) },
      { path: "schedule", element: WithSuspense(<CalendarPage />) },
      { path: "messages", element: WithSuspense(<ChatPage />) },
    ],
  },

  // Advisor routes
  {
    path: "/advise",
    element: (
      <ProtectedRoute requiredAnyRole={["advisor"]}>
        {WithSuspense(<AdvisorLayout />)}
      </ProtectedRoute>
    ),
    errorElement: <RouteErrorBoundary />,
    children: [
      { index: true, element: WithSuspense(<AdvisorDashboard />) },
      { path: "students", element: WithSuspense(<StudentsMonitorPage />) },
      { path: "students/:id", element: WithSuspense(<StudentDetailPage />) },
      { path: "documents", element: WithSuspense(<AdvisorInbox />) },
      { path: "approvals", element: WithSuspense(<ApprovalQueuePage />) },
      { path: "thesis-tracker", element: WithSuspense(<ThesisTrackerPage />) },
      { path: "committees", element: WithSuspense(<CommitteeSchedulePage />) },
      { path: "messages", element: WithSuspense(<ChatPage />) },
    ],
  },

  // HR routes
  {
    path: "/hr",
    element: (
      <ProtectedRoute requiredAnyRole={["hr_admin"]}>
        {WithSuspense(<HRLayout />)}
      </ProtectedRoute>
    ),
    errorElement: <RouteErrorBoundary />,
    children: [
      { index: true, element: WithSuspense(<HRDashboard />) },
      { path: "users", element: WithSuspense(<AdminUsersPage />) },
      { path: "users/new", element: WithSuspense(<CreateUsers />) },
      { path: "import", element: WithSuspense(<BulkImportPage />) },
      { path: "roles", element: WithSuspense(<RoleManagementPage />) },
      {
        path: "access-requests",
        element: WithSuspense(<AccessRequestsPage />),
      },
      { path: "audit", element: WithSuspense(<AuditLogPage />) },
    ],
  },

  // Facility routes
  {
    path: "/facilities",
    element: (
      <ProtectedRoute requiredAnyRole={["facility_manager"]}>
        {WithSuspense(<FacilityLayout />)}
      </ProtectedRoute>
    ),
    errorElement: <RouteErrorBoundary />,
    children: [
      { index: true, element: WithSuspense(<FacilityDashboard />) },
      { path: "buildings", element: WithSuspense(<BuildingsPage />) },
      { path: "rooms", element: WithSuspense(<RoomsPage />) },
      { path: "equipment", element: WithSuspense(<EquipmentPage />) },
      { path: "bookings", element: WithSuspense(<BookingsPage />) },
      { path: "maintenance", element: WithSuspense(<MaintenancePage />) },
      { path: "reports", element: WithSuspense(<UtilizationReportsPage />) },
    ],
  },

  // Scheduler routes
  {
    path: "/scheduling",
    element: (
      <ProtectedRoute requiredAnyRole={["scheduler_admin"]}>
        {WithSuspense(<SchedulerLayout />)}
      </ProtectedRoute>
    ),
    errorElement: <RouteErrorBoundary />,
    children: [
      { index: true, element: WithSuspense(<SchedulerDashboard />) },
      { path: "calendar", element: WithSuspense(<MasterSchedulePage />) },
      { path: "create", element: WithSuspense(<CreateSchedulePage />) },
      { path: "conflicts", element: WithSuspense(<ConflictsPage />) },
      { path: "approvals", element: WithSuspense(<ScheduleApprovalsPage />) },
      {
        path: "instructor-load",
        element: WithSuspense(<InstructorLoadPage />),
      },
      { path: "time-slots", element: WithSuspense(<TimeSlotsPage />) },
    ],
  },

  // Academic/Dean routes
  {
    path: "/academic",
    element: (
      <ProtectedRoute requiredAnyRole={["dean", "chair"]}>
        {WithSuspense(<DeanLayout />)}
      </ProtectedRoute>
    ),
    errorElement: <RouteErrorBoundary />,
    children: [
      { index: true, element: WithSuspense(<AcademicDashboard />) },
      { path: "approvals", element: WithSuspense(<ApprovalQueuePage />) },
      { path: "courses", element: WithSuspense(<CoursesPage />) },
      { path: "programs", element: WithSuspense(<ProgramsPage />) },
      { path: "committees", element: WithSuspense(<CommitteesPage />) },
      { path: "analytics", element: WithSuspense(<AnalyticsDashboard />) },
      { path: "policies", element: WithSuspense(<PoliciesPage />) },
    ],
  },
]);
```

### 3.11 Updated AuthContext

**File**: `frontend/src/contexts/AuthContext.tsx` (modifications)

```tsx
interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;

  // Multi-role support
  activeRole: string;
  availableRoles: string[];
  setActiveRole: (role: string, token?: string) => void;

  // Convenience methods
  hasRole: (role: string) => boolean;
  hasAnyRole: (roles: string[]) => boolean;
  canAccessRoute: (requiredRoles: string[]) => boolean;
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [activeRole, setActiveRoleState] = useState<string>("");
  const [availableRoles, setAvailableRoles] = useState<string[]>([]);

  // Initialize from token
  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      try {
        const decoded = jwtDecode<JWTClaims>(token);
        setUser(decoded.user);
        setActiveRoleState(decoded.active_role || decoded.role);
        setAvailableRoles(decoded.available_roles || [decoded.role]);
      } catch {
        localStorage.removeItem("token");
      }
    }
    setIsLoading(false);
  }, []);

  const setActiveRole = (role: string, token?: string) => {
    if (token) {
      localStorage.setItem("token", token);
    }
    setActiveRoleState(role);
  };

  const hasRole = (role: string) => availableRoles.includes(role);

  const hasAnyRole = (roles: string[]) =>
    roles.some((role) => availableRoles.includes(role));

  const canAccessRoute = (requiredRoles: string[]) => {
    // Admin and superadmin can access everything
    if (hasRole("admin") || hasRole("superadmin")) return true;
    return hasAnyRole(requiredRoles);
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        login,
        logout,
        activeRole,
        availableRoles,
        setActiveRole,
        hasRole,
        hasAnyRole,
        canAccessRoute,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}
```

---

## Phase 3 Implementation Checklist

### Layout Files

- [ ] **3.1** Create `frontend/src/layouts/shared/LayoutShell.tsx`
- [ ] **3.2** Create `frontend/src/layouts/shared/RoleSwitcher.tsx`
- [ ] **3.3** Create `frontend/src/layouts/shared/ApprovalBadge.tsx`
- [ ] **3.4** Create `frontend/src/layouts/InstructorLayout.tsx`
- [ ] **3.5** Create `frontend/src/layouts/AdvisorLayout.tsx`
- [ ] **3.6** Create `frontend/src/layouts/HRLayout.tsx`
- [ ] **3.7** Create `frontend/src/layouts/FacilityLayout.tsx`
- [ ] **3.8** Create `frontend/src/layouts/SchedulerLayout.tsx`
- [ ] **3.9** Create `frontend/src/layouts/DeanLayout.tsx`

### Routing

- [ ] **3.10** Update `frontend/src/routes/index.tsx` with new route groups
- [ ] **3.11** Remove teacher routes from `/admin/*`
- [ ] **3.12** Update `AdminLayout` to remove teacher navigation items
- [ ] **3.13** Add redirects from old routes to new routes

### Auth Context

- [ ] **3.14** Update `AuthContext` with multi-role support
- [ ] **3.15** Update `ProtectedRoute` to use `activeRole`
- [ ] **3.16** Add landing page redirect based on role

### Pages (placeholders or migrate existing)

- [ ] **3.17** Create/migrate Instructor pages
- [ ] **3.18** Create/migrate Advisor pages
- [ ] **3.19** Create HR pages
- [ ] **3.20** Create Facility pages
- [ ] **3.21** Create Scheduler pages
- [ ] **3.22** Create Academic/Dean pages

### Testing

- [ ] **3.23** Test role switching between layouts
- [ ] **3.24** Test navigation in each layout
- [ ] **3.25** Test protected routes
- [ ] **3.26** Test backward compatibility

---

## Phase 4: LTI 1.3 and Standards Compliance

### 4.1 Overview

This phase completes the LTI 1.3 integration and implements additional EdTech standards to enable seamless integration with Learning Management Systems (Canvas, Moodle, Blackboard, etc.) and learning analytics platforms.

### 4.2 LTI 1.3 Implementation

#### Current State

The application has partial LTI 1.3 support in `backend/internal/handlers/lti.go`:

- ✅ Tool Registration endpoints
- ✅ JWKS endpoint
- ✅ Login initiation
- ❌ Launch validation (returns `StatusNotImplemented`)
- ❌ Assignment and Grade Services (AGS)
- ❌ Names and Role Provisioning Services (NRPS)
- ❌ Deep Linking 2.0

#### LTI 1.3 Launch Implementation

**File**: `backend/internal/handlers/lti.go` (complete implementation)

```go
package handlers

import (
    "crypto/rsa"
    "encoding/json"
    "net/http"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

// LTI 1.3 Claims structure
type LTI13Claims struct {
    jwt.RegisteredClaims

    // Required LTI claims
    MessageType       string `json:"https://purl.imsglobal.org/spec/lti/claim/message_type"`
    Version           string `json:"https://purl.imsglobal.org/spec/lti/claim/version"`
    DeploymentID      string `json:"https://purl.imsglobal.org/spec/lti/claim/deployment_id"`
    TargetLinkURI     string `json:"https://purl.imsglobal.org/spec/lti/claim/target_link_uri"`
    ResourceLink      *ResourceLinkClaim `json:"https://purl.imsglobal.org/spec/lti/claim/resource_link,omitempty"`

    // User claims
    Sub               string   `json:"sub"` // User ID from platform
    Name              string   `json:"name"`
    Email             string   `json:"email"`
    GivenName         string   `json:"given_name"`
    FamilyName        string   `json:"family_name"`

    // Role claims (IMS LIS vocabulary)
    Roles             []string `json:"https://purl.imsglobal.org/spec/lti/claim/roles"`

    // Context claims (course/section)
    Context           *ContextClaim `json:"https://purl.imsglobal.org/spec/lti/claim/context,omitempty"`

    // Platform claims
    ToolPlatform      *ToolPlatformClaim `json:"https://purl.imsglobal.org/spec/lti/claim/tool_platform,omitempty"`

    // LIS claims
    LIS               *LISClaim `json:"https://purl.imsglobal.org/spec/lti/claim/lis,omitempty"`

    // Custom parameters
    Custom            map[string]string `json:"https://purl.imsglobal.org/spec/lti/claim/custom,omitempty"`

    // Service endpoints
    AGS               *AGSClaim `json:"https://purl.imsglobal.org/spec/lti-ags/claim/endpoint,omitempty"`
    NRPS              *NRPSClaim `json:"https://purl.imsglobal.org/spec/lti-nrps/claim/namesroleservice,omitempty"`
}

type ResourceLinkClaim struct {
    ID          string `json:"id"`
    Title       string `json:"title,omitempty"`
    Description string `json:"description,omitempty"`
}

type ContextClaim struct {
    ID    string   `json:"id"`
    Label string   `json:"label,omitempty"`
    Title string   `json:"title,omitempty"`
    Type  []string `json:"type,omitempty"`
}

type ToolPlatformClaim struct {
    GUID            string `json:"guid"`
    Name            string `json:"name,omitempty"`
    Version         string `json:"version,omitempty"`
    ProductFamilyCode string `json:"product_family_code,omitempty"`
}

type LISClaim struct {
    PersonSourcedID   string `json:"person_sourcedid,omitempty"`
    CourseOfferingSourcedID string `json:"course_offering_sourcedid,omitempty"`
    CourseSectionSourcedID string `json:"course_section_sourcedid,omitempty"`
}

type AGSClaim struct {
    Scope     []string `json:"scope"`
    LineItems string   `json:"lineitems,omitempty"`
    LineItem  string   `json:"lineitem,omitempty"`
}

type NRPSClaim struct {
    ContextMembershipsURL string   `json:"context_memberships_url"`
    ServiceVersions       []string `json:"service_versions,omitempty"`
}

// Launch handles the LTI 1.3 launch request
// POST /api/lti/launch
func (h *LTIHandler) Launch(c *gin.Context) {
    // 1. Get id_token from request
    idToken := c.PostForm("id_token")
    state := c.PostForm("state")

    if idToken == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id_token"})
        return
    }

    // 2. Validate state (CSRF protection)
    if !h.validateState(state) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
        return
    }

    // 3. Parse and validate JWT
    claims, err := h.parseAndValidateLTIToken(idToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid id_token: " + err.Error()})
        return
    }

    // 4. Validate required claims
    if err := h.validateRequiredClaims(claims); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 5. Get or create tenant from platform
    tenant, err := h.getOrCreateTenant(c.Request.Context(), claims.ToolPlatform)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve tenant"})
        return
    }

    // 6. Get or create user
    user, err := h.getOrCreateUser(c.Request.Context(), claims, tenant.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve user"})
        return
    }

    // 7. Map LIS roles to internal roles
    internalRoles := models.MapLISRolesToInternal(claims.Roles)

    // 8. Update user roles if changed
    if err := h.updateUserRoles(c.Request.Context(), user.ID, internalRoles); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update roles"})
        return
    }

    // 9. Handle context (course enrollment if applicable)
    if claims.Context != nil {
        if err := h.handleCourseContext(c.Request.Context(), user.ID, claims.Context, internalRoles); err != nil {
            // Log but don't fail - context handling is optional
            h.logger.Warn("Failed to handle course context", "error", err)
        }
    }

    // 10. Store AGS/NRPS endpoints for later use
    if claims.AGS != nil || claims.NRPS != nil {
        h.storeLTIServices(c.Request.Context(), claims, user.ID)
    }

    // 11. Generate session token
    token, expiresAt, err := h.authService.GenerateTokenWithRole(
        c.Request.Context(),
        user.ID.String(),
        string(internalRoles[0]), // Primary role
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
        return
    }

    // 12. Determine redirect URL based on role and target
    redirectURL := h.determineRedirectURL(claims, internalRoles)

    // 13. Set cookie and redirect (or return JSON for SPA)
    c.SetCookie("auth_token", token, int(expiresAt), "/", "", true, true)

    // For SPA, return JSON with redirect info
    c.JSON(http.StatusOK, gin.H{
        "token":        token,
        "user":         user,
        "roles":        internalRoles,
        "redirect_url": redirectURL,
        "expires_at":   expiresAt,
    })
}

// parseAndValidateLTIToken validates the LTI 1.3 id_token JWT
func (h *LTIHandler) parseAndValidateLTIToken(tokenString string) (*LTI13Claims, error) {
    // 1. Decode header to get kid
    parts := strings.Split(tokenString, ".")
    if len(parts) != 3 {
        return nil, errors.New("invalid token format")
    }

    // 2. Get platform's public key from JWKS
    token, err := jwt.ParseWithClaims(tokenString, &LTI13Claims{}, func(token *jwt.Token) (interface{}, error) {
        // Verify algorithm
        if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }

        // Get kid from header
        kid, ok := token.Header["kid"].(string)
        if !ok {
            return nil, errors.New("missing kid in token header")
        }

        // Fetch platform's public key
        return h.getPlatformPublicKey(token.Claims.(*LTI13Claims).Issuer, kid)
    })

    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*LTI13Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token claims")
    }

    // 3. Validate nonce (replay protection)
    if !h.validateNonce(claims.ID) {
        return nil, errors.New("nonce already used")
    }

    return claims, nil
}

// validateRequiredClaims checks all required LTI 1.3 claims
func (h *LTIHandler) validateRequiredClaims(claims *LTI13Claims) error {
    if claims.MessageType != "LtiResourceLinkRequest" {
        return errors.New("unsupported message type")
    }
    if claims.Version != "1.3.0" {
        return errors.New("unsupported LTI version")
    }
    if claims.DeploymentID == "" {
        return errors.New("missing deployment_id")
    }
    if claims.Sub == "" {
        return errors.New("missing sub (user id)")
    }
    if len(claims.Roles) == 0 {
        return errors.New("missing roles")
    }
    return nil
}
```

### 4.3 Assignment and Grade Services (AGS)

**File**: `backend/internal/services/lti_ags.go` (new file)

```go
package services

import (
    "context"
    "encoding/json"
    "net/http"
)

// LTIAGSService handles grade passback to LMS
type LTIAGSService interface {
    // Line Items (assignments)
    GetLineItems(ctx context.Context, endpoint string) ([]LineItem, error)
    CreateLineItem(ctx context.Context, endpoint string, item LineItem) (*LineItem, error)
    UpdateLineItem(ctx context.Context, endpoint string, item LineItem) error
    DeleteLineItem(ctx context.Context, endpoint string) error

    // Scores (grades)
    PostScore(ctx context.Context, lineItemURL string, score Score) error
    GetResults(ctx context.Context, lineItemURL string) ([]Result, error)
}

// LineItem represents an LTI AGS line item (gradebook column)
type LineItem struct {
    ID              string   `json:"id,omitempty"`
    ScoreMaximum    float64  `json:"scoreMaximum"`
    Label           string   `json:"label"`
    ResourceID      string   `json:"resourceId,omitempty"`
    ResourceLinkID  string   `json:"resourceLinkId,omitempty"`
    Tag             string   `json:"tag,omitempty"`
    StartDateTime   string   `json:"startDateTime,omitempty"`
    EndDateTime     string   `json:"endDateTime,omitempty"`
}

// Score represents a grade submission
type Score struct {
    UserID          string  `json:"userId"`
    ScoreGiven      float64 `json:"scoreGiven,omitempty"`
    ScoreMaximum    float64 `json:"scoreMaximum,omitempty"`
    Comment         string  `json:"comment,omitempty"`
    Timestamp       string  `json:"timestamp"`
    ActivityProgress string `json:"activityProgress"` // Initialized, Started, InProgress, Submitted, Completed
    GradingProgress  string `json:"gradingProgress"`  // FullyGraded, Pending, PendingManual, Failed, NotReady
}

// Result represents a grade result from LMS
type Result struct {
    ID              string  `json:"id"`
    ScoreOf         string  `json:"scoreOf"`
    UserID          string  `json:"userId"`
    ResultScore     float64 `json:"resultScore,omitempty"`
    ResultMaximum   float64 `json:"resultMaximum,omitempty"`
    Comment         string  `json:"comment,omitempty"`
}

// PostScore sends a grade to the LMS
func (s *ltiAGSServiceImpl) PostScore(ctx context.Context, lineItemURL string, score Score) error {
    // 1. Get OAuth2 access token for AGS
    token, err := s.getAccessToken(ctx, []string{
        "https://purl.imsglobal.org/spec/lti-ags/scope/score",
    })
    if err != nil {
        return err
    }

    // 2. Build request
    scoreURL := lineItemURL + "/scores"
    payload, _ := json.Marshal(score)

    req, _ := http.NewRequestWithContext(ctx, "POST", scoreURL, bytes.NewReader(payload))
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/vnd.ims.lis.v1.score+json")

    // 3. Send request
    resp, err := s.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        return fmt.Errorf("AGS score submission failed: %d", resp.StatusCode)
    }

    return nil
}
```

### 4.4 Names and Role Provisioning Services (NRPS)

**File**: `backend/internal/services/lti_nrps.go` (new file)

```go
package services

import (
    "context"
    "encoding/json"
    "net/http"
)

// LTINRPSService handles roster sync from LMS
type LTINRPSService interface {
    GetMemberships(ctx context.Context, contextMembershipsURL string) (*MembershipContainer, error)
    SyncCourseRoster(ctx context.Context, courseID uuid.UUID, memberships []Member) error
}

// MembershipContainer represents NRPS response
type MembershipContainer struct {
    ID        string   `json:"id"`
    Context   NRPSContext `json:"context"`
    Members   []Member `json:"members"`
}

type NRPSContext struct {
    ID    string `json:"id"`
    Label string `json:"label"`
    Title string `json:"title"`
}

// Member represents a course member from NRPS
type Member struct {
    UserID         string   `json:"user_id"`
    Status         string   `json:"status"` // Active, Inactive, Deleted
    Roles          []string `json:"roles"`
    Name           string   `json:"name,omitempty"`
    Email          string   `json:"email,omitempty"`
    GivenName      string   `json:"given_name,omitempty"`
    FamilyName     string   `json:"family_name,omitempty"`
    LISPersonSourcedID string `json:"lis_person_sourcedid,omitempty"`
}

// GetMemberships fetches course roster from LMS
func (s *ltiNRPSServiceImpl) GetMemberships(ctx context.Context, url string) (*MembershipContainer, error) {
    // 1. Get OAuth2 access token for NRPS
    token, err := s.getAccessToken(ctx, []string{
        "https://purl.imsglobal.org/spec/lti-nrps/scope/contextmembership.readonly",
    })
    if err != nil {
        return nil, err
    }

    // 2. Build request
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Accept", "application/vnd.ims.lti-nrps.v2.membershipcontainer+json")

    // 3. Send request
    resp, err := s.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        return nil, fmt.Errorf("NRPS request failed: %d", resp.StatusCode)
    }

    // 4. Parse response
    var container MembershipContainer
    if err := json.NewDecoder(resp.Body).Decode(&container); err != nil {
        return nil, err
    }

    return &container, nil
}

// SyncCourseRoster updates local enrollments from NRPS data
func (s *ltiNRPSServiceImpl) SyncCourseRoster(ctx context.Context, courseID uuid.UUID, members []Member) error {
    for _, member := range members {
        // 1. Get or create user
        user, err := s.userService.GetOrCreateByLTIID(ctx, member.UserID, member.Email, member.Name)
        if err != nil {
            s.logger.Error("Failed to sync user", "lti_id", member.UserID, "error", err)
            continue
        }

        // 2. Map LIS roles
        internalRoles := models.MapLISRolesToInternal(member.Roles)

        // 3. Update enrollment
        enrollment := models.CourseEnrollment{
            UserID:   user.ID,
            CourseID: courseID,
            Status:   mapMemberStatus(member.Status),
            Role:     determineCourseRole(internalRoles),
        }

        if err := s.enrollmentRepo.Upsert(ctx, &enrollment); err != nil {
            s.logger.Error("Failed to sync enrollment", "user_id", user.ID, "error", err)
        }
    }

    return nil
}
```

### 4.5 FERPA/GDPR Compliance

**File**: `backend/internal/models/data_retention.go` (new file)

```go
package models

import (
    "time"
    "github.com/google/uuid"
)

// DataClassification defines sensitivity levels
type DataClassification string

const (
    ClassPublic       DataClassification = "public"
    ClassInternal     DataClassification = "internal"
    ClassConfidential DataClassification = "confidential" // FERPA protected
    ClassRestricted   DataClassification = "restricted"   // GDPR special categories
)

// DataRetentionPolicy defines how long data is kept
type DataRetentionPolicy struct {
    ID              uuid.UUID          `json:"id" gorm:"type:uuid;primary_key"`
    TenantID        *uuid.UUID         `json:"tenant_id" gorm:"type:uuid"`
    EntityType      string             `json:"entity_type"` // user_data, activity_logs, grades, etc.
    Classification  DataClassification `json:"classification"`
    RetentionDays   int                `json:"retention_days"` // 0 = indefinite
    AnonymizeAfter  int                `json:"anonymize_after_days"` // Days before anonymization
    LegalBasis      string             `json:"legal_basis"` // consent, contract, legal_obligation, etc.
    Description     string             `json:"description"`
    IsActive        bool               `json:"is_active"`
    CreatedAt       time.Time          `json:"created_at"`
}

// ConsentRecord tracks user consent for GDPR
type ConsentRecord struct {
    ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
    UserID       uuid.UUID  `json:"user_id" gorm:"type:uuid;index"`
    ConsentType  string     `json:"consent_type"` // marketing, analytics, third_party_sharing
    IsGranted    bool       `json:"is_granted"`
    GrantedAt    *time.Time `json:"granted_at"`
    RevokedAt    *time.Time `json:"revoked_at"`
    IPAddress    string     `json:"ip_address"`
    UserAgent    string     `json:"user_agent"`
    Version      string     `json:"version"` // Policy version consented to
    CreatedAt    time.Time  `json:"created_at"`
}

// DataExportRequest for GDPR data portability
type DataExportRequest struct {
    ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
    UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;index"`
    Status      string     `json:"status"` // pending, processing, completed, failed
    RequestedAt time.Time  `json:"requested_at"`
    CompletedAt *time.Time `json:"completed_at"`
    ExportURL   string     `json:"export_url"` // Temporary download URL
    ExpiresAt   *time.Time `json:"expires_at"`
    Format      string     `json:"format"` // json, csv
    CreatedAt   time.Time  `json:"created_at"`
}

// DataDeletionRequest for GDPR right to erasure
type DataDeletionRequest struct {
    ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
    UserID        uuid.UUID  `json:"user_id" gorm:"type:uuid;index"`
    Status        string     `json:"status"` // pending, approved, processing, completed, rejected
    RequestedAt   time.Time  `json:"requested_at"`
    ApprovedBy    *uuid.UUID `json:"approved_by" gorm:"type:uuid"`
    ApprovedAt    *time.Time `json:"approved_at"`
    CompletedAt   *time.Time `json:"completed_at"`
    RejectionReason string   `json:"rejection_reason"`
    RetainedData  []string   `json:"retained_data"` // Data kept for legal reasons
    CreatedAt     time.Time  `json:"created_at"`
}
```

### 4.6 Enrollment-Based Authorization Integration

**File**: `backend/internal/middleware/rbac.go` (additions)

```go
// RequireEnrollment checks if user has enrollment-based access to a course
func RequireEnrollment(requiredRole string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("user_id")
        courseID := c.Param("courseId")

        if courseID == "" {
            courseID = c.Param("course_id")
        }

        if courseID == "" {
            c.Next() // No course context, skip enrollment check
            return
        }

        // Check enrollment
        enrollment, err := enrollmentRepo.GetByUserAndCourse(c.Request.Context(), userID, courseID)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                "error": "Not enrolled in this course",
            })
            return
        }

        // Check role in course if specified
        if requiredRole != "" && enrollment.Role != requiredRole {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                "error": "Insufficient role in course",
            })
            return
        }

        // Set enrollment info in context
        c.Set("enrollment", enrollment)
        c.Set("course_role", enrollment.Role)

        c.Next()
    }
}

// CombinedRBACMiddleware checks both global RBAC and enrollment-based access
func CombinedRBACMiddleware(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("user_id")
        activeRole := c.GetString("active_role")

        // First check global/tenant permissions
        hasGlobalAccess, _ := authzService.HasPermission(c.Request.Context(), userID, permission, "global", "")
        if hasGlobalAccess {
            c.Next()
            return
        }

        // Then check course-context permissions
        courseID := c.Param("courseId")
        if courseID != "" {
            hasCourseAccess, _ := authzService.HasPermission(c.Request.Context(), userID, permission, "course", courseID)
            if hasCourseAccess {
                c.Next()
                return
            }
        }

        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
            "error": "Permission denied",
        })
    }
}
```

---

## Phase 4 Implementation Checklist

### LTI 1.3 Core

- [ ] **4.1** Complete `Launch()` handler with full validation
- [ ] **4.2** Implement `parseAndValidateLTIToken()`
- [ ] **4.3** Implement nonce validation and storage
- [ ] **4.4** Implement platform public key fetching (JWKS)
- [ ] **4.5** Implement user provisioning from LTI claims
- [ ] **4.6** Store LTI context in session/database

### LTI Services

- [ ] **4.7** Create `backend/internal/services/lti_ags.go`
- [ ] **4.8** Implement `PostScore()` for grade passback
- [ ] **4.9** Implement `GetLineItems()` for assignment sync
- [ ] **4.10** Create `backend/internal/services/lti_nrps.go`
- [ ] **4.11** Implement `GetMemberships()` for roster sync
- [ ] **4.12** Create background job for roster sync

### Compliance

- [ ] **4.13** Create `backend/internal/models/data_retention.go`
- [ ] **4.14** Implement data export endpoint (GDPR)
- [ ] **4.15** Implement data deletion request workflow
- [ ] **4.16** Add consent management endpoints
- [ ] **4.17** Create data anonymization job

### Authorization Integration

- [ ] **4.18** Update RBAC middleware with enrollment check
- [ ] **4.19** Implement `CombinedRBACMiddleware`
- [ ] **4.20** Update routes to use combined authorization

### Testing

- [ ] **4.21** Integration tests with LTI Reference Implementation
- [ ] **4.22** Test grade passback to Canvas/Moodle
- [ ] **4.23** Test roster sync
- [ ] **4.24** Test data export functionality
- [ ] **4.25** Security audit for LTI endpoints

---

## Migration Strategy

### Phase 1 Migration (RBAC)

1. Deploy new role definitions (backward compatible)
2. Run migration to create `user_roles` table
3. Migrate existing single roles to `user_roles`
4. Deploy role switching API
5. Frontend: Add role switcher, test with existing users
6. Gradually assign additional roles to multi-role users

### Phase 2 Migration (Workflows)

1. Deploy workflow tables and service
2. Create default workflow templates
3. Enable workflows for new entities first (opt-in)
4. Migrate existing approval processes to workflows
5. Configure tenant-specific workflow templates

### Phase 3 Migration (Layouts)

1. Deploy new layouts alongside existing AdminLayout
2. Add feature flag for new navigation
3. Redirect users to new layouts based on role
4. Remove teacher routes from AdminLayout
5. Monitor usage and fix issues
6. Remove feature flag after stabilization

### Phase 4 Migration (LTI)

1. Complete LTI 1.3 launch implementation
2. Test with sandbox LMS instances
3. Enable AGS for grade passback
4. Enable NRPS for roster sync
5. Coordinate with partner institutions for production testing

---

## Implementation Timeline

| Phase              | Duration  | Dependencies | Team               |
| ------------------ | --------- | ------------ | ------------------ |
| Phase 1: RBAC      | 2-3 weeks | None         | Backend + Frontend |
| Phase 2: Workflows | 3-4 weeks | Phase 1      | Backend + Frontend |
| Phase 3: Layouts   | 2-3 weeks | Phase 1      | Frontend           |
| Phase 4: LTI       | 3-4 weeks | Phase 1      | Backend            |

**Total Estimated Duration**: 8-12 weeks (with parallel work on Phases 2-4)

### Milestones

1. **Week 2**: Phase 1 complete, role switching functional
2. **Week 5**: Phase 2 complete, workflow engine operational
3. **Week 7**: Phase 3 complete, all layouts deployed
4. **Week 10**: Phase 4 complete, LTI 1.3 certified

---

## Appendix A: Role-Permission Matrix

| Permission       | superadmin | admin | hr_admin | facility_manager | scheduler_admin | dean | chair | instructor | advisor | student |
| ---------------- | ---------- | ----- | -------- | ---------------- | --------------- | ---- | ----- | ---------- | ------- | ------- |
| user.view        | ✅         | ✅    | ✅       | ❌               | ❌              | ❌   | ❌    | ❌         | ❌      | ❌      |
| user.create      | ✅         | ❌    | ✅       | ❌               | ❌              | ❌   | ❌    | ❌         | ❌      | ❌      |
| user.edit        | ✅         | ❌    | ✅       | ❌               | ❌              | ❌   | ❌    | ❌         | ❌      | ❌      |
| user.delete      | ✅         | ❌    | ✅       | ❌               | ❌              | ❌   | ❌    | ❌         | ❌      | ❌      |
| user.assign_role | ✅         | ❌    | ✅       | ❌               | ❌              | ❌   | ❌    | ❌         | ❌      | ❌      |
| course.view      | ✅         | ✅    | ❌       | ❌               | ✅              | ✅   | ✅    | ✅         | ✅      | ✅      |
| course.create    | ✅         | ✅    | ❌       | ❌               | ❌              | ❌   | ❌    | ✅         | ✅      | ❌      |
| course.approve   | ✅         | ❌    | ❌       | ❌               | ❌              | ✅   | ✅    | ❌         | ❌      | ❌      |
| course.teach     | ✅         | ❌    | ❌       | ❌               | ❌              | ❌   | ❌    | ✅         | ✅      | ❌      |
| grade.view       | ✅         | ✅    | ❌       | ❌               | ❌              | ✅   | ✅    | ✅         | ✅      | ✅\*    |
| grade.edit       | ✅         | ❌    | ❌       | ❌               | ❌              | ❌   | ❌    | ✅         | ✅      | ❌      |
| grade.approve    | ✅         | ❌    | ❌       | ❌               | ❌              | ✅   | ❌    | ❌         | ❌      | ❌      |
| schedule.view    | ✅         | ✅    | ❌       | ✅               | ✅              | ✅   | ✅    | ✅         | ✅      | ✅      |
| schedule.create  | ✅         | ❌    | ❌       | ❌               | ✅              | ❌   | ❌    | ❌         | ❌      | ❌      |
| schedule.approve | ✅         | ❌    | ❌       | ❌               | ❌              | ✅   | ❌    | ❌         | ❌      | ❌      |
| facility.view    | ✅         | ✅    | ❌       | ✅               | ✅              | ❌   | ❌    | ❌         | ❌      | ❌      |
| facility.create  | ✅         | ❌    | ❌       | ✅               | ❌              | ❌   | ❌    | ❌         | ❌      | ❌      |
| facility.book    | ✅         | ❌    | ❌       | ✅               | ✅              | ❌   | ❌    | ❌         | ❌      | ❌      |
| workflow.view    | ✅         | ✅    | ✅       | ✅               | ✅              | ✅   | ✅    | ❌         | ✅      | ❌      |
| workflow.approve | ✅         | ❌    | ✅       | ✅               | ✅              | ✅   | ✅    | ❌         | ✅      | ❌      |
| analytics.view   | ✅         | ✅    | ❌       | ❌               | ❌              | ✅   | ✅    | ❌         | ❌      | ❌      |
| system.settings  | ✅         | ✅    | ❌       | ❌               | ❌              | ❌   | ❌    | ❌         | ❌      | ❌      |

\*Student can only view own grades

---

## Appendix B: API Endpoints Summary

### New Endpoints

```
# Role Management
POST   /api/auth/switch-role
GET    /api/users/:id/roles
POST   /api/users/:id/roles
DELETE /api/users/:id/roles/:roleId

# Workflows
GET    /api/workflows/templates
POST   /api/workflows/templates
GET    /api/workflows/pending
GET    /api/workflows/pending/count
POST   /api/workflows/start
POST   /api/workflows/approve
GET    /api/workflows/instances/:id
POST   /api/workflows/delegations
DELETE /api/workflows/delegations/:id

# LTI 1.3
POST   /api/lti/launch
GET    /api/lti/jwks
POST   /api/lti/login
GET    /api/lti/config/:toolId

# Compliance
GET    /api/gdpr/export
POST   /api/gdpr/export
POST   /api/gdpr/delete-request
GET    /api/gdpr/consent
POST   /api/gdpr/consent
```

---

_Document Version: 1.0_  
_Last Updated: January 8, 2026_
