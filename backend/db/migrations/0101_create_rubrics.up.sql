CREATE TABLE rubrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_offering_id UUID NOT NULL REFERENCES course_offerings(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    is_global BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rubric_criteria (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rubric_id UUID NOT NULL REFERENCES rubrics(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL, -- e.g. "Grammar"
    description TEXT,
    weight FLOAT NOT NULL DEFAULT 1.0, -- Multiplier
    position INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rubric_levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    criterion_id UUID NOT NULL REFERENCES rubric_criteria(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL, -- e.g. "Excellent"
    description TEXT,
    points FLOAT NOT NULL, -- e.g. 5.0
    position INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rubric_grades (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id UUID NOT NULL REFERENCES activity_submissions(id) ON DELETE CASCADE,
    rubric_id UUID NOT NULL REFERENCES rubrics(id) ON DELETE CASCADE,
    grader_id UUID REFERENCES users(id) ON DELETE SET NULL,
    total_score FLOAT NOT NULL DEFAULT 0,
    comments TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rubric_grade_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rubric_grade_id UUID NOT NULL REFERENCES rubric_grades(id) ON DELETE CASCADE,
    criterion_id UUID NOT NULL REFERENCES rubric_criteria(id) ON DELETE CASCADE,
    level_id UUID REFERENCES rubric_levels(id) ON DELETE SET NULL,
    points_awarded FLOAT NOT NULL DEFAULT 0,
    comments TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rubrics_course ON rubrics(course_offering_id);
CREATE INDEX idx_rubric_grades_submission ON rubric_grades(submission_id);
