-- Add meeting type and related fields for online/offline events
ALTER TABLE events ADD COLUMN IF NOT EXISTS meeting_type VARCHAR(20) DEFAULT 'offline';
-- 'online' | 'offline'

ALTER TABLE events ADD COLUMN IF NOT EXISTS meeting_url TEXT;
-- For Zoom/Google Meet links (used when meeting_type = 'online')

ALTER TABLE events ADD COLUMN IF NOT EXISTS physical_address TEXT;
-- For physical location details (used when meeting_type = 'offline')

-- Add color field for custom event colors
ALTER TABLE events ADD COLUMN IF NOT EXISTS color VARCHAR(20);
-- Stores hex color or color key like 'blue', 'red', 'green', 'purple'
