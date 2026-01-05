-- Migration: Seed demo.university with scheduling-ready catalog data
-- Adds additional faculty users, programs, and courses (with requirements) to better demo the scheduler.

-- Constants
-- Demo tenant: dd000000-0000-0000-0000-d00000000001 (from migration 0047)
-- Demo admin:  dd000001-0000-0000-0001-000000000001

-- =========================================================
-- 1) Faculty / Instructors (role=advisor; used as INSTRUCTOR)
-- =========================================================
-- Password: demopassword123! (same hash as demo advisors/students in 0047)
INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active, created_at)
VALUES
  ('dd000007-0000-0000-0006-000000000006', 'prof.aliyev', 'aliyev@demo.university.edu', 'Nurlan', 'Aliyev', 'advisor', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd000008-0000-0000-0007-000000000007', 'prof.kim', 'kim@demo.university.edu', 'Min', 'Kim', 'advisor', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd000009-0000-0000-0008-000000000008', 'prof.patel', 'patel@demo.university.edu', 'Anaya', 'Patel', 'advisor', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd000010-0000-0000-0009-000000000009', 'prof.sato', 'sato@demo.university.edu', 'Haruto', 'Sato', 'advisor', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd000011-0000-0000-0010-000000000010', 'prof.smits', 'smits@demo.university.edu', 'Eva', 'Smits', 'advisor', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW())
ON CONFLICT DO NOTHING;

INSERT INTO user_tenant_memberships (user_id, tenant_id, role, roles, is_primary)
VALUES
  ('dd000007-0000-0000-0006-000000000006', 'dd000000-0000-0000-0000-d00000000001', 'advisor', ARRAY['advisor'], true),
  ('dd000008-0000-0000-0007-000000000007', 'dd000000-0000-0000-0000-d00000000001', 'advisor', ARRAY['advisor'], true),
  ('dd000009-0000-0000-0008-000000000008', 'dd000000-0000-0000-0000-d00000000001', 'advisor', ARRAY['advisor'], true),
  ('dd000010-0000-0000-0009-000000000009', 'dd000000-0000-0000-0000-d00000000001', 'advisor', ARRAY['advisor'], true),
  ('dd000011-0000-0000-0010-000000000010', 'dd000000-0000-0000-0000-d00000000001', 'advisor', ARRAY['advisor'], true)
ON CONFLICT (user_id, tenant_id) DO NOTHING;

-- =========================
-- 2) Additional Demo Programs
-- =========================
INSERT INTO programs (id, tenant_id, code, name, title, description, credits, duration_months, is_active, created_at, updated_at)
VALUES
  (
    'dd200006-0000-0000-0006-000000000006',
    'dd000000-0000-0000-0000-d00000000001',
    'PHDHINF',
    'PhD in Health Informatics',
    '{"en":"PhD in Health Informatics","ru":"PhD по медицинской информатике"}'::jsonb,
    '{"en":"Data-driven healthcare, clinical data systems, and digital health.","ru":"Данные в здравоохранении, клинические ИС и цифровая медицина."}'::jsonb,
    180,
    48,
    true,
    NOW(),
    NOW()
  ),
  (
    'dd200007-0000-0000-0007-000000000007',
    'dd000000-0000-0000-0000-d00000000001',
    'PHDCEPI',
    'PhD in Clinical Epidemiology',
    '{"en":"PhD in Clinical Epidemiology","ru":"PhD по клинической эпидемиологии"}'::jsonb,
    '{"en":"Advanced methods for evidence-based medicine and clinical research.","ru":"Методы доказательной медицины и клинических исследований."}'::jsonb,
    180,
    48,
    true,
    NOW(),
    NOW()
  ),
  (
    'dd200008-0000-0000-0008-000000000008',
    'dd000000-0000-0000-0000-d00000000001',
    'PHDENVP',
    'PhD in Environmental Policy',
    '{"en":"PhD in Environmental Policy","ru":"PhD по экологической политике"}'::jsonb,
    '{"en":"Policy, governance, and environmental health regulation.","ru":"Политика, управление и регулирование в области экологического здоровья."}'::jsonb,
    180,
    48,
    true,
    NOW(),
    NOW()
  )
ON CONFLICT DO NOTHING;

