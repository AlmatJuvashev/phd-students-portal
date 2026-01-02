DROP INDEX IF EXISTS idx_rooms_building_name;

ALTER TABLE rooms
DROP COLUMN IF EXISTS created_by,
DROP COLUMN IF EXISTS updated_by,
DROP COLUMN IF EXISTS deleted_at;

ALTER TABLE buildings
DROP COLUMN IF EXISTS created_by,
DROP COLUMN IF EXISTS updated_by,
DROP COLUMN IF EXISTS deleted_at;
