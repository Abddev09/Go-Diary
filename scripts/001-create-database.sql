CREATE DATABASE IF NOT EXISTS blog_system;
USE blog_system;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    image_url VARCHAR(500),
    image_file VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    author_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    author_username VARCHAR(50) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_posts_author ON posts(author_id);
CREATE INDEX IF NOT EXISTS idx_posts_created ON posts(created_at DESC);



\c blog_system;

GRANT ALL PRIVILEGES ON DATABASE blog_system TO blog_user;
GRANT ALL PRIVILEGES ON SCHEMA public TO blog_user;
GRANT CREATE ON SCHEMA public TO blog_user;

ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO blog_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO blog_user;