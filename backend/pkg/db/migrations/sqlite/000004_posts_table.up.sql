-- Migration: Create posts table
-- User posts with text, optional images, and privacy settings

CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    group_id INTEGER,
    content TEXT NOT NULL,
    image_uuid TEXT NOT NULL UNIQUE, 
    privacy TEXT NOT NULL CHECK(privacy IN ('public', 'followers', 'selected', 'group')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