-- =========================
-- 3) Additional Demo Courses
-- =========================
-- Department IDs seeded in migration 0047:
--   Epidemiology:          dd400001-0000-0000-0001-000000000001
--   Health Policy:         dd400002-0000-0000-0002-000000000002
--   Biostatistics:         dd400003-0000-0000-0003-000000000003
--   Environmental Health:  dd400004-0000-0000-0004-000000000004
--   Global Health:         dd400005-0000-0000-0005-000000000005
INSERT INTO courses (id, tenant_id, program_id, department_id, code, title, description, credits, workload_hours, is_active, created_at, updated_at)
VALUES
  (
    'dd210001-0000-0000-0001-000000000001',
    'dd000000-0000-0000-0000-d00000000001',
    'dd200006-0000-0000-0006-000000000006',
    'dd400003-0000-0000-0003-000000000003',
    'HINF-510',
    '{"en":"Applied Statistical Computing for Health Data","ru":"Прикладные вычисления и статистика для данных здравоохранения"}'::jsonb,
    '{"en":"Hands-on computing for biostatistics and health informatics (R/Python).","ru":"Практикум по вычислениям для биостатистики и информатики (R/Python)."}'::jsonb,
    5,
    150,
    true,
    NOW(),
    NOW()
  ),
  (
    'dd210002-0000-0000-0002-000000000002',
    'dd000000-0000-0000-0000-d00000000001',
    'dd200007-0000-0000-0007-000000000007',
    'dd400001-0000-0000-0001-000000000001',
    'EPI-520',
    '{"en":"Outbreak Investigation Lab","ru":"Лаборатория расследования вспышек"}'::jsonb,
    '{"en":"Case-based epidemiology with lab-style group work and presentations.","ru":"Кейс-эпидемиология с групповыми заданиями и презентациями."}'::jsonb,
    4,
    120,
    true,
    NOW(),
    NOW()
  ),
  (
    'dd210003-0000-0000-0003-000000000003',
    'dd000000-0000-0000-0000-d00000000001',
    'dd200008-0000-0000-0008-000000000008',
    'dd400002-0000-0000-0002-000000000002',
    'HPM-510',
    '{"en":"Healthcare Systems & Policy Analysis","ru":"Системы здравоохранения и анализ политики"}'::jsonb,
    '{"en":"Policy frameworks, financing, and comparative health systems.","ru":"Политика, финансирование и сравнительный анализ систем здравоохранения."}'::jsonb,
    4,
    120,
    true,
    NOW(),
    NOW()
  ),
  (
    'dd210004-0000-0000-0004-000000000004',
    'dd000000-0000-0000-0000-d00000000001',
    NULL,
    'dd400005-0000-0000-0005-000000000005',
    'GLH-530',
    '{"en":"Global Health Field Methods","ru":"Полевые методы в глобальном здравоохранении"}'::jsonb,
    '{"en":"Mixed-methods data collection, ethics, and field operations.","ru":"Смешанные методы сбора данных, этика и полевые исследования."}'::jsonb,
    3,
    90,
    true,
    NOW(),
    NOW()
  ),
  (
    'dd210005-0000-0000-0005-000000000005',
    'dd000000-0000-0000-0000-d00000000001',
    NULL,
    'dd400004-0000-0000-0004-000000000004',
    'EHS-540',
    '{"en":"Environmental Exposure Assessment","ru":"Оценка экологических воздействий"}'::jsonb,
    '{"en":"Sampling, monitoring, and exposure modeling with lab demonstrations.","ru":"Отбор проб, мониторинг и моделирование воздействий с лабораторными демонстрациями."}'::jsonb,
    4,
    120,
    true,
    NOW(),
    NOW()
  ),
  (
    'dd210006-0000-0000-0006-000000000006',
    'dd000000-0000-0000-0000-d00000000001',
    NULL,
    NULL,
    'UNI-500',
    '{"en":"Academic Writing & Research Integrity","ru":"Академическое письмо и исследовательская добросовестность"}'::jsonb,
    '{"en":"Writing, citation practice, reproducibility, and ethics.","ru":"Письмо, цитирование, воспроизводимость и этика."}'::jsonb,
    3,
    90,
    true,
    NOW(),
    NOW()
  ),
  (
    'dd210007-0000-0000-0007-000000000007',
    'dd000000-0000-0000-0000-d00000000001',
    NULL,
    NULL,
    'UNI-510',
    '{"en":"Teaching Practicum for PhD Students","ru":"Педагогическая практика для докторантов"}'::jsonb,
    '{"en":"Micro-teaching, course design, and assessment basics.","ru":"Микро-уроки, дизайн курса и основы оценивания."}'::jsonb,
    2,
    60,
    true,
    NOW(),
    NOW()
  )
