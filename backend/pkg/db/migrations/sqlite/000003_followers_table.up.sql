-- Migration: Create followers table
-- Tracks follow requests and follow relationships between users with status

CREATE TABLE IF NOT EXISTS follow_request (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    from_user_id INTEGER NOT NULL,
    to_user_id INTEGER NOT NULL,
    status TEXT NOT NULL CHECK(status IN ('pending', 'accepted', 'declined')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(from_user_id, to_user_id),
    FOREIGN KEY(from_user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(to_user_id) REFERENCES users(id) ON DELETE CASCADE
);
