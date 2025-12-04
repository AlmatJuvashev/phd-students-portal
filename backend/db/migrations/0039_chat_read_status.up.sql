-- Track when each user last read messages in a room
CREATE TABLE IF NOT EXISTS chat_room_read_status (
    room_id UUID NOT NULL REFERENCES chat_rooms(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_read_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (room_id, user_id)
);

-- Add index for efficient lookup
CREATE INDEX IF NOT EXISTS idx_chat_room_read_status_user ON chat_room_read_status(user_id);