ON CONFLICT DO NOTHING;

-- ===================================
-- 4) Course Requirements (Room Matching)
-- ===================================
INSERT INTO course_requirements (course_id, key, value)
VALUES
  ('dd210001-0000-0000-0001-000000000001', 'EQUIPMENT', 'COMPUTERS'),
  ('dd210001-0000-0000-0001-000000000001', 'SOFTWARE', 'RSTUDIO'),
  ('dd210002-0000-0000-0002-000000000002', 'EQUIPMENT', 'PROJECTOR'),
  ('dd210005-0000-0000-0005-000000000005', 'EQUIPMENT', 'FUME_HOOD')
ON CONFLICT DO NOTHING;

-- ==========================================
-- 5) Buildings, Rooms, and Room Attributes
-- ==========================================
-- These resources are used to demo scheduling constraints:
-- - specialized labs (COMPUTERS/RSTUDIO, FUME_HOOD)
-- - universal rooms (department_id NULL) with various capacities

INSERT INTO buildings (id, tenant_id, name, address, description, is_active, created_by, updated_by, created_at, updated_at)
VALUES
  (
    'dd600001-0000-0000-0001-000000000001',
    'dd000000-0000-0000-0000-d00000000001',
    'Main Teaching Center',
    '1 University Ave',
    '{"en":"Shared teaching spaces for lectures and seminars.","ru":"Общие аудитории для лекций и семинаров."}'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  ),
  (
    'dd600002-0000-0000-0002-000000000002',
    'dd000000-0000-0000-0000-d00000000001',
    'Biostatistics & Informatics Lab',
    '2 University Ave',
    '{"en":"Computer labs for data science and biostatistics.","ru":"Компьютерные классы для анализа данных и биостатистики."}'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  ),
  (
    'dd600003-0000-0000-0003-000000000003',
    'dd000000-0000-0000-0000-d00000000001',
    'Epidemiology & Policy Building',
    '3 University Ave',
    '{"en":"Departmental teaching rooms for epidemiology and policy.","ru":"Аудитории кафедр эпидемиологии и политики здравоохранения."}'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  ),
  (
    'dd600004-0000-0000-0004-000000000004',
    'dd000000-0000-0000-0000-d00000000001',
    'Environmental Health Labs',
    '4 University Ave',
    '{"en":"Wet labs and exposure assessment facilities.","ru":"Лаборатории для экологических исследований и оценки воздействий."}'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  )
ON CONFLICT DO NOTHING;

