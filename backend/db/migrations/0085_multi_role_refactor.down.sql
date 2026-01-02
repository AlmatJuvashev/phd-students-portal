-- Revert: Copy first role back to single column (lossy if multiple roles exist, but acceptable for down)
UPDATE user_tenant_memberships SET role = roles[1] WHERE roles IS NOT NULL AND array_length(roles, 1) > 0;

-- Restore NOT NULL constraint (will fail if role is null, so we assume migration worked or data is clean)
-- Actually, better to just drop the array.
ALTER TABLE user_tenant_memberships DROP COLUMN roles;

-- We can't easily restore NOT NULL on 'role' without ensuring data integrity, 
-- but we can try if we assume rollback is immediate.
ALTER TABLE user_tenant_memberships ALTER COLUMN role SET NOT NULL;
