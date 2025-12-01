-- First, fill any null emails with a placeholder to avoid constraint violation
UPDATE users 
SET email = 'missing-' || id || '@placeholder.local' 
WHERE email IS NULL;

-- Now make the column required
ALTER TABLE users ALTER COLUMN email SET NOT NULL;
