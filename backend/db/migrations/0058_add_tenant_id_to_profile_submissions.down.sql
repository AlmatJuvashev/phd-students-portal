DROP INDEX IF EXISTS idx_profile_submissions_tenant;
ALTER TABLE profile_submissions DROP COLUMN IF EXISTS tenant_id;
