ALTER TABLE users
  ADD COLUMN IF NOT EXISTS phone text,
  ADD COLUMN IF NOT EXISTS program text,
  ADD COLUMN IF NOT EXISTS department text,
  ADD COLUMN IF NOT EXISTS cohort text;

-- Make email optional to allow student accounts without email
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;

-- Link students to advisors (many-to-many)
CREATE TABLE IF NOT EXISTS student_advisors (
  student_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  advisor_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  PRIMARY KEY (student_id, advisor_id)
);

