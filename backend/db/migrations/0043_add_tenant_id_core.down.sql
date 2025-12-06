-- Rollback: Remove tenant_id from core tables
DROP INDEX IF EXISTS idx_reminders_tenant;
ALTER TABLE reminders DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_node_deadlines_tenant;
ALTER TABLE node_deadlines DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_student_steps_tenant;
ALTER TABLE student_steps DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_checklist_steps_tenant;
ALTER TABLE checklist_steps DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_checklist_modules_tenant;
ALTER TABLE checklist_modules DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_journey_states_tenant;
ALTER TABLE journey_states DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_node_instance_slots_tenant;
ALTER TABLE node_instance_slots DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_node_instances_tenant;
ALTER TABLE node_instances DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_document_versions_tenant;
ALTER TABLE document_versions DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_documents_tenant;
ALTER TABLE documents DROP COLUMN IF EXISTS tenant_id;
