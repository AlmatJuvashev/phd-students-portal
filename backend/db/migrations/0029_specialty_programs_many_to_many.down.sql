-- Re-add program_id column to specialties
ALTER TABLE specialties ADD COLUMN program_id uuid REFERENCES programs(id) ON DELETE SET NULL;

-- Migrate data back (keep only the first program if multiple)
UPDATE specialties s
SET program_id = sp.program_id
FROM (
  SELECT DISTINCT ON (specialty_id) specialty_id, program_id
  FROM specialty_programs
  ORDER BY specialty_id, created_at
) sp
WHERE s.id = sp.specialty_id;

-- Drop junction table
DROP TABLE IF EXISTS specialty_programs;

-- Recreate the index
CREATE INDEX idx_specialties_program_id ON specialties(program_id);
