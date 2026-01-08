CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL,
    tenant_id UUID, -- Optional: for tenant-scoped roles (nullable implies global)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, role, tenant_id)
);

CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role ON user_roles(role);

CREATE TABLE IF NOT EXISTS role_metadata (
    role_name VARCHAR(50) PRIMARY KEY,
    description TEXT,
    category VARCHAR(50),
    is_system_role BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS permission_change_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_name VARCHAR(50) NOT NULL,
    permission VARCHAR(100) NOT NULL,
    action VARCHAR(20) NOT NULL, -- 'GRANT', 'REVOKE'
    changed_by UUID REFERENCES users(id),
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    reason TEXT
);

-- Seed default basic role metadata (comprehensive list in later seed or via admin UI)
INSERT INTO role_metadata (role_name, category, description) VALUES
('superadmin', 'platform', 'Global System Superadmin'),
('admin', 'admin', 'Administrator'),
('student', 'student', 'Learner'),
('instructor', 'teaching', 'Course Instructor'),
('advisor', 'academic', 'Academic Advisor'),
('dean', 'academic', 'Faculty Dean'),
('registrar', 'administrative', 'Registrar Staff')
ON CONFLICT DO NOTHING;
