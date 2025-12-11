ALTER TABLE chat_rooms DROP CONSTRAINT chat_rooms_type_check;
-- WARNING: This will fail if there are rows with 'group' or 'channel'. 
-- We accept that risk for local dev down migration.
ALTER TABLE chat_rooms ADD CONSTRAINT chat_rooms_type_check CHECK (type IN ('cohort','advisory','other'));
