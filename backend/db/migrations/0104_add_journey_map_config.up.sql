-- Add config column to journey_maps for storing layout/phases
ALTER TABLE journey_maps ADD COLUMN IF NOT EXISTS config JSONB DEFAULT '{}'::jsonb;
