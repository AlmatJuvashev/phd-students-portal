CREATE TABLE IF NOT EXISTS grading_schemas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    scale JSONB NOT NULL DEFAULT '[]',
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for tenant lookups
CREATE INDEX IF NOT EXISTS idx_grading_schemas_tenant ON grading_schemas(tenant_id);


CREATE TABLE IF NOT EXISTS gradebook_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_offering_id UUID NOT NULL REFERENCES course_offerings(id) ON DELETE CASCADE,
    activity_id UUID NOT NULL REFERENCES course_activities(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    score NUMERIC(5, 2) NOT NULL DEFAULT 0.00,
    max_score NUMERIC(5, 2) NOT NULL DEFAULT 100.00,
    grade VARCHAR(50) DEFAULT NULL,
    feedback TEXT DEFAULT '',
    graded_by_id UUID REFERENCES users(id) ON DELETE SET NULL,
    graded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ensure one entry per student per activity per offering
    UNIQUE(course_offering_id, activity_id, student_id)
);

CREATE INDEX IF NOT EXISTS idx_gradebook_offering ON gradebook_entries(course_offering_id);
CREATE INDEX IF NOT EXISTS idx_gradebook_student ON gradebook_entries(student_id);
