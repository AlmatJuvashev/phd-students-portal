-- Migration: Extend demo.university scheduling catalog data
-- Adds more programs, courses, instructors, and specialized/universal rooms
-- to stress-test the scheduling algorithm (department + attribute matching).
--
-- Demo tenant: dd000000-0000-0000-0000-d00000000001
-- Demo admin:  dd000001-0000-0000-0001-000000000001

-- =========================================================
-- 1) Additional Faculty / Instructors (role=advisor)
-- =========================================================
-- Password: demopassword123! (same hash as demo advisors/students in 0047)
INSERT INTO users (id, username, email, first_name, last_name, role, password_hash, is_active, created_at)
VALUES
  ('dd000012-0000-0000-0011-000000000011', 'prof.chen', 'chen@demo.university.edu', 'Lina', 'Chen', 'advisor', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd000013-0000-0000-0012-000000000012', 'prof.nakamura', 'nakamura@demo.university.edu', 'Kenji', 'Nakamura', 'advisor', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd000014-0000-0000-0013-000000000013', 'prof.singh', 'singh@demo.university.edu', 'Riya', 'Singh', 'advisor', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd000015-0000-0000-0014-000000000014', 'prof.garcia', 'garcia@demo.university.edu', 'Mateo', 'Garcia', 'advisor', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW()),
  ('dd000016-0000-0000-0015-000000000015', 'prof.kuznetsova', 'kuznetsova@demo.university.edu', 'Elena', 'Kuznetsova', 'advisor', '$2a$10$Wz5yHrQmXhKLJxVKxUQXOeJ7.N6GYj9bIxJZ3vMVq7qh0Kz5D1DmC', true, NOW())
ON CONFLICT DO NOTHING;

INSERT INTO user_tenant_memberships (user_id, tenant_id, role, roles, is_primary)
VALUES
  ('dd000012-0000-0000-0011-000000000011', 'dd000000-0000-0000-0000-d00000000001', 'advisor', ARRAY['advisor'], true),
  ('dd000013-0000-0000-0012-000000000012', 'dd000000-0000-0000-0000-d00000000001', 'advisor', ARRAY['advisor'], true),
  ('dd000014-0000-0000-0013-000000000013', 'dd000000-0000-0000-0000-d00000000001', 'advisor', ARRAY['advisor'], true),
  ('dd000015-0000-0000-0014-000000000014', 'dd000000-0000-0000-0000-d00000000001', 'advisor', ARRAY['advisor'], true),
  ('dd000016-0000-0000-0015-000000000015', 'dd000000-0000-0000-0000-d00000000001', 'advisor', ARRAY['advisor'], true)
ON CONFLICT (user_id, tenant_id) DO NOTHING;

-- =========================
-- 2) Additional Demo Programs
-- =========================
INSERT INTO programs (id, tenant_id, code, name, title, description, credits, duration_months, is_active, created_at, updated_at)
VALUES
  (
    'dd200009-0000-0000-0009-000000000009',
    'dd000000-0000-0000-0000-d00000000001',
    'PHDBIOM',
    'PhD in Biomedical Sciences',
    '{"en":"PhD in Biomedical Sciences","ru":"PhD по биомедицинским наукам"}'::jsonb,
    '{"en":"Molecular methods, translational research, and lab-based investigation.","ru":"Молекулярные методы, трансляционные исследования и лабораторная работа."}'::jsonb,
    180,
    48,
    true,
    NOW(),
    NOW()
  ),
  (
    'dd200010-0000-0000-0010-000000000010',
    'dd000000-0000-0000-0000-d00000000001',
    'PHDCLSIM',
    'PhD in Clinical Simulation Science',
    '{"en":"PhD in Clinical Simulation Science","ru":"PhD по клиническим симуляциям"}'::jsonb,
    '{"en":"Simulation-based education, OSCE design, and clinical skills training research.","ru":"Симуляционное обучение, дизайн OSCE и исследования клинических навыков."}'::jsonb,
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
INSERT INTO courses (id, tenant_id, program_id, department_id, code, title, description, credits, workload_hours, is_active, created_at, updated_at)
VALUES
  (
    'dd210008-0000-0000-0008-000000000008',
    'dd000000-0000-0000-0000-d00000000001',
    'dd200009-0000-0000-0009-000000000009',
    'dd400004-0000-0000-0004-000000000004',
    'BIO-610',
    '{"en":"Molecular Methods for Environmental Health","ru":"Молекулярные методы в экологическом здоровье"}'::jsonb,
    '{"en":"PCR, sequencing basics, and lab workflows for exposure science.","ru":"ПЦР, основы секвенирования и лабораторные процессы для наук о воздействиях."}'::jsonb,
    4,
    120,
    true,
    NOW(),
    NOW()
  ),
  (
    'dd210009-0000-0000-0009-000000000009',
    'dd000000-0000-0000-0000-d00000000001',
    'dd200009-0000-0000-0009-000000000009',
    'dd400004-0000-0000-0004-000000000004',
    'BIO-620',
    '{"en":"Biosafety & Laboratory Practice","ru":"Биобезопасность и лабораторная практика"}'::jsonb,
    '{"en":"Safe lab operations, containment, and documentation standards.","ru":"Безопасная работа в лаборатории, контейнмент и стандарты документации."}'::jsonb,
    3,
    90,
    true,
    NOW(),
    NOW()
  ),
  (
    'dd210010-0000-0000-0010-000000000010',
    'dd000000-0000-0000-0000-d00000000001',
    'dd200010-0000-0000-0010-000000000010',
    NULL,
    'SIM-600',
    '{"en":"Clinical Simulation for Researchers","ru":"Клиническая симуляция для исследователей"}'::jsonb,
    '{"en":"Simulation scenarios, debriefing, and evaluation for academic research.","ru":"Сценарии симуляции, дебрифинг и оценивание для исследований."}'::jsonb,
    3,
    90,
    true,
    NOW(),
    NOW()
  ),
  (
    'dd210011-0000-0000-0011-000000000011',
    'dd000000-0000-0000-0000-d00000000001',
    'dd200010-0000-0000-0010-000000000010',
    NULL,
    'OSCE-610',
    '{"en":"OSCE Station Design & Assessment","ru":"Дизайн станций OSCE и оценивание"}'::jsonb,
    '{"en":"Blueprinting, standard setting, and reliable OSCE assessment design.","ru":"Планирование, стандартизация и надежный дизайн оценивания OSCE."}'::jsonb,
    2,
    60,
    true,
    NOW(),
    NOW()
  ),
  (
    'dd210012-0000-0000-0012-000000000012',
    'dd000000-0000-0000-0000-d00000000001',
    NULL,
    'dd400005-0000-0000-0005-000000000005',
    'GLH-540',
    '{"en":"Global Health Virtual Collaboration","ru":"Виртуальное сотрудничество в глобальном здравоохранении"}'::jsonb,
    '{"en":"Remote collaboration patterns, ethics, and cross-border teamwork.","ru":"Дистанционное взаимодействие, этика и межстрановые команды."}'::jsonb,
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
  ('dd210008-0000-0000-0008-000000000008', 'EQUIPMENT', 'MICROSCOPE'),
  ('dd210009-0000-0000-0009-000000000009', 'EQUIPMENT', 'BIOSAFETY_CABINET'),
  ('dd210010-0000-0000-0010-000000000010', 'EQUIPMENT', 'SIMULATION'),
  ('dd210011-0000-0000-0011-000000000011', 'EQUIPMENT', 'SIMULATION'),
  ('dd210012-0000-0000-0012-000000000012', 'EQUIPMENT', 'VIDEO_CONF')
ON CONFLICT DO NOTHING;

-- ==========================================
-- 5) Additional Buildings and Rooms
-- ==========================================
INSERT INTO buildings (id, tenant_id, name, address, description, is_active, created_by, updated_by, created_at, updated_at)
VALUES
  (
    'dd600005-0000-0000-0005-000000000005',
    'dd000000-0000-0000-0000-d00000000001',
    'Clinical Simulation Center',
    '5 University Ave',
    '{"en":"Simulation labs and OSCE rooms for skills training.","ru":"Симуляционные лаборатории и аудитории OSCE для тренинга навыков."}'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  ),
  (
    'dd600006-0000-0000-0006-000000000006',
    'dd000000-0000-0000-0000-d00000000001',
    'Biomedical Research Wing',
    '6 University Ave',
    '{"en":"Wet labs and research seminar rooms for biomedical studies.","ru":"Лаборатории и семинарские аудитории для биомедицинских исследований."}'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  )
ON CONFLICT DO NOTHING;

INSERT INTO rooms (id, building_id, name, capacity, floor, department_id, type, features, is_active, created_by, updated_by, created_at, updated_at)
VALUES
  -- Clinical Simulation Center (universal, department_id NULL)
  (
    'dd650001-0000-0000-0001-000000000001',
    'dd600005-0000-0000-0005-000000000005',
    'Sim Lab CS-1',
    24,
    1,
    NULL,
    'lab',
    '["Simulation","Projector"]'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  ),
  (
    'dd650002-0000-0000-0002-000000000002',
    'dd600005-0000-0000-0005-000000000005',
    'OSCE Room CS-2',
    16,
    1,
    NULL,
    'lab',
    '["Simulation"]'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  ),
  (
    'dd650003-0000-0000-0003-000000000003',
    'dd600005-0000-0000-0005-000000000005',
    'Skills Workshop CS-3',
    30,
    2,
    NULL,
    'classroom',
    '["Projector","Whiteboard"]'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  ),

  -- Biomedical Research Wing (Environmental Health department)
  (
    'dd660001-0000-0000-0001-000000000001',
    'dd600006-0000-0000-0006-000000000006',
    'Molecular Lab BR-1',
    18,
    1,
    'dd400004-0000-0000-0004-000000000004',
    'lab',
    '["Microscopes","LabBench"]'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  ),
  (
    'dd660002-0000-0000-0002-000000000002',
    'dd600006-0000-0000-0006-000000000006',
    'Biosafety Lab BR-2',
    12,
    1,
    'dd400004-0000-0000-0004-000000000004',
    'lab',
    '["BiosafetyCabinet"]'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  ),
  (
    'dd660003-0000-0000-0003-000000000003',
    'dd600006-0000-0000-0006-000000000006',
    'Research Seminar BR-3',
    35,
    2,
    'dd400004-0000-0000-0004-000000000004',
    'seminar_room',
    '["Projector","Whiteboard"]'::jsonb,
    true,
    'dd000001-0000-0000-0001-000000000001',
    'dd000001-0000-0000-0001-000000000001',
    NOW(),
    NOW()
  )
ON CONFLICT DO NOTHING;

INSERT INTO room_attributes (room_id, key, value)
VALUES
  -- Simulation Center
  ('dd650001-0000-0000-0001-000000000001', 'EQUIPMENT', 'SIMULATION'),
  ('dd650002-0000-0000-0002-000000000002', 'EQUIPMENT', 'SIMULATION'),

  -- Biomedical Research Wing
  ('dd660001-0000-0000-0001-000000000001', 'EQUIPMENT', 'MICROSCOPE'),
  ('dd660002-0000-0000-0002-000000000002', 'EQUIPMENT', 'BIOSAFETY_CABINET')
ON CONFLICT DO NOTHING;

