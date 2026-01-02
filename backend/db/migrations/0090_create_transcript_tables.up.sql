CREATE TABLE IF NOT EXISTS term_grades (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID NOT NULL REFERENCES users(id),
    term_id UUID NOT NULL REFERENCES academic_terms(id),
    course_offering_id UUID NOT NULL REFERENCES course_offerings(id),
    
    -- Snapshots
    course_title TEXT NOT NULL,
    course_code TEXT NOT NULL,
    credits NUMERIC(5,2) NOT NULL DEFAULT 0,
    
    -- Grade Info
    grade VARCHAR(5) NOT NULL, -- "A", "B+", "W"
    grade_points NUMERIC(4,2) NOT NULL DEFAULT 0, -- 4.0, 3.33
    percentage NUMERIC(5,2) NOT NULL DEFAULT 0,
    is_passed BOOLEAN NOT NULL DEFAULT FALSE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    UNIQUE(student_id, course_offering_id)
);

CREATE INDEX idx_term_grades_student ON term_grades(student_id);
CREATE INDEX idx_term_grades_term ON term_grades(term_id);
