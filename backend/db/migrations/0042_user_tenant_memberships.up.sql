-- Multitenancy: Create user-tenant memberships for multi-tenant users
-- This allows users (especially advisors) to belong to multiple tenants

CREATE TABLE user_tenant_memberships (
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  tenant_id uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  role user_role NOT NULL,                -- Role within this tenant
  is_primary boolean DEFAULT false,       -- User's primary/default tenant
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  PRIMARY KEY (user_id, tenant_id)
);

-- Indexes for common queries
CREATE INDEX idx_user_tenant_user ON user_tenant_memberships(user_id);
CREATE INDEX idx_user_tenant_tenant ON user_tenant_memberships(tenant_id);
CREATE INDEX idx_user_tenant_primary ON user_tenant_memberships(user_id) WHERE is_primary = true;

-- Migrate existing users to default tenant (KazNMU)
INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary)
SELECT id, '00000000-0000-0000-0000-000000000001', role, true
FROM users
WHERE NOT EXISTS (
  SELECT 1 FROM user_tenant_memberships utm WHERE utm.user_id = users.id
);

-- Trigger to update updated_at
CREATE OR REPLACE FUNCTION update_user_tenant_memberships_modified()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_tenant_memberships_updated
  BEFORE UPDATE ON user_tenant_memberships
  FOR EACH ROW
  EXECUTE FUNCTION update_user_tenant_memberships_modified();
