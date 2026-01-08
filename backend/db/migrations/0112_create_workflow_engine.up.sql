-- Workflow Templates
CREATE TABLE IF NOT EXISTS workflow_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID, -- REFERENCES tenants(id), -- Add FK if tenants table guaranteed to exist in all setups
    name VARCHAR(100) NOT NULL,
    description TEXT,
    entity_type VARCHAR(50) NOT NULL, -- course_approval, schedule_approval, thesis_stage, user_access
    is_active BOOLEAN DEFAULT true,
    is_system_template BOOLEAN DEFAULT false, -- Cannot be deleted
    created_by UUID, -- REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(tenant_id, name)
);

-- Workflow Steps (ordered within template)
CREATE TABLE IF NOT EXISTS workflow_steps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    template_id UUID NOT NULL REFERENCES workflow_templates(id) ON DELETE CASCADE,
    step_order INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    
    -- Who can approve this step
    required_role VARCHAR(50), -- Role required to approve
    required_permission VARCHAR(100), -- Or specific permission
    specific_user_id UUID, -- REFERENCES users(id), -- Or specific user
    
    -- Step configuration
    is_optional BOOLEAN DEFAULT false,
    allow_delegation BOOLEAN DEFAULT true,
    parallel_with_previous BOOLEAN DEFAULT false, -- Can run in parallel with previous step
    
    -- Timeout handling
    timeout_days INT DEFAULT 7,
    auto_approve_on_timeout BOOLEAN DEFAULT false,
    auto_reject_on_timeout BOOLEAN DEFAULT false,
    escalation_role VARCHAR(50), -- Role to escalate to on timeout
    
    -- Notifications
    notify_on_pending BOOLEAN DEFAULT true,
    reminder_days INT DEFAULT 3, -- Days before timeout to send reminder
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(template_id, step_order)
);

-- Workflow Instances (actual running workflows)
CREATE TABLE IF NOT EXISTS workflow_instances (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    template_id UUID NOT NULL REFERENCES workflow_templates(id),
    tenant_id UUID, -- REFERENCES tenants(id),
    
    -- What is being approved
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    entity_name VARCHAR(200), -- Cached for display
    
    -- Initiator
    initiated_by UUID NOT NULL, -- REFERENCES users(id),
    initiated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Current state
    current_step_id UUID REFERENCES workflow_steps(id),
    current_step_order INT DEFAULT 1,
    status VARCHAR(30) DEFAULT 'pending', -- pending, approved, rejected, cancelled, expired
    
    -- Completion
    completed_at TIMESTAMP WITH TIME ZONE,
    final_decision VARCHAR(30), -- approved, rejected
    final_comment TEXT,
    
    -- Metadata
    metadata JSONB, -- Additional context data
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_workflow_instances_entity ON workflow_instances(entity_type, entity_id);
CREATE INDEX idx_workflow_instances_status ON workflow_instances(status);
CREATE INDEX idx_workflow_instances_pending ON workflow_instances(status, current_step_id) WHERE status = 'pending';

-- Workflow Approvals (decisions on each step)
CREATE TABLE IF NOT EXISTS workflow_approvals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    instance_id UUID NOT NULL REFERENCES workflow_instances(id) ON DELETE CASCADE,
    step_id UUID NOT NULL REFERENCES workflow_steps(id),
    
    -- Who made the decision
    approver_id UUID, -- REFERENCES users(id),
    approver_role VARCHAR(50),
    delegated_from UUID, -- REFERENCES users(id), -- If delegated
    
    -- Decision
    decision VARCHAR(30) DEFAULT '', -- approved, rejected, returned, delegated
    comment TEXT,
    
    -- Timing
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    decided_at TIMESTAMP WITH TIME ZONE,
    due_at TIMESTAMP WITH TIME ZONE, -- Based on timeout_days
    
    -- Notifications sent
    notification_sent_at TIMESTAMP WITH TIME ZONE,
    reminder_sent_at TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_workflow_approvals_instance ON workflow_approvals(instance_id);
CREATE INDEX idx_workflow_approvals_pending ON workflow_approvals(approver_id, decided_at) WHERE decided_at IS NULL;

-- Workflow Delegation
CREATE TABLE IF NOT EXISTS workflow_delegations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    delegator_id UUID NOT NULL, -- REFERENCES users(id),
    delegate_id UUID NOT NULL, -- REFERENCES users(id),
    
    -- Scope of delegation
    workflow_type VARCHAR(50), -- NULL = all workflows
    role VARCHAR(50), -- Delegate acts as this role
    
    -- Validity period
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    reason TEXT,
    
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CHECK(end_date >= start_date)
);

CREATE INDEX idx_workflow_delegations_active ON workflow_delegations(delegate_id, is_active) WHERE is_active = true;

-- Insert default workflow templates
INSERT INTO workflow_templates (name, description, entity_type, is_system_template) VALUES
    ('Course Approval', 'Standard course approval workflow', 'course_approval', true),
    ('Schedule Approval', 'Multi-party schedule approval', 'schedule_approval', true),
    ('Thesis Stage Approval', 'PhD thesis stage progression', 'thesis_stage', true),
    ('User Access Request', 'Request for elevated access', 'user_access', true)
ON CONFLICT DO NOTHING;

-- Insert steps for Course Approval template
-- Note: UUIDs are dynamic in standard insert, we need CTE or lookup.
-- For simplicity in migration, we can rely on application seeding or do a sub-select insert.

INSERT INTO workflow_steps (template_id, step_order, name, required_role, timeout_days, reminder_days) 
SELECT id, 1, 'Department Chair Review', 'chair', 5, 2
FROM workflow_templates WHERE name = 'Course Approval'
ON CONFLICT DO NOTHING;

INSERT INTO workflow_steps (template_id, step_order, name, required_role, timeout_days, reminder_days) 
SELECT id, 2, 'Dean Approval', 'dean', 7, 3
FROM workflow_templates WHERE name = 'Course Approval'
ON CONFLICT DO NOTHING;
