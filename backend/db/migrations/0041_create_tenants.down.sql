-- Rollback: Drop tenants table
DROP INDEX IF EXISTS idx_tenants_active;
DROP INDEX IF EXISTS idx_tenants_domain;
DROP INDEX IF EXISTS idx_tenants_slug;
DROP TABLE IF EXISTS tenants;
ALTER TABLE users DROP COLUMN IF EXISTS is_superadmin;
