-- +goose Up
SELECT 'up SQL query';
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_name VARCHAR(50),
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT,
    role VARCHAR(20) CHECK (role IN ('user', 'admin')) DEFAULT 'user',
    source VARCHAR(20) CHECK (
        source IN ('google', 'facebook', 'general')
    ) DEFAULT 'general',
    country_code VARCHAR(5),
    is_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    image TEXT,
    gender VARCHAR(10) CHECK (gender IN ('male', 'female', 'other')),
    refresh_token TEXT,
    is_ban BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
SELECT 'down SQL query';
DROP TABLE IF EXISTS users;

