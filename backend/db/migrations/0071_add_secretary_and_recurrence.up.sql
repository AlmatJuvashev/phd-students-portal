-- Add secretary to user_role enum
ALTER TYPE user_role ADD VALUE 'secretary';

-- Add recurrence fields to events
ALTER TABLE events
  ADD COLUMN recurrence_type VARCHAR(50),
  ADD COLUMN recurrence_end TIMESTAMPTZ;
