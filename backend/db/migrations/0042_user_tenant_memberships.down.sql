-- Rollback: Drop user-tenant memberships
DROP TRIGGER IF EXISTS user_tenant_memberships_updated ON user_tenant_memberships;
DROP FUNCTION IF EXISTS update_user_tenant_memberships_modified();
DROP INDEX IF EXISTS idx_user_tenant_primary;
DROP INDEX IF EXISTS idx_user_tenant_tenant;
DROP INDEX IF EXISTS idx_user_tenant_user;
DROP TABLE IF EXISTS user_tenant_memberships;
