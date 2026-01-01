CREATE TABLE IF NOT EXISTS question_banks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_question_banks_tenant ON question_banks(tenant_id);

CREATE TABLE IF NOT EXISTS question_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bank_id UUID NOT NULL REFERENCES question_banks(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    content JSONB NOT NULL DEFAULT '{}',
    difficulty INTEGER DEFAULT 1,
    tags TEXT[] DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_question_items_bank ON question_items(bank_id);
CREATE INDEX IF NOT EXISTS idx_question_items_tags ON question_items USING GIN (tags);
