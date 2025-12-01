-- Create test student user
-- This script creates a test student account for E2E testing

INSERT INTO users (username, password_hash, role, email, full_name, created_at, updated_at)
VALUES (
  'tu6260',
  -- Password: thunder-pluto-river72
  -- You'll need to hash this using your backend's hashing function
  'REPLACE_WITH_HASHED_PASSWORD',
  'student',
  'tu6260@test.kaznmu.edu.kz',
  'Test Student',
  NOW(),
  NOW()
)
ON CONFLICT (username) DO NOTHING;

-- Note: You need to generate the password hash using your backend's password hashing function
-- For example, if using bcrypt in Go:
-- bcrypt.GenerateFromPassword([]byte("thunder-pluto-river72"), bcrypt.DefaultCost)
