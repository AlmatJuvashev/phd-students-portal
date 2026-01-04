-- Down Migration: Revert program version tables back to journey map tables

-- Drop indexes/constraints added in up migration
DROP INDEX IF EXISTS idx_program_versions_one_active;
DROP INDEX IF EXISTS idx_program_versions_program_id;
DROP INDEX IF EXISTS idx_program_versions_created_at;

ALTER TABLE IF EXISTS program_versions
    DROP CONSTRAINT IF EXISTS program_versions_program_id_version_unique;

-- Note: We intentionally keep data; we only rename back.
-- Revert column/table names
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'program_version_node_definitions' AND column_name = 'program_version_id'
    ) THEN
        ALTER TABLE program_version_node_definitions RENAME COLUMN program_version_id TO journey_map_id;
    END IF;
END $$;

ALTER TABLE IF EXISTS program_version_node_definitions RENAME TO journey_node_definitions;
ALTER TABLE IF EXISTS program_versions RENAME TO journey_maps;

