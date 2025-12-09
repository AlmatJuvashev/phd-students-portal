-- Multitenancy: Add tenant_id to chat, events, notifications, and reference tables

-- Chat rooms
ALTER TABLE chat_rooms ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE chat_rooms SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
ALTER TABLE chat_rooms ALTER COLUMN tenant_id SET NOT NULL;
CREATE INDEX idx_chat_rooms_tenant ON chat_rooms(tenant_id);

-- Chat messages (tenant from room, but add for RLS)
ALTER TABLE chat_messages ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE chat_messages SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
ALTER TABLE chat_messages ALTER COLUMN tenant_id SET NOT NULL;
CREATE INDEX idx_chat_messages_tenant ON chat_messages(tenant_id);

-- Chat room members
ALTER TABLE chat_room_members ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE chat_room_members SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
ALTER TABLE chat_room_members ALTER COLUMN tenant_id SET NOT NULL;
CREATE INDEX idx_chat_room_members_tenant ON chat_room_members(tenant_id);

-- Events (calendar)
ALTER TABLE events ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE events SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
ALTER TABLE events ALTER COLUMN tenant_id SET NOT NULL;
CREATE INDEX idx_events_tenant ON events(tenant_id);

-- Event attendees
ALTER TABLE event_attendees ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE event_attendees SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
ALTER TABLE event_attendees ALTER COLUMN tenant_id SET NOT NULL;
CREATE INDEX idx_event_attendees_tenant ON event_attendees(tenant_id);

-- Notifications
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE notifications SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
ALTER TABLE notifications ALTER COLUMN tenant_id SET NOT NULL;
CREATE INDEX idx_notifications_tenant ON notifications(tenant_id);

-- Contacts
ALTER TABLE contacts ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE contacts SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
ALTER TABLE contacts ALTER COLUMN tenant_id SET NOT NULL;
CREATE INDEX idx_contacts_tenant ON contacts(tenant_id);

-- Comments
ALTER TABLE comments ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE comments SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
ALTER TABLE comments ALTER COLUMN tenant_id SET NOT NULL;
CREATE INDEX idx_comments_tenant ON comments(tenant_id);

-- Playbook versions (tenant-specific playbooks)
ALTER TABLE playbook_versions ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE playbook_versions SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
ALTER TABLE playbook_versions ALTER COLUMN tenant_id SET NOT NULL;
CREATE INDEX idx_playbook_versions_tenant ON playbook_versions(tenant_id);

-- Playbook active version (make tenant-specific)
ALTER TABLE playbook_active_version ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE playbook_active_version SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
ALTER TABLE playbook_active_version ALTER COLUMN tenant_id SET NOT NULL;

-- Reference data: programs, specialties, cohorts, departments
ALTER TABLE programs ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE programs SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
CREATE INDEX idx_programs_tenant ON programs(tenant_id);

ALTER TABLE specialties ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE specialties SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
CREATE INDEX idx_specialties_tenant ON specialties(tenant_id);

ALTER TABLE cohorts ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE cohorts SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
CREATE INDEX idx_cohorts_tenant ON cohorts(tenant_id);

ALTER TABLE departments ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);
UPDATE departments SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;
CREATE INDEX idx_departments_tenant ON departments(tenant_id);
