ALTER TABLE chat_messages ADD COLUMN attachments JSONB DEFAULT '[]'::jsonb;
