CREATE TABLE IF NOT EXISTS rooms_messages (
    id SERIAL PRIMARY KEY,
    room_id BIGINT NOT NULL,
    sender_id BIGINT NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_room_messages_user
        FOREIGN KEY(sender_id) 
            REFERENCES users(id)
)
