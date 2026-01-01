-- Remove recurrence fields from events
ALTER TABLE events
  DROP COLUMN recurrence_type,
  DROP COLUMN recurrence_end;

-- Note: Removing enum value 'secretary' is not trivial in Postgres and usually not done in simple down migrations.
-- We will leave it as is or requires creating a new type and swapping.
-- For now, just dropping columns.
