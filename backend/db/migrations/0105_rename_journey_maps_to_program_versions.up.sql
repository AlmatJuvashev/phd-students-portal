-- Migration: Rename journey map tables to program version tables
-- Rationale: "Journey Map" is a program template/version built in the Program Builder.
-- We keep the API/UI naming flexible, but store versions under program_versions.

-- 1) Rename tables
ALTER TABLE IF EXISTS journey_maps RENAME TO program_versions;
ALTER TABLE IF EXISTS journey_node_definitions RENAME TO program_version_node_definitions;

-- 2) Rename FK column on nodes
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'program_version_node_definitions' AND column_name = 'journey_map_id'
    ) THEN
        ALTER TABLE program_version_node_definitions RENAME COLUMN journey_map_id TO program_version_id;
    END IF;
END $$;

-- 3) Ensure program_versions has config + updated_at (older DBs may not have 0104 applied)
ALTER TABLE program_versions
    ADD COLUMN IF NOT EXISTS config JSONB DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT NOW();

-- 3.1) Ensure node definitions table has updated_at for edits from the builder
ALTER TABLE program_version_node_definitions
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT NOW();

-- 4) Normalize version and enforce uniqueness per program
UPDATE program_versions SET version = '0.0.0' WHERE version IS NULL OR version = '';
ALTER TABLE program_versions ALTER COLUMN version SET NOT NULL;
ALTER TABLE program_versions ALTER COLUMN is_active SET DEFAULT false;

-- Ensure at most one active version per program before adding the partial unique index.
-- Older data may have multiple active=true rows because journey_maps had no such constraint.
WITH ranked AS (
    SELECT
        id,
        program_id,
        ROW_NUMBER() OVER (PARTITION BY program_id ORDER BY created_at DESC, id DESC) AS rn
    FROM program_versions
    WHERE is_active = true
)
UPDATE program_versions pv
SET is_active = false
FROM ranked r
WHERE pv.id = r.id AND r.rn > 1;

ALTER TABLE program_versions
    ADD CONSTRAINT program_versions_program_id_version_unique UNIQUE (program_id, version);

-- 5) Enforce at most one active version per program
CREATE UNIQUE INDEX IF NOT EXISTS idx_program_versions_one_active
    ON program_versions (program_id)
    WHERE is_active = true;

-- 6) Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_program_versions_program_id ON program_versions(program_id);
CREATE INDEX IF NOT EXISTS idx_program_versions_created_at ON program_versions(created_at);
