-- Superadmin Portal: Extended tenant fields, activity logs, global settings

-- Add tenant type for different institution types
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS tenant_type text DEFAULT 'university';
-- Types: 'university', 'college', 'vocational', 'school'

-- Add branding/customization fields
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS app_name text;           -- Custom app name per tenant
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS primary_color text DEFAULT '#3b82f6';
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS secondary_color text DEFAULT '#1e40af';

-- Activity logs table for tracking all actions
CREATE TABLE activity_logs (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  tenant_id uuid REFERENCES tenants(id) ON DELETE SET NULL,
  user_id uuid REFERENCES users(id) ON DELETE SET NULL,
  action text NOT NULL,              -- 'login', 'logout', 'create', 'update', 'delete', etc.
  entity_type text,                  -- 'user', 'document', 'node', 'tenant', etc.
  entity_id uuid,                    -- ID of the affected entity
  description text,                  -- Human-readable description
  metadata jsonb DEFAULT '{}',       -- Additional context (old/new values, etc.)
  ip_address inet,                   -- Client IP address
  user_agent text,                   -- Browser/client info
  created_at timestamptz DEFAULT now()
);

-- Indexes for efficient querying
CREATE INDEX idx_activity_logs_tenant ON activity_logs(tenant_id);
CREATE INDEX idx_activity_logs_user ON activity_logs(user_id);
CREATE INDEX idx_activity_logs_action ON activity_logs(action);
CREATE INDEX idx_activity_logs_entity ON activity_logs(entity_type, entity_id);
CREATE INDEX idx_activity_logs_created ON activity_logs(created_at DESC);

-- Global settings table (platform-wide configuration)
CREATE TABLE global_settings (
  key text PRIMARY KEY,
  value jsonb NOT NULL,
  description text,                  -- What this setting does
  category text DEFAULT 'general',   -- For grouping in UI
  updated_at timestamptz DEFAULT now(),
  updated_by uuid REFERENCES users(id) ON DELETE SET NULL
);

-- Insert default settings
INSERT INTO global_settings (key, value, description, category) VALUES
  ('maintenance_mode', 'false', 'Enable maintenance mode for the entire platform', 'system'),
  ('allow_new_tenants', 'true', 'Allow creation of new tenants', 'system'),
  ('default_tenant_type', '"university"', 'Default type for new tenants', 'defaults'),
  ('max_file_size_mb', '50', 'Maximum file upload size in MB', 'limits'),
  ('session_timeout_hours', '24', 'User session timeout in hours', 'security')
ON CONFLICT (key) DO NOTHING;
