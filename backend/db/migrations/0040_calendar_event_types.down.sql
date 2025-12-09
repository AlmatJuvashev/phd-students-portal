-- Rollback meeting type and related fields
ALTER TABLE events DROP COLUMN IF EXISTS meeting_type;
ALTER TABLE events DROP COLUMN IF EXISTS meeting_url;
ALTER TABLE events DROP COLUMN IF EXISTS physical_address;
ALTER TABLE events DROP COLUMN IF EXISTS color;
