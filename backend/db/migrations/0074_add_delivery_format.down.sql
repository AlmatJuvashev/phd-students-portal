-- Rollback: Remove delivery format support

ALTER TABLE class_sessions 
    DROP COLUMN IF EXISTS meeting_url;

ALTER TABLE class_sessions 
    DROP COLUMN IF EXISTS session_format;

ALTER TABLE course_offerings 
    DROP COLUMN IF EXISTS meeting_url;

ALTER TABLE course_offerings 
    DROP COLUMN IF EXISTS virtual_capacity;

DROP INDEX IF EXISTS idx_course_offerings_delivery_format;

ALTER TABLE course_offerings 
    DROP COLUMN IF EXISTS delivery_format;
