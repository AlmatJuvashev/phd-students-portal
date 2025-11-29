-- Add creator role to chat rooms (safe to apply after 0022 ran without this column)

ALTER TABLE chat_rooms
ADD COLUMN IF NOT EXISTS created_by_role user_role;

-- Backfill from users table where missing
UPDATE chat_rooms r
SET created_by_role = u.role
FROM users u
WHERE u.id = r.created_by
  AND r.created_by_role IS NULL;

ALTER TABLE chat_rooms
ALTER COLUMN created_by_role SET NOT NULL;
