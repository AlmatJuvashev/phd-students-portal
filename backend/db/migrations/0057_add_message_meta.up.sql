ALTER TABLE chat_messages ADD COLUMN meta jsonb DEFAULT '{}'::jsonb;
