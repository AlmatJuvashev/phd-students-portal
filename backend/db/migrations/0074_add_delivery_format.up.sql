-- Migration: Add delivery format support to course offerings
-- Supports: IN_PERSON, ONLINE_SYNC, ONLINE_ASYNC, HYBRID

-- Add delivery_format to course_offerings
ALTER TABLE course_offerings 
    ADD COLUMN IF NOT EXISTS delivery_format VARCHAR(20) DEFAULT 'IN_PERSON' NOT NULL;

-- Add virtual_capacity for online courses (optional, NULL means unlimited or same as max_capacity)
ALTER TABLE course_offerings 
    ADD COLUMN IF NOT EXISTS virtual_capacity INTEGER;

-- Add default meeting URL for online offerings
ALTER TABLE course_offerings 
    ADD COLUMN IF NOT EXISTS meeting_url VARCHAR(512);

-- Add session-level format override for HYBRID courses
ALTER TABLE class_sessions 
    ADD COLUMN IF NOT EXISTS session_format VARCHAR(20);

-- Add session-specific meeting URL
ALTER TABLE class_sessions 
    ADD COLUMN IF NOT EXISTS meeting_url VARCHAR(512);

-- Index for filtering offerings by format
CREATE INDEX IF NOT EXISTS idx_course_offerings_delivery_format 
    ON course_offerings(delivery_format);

-- Comment on columns
COMMENT ON COLUMN course_offerings.delivery_format IS 'Delivery format: IN_PERSON, ONLINE_SYNC, ONLINE_ASYNC, HYBRID';
COMMENT ON COLUMN course_offerings.virtual_capacity IS 'Max participants for online sessions (NULL = unlimited)';
COMMENT ON COLUMN course_offerings.meeting_url IS 'Default meeting link for online sessions (Zoom, Teams, etc)';
COMMENT ON COLUMN class_sessions.session_format IS 'Override format for HYBRID courses (IN_PERSON or ONLINE_SYNC)';
COMMENT ON COLUMN class_sessions.meeting_url IS 'Session-specific meeting link';
