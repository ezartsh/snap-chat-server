CREATE TABLE IF NOT EXISTS groups (
    id SERIAL PRIMARY KEY,
    room_id BIGINT NOT NULL,
    key VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_group_rooms
        FOREIGN KEY(room_id) 
            REFERENCES rooms(id)
)
