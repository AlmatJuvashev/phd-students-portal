-- Down migration: Remove demo.university tenant and related data
-- This removes all demo data created in the up migration

-- Remove specialty-program links
DELETE FROM specialty_programs WHERE specialty_id LIKE 'dd100%';

-- Remove departments
DELETE FROM departments WHERE tenant_id = 'dd000000-0000-0000-0000-demo00000001';

-- Remove cohorts
DELETE FROM cohorts WHERE tenant_id = 'dd000000-0000-0000-0000-demo00000001';

-- Remove programs
DELETE FROM programs WHERE tenant_id = 'dd000000-0000-0000-0000-demo00000001';

-- Remove specialties
DELETE FROM specialties WHERE tenant_id = 'dd000000-0000-0000-0000-demo00000001';

-- Remove user memberships for demo tenant
DELETE FROM user_tenant_memberships WHERE tenant_id = 'dd000000-0000-0000-0000-demo00000001';

-- Remove demo users (students and advisors)
DELETE FROM users WHERE id LIKE 'dd000%' OR id LIKE 'dd001%';

-- Remove demo tenant
DELETE FROM tenants WHERE id = 'dd000000-0000-0000-0000-demo00000001';

