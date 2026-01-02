-- Add roles array column
ALTER TABLE user_tenant_memberships ADD COLUMN roles text[];

-- Migrate existing single role to array
UPDATE user_tenant_memberships SET roles = ARRAY[role];

-- Make single role nullable (deprecate it)
ALTER TABLE user_tenant_memberships ALTER COLUMN role DROP NOT NULL;
