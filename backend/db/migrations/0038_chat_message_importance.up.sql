-- Add importance column for admin messages (alert/warning)
ALTER TABLE chat_messages ADD COLUMN IF NOT EXISTS importance VARCHAR(20);
