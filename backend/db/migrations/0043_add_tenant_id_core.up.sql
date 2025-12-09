-- Multitenancy: Add tenant_id to core data tables
-- All tenant-scoped data gets a tenant_id foreign key

-- Documents
ALTER TABLE documents ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE documents d SET tenant_id = '00000000-0000-0000-0000-000000000001' 
WHERE tenant_id IS NULL;
ALTER TABLE documents ALTER COLUMN tenant_id SET NOT NULL;
CREATE INDEX idx_documents_tenant ON documents(tenant_id);

-- Document versions (inherits tenant from document, but add for RLS)
ALTER TABLE document_versions ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE document_versions dv SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
ALTER TABLE document_versions ALTER COLUMN tenant_id SET NOT NULL;
CREATE INDEX idx_document_versions_tenant ON document_versions(tenant_id);

-- Node instances (journey progress)
ALTER TABLE node_instances ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE node_instances SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
ALTER TABLE node_instances ALTER COLUMN tenant_id SET NOT NULL;
CREATE INDEX idx_node_instances_tenant ON node_instances(tenant_id);

-- Node instance slots (node attachments)
ALTER TABLE node_instance_slots ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE node_instance_slots SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
ALTER TABLE node_instance_slots ALTER COLUMN tenant_id SET NOT NULL;
CREATE INDEX idx_node_instance_slots_tenant ON node_instance_slots(tenant_id);

-- Journey states
ALTER TABLE journey_states ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE journey_states SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
ALTER TABLE journey_states ALTER COLUMN tenant_id SET NOT NULL;
CREATE INDEX idx_journey_states_tenant ON journey_states(tenant_id);

-- Checklist modules (tenant-specific curriculum)
ALTER TABLE checklist_modules ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE checklist_modules SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
CREATE INDEX idx_checklist_modules_tenant ON checklist_modules(tenant_id);

-- Checklist steps
ALTER TABLE checklist_steps ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE checklist_steps SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
CREATE INDEX idx_checklist_steps_tenant ON checklist_steps(tenant_id);

-- Student steps
ALTER TABLE student_steps ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE student_steps SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
CREATE INDEX idx_student_steps_tenant ON student_steps(tenant_id);

-- Node deadlines
ALTER TABLE node_deadlines ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE node_deadlines SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
CREATE INDEX idx_node_deadlines_tenant ON node_deadlines(tenant_id);

-- Reminders
ALTER TABLE reminders ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE reminders SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
CREATE INDEX idx_reminders_tenant ON reminders(tenant_id);
