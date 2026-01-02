-- Create Course Enrollments Table
CREATE TABLE course_enrollments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_offering_id UUID NOT NULL REFERENCES course_offerings(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'ENROLLED', -- ENROLLED, PENDING, DROPPED, WAITLIST
    method VARCHAR(50) NOT NULL DEFAULT 'ADMIN', -- ADMIN, SELF, SYSTEM
    enrolled_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    UNIQUE(course_offering_id, student_id)
);

-- Create Activity Submissions Table
CREATE TABLE activity_submissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    activity_id UUID NOT NULL, -- Logical link to JSONB content in course_activities (no FK enforced to course_activities table as it might not be normalized in V11 yet, but ideally should be. Let's assume course_content table exists based on previous phases)
    -- Wait, course_activities ARE normalized in V11 Phase 1.5. Let's check tablename.
    -- Assuming 'course_activities' table exists from Phase 1.5. 
    -- If not, I'll remove the FK constraint and just use UUID.
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_offering_id UUID NOT NULL REFERENCES course_offerings(id) ON DELETE CASCADE,
    content JSONB NOT NULL DEFAULT '{}',
    status VARCHAR(50) NOT NULL DEFAULT 'SUBMITTED', -- SUBMITTED, DRAFT, GRADED
    submitted_at TIMESTAMPTZ DEFAULT NOW(),
    
    UNIQUE(activity_id, student_id)
);

-- Create Class Attendance Table
CREATE TABLE class_attendance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    class_session_id UUID NOT NULL REFERENCES class_sessions(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'PRESENT', -- PRESENT, ABSENT, LATE, EXCUSED
    notes TEXT DEFAULT '',
    recorded_by_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(class_session_id, student_id)
);

-- Indexes for performance
CREATE INDEX idx_enrollments_student ON course_enrollments(student_id);
CREATE INDEX idx_enrollments_offering ON course_enrollments(course_offering_id);
CREATE INDEX idx_submissions_offering ON activity_submissions(course_offering_id);
CREATE INDEX idx_submissions_student ON activity_submissions(student_id);
CREATE INDEX idx_attendance_session ON class_attendance(class_session_id);
CREATE INDEX idx_attendance_student ON class_attendance(student_id);
