CREATE TABLE IF NOT EXISTS contacts (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    contact_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_contact_user
        FOREIGN KEY(user_id) 
            REFERENCES users(id),
    CONSTRAINT fk_contact_party_user
        FOREIGN KEY(user_id) 
            REFERENCES users(id)
)
