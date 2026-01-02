-- Create Question Banks Table
CREATE TABLE question_banks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL, -- Logical tenant isolation
    title VARCHAR(255) NOT NULL,
    description TEXT,
    subject VARCHAR(100), -- Anatomy, Histology, etc.
    blooms_taxonomy VARCHAR(50), -- Knowledge, Comprehension, Application, etc.
    is_public BOOLEAN DEFAULT FALSE,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create Questions Table
CREATE TABLE questions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bank_id UUID NOT NULL REFERENCES question_banks(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- MCQ, MRQ, TRUE_FALSE, TEXT, LIKERT
    stem TEXT NOT NULL, -- The actual question text
    media_url VARCHAR(255), -- Optional image/video
    points_default FLOAT DEFAULT 1.0,
    difficulty_level VARCHAR(50), -- EASY, MEDIUM, HARD
    learning_outcome_id UUID, -- Link to defined learning outcome (nullable for now)
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create Question Options Table (for objective questions)
CREATE TABLE question_options (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    is_correct BOOLEAN DEFAULT FALSE,
    sort_order INT NOT NULL DEFAULT 0,
    feedback TEXT -- Explanation for why this is correct/incorrect
);

-- Create Assessments Table (The "Exam" configuration)
CREATE TABLE assessments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    course_offering_id UUID NOT NULL REFERENCES course_offerings(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    time_limit_minutes INT, -- NULL means no limit
    available_from TIMESTAMPTZ,
    available_until TIMESTAMPTZ,
    shuffle_questions BOOLEAN DEFAULT FALSE,
    grading_policy VARCHAR(50) DEFAULT 'AUTOMATIC', -- AUTOMATIC, MANUAL_REVIEW
    passing_score FLOAT DEFAULT 0.0,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create Assessment Sections (Optional grouping, e.g. Part 1, Part 2)
CREATE TABLE assessment_sections (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    title VARCHAR(255),
    instructions TEXT,
    sort_order INT NOT NULL DEFAULT 0
);

-- Create Assessment Items (Linking Questions to Assessments)
CREATE TABLE assessment_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    section_id UUID REFERENCES assessment_sections(id) ON DELETE SET NULL,
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    points_override FLOAT, -- Override default points for this specific exam
    sort_order INT NOT NULL DEFAULT 0,
    
    UNIQUE(assessment_id, question_id) 
);

-- Create Assessment Attempts (Student taking the test)
CREATE TABLE assessment_attempts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    started_at TIMESTAMPTZ DEFAULT NOW(),
    finished_at TIMESTAMPTZ,
    score FLOAT DEFAULT 0.0,
    status VARCHAR(50) DEFAULT 'IN_PROGRESS', -- IN_PROGRESS, SUBMITTED, GRADED
    
    UNIQUE(assessment_id, student_id, started_at) -- Allow multiple attempts but different times
);

-- Create Item Responses (Granular answers)
CREATE TABLE item_responses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    attempt_id UUID NOT NULL REFERENCES assessment_attempts(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    selected_option_id UUID REFERENCES question_options(id), -- For MCQ
    text_response TEXT, -- For Essay/Short Answer
    score FLOAT DEFAULT 0.0,
    is_correct BOOLEAN DEFAULT FALSE,
    graded_at TIMESTAMPTZ,
    
    UNIQUE(attempt_id, question_id)
);

-- Indexes for performance
CREATE INDEX idx_questions_bank ON questions(bank_id);
CREATE INDEX idx_options_question ON question_options(question_id);
CREATE INDEX idx_assessments_offering ON assessments(course_offering_id);
CREATE INDEX idx_items_assessment ON assessment_items(assessment_id);
CREATE INDEX idx_attempts_student ON assessment_attempts(student_id);
CREATE INDEX idx_responses_attempt ON item_responses(attempt_id);
