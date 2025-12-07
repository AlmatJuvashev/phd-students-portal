-- Migration: Create demo.university tenant with comprehensive sample data
-- This migration seeds a demo university with:
-- - 1 demo tenant (demo.university) with all services enabled
-- - 4-5 advisors
-- - 20+ students across different stages
-- - Public health/healthcare focused dictionaries (specialties, programs, departments)

-- Create demo.university tenant
INSERT INTO tenants (id, slug, name, tenant_type, is_active, enabled_services, primary_color, secondary_color, app_name)
VALUES (
  'dd000000-demo-demo-demo-demouniversity',
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
  'dd000001-demo-admi-0001-demoadmin001',
  'demo.admin',
  'admin@demo.university.edu',
  'Demo',
  'Admin',
  'admin',
  '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', -- password: demopassword123!
  true,
  NOW()
) ON CONFLICT (username) DO NOTHING;

-- Create membership for demo admin
INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary)
VALUES (
  'dd000001-demo-admi-0001-demoadmin001',
  'dd000000-demo-demo-demo-demouniversity',
  'admin',
  true
) ON CONFLICT (user_id, tenant_id) DO NOTHING;

-- Create 5 Advisor users for demo tenant (same password: demopassword123!)
INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active, created_at)
VALUES
  ('dd000002-demo-advi-0001-demoadvisor1', 'dr.johnson', 'johnson@demo.university.edu', 'Sarah', 'Johnson', 'advisor', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd000003-demo-advi-0002-demoadvisor2', 'dr.williams', 'williams@demo.university.edu', 'Michael', 'Williams', 'advisor', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd000004-demo-advi-0003-demoadvisor3', 'dr.chen', 'chen@demo.university.edu', 'Wei', 'Chen', 'advisor', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd000005-demo-advi-0004-demoadvisor4', 'dr.martinez', 'martinez@demo.university.edu', 'Elena', 'Martinez', 'advisor', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd000006-demo-advi-0005-demoadvisor5', 'dr.thompson', 'thompson@demo.university.edu', 'James', 'Thompson', 'advisor', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW())
ON CONFLICT (username) DO NOTHING;

-- Create advisor memberships
INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary)
VALUES
  ('dd000002-demo-advi-0001-demoadvisor1', 'dd000000-demo-demo-demo-demouniversity', 'advisor', true),
  ('dd000003-demo-advi-0002-demoadvisor2', 'dd000000-demo-demo-demo-demouniversity', 'advisor', true),
  ('dd000004-demo-advi-0003-demoadvisor3', 'dd000000-demo-demo-demo-demouniversity', 'advisor', true),
  ('dd000005-demo-advi-0004-demoadvisor4', 'dd000000-demo-demo-demo-demouniversity', 'advisor', true),
  ('dd000006-demo-advi-0005-demoadvisor5', 'dd000000-demo-demo-demo-demouniversity', 'advisor', true)
ON CONFLICT (user_id, tenant_id) DO NOTHING;

