-- Add tenant_id to profile_submissions
ALTER TABLE profile_submissions ADD COLUMN IF NOT EXISTS tenant_id uuid REFERENCES tenants(id);

-- Backfill existing records with default tenant if any exist
UPDATE profile_submissions SET tenant_id = '00000000-0000-0000-0000-000000000001'
WHERE tenant_id IS NULL;

-- Enforce NOT NULL after backfill
ALTER TABLE profile_submissions ALTER COLUMN tenant_id SET NOT NULL;

-- Add index for performance
CREATE INDEX idx_profile_submissions_tenant ON profile_submissions(tenant_id);
