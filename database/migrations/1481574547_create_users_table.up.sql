CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name varchar(255) NOT NULl,
    username varchar(255) NOT NULl,
    password text NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
)
