-- Add audit fields to buildings
ALTER TABLE buildings
ADD COLUMN created_by UUID,
ADD COLUMN updated_by UUID,
ADD COLUMN deleted_at TIMESTAMPTZ;

-- Add audit fields to rooms
ALTER TABLE rooms
ADD COLUMN created_by UUID,
ADD COLUMN updated_by UUID,
ADD COLUMN deleted_at TIMESTAMPTZ;

-- Add unique constraint for room name per building (only for active/non-deleted rooms)
CREATE UNIQUE INDEX idx_rooms_building_name ON rooms (building_id, name) WHERE deleted_at IS NULL;
