-- Migration: Repair demo seed (idempotent re-run of 0047 fixes)

-- Create demo.university tenant
INSERT INTO tenants (id, slug, name, tenant_type, is_active, enabled_services, primary_color, secondary_color, app_name)
VALUES (
  'dd000000-0000-0000-0000-d00000000001',
  'demo',
  'Demo University',
  'university',
  true,
  ARRAY['chat', 'calendar'],
  '#059669', -- Emerald green
  '#047857',
  'Demo University PhD Portal'
) ON CONFLICT (slug) DO UPDATE SET
  enabled_services = ARRAY['chat', 'calendar'],
  is_active = true;

-- Get the demo tenant ID for later use
-- We use a fixed UUID for consistent referencing

-- Create demo admin user
-- Universal demo password: demopassword123! (bcrypt hash below)
INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active, created_at)
VALUES (
  'dd000001-0000-0000-0001-000000000001',
  'demo.admin',
  'admin@demo.university.edu',
  'Demo',
  'Admin',
  'admin',
  '$2a$10$3r/ZsOJoxgNK6mMj4Zl9FeYdq8pokmxAqk975/vu6JbEjeyXk6cR6', -- password: demopassword123!
  true,
  NOW()
) ON CONFLICT (username) DO NOTHING;

-- Create membership for demo admin
INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary)
VALUES (
  'dd000001-0000-0000-0001-000000000001',
  'dd000000-0000-0000-0000-d00000000001',
  'admin',
  true
) ON CONFLICT (user_id, tenant_id) DO NOTHING;

-- Create 5 Advisor users for demo tenant (same password: demopassword123!)
INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active, created_at)
VALUES
  ('dd000002-0000-0000-0001-000000000001', 'dr.johnson', 'johnson@demo.university.edu', 'Sarah', 'Johnson', 'advisor', '$2a$10$3r/ZsOJoxgNK6mMj4Zl9FeYdq8pokmxAqk975/vu6JbEjeyXk6cR6', true, NOW()),
  ('dd000003-0000-0000-0002-000000000002', 'dr.williams', 'williams@demo.university.edu', 'Michael', 'Williams', 'advisor', '$2a$10$3r/ZsOJoxgNK6mMj4Zl9FeYdq8pokmxAqk975/vu6JbEjeyXk6cR6', true, NOW()),
  ('dd000004-0000-0000-0003-000000000003', 'dr.chen', 'chen@demo.university.edu', 'Wei', 'Chen', 'advisor', '$2a$10$3r/ZsOJoxgNK6mMj4Zl9FeYdq8pokmxAqk975/vu6JbEjeyXk6cR6', true, NOW()),
  ('dd000005-0000-0000-0004-000000000004', 'dr.martinez', 'martinez@demo.university.edu', 'Elena', 'Martinez', 'advisor', '$2a$10$3r/ZsOJoxgNK6mMj4Zl9FeYdq8pokmxAqk975/vu6JbEjeyXk6cR6', true, NOW()),
  ('dd000006-0000-0000-0005-000000000005', 'dr.thompson', 'thompson@demo.university.edu', 'James', 'Thompson', 'advisor', '$2a$10$3r/ZsOJoxgNK6mMj4Zl9FeYdq8pokmxAqk975/vu6JbEjeyXk6cR6', true, NOW())
ON CONFLICT (username) DO NOTHING;

-- Create advisor memberships
INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary)
VALUES
  ('dd000002-0000-0000-0001-000000000001', 'dd000000-0000-0000-0000-d00000000001', 'advisor', true),
  ('dd000003-0000-0000-0002-000000000002', 'dd000000-0000-0000-0000-d00000000001', 'advisor', true),
  ('dd000004-0000-0000-0003-000000000003', 'dd000000-0000-0000-0000-d00000000001', 'advisor', true),
  ('dd000005-0000-0000-0004-000000000004', 'dd000000-0000-0000-0000-d00000000001', 'advisor', true),
  ('dd000006-0000-0000-0005-000000000005', 'dd000000-0000-0000-0000-d00000000001', 'advisor', true)
ON CONFLICT (user_id, tenant_id) DO NOTHING;

-- Create 24 Student users (distributed across stages) - password: demopassword123!
INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active, created_at)
VALUES
  -- Year 1 Students (Stage: Coursework - 6 students)
  ('dd001001-0000-0000-0001-000000000001', 'demo.student1', 'student1@demo.university.edu', 'Emma', 'Brown', 'student', '$2a$10$3r/ZsOJoxgNK6mMj4Zl9FeYdq8pokmxAqk975/vu6JbEjeyXk6cR6', true, NOW()),
  ('dd001002-0000-0000-0002-000000000002', 'demo.student2', 'student2@demo.university.edu', 'Liam', 'Davis', 'student', '$2a$10$3r/ZsOJoxgNK6mMj4Zl9FeYdq8pokmxAqk975/vu6JbEjeyXk6cR6', true, NOW()),
  ('dd001003-0000-0000-0003-000000000003', 'demo.student3', 'student3@demo.university.edu', 'Sophia', 'Garcia', 'student', '$2a$10$3r/ZsOJoxgNK6mMj4Zl9FeYdq8pokmxAqk975/vu6JbEjeyXk6cR6', true, NOW()),
  ('dd001004-0000-0000-0004-000000000004', 'demo.student4', 'student4@demo.university.edu', 'Noah', 'Miller', 'student', '$2a$10$3r/ZsOJoxgNK6mMj4Zl9FeYdq8pokmxAqk975/vu6JbEjeyXk6cR6', true, NOW()),
  ('dd001005-0000-0000-0005-000000000005', 'demo.student5', 'student5@demo.university.edu', 'Ava', 'Wilson', 'student', '$2a$10$3r/ZsOJoxgNK6mMj4Zl9FeYdq8pokmxAqk975/vu6JbEjeyXk6cR6', true, NOW()),
  ('dd001002-0000-0000-0002-000000000002', 'demo.student2', 'student2@demo.university.edu', 'Liam', 'Davis', 'student', '$2a$10$3r/ZsOJoxgNK6mMj4Zl9FeYdq8pokmxAqk975/vu6JbEjeyXk6cR6', true, NOW())
ON CONFLICT (username) DO NOTHING;

-- Create student memberships
INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary)
SELECT id, 'dd000000-0000-0000-0000-d00000000001', 'student', true
FROM users 
WHERE id::text LIKE 'dd001%'
ON CONFLICT (user_id, tenant_id) DO NOTHING;
