-- Learning Outcomes table
CREATE TABLE learning_outcomes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    program_id UUID REFERENCES programs(id),
    course_id UUID REFERENCES courses(id),
    code VARCHAR(20) NOT NULL,
    title JSONB NOT NULL,
    description JSONB,
    category VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Link outcomes to assessments (journey nodes)
CREATE TABLE outcome_assessments (
    outcome_id UUID REFERENCES learning_outcomes(id) ON DELETE CASCADE,
    node_definition_id UUID REFERENCES journey_node_definitions(id) ON DELETE CASCADE,
    weight DECIMAL(3,2) DEFAULT 1.0,
    PRIMARY KEY (outcome_id, node_definition_id)
);

-- Curriculum Change Log (audit history)
CREATE TABLE curriculum_change_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    action VARCHAR(20) NOT NULL,
    old_value JSONB,
    new_value JSONB,
    changed_by UUID REFERENCES users(id),
    changed_at TIMESTAMPTZ DEFAULT now()
);

-- Audit access tokens for external users
CREATE TABLE audit_access_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    user_id UUID REFERENCES users(id),
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    scope VARCHAR(100) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Indexes
CREATE INDEX idx_learning_outcomes_tenant ON learning_outcomes(tenant_id);
CREATE INDEX idx_learning_outcomes_program ON learning_outcomes(program_id) WHERE program_id IS NOT NULL;
CREATE INDEX idx_learning_outcomes_course ON learning_outcomes(course_id) WHERE course_id IS NOT NULL;
CREATE INDEX idx_curriculum_change_log_tenant ON curriculum_change_log(tenant_id);
CREATE INDEX idx_curriculum_change_log_entity ON curriculum_change_log(entity_type, entity_id);
CREATE INDEX idx_audit_access_tokens_user ON audit_access_tokens(user_id);
