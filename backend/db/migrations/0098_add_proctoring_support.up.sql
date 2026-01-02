-- Add security settings to assessment configuration
ALTER TABLE assessments ADD COLUMN security_settings JSONB DEFAULT '{}';

-- Create table for proctoring logs
CREATE TABLE proctoring_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    attempt_id UUID NOT NULL REFERENCES assessment_attempts(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL, -- TAB_SWITCH, BLUR, FOCUS_LOST, MULTIPLE_FACES, NO_FACE
    occurred_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'
);

CREATE INDEX idx_proctoring_logs_attempt ON proctoring_logs(attempt_id);