-- Create 24 Student users (distributed across stages) - password: demopassword123!
INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active, created_at)
VALUES
  -- Year 1 Students (Stage: Coursework - 6 students)
  ('dd001001-demo-stud-0001-demostudent1', 'demo.student1', 'student1@demo.university.edu', 'Emma', 'Brown', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001002-demo-stud-0002-demostudent2', 'demo.student2', 'student2@demo.university.edu', 'Liam', 'Davis', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001003-demo-stud-0003-demostudent3', 'demo.student3', 'student3@demo.university.edu', 'Sophia', 'Garcia', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001004-demo-stud-0004-demostudent4', 'demo.student4', 'student4@demo.university.edu', 'Noah', 'Miller', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001005-demo-stud-0005-demostudent5', 'demo.student5', 'student5@demo.university.edu', 'Ava', 'Wilson', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001006-demo-stud-0006-demostudent6', 'demo.student6', 'student6@demo.university.edu', 'William', 'Moore', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  
  -- Year 2 Students (Stage: Qualifying Exams - 5 students)
  ('dd001007-demo-stud-0007-demostudent7', 'demo.student7', 'student7@demo.university.edu', 'Isabella', 'Taylor', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001008-demo-stud-0008-demostudent8', 'demo.student8', 'student8@demo.university.edu', 'James', 'Anderson', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001009-demo-stud-0009-demostudent9', 'demo.student9', 'student9@demo.university.edu', 'Mia', 'Thomas', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001010-demo-stud-0010-demostude10', 'demo.student10', 'student10@demo.university.edu', 'Benjamin', 'Jackson', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001011-demo-stud-0011-demostude11', 'demo.student11', 'student11@demo.university.edu', 'Charlotte', 'White', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  
  -- Year 3 Students (Stage: Proposal Writing - 5 students)
  ('dd001012-demo-stud-0012-demostude12', 'demo.student12', 'student12@demo.university.edu', 'Elijah', 'Harris', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001013-demo-stud-0013-demostude13', 'demo.student13', 'student13@demo.university.edu', 'Amelia', 'Martin', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001014-demo-stud-0014-demostude14', 'demo.student14', 'student14@demo.university.edu', 'Oliver', 'Lee', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001015-demo-stud-0015-demostude15', 'demo.student15', 'student15@demo.university.edu', 'Harper', 'Perez', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001016-demo-stud-0016-demostude16', 'demo.student16', 'student16@demo.university.edu', 'Ethan', 'Thompson', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  
  -- Year 4 Students (Stage: Research/Dissertation - 5 students)
  ('dd001017-demo-stud-0017-demostude17', 'demo.student17', 'student17@demo.university.edu', 'Evelyn', 'Garcia', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001018-demo-stud-0018-demostude18', 'demo.student18', 'student18@demo.university.edu', 'Aiden', 'Martinez', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001019-demo-stud-0019-demostude19', 'demo.student19', 'student19@demo.university.edu', 'Luna', 'Robinson', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001020-demo-stud-0020-demostude20', 'demo.student20', 'student20@demo.university.edu', 'Lucas', 'Clark', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001021-demo-stud-0021-demostude21', 'demo.student21', 'student21@demo.university.edu', 'Chloe', 'Rodriguez', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  
  -- Year 5 Students (Stage: Defense Preparation - 3 students)
  ('dd001022-demo-stud-0022-demostude22', 'demo.student22', 'student22@demo.university.edu', 'Mason', 'Lewis', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001023-demo-stud-0023-demostude23', 'demo.student23', 'student23@demo.university.edu', 'Ella', 'Walker', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd001024-demo-stud-0024-demostude24', 'demo.student24', 'student24@demo.university.edu', 'Jacob', 'Hall', 'student', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW())
ON CONFLICT (username) DO NOTHING;

-- Create student memberships
INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary)
SELECT id, 'dd000000-demo-demo-demo-demouniversity', 'student', true
FROM users 
WHERE id LIKE 'dd001%'
ON CONFLICT (user_id, tenant_id) DO NOTHING;

-- Create Public Health & Healthcare focused Specialties for demo tenant
INSERT INTO specialties (id, name, code, tenant_id, is_active, created_at)
VALUES
  ('dd100001-demo-spec-0001-epidemiology', 'Epidemiology', 'EPD', 'dd000000-demo-demo-demo-demouniversity', true, NOW()),
  ('dd100002-demo-spec-0002-publichealth', 'Public Health', 'PBH', 'dd000000-demo-demo-demo-demouniversity', true, NOW()),
  ('dd100003-demo-spec-0003-healthpolicy', 'Health Policy & Management', 'HPM', 'dd000000-demo-demo-demo-demouniversity', true, NOW()),
  ('dd100004-demo-spec-0004-globalhealth', 'Global Health', 'GLH', 'dd000000-demo-demo-demo-demouniversity', true, NOW()),
  ('dd100005-demo-spec-0005-environhlth', 'Environmental Health Sciences', 'EHS', 'dd000000-demo-demo-demo-demouniversity', true, NOW()),
  ('dd100006-demo-spec-0006-biostatistcs', 'Biostatistics', 'BST', 'dd000000-demo-demo-demo-demouniversity', true, NOW()),
  ('dd100007-demo-spec-0007-healthbehav', 'Health Behavior', 'HBV', 'dd000000-demo-demo-demo-demouniversity', true, NOW()),
  ('dd100008-demo-spec-0008-clinicalres', 'Clinical Research', 'CLR', 'dd000000-demo-demo-demo-demouniversity', true, NOW())
ON CONFLICT (id) DO NOTHING;

-- Create Programs for demo tenant
INSERT INTO programs (id, name, code, tenant_id, is_active, created_at)
VALUES
  ('dd200001-demo-prog-0001-drph', 'Doctor of Public Health (DrPH)', 'DRPH', 'dd000000-demo-demo-demo-demouniversity', true, NOW()),
  ('dd200002-demo-prog-0002-phdph', 'PhD in Public Health', 'PHDPH', 'dd000000-demo-demo-demo-demouniversity', true, NOW()),
  ('dd200003-demo-prog-0003-phdhsm', 'PhD in Health Services Management', 'PHDHSM', 'dd000000-demo-demo-demo-demouniversity', true, NOW()),
  ('dd200004-demo-prog-0004-phdepid', 'PhD in Epidemiology', 'PHDEPD', 'dd000000-demo-demo-demo-demouniversity', true, NOW()),
  ('dd200005-demo-prog-0005-phdbios', 'PhD in Biostatistics', 'PHDBST', 'dd000000-demo-demo-demo-demouniversity', true, NOW())
ON CONFLICT (id) DO NOTHING;

-- Create Cohorts for demo tenant
INSERT INTO cohorts (id, name, year, tenant_id, is_active, created_at)
VALUES
  ('dd300001-demo-coho-0001-2020', 'Cohort 2020', 2020, 'dd000000-demo-demo-demo-demouniversity', true, NOW()),
  ('dd300002-demo-coho-0002-2021', 'Cohort 2021', 2021, 'dd000000-demo-demo-demo-demouniversity', true, NOW()),
  ('dd300003-demo-coho-0003-2022', 'Cohort 2022', 2022, 'dd000000-demo-demo-demo-demouniversity', true, NOW()),
  ('dd300004-demo-coho-0004-2023', 'Cohort 2023', 2023, 'dd000000-demo-demo-demo-demouniversity', true, NOW()),
  ('dd300005-demo-coho-0005-2024', 'Cohort 2024', 2024, 'dd000000-demo-demo-demo-demouniversity', true, NOW())
ON CONFLICT (id) DO NOTHING;

-- Create Departments for demo tenant
INSERT INTO departments (id, name, code, tenant_id, is_active, created_at)
VALUES
  ('dd400001-demo-dept-0001-epidemhlth', 'Department of Epidemiology', 'EPID', 'dd000000-demo-demo-demo-demouniversity', true, NOW()),
  ('dd400002-demo-dept-0002-healthpol', 'Department of Health Policy', 'HPOL', 'dd000000-demo-demo-demo-demouniversity', true, NOW()),
  ('dd400003-demo-dept-0003-biostat', 'Department of Biostatistics', 'BIOS', 'dd000000-demo-demo-demo-demouniversity', true, NOW()),
  ('dd400004-demo-dept-0004-envhealth', 'Department of Environmental Health', 'ENVH', 'dd000000-demo-demo-demo-demouniversity', true, NOW()),
  ('dd400005-demo-dept-0005-globalhlth', 'Department of Global Health', 'GLBH', 'dd000000-demo-demo-demo-demouniversity', true, NOW()),
  ('dd400006-demo-dept-0006-behavorsc', 'Department of Behavioral Sciences', 'BHVS', 'dd000000-demo-demo-demo-demouniversity', true, NOW())
ON CONFLICT (id) DO NOTHING;

-- Link Programs to Specialties (many-to-many relationship)
INSERT INTO specialty_programs (specialty_id, program_id)
VALUES
  ('dd100001-demo-spec-0001-epidemiology', 'dd200004-demo-prog-0004-phdepid'),
  ('dd100002-demo-spec-0002-publichealth', 'dd200001-demo-prog-0001-drph'),
  ('dd100002-demo-spec-0002-publichealth', 'dd200002-demo-prog-0002-phdph'),
  ('dd100003-demo-spec-0003-healthpolicy', 'dd200003-demo-prog-0003-phdhsm'),
  ('dd100004-demo-spec-0004-globalhealth', 'dd200002-demo-prog-0002-phdph'),
  ('dd100005-demo-spec-0005-environhlth', 'dd200002-demo-prog-0002-phdph'),
  ('dd100006-demo-spec-0006-biostatistcs', 'dd200005-demo-prog-0005-phdbios'),
  ('dd100007-demo-spec-0007-healthbehav', 'dd200001-demo-prog-0001-drph'),
  ('dd100008-demo-spec-0008-clinicalres', 'dd200002-demo-prog-0002-phdph')
ON CONFLICT DO NOTHING;
