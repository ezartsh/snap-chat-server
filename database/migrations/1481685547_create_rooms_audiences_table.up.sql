CREATE TABLE IF NOT EXISTS rooms_audiences (
    id SERIAL PRIMARY KEY,
    room_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_rooms_audiences_room
        FOREIGN KEY(room_id) 
            REFERENCES rooms(id),
    CONSTRAINT fk_rooms_audiences_user
        FOREIGN KEY(user_id) 
            REFERENCES users(id)
)
