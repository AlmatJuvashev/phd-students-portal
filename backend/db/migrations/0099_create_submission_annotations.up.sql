CREATE TABLE submission_annotations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    submission_id UUID NOT NULL REFERENCES activity_submissions(id) ON DELETE CASCADE,
    author_id UUID NOT NULL REFERENCES users(id),
    
    page_number INT NOT NULL,
    annotation_type VARCHAR(50) NOT NULL, -- 'highlight', 'text', 'drawing', 'strikeout'
    
    -- Coordinates (Normalized 0-100% to handle resize)
    x_percent FLOAT NOT NULL,
    y_percent FLOAT NOT NULL,
    width_percent FLOAT,
    height_percent FLOAT,
    
    content TEXT, -- For text comments
    color VARCHAR(20) DEFAULT '#FF0000',
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_annotations_submission ON submission_annotations(submission_id);
