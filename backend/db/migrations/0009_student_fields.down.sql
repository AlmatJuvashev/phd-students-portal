-- Revert student-advisor mapping and extra fields
DROP TABLE IF EXISTS student_advisors;

-- Cannot set NOT NULL back on email safely without data migration; skip in down migration.
-- ALTER TABLE users ALTER COLUMN email SET NOT NULL;

ALTER TABLE users
  DROP COLUMN IF EXISTS cohort,
  DROP COLUMN IF EXISTS department,
  DROP COLUMN IF EXISTS program,
  DROP COLUMN IF EXISTS phone;

