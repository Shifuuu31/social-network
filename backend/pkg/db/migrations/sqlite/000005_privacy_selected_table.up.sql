-- Migration: Create post_privacy_selected table
-- Links posts with specific users allowed to see posts marked 'selected'

CREATE TABLE IF NOT EXISTS privacy_selected (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    selected_user_id INTEGER NOT NULL,
    UNIQUE(post_id, user_id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
    FOREIGN KEY(selected_user_id) REFERENCES users(id) ON DELETE CASCADE,
);
