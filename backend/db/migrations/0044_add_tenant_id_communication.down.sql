-- Rollback: Remove tenant_id from communication and reference tables
DROP INDEX IF EXISTS idx_departments_tenant;
ALTER TABLE departments DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_cohorts_tenant;
ALTER TABLE cohorts DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_specialties_tenant;
ALTER TABLE specialties DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_programs_tenant;
ALTER TABLE programs DROP COLUMN IF EXISTS tenant_id;

ALTER TABLE playbook_active_version DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_playbook_versions_tenant;
ALTER TABLE playbook_versions DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_comments_tenant;
ALTER TABLE comments DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_contacts_tenant;
ALTER TABLE contacts DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_notifications_tenant;
ALTER TABLE notifications DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_event_attendees_tenant;
ALTER TABLE event_attendees DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_events_tenant;
ALTER TABLE events DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_chat_room_members_tenant;
ALTER TABLE chat_room_members DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_chat_messages_tenant;
ALTER TABLE chat_messages DROP COLUMN IF EXISTS tenant_id;

DROP INDEX IF EXISTS idx_chat_rooms_tenant;
ALTER TABLE chat_rooms DROP COLUMN IF EXISTS tenant_id;
