-- Permissions: Capability definitions (e.g., 'course.view', 'grade.edit')
CREATE TABLE permissions (
    slug VARCHAR(100) PRIMARY KEY,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Roles: Named collections of permissions (e.g., 'Instructor', 'Student')
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_system_role BOOLEAN DEFAULT false, -- If true, cannot be deleted/edited by users
    tenant_id UUID, -- Optional: custom roles scoped to a tenant. Null for system global roles.
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(name, tenant_id) -- Unique role name per tenant (or global if null)
);

-- Role -> Permissions (Many-to-Many)
CREATE TABLE role_permissions (
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    permission_slug VARCHAR(100) REFERENCES permissions(slug) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_slug)
);

-- User -> Context -> Role (The Core of Contextual RBAC)
-- context_type: 'global', 'tenant', 'department', 'course'
-- context_id: UUID of the related entity (or null/zero-uuid for global)
CREATE TABLE user_context_roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    context_type VARCHAR(50) NOT NULL, 
    context_id UUID NOT NULL, -- Logical ID of target (TenantID, CourseID, etc.)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Prevent duplicate role assignment for same user in same context
    UNIQUE(user_id, role_id, context_type, context_id) 
);

-- Indexes for fast lookup during authorization checks
CREATE INDEX idx_ucr_user_context ON user_context_roles(user_id, context_type, context_id);
CREATE INDEX idx_ucr_user_global ON user_context_roles(user_id) WHERE context_type = 'global';
