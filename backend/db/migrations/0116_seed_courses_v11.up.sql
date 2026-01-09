-- Seed Global Courses (matching actual courses table schema)
INSERT INTO courses (id, tenant_id, code, title, description, credits, workload_hours, is_active, created_at, updated_at)
VALUES
  ('11111111-0000-0000-0001-000000000001'::uuid, 'dd000000-0000-0000-0000-d00000000001'::uuid, 'RES-101', '{"en": "Research Methodology", "kk": "Зерттеу әдіснамасы", "ru": "Методология исследований"}'::jsonb, '{"en": "Fundamental research methods including qualitative and quantitative analysis."}'::jsonb, 5, 45, true, NOW(), NOW()),
  ('11111111-0000-0000-0001-000000000002'::uuid, 'dd000000-0000-0000-0000-d00000000001'::uuid, 'WRT-202', '{"en": "Academic Writing", "kk": "Академиялық жазу", "ru": "Академическое письмо"}'::jsonb, '{"en": "Advanced academic writing skills for thesis and publication."}'::jsonb, 3, 30, true, NOW(), NOW()),
  ('11111111-0000-0000-0001-000000000003'::uuid, 'dd000000-0000-0000-0000-d00000000001'::uuid, 'STAT-300', '{"en": "Advanced Statistics", "kk": "Жоғары статистика", "ru": "Продвинутая статистика"}'::jsonb, '{"en": "Statistical analysis methods for doctoral research."}'::jsonb, 5, 60, true, NOW(), NOW()),
  ('11111111-0000-0000-0001-000000000004'::uuid, 'dd000000-0000-0000-0000-d00000000001'::uuid, 'ETH-100', '{"en": "Research Ethics", "kk": "Зерттеу этикасы", "ru": "Этика исследований"}'::jsonb, '{"en": "Ethical considerations in academic and medical research."}'::jsonb, 2, 20, true, NOW(), NOW()),
  ('11111111-0000-0000-0001-000000000005'::uuid, 'dd000000-0000-0000-0000-d00000000001'::uuid, 'AI-500', '{"en": "Intro to AI Systems", "kk": "AI жүйелеріне кіріспе", "ru": "Введение в системы ИИ"}'::jsonb, '{"en": "Introduction to artificial intelligence systems and applications."}'::jsonb, 4, 50, true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
