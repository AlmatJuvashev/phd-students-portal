CREATE TABLE profile_submissions (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    form_data JSONB NOT NULL,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER profile_submissions_touch_updated_at
    BEFORE UPDATE ON profile_submissions
    FOR EACH ROW
    EXECUTE FUNCTION touch_updated_at();
