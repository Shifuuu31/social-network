-- Migration: Create group_members table
-- Tracks membership status (invited/requested/member/declined) in groups

CREATE TABLE IF NOT EXISTS group_members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    prev_status TEXT NOT NULL CHECK(prev_status IN ('none','invited', 'requested', 'member', 'declined')),
    status TEXT NOT NULL CHECK(status IN ('invited', 'requested', 'member', 'declined')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(group_id, user_id),
    FOREIGN KEY(group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