INSERT INTO rooms (id, building_id, name, capacity, floor, department_id, type, features, is_active, created_by, updated_by, created_at, updated_at)
VALUES
  -- Main Teaching Center (universal rooms, department_id NULL)
  (
    'dd610001-0000-0000-0001-000000000001',
    'dd600001-0000-0000-0001-000000000001',
    'Auditorium A',
    220,
    1,
    NULL,
    'lecture_hall',
    '["Projector","Audio"]'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  ),
  (
    'dd610002-0000-0000-0002-000000000002',
    'dd600001-0000-0000-0001-000000000001',
    'Room 101',
    60,
    1,
    NULL,
    'classroom',
    '["Projector","Whiteboard"]'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  ),
  (
    'dd610003-0000-0000-0003-000000000003',
    'dd600001-0000-0000-0001-000000000001',
    'Room 102',
    40,
    1,
    NULL,
    'classroom',
    '["Projector","Whiteboard"]'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  ),
  (
    'dd610004-0000-0000-0004-000000000004',
    'dd600001-0000-0000-0001-000000000001',
    'Seminar 201',
    28,
    2,
    NULL,
    'seminar_room',
    '["Whiteboard"]'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  ),
  (
    'dd610005-0000-0000-0005-000000000005',
    'dd600001-0000-0000-0001-000000000001',
    'Seminar 202',
    28,
    2,
    NULL,
    'seminar_room',
    '["Whiteboard"]'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  ),

  -- Biostatistics & Informatics Lab (specialized)
  (
    'dd620001-0000-0000-0001-000000000001',
    'dd600002-0000-0000-0002-000000000002',
    'Computer Lab BL-1',
    40,
    1,
    'dd400003-0000-0000-0003-000000000003',
    'lab',
    '["Computers","RStudio"]'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  ),
  (
    'dd620002-0000-0000-0002-000000000002',
    'dd600002-0000-0000-0002-000000000002',
    'Analytics Classroom BL-2',
    30,
    1,
    'dd400003-0000-0000-0003-000000000003',
    'classroom',
    '["Projector","Whiteboard"]'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  ),

  -- Epidemiology & Policy Building (departmental)
  (
    'dd630001-0000-0000-0001-000000000001',
    'dd600003-0000-0000-0003-000000000003',
    'EPI Lecture Hall 1',
    80,
    1,
    'dd400001-0000-0000-0001-000000000001',
    'lecture_hall',
    '["Projector","Audio"]'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  ),
  (
    'dd630002-0000-0000-0002-000000000002',
    'dd600003-0000-0000-0003-000000000003',
    'Policy Seminar Room',
    35,
    2,
    'dd400002-0000-0000-0002-000000000002',
    'seminar_room',
    '["Projector","Whiteboard"]'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  ),
  (
    'dd630003-0000-0000-0003-000000000003',
    'dd600003-0000-0000-0003-000000000003',
    'Global Health Meeting Room',
    25,
    3,
    'dd400005-0000-0000-0005-000000000005',
    'seminar_room',
    '["Projector","VideoConference"]'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  ),

  -- Environmental Health Labs (specialized)
  (
    'dd640001-0000-0000-0001-000000000001',
    'dd600004-0000-0000-0004-000000000004',
    'Wet Lab EH-1',
    20,
    1,
    'dd400004-0000-0000-0004-000000000004',
    'lab',
    '["FumeHood","LabBench"]'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  ),
  (
    'dd640002-0000-0000-0002-000000000002',
    'dd600004-0000-0000-0004-000000000004',
    'Sample Prep Lab EH-2',
    16,
    1,
    'dd400004-0000-0000-0004-000000000004',
    'lab',
    '["FumeHood"]'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  )
ON CONFLICT DO NOTHING;

INSERT INTO room_attributes (room_id, key, value)
VALUES
  -- Universal teaching rooms
  ('dd610001-0000-0000-0001-000000000001', 'EQUIPMENT', 'PROJECTOR'),
  ('dd610002-0000-0000-0002-000000000002', 'EQUIPMENT', 'PROJECTOR'),
  ('dd610003-0000-0000-0003-000000000003', 'EQUIPMENT', 'PROJECTOR'),
  ('dd610004-0000-0000-0004-000000000004', 'EQUIPMENT', 'WHITEBOARD'),
  ('dd610005-0000-0000-0005-000000000005', 'EQUIPMENT', 'WHITEBOARD'),

  -- Biostatistics & Informatics
  ('dd620001-0000-0000-0001-000000000001', 'EQUIPMENT', 'COMPUTERS'),
  ('dd620001-0000-0000-0001-000000000001', 'SOFTWARE', 'RSTUDIO'),
  ('dd620002-0000-0000-0002-000000000002', 'EQUIPMENT', 'PROJECTOR'),

  -- Epidemiology & Policy
  ('dd630001-0000-0000-0001-000000000001', 'EQUIPMENT', 'PROJECTOR'),
  ('dd630002-0000-0000-0002-000000000002', 'EQUIPMENT', 'PROJECTOR'),
  ('dd630003-0000-0000-0003-000000000003', 'EQUIPMENT', 'VIDEO_CONF'),

  -- Environmental Health
  ('dd640001-0000-0000-0001-000000000001', 'EQUIPMENT', 'FUME_HOOD'),
  ('dd640002-0000-0000-0002-000000000002', 'EQUIPMENT', 'FUME_HOOD')
ON CONFLICT DO NOTHING;
