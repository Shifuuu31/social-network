-- Migration: Create groups table
-- User-created groups with title and description

CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    creator_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    image_uuid TEXT UNIQUE, 
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(creator_id) REFERENCES users(id) ON DELETE CASCADE
);
