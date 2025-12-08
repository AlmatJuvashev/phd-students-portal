-- Migration: Create demo.university tenant with comprehensive sample data
-- This migration seeds a demo university with:
-- - 1 demo tenant (demo.university) with all services enabled
-- - 4-5 advisors
-- - 20+ students across different stages
-- - Public health/healthcare focused dictionaries (specialties, programs, departments)

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
  ('dd000002-0000-0000-0001-000000000001', 'dr.johnson', 'johnson@demo.university.edu', 'Sarah', 'Johnson', 'advisor', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd000003-0000-0000-0002-000000000002', 'dr.williams', 'williams@demo.university.edu', 'Michael', 'Williams', 'advisor', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd000004-0000-0000-0003-000000000003', 'dr.chen', 'chen@demo.university.edu', 'Wei', 'Chen', 'advisor', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd000005-0000-0000-0004-000000000004', 'dr.martinez', 'martinez@demo.university.edu', 'Elena', 'Martinez', 'advisor', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd000006-0000-0000-0005-000000000005', 'dr.thompson', 'thompson@demo.university.edu', 'James', 'Thompson', 'advisor', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW())
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
  ('dd001001-0000-0000-0001-000000000001', 'demo.student1', 'student1@demo.university.edu', 'Emma', 'Brown', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001002-0000-0000-0002-000000000002', 'demo.student2', 'student2@demo.university.edu', 'Liam', 'Davis', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001003-0000-0000-0003-000000000003', 'demo.student3', 'student3@demo.university.edu', 'Sophia', 'Garcia', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001004-0000-0000-0004-000000000004', 'demo.student4', 'student4@demo.university.edu', 'Noah', 'Miller', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001005-0000-0000-0005-000000000005', 'demo.student5', 'student5@demo.university.edu', 'Ava', 'Wilson', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001006-0000-0000-0006-000000000006', 'demo.student6', 'student6@demo.university.edu', 'William', 'Moore', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  
  -- Year 2 Students (Stage: Qualifying Exams - 5 students)
  ('dd001007-0000-0000-0007-000000000007', 'demo.student7', 'student7@demo.university.edu', 'Isabella', 'Taylor', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001008-0000-0000-0008-000000000008', 'demo.student8', 'student8@demo.university.edu', 'James', 'Anderson', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001009-0000-0000-0009-000000000009', 'demo.student9', 'student9@demo.university.edu', 'Mia', 'Thomas', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001010-0000-0000-0010-000000000010', 'demo.student10', 'student10@demo.university.edu', 'Benjamin', 'Jackson', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001011-0000-0000-0011-000000000011', 'demo.student11', 'student11@demo.university.edu', 'Charlotte', 'White', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  
  -- Year 3 Students (Stage: Proposal Writing - 5 students)
  ('dd001012-0000-0000-0012-000000000012', 'demo.student12', 'student12@demo.university.edu', 'Elijah', 'Harris', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001013-0000-0000-0013-000000000013', 'demo.student13', 'student13@demo.university.edu', 'Amelia', 'Martin', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001014-0000-0000-0014-000000000014', 'demo.student14', 'student14@demo.university.edu', 'Oliver', 'Lee', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001015-0000-0000-0015-000000000015', 'demo.student15', 'student15@demo.university.edu', 'Harper', 'Perez', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001016-0000-0000-0016-000000000016', 'demo.student16', 'student16@demo.university.edu', 'Ethan', 'Thompson', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  
  -- Year 4 Students (Stage: Research/Dissertation - 5 students)
  ('dd001017-0000-0000-0017-000000000017', 'demo.student17', 'student17@demo.university.edu', 'Evelyn', 'Garcia', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001018-0000-0000-0018-000000000018', 'demo.student18', 'student18@demo.university.edu', 'Aiden', 'Martinez', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001019-0000-0000-0019-000000000019', 'demo.student19', 'student19@demo.university.edu', 'Luna', 'Robinson', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001020-0000-0000-0020-000000000020', 'demo.student20', 'student20@demo.university.edu', 'Lucas', 'Clark', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001021-0000-0000-0021-000000000021', 'demo.student21', 'student21@demo.university.edu', 'Chloe', 'Rodriguez', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  
  -- Year 5 Students (Stage: Defense Preparation - 3 students)
  ('dd001022-0000-0000-0022-000000000022', 'demo.student22', 'student22@demo.university.edu', 'Mason', 'Lewis', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001023-0000-0000-0023-000000000023', 'demo.student23', 'student23@demo.university.edu', 'Ella', 'Walker', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001024-0000-0000-0024-000000000024', 'demo.student24', 'student24@demo.university.edu', 'Jacob', 'Hall', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW())
ON CONFLICT (username) DO NOTHING;

-- Create student memberships
INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary)
SELECT id, 'dd000000-0000-0000-0000-d00000000001', 'student', true
FROM users 
WHERE id::text LIKE 'dd001%'
ON CONFLICT (user_id, tenant_id) DO NOTHING;

-- Create Public Health & Healthcare focused Specialties for demo tenant
INSERT INTO specialties (id, name, code, tenant_id, is_active, created_at)
VALUES
  ('dd100001-0000-0000-0001-000000000001', 'Epidemiology', 'EPD', 'dd000000-0000-0000-0000-d00000000001', true, NOW()),
  ('dd100002-0000-0000-0002-000000000002', 'Public Health', 'PBH', 'dd000000-0000-0000-0000-d00000000001', true, NOW()),
  ('dd100003-0000-0000-0003-000000000003', 'Health Policy & Management', 'HPM', 'dd000000-0000-0000-0000-d00000000001', true, NOW()),
  ('dd100004-0000-0000-0004-000000000004', 'Global Health', 'GLH', 'dd000000-0000-0000-0000-d00000000001', true, NOW()),
  ('dd100005-0000-0000-0005-000000000005', 'Environmental Health Sciences', 'EHS', 'dd000000-0000-0000-0000-d00000000001', true, NOW()),
  ('dd100006-0000-0000-0006-000000000006', 'Biostatistics', 'BST', 'dd000000-0000-0000-0000-d00000000001', true, NOW()),
  ('dd100007-0000-0000-0007-000000000007', 'Health Behavior', 'HBV', 'dd000000-0000-0000-0000-d00000000001', true, NOW()),
  ('dd100008-0000-0000-0008-000000000008', 'Clinical Research', 'CLR', 'dd000000-0000-0000-0000-d00000000001', true, NOW())
ON CONFLICT (id) DO NOTHING;

-- Create Programs for demo tenant
INSERT INTO programs (id, name, code, tenant_id, is_active, created_at)
VALUES
  ('dd200001-0000-0000-0001-000000000001', 'Doctor of Public Health (DrPH)', 'DRPH', 'dd000000-0000-0000-0000-d00000000001', true, NOW()),
  ('dd200002-0000-0000-0002-000000000002', 'PhD in Public Health', 'PHDPH', 'dd000000-0000-0000-0000-d00000000001', true, NOW()),
  ('dd200003-0000-0000-0003-000000000003', 'PhD in Health Services Management', 'PHDHSM', 'dd000000-0000-0000-0000-d00000000001', true, NOW()),
  ('dd200004-0000-0000-0004-000000000004', 'PhD in Epidemiology', 'PHDEPD', 'dd000000-0000-0000-0000-d00000000001', true, NOW()),
  ('dd200005-0000-0000-0005-000000000005', 'PhD in Biostatistics', 'PHDBST', 'dd000000-0000-0000-0000-d00000000001', true, NOW())
ON CONFLICT (id) DO NOTHING;

-- Create Cohorts for demo tenant
INSERT INTO cohorts (id, name, start_date, tenant_id, is_active, created_at)
VALUES
  ('dd300001-0000-0000-0001-000000002020', 'Cohort 2020', '2020-09-01', 'dd000000-0000-0000-0000-d00000000001', true, NOW()),
  ('dd300002-0000-0000-0002-000000002021', 'Cohort 2021', '2021-09-01', 'dd000000-0000-0000-0000-d00000000001', true, NOW()),
  ('dd300003-0000-0000-0003-000000002022', 'Cohort 2022', '2022-09-01', 'dd000000-0000-0000-0000-d00000000001', true, NOW()),
  ('dd300004-0000-0000-0004-000000002023', 'Cohort 2023', '2023-09-01', 'dd000000-0000-0000-0000-d00000000001', true, NOW()),
  ('dd300005-0000-0000-0005-000000002024', 'Cohort 2024', '2024-09-01', 'dd000000-0000-0000-0000-d00000000001', true, NOW())
ON CONFLICT (id) DO NOTHING;

-- Create Departments for demo tenant
INSERT INTO departments (id, name, code, tenant_id, is_active, created_at)
VALUES
  ('dd400001-0000-0000-0001-000000000001', 'Department of Epidemiology', 'EPID', 'dd000000-0000-0000-0000-d00000000001', true, NOW()),
  ('dd400002-0000-0000-0002-000000000002', 'Department of Health Policy', 'HPOL', 'dd000000-0000-0000-0000-d00000000001', true, NOW()),
  ('dd400003-0000-0000-0003-000000000003', 'Department of Biostatistics', 'BIOS', 'dd000000-0000-0000-0000-d00000000001', true, NOW()),
  ('dd400004-0000-0000-0004-000000000004', 'Department of Environmental Health', 'ENVH', 'dd000000-0000-0000-0000-d00000000001', true, NOW()),
  ('dd400005-0000-0000-0005-000000000005', 'Department of Global Health', 'GLBH', 'dd000000-0000-0000-0000-d00000000001', true, NOW()),
  ('dd400006-0000-0000-0006-000000000006', 'Department of Behavioral Sciences', 'BHVS', 'dd000000-0000-0000-0000-d00000000001', true, NOW())
ON CONFLICT (id) DO NOTHING;

-- Link Programs to Specialties (many-to-many relationship)
INSERT INTO specialty_programs (specialty_id, program_id)
VALUES
  ('dd100001-0000-0000-0001-000000000001', 'dd200004-0000-0000-0004-000000000004'),
  ('dd100002-0000-0000-0002-000000000002', 'dd200001-0000-0000-0001-000000000001'),
  ('dd100002-0000-0000-0002-000000000002', 'dd200002-0000-0000-0002-000000000002'),
  ('dd100003-0000-0000-0003-000000000003', 'dd200003-0000-0000-0003-000000000003'),
  ('dd100004-0000-0000-0004-000000000004', 'dd200002-0000-0000-0002-000000000002'),
  ('dd100005-0000-0000-0005-000000000005', 'dd200002-0000-0000-0002-000000000002'),
  ('dd100006-0000-0000-0006-000000000006', 'dd200005-0000-0000-0005-000000000005'),
  ('dd100007-0000-0000-0007-000000000007', 'dd200001-0000-0000-0001-000000000001'),
  ('dd100008-0000-0000-0008-000000000008', 'dd200002-0000-0000-0002-000000000002')
ON CONFLICT DO NOTHING;
