-- 1. Academic Terms
CREATE TABLE academic_terms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    is_active BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(tenant_id, code) -- Prevent duplicate term codes per tenant
);

-- 2. Course Offerings (Instances)
CREATE TABLE course_offerings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    term_id UUID NOT NULL REFERENCES academic_terms(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    section VARCHAR(50) NOT NULL, -- e.g. "001", "A"
    max_capacity INTEGER DEFAULT 30,
    current_enrolled INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    status VARCHAR(50) DEFAULT 'DRAFT', -- DRAFT, PUBLISHED, ARCHIVED
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(term_id, course_id, section) -- Prevent duplicate sections for same course in a term
);

-- 3. Course Staff (Instructors, TAs)
CREATE TABLE course_staff (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_offering_id UUID NOT NULL REFERENCES course_offerings(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL, -- INSTRUCTOR, TA, GRADER
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(course_offering_id, user_id)
);

-- 4. Class Sessions (Calendar Events)
CREATE TABLE class_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_offering_id UUID NOT NULL REFERENCES course_offerings(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    start_time VARCHAR(5) NOT NULL, -- HH:MM (24h)
    end_time VARCHAR(5) NOT NULL,   -- HH:MM (24h)
    room_id UUID REFERENCES rooms(id) ON DELETE SET NULL,
    instructor_id UUID REFERENCES users(id) ON DELETE SET NULL, -- Override default instructor
    type VARCHAR(50) DEFAULT 'LECTURE', -- LECTURE, LAB, EXAM
    is_cancelled BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_academic_terms_tenant ON academic_terms(tenant_id);
CREATE INDEX idx_course_offerings_term ON course_offerings(term_id);
CREATE INDEX idx_course_offerings_course ON course_offerings(course_id);
CREATE INDEX idx_course_staff_user ON course_staff(user_id);
CREATE INDEX idx_class_sessions_offering ON class_sessions(course_offering_id);
CREATE INDEX idx_class_sessions_date ON class_sessions(date);
CREATE INDEX idx_class_sessions_room ON class_sessions(room_id);
