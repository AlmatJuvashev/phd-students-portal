-- Rollback: Remove enabled_services column
ALTER TABLE tenants DROP COLUMN IF EXISTS enabled_services;
