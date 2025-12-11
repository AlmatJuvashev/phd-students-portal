ALTER TABLE chat_rooms DROP CONSTRAINT chat_rooms_type_check;
ALTER TABLE chat_rooms ADD CONSTRAINT chat_rooms_type_check CHECK (type IN ('cohort','advisory','other','group','channel'));
