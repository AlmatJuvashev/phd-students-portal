-- Course Offering -> Cohorts (Many-to-Many)
-- Used for Conflict Detection: Prevents scheduling two offerings for the same cohort at the same time.
CREATE TABLE course_offering_cohorts (
    course_offering_id UUID REFERENCES course_offerings(id) ON DELETE CASCADE,
    cohort_id UUID NOT NULL, -- references dictionaries.cohorts(id) logically
    PRIMARY KEY (course_offering_id, cohort_id)
);

-- Room Attributes (e.g. "EQUIPMENT": "Projector", "TYPE": "Computer Lab")
CREATE TABLE room_attributes (
    room_id UUID REFERENCES rooms(id) ON DELETE CASCADE,
    key VARCHAR(50) NOT NULL,
    value VARCHAR(100) NOT NULL,
    PRIMARY KEY (room_id, key, value)
);

-- Course Requirements (e.g. "REQUIRES_EQUIPMENT": "Projector")
CREATE TABLE course_requirements (
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    key VARCHAR(50) NOT NULL,
    value VARCHAR(100) NOT NULL,
    PRIMARY KEY (course_id, key, value)
);
