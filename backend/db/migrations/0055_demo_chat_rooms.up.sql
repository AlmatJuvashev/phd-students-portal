-- Migration: Add chat rooms for demo.university tenant and add members
-- This migration creates chat rooms and adds demo.admin and demo students as members

-- Define the demo tenant ID
-- dd000000-0000-0000-0000-d00000000001 = demo.university tenant
-- dd000001-0000-0000-0001-000000000001 = demo.admin user
-- dd001001-0000-0000-0001-000000000001 = demo.student1 user
-- dd001002-0000-0000-0002-000000000002 = demo.student2 user

-- Create chat rooms for demo tenant
INSERT INTO chat_rooms (id, name, type, created_by, created_by_role, tenant_id, is_archived, meta, created_at)
VALUES
  ('dd500001-0000-0000-0001-000000000001', 'Public Health Cohort 2024', 'cohort', 'dd000001-0000-0000-0001-000000000001', 'admin', 'dd000000-0000-0000-0000-d00000000001', false, '{}', NOW()),
  ('dd500002-0000-0000-0002-000000000002', 'Advisory: Dr. Johnson', 'advisory', 'dd000001-0000-0000-0001-000000000001', 'admin', 'dd000000-0000-0000-0000-d00000000001', false, '{}', NOW()),
  ('dd500003-0000-0000-0003-000000000003', 'Epidemiology Specialists', 'other', 'dd000001-0000-0000-0001-000000000001', 'admin', 'dd000000-0000-0000-0000-d00000000001', false, '{}', NOW())
ON CONFLICT (id) DO NOTHING;

-- Add demo.admin to all demo rooms
INSERT INTO chat_room_members (room_id, user_id, role_in_room, tenant_id, joined_at)
VALUES
  ('dd500001-0000-0000-0001-000000000001', 'dd000001-0000-0000-0001-000000000001', 'admin', 'dd000000-0000-0000-0000-d00000000001', NOW()),
  ('dd500002-0000-0000-0002-000000000002', 'dd000001-0000-0000-0001-000000000001', 'admin', 'dd000000-0000-0000-0000-d00000000001', NOW()),
  ('dd500003-0000-0000-0003-000000000003', 'dd000001-0000-0000-0001-000000000001', 'admin', 'dd000000-0000-0000-0000-d00000000001', NOW())
ON CONFLICT (room_id, user_id) DO NOTHING;

-- Add demo.student1 to first two rooms
INSERT INTO chat_room_members (room_id, user_id, role_in_room, tenant_id, joined_at)
VALUES
  ('dd500001-0000-0000-0001-000000000001', 'dd001001-0000-0000-0001-000000000001', 'member', 'dd000000-0000-0000-0000-d00000000001', NOW()),
  ('dd500002-0000-0000-0002-000000000002', 'dd001001-0000-0000-0001-000000000001', 'member', 'dd000000-0000-0000-0000-d00000000001', NOW())
ON CONFLICT (room_id, user_id) DO NOTHING;

-- Add demo.student2 to first and third rooms
INSERT INTO chat_room_members (room_id, user_id, role_in_room, tenant_id, joined_at)
VALUES
  ('dd500001-0000-0000-0001-000000000001', 'dd001002-0000-0000-0002-000000000002', 'member', 'dd000000-0000-0000-0000-d00000000001', NOW()),
  ('dd500003-0000-0000-0003-000000000003', 'dd001002-0000-0000-0002-000000000002', 'member', 'dd000000-0000-0000-0000-d00000000001', NOW())
ON CONFLICT (room_id, user_id) DO NOTHING;

-- Add advisors to advisory room
INSERT INTO chat_room_members (room_id, user_id, role_in_room, tenant_id, joined_at)
SELECT 
  'dd500002-0000-0000-0002-000000000002',
  id,
  'member',
  'dd000000-0000-0000-0000-d00000000001',
  NOW()
FROM users
WHERE id IN (
  'dd000002-0000-0000-0001-000000000001',  -- dr.johnson
  'dd000003-0000-0000-0002-000000000002'   -- dr.williams
)
ON CONFLICT (room_id, user_id) DO NOTHING;

-- Add some welcome messages to the cohort room
INSERT INTO chat_messages (id, room_id, sender_id, body, tenant_id, created_at)
VALUES
  ('dd600001-0000-0000-0001-000000000001', 'dd500001-0000-0000-0001-000000000001', 'dd000001-0000-0000-0001-000000000001', 'Welcome to the Public Health Cohort 2024 chat room! 🎓', 'dd000000-0000-0000-0000-d00000000001', NOW() - INTERVAL '1 day'),
  ('dd600002-0000-0000-0002-000000000002', 'dd500001-0000-0000-0001-000000000001', 'dd000001-0000-0000-0001-000000000001', 'Please use this room for program announcements and cohort-wide discussions.', 'dd000000-0000-0000-0000-d00000000001', NOW() - INTERVAL '1 day' + INTERVAL '1 minute')
ON CONFLICT (id) DO NOTHING;
