-- Create junction table for many-to-many relationship between specialties and programs
CREATE TABLE specialty_programs (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  specialty_id uuid NOT NULL REFERENCES specialties(id) ON DELETE CASCADE,
  program_id uuid NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE(specialty_id, program_id)
);

CREATE INDEX idx_specialty_programs_specialty ON specialty_programs(specialty_id);
CREATE INDEX idx_specialty_programs_program ON specialty_programs(program_id);

-- Migrate existing data from specialties.program_id to specialty_programs
INSERT INTO specialty_programs (specialty_id, program_id)
SELECT id, program_id
FROM specialties
WHERE program_id IS NOT NULL;

-- Remove program_id column from specialties (it's now in the junction table)
ALTER TABLE specialties DROP COLUMN program_id;
