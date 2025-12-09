-- Rollback: Remove superadmin tables and tenant extensions

DROP TABLE IF EXISTS global_settings;
DROP TABLE IF EXISTS activity_logs;

ALTER TABLE tenants DROP COLUMN IF EXISTS tenant_type;
ALTER TABLE tenants DROP COLUMN IF EXISTS app_name;
ALTER TABLE tenants DROP COLUMN IF EXISTS primary_color;
ALTER TABLE tenants DROP COLUMN IF EXISTS secondary_color;
