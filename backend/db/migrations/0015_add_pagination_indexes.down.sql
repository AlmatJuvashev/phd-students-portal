-- Rollback pagination indexes

DROP INDEX IF EXISTS idx_users_last_name_active;
DROP INDEX IF EXISTS idx_users_role_active;
DROP INDEX IF EXISTS idx_profile_submissions_user_submitted;
