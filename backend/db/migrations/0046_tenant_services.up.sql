-- Add enabled_services to tenants for optional features
-- Default services are always enabled: journey, contacts, notifications, uploads
-- Optional services: chat, calendar

ALTER TABLE tenants ADD COLUMN IF NOT EXISTS enabled_services text[] DEFAULT ARRAY['chat', 'calendar'];

-- All existing tenants get all services enabled by default
UPDATE tenants SET enabled_services = ARRAY['chat', 'calendar'] WHERE enabled_services IS NULL;

COMMENT ON COLUMN tenants.enabled_services IS 'Optional services enabled for this tenant. Valid values: chat, calendar. Core services (journey, contacts, notifications, uploads) are always enabled.';
