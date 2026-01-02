CREATE TABLE student_risk_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    risk_score FLOAT NOT NULL DEFAULT 0, -- 0 to 100
    risk_factors JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_risk_snapshots_student ON student_risk_snapshots(student_id);
CREATE INDEX idx_risk_snapshots_score ON student_risk_snapshots(risk_score);
