-- Migration: Create event_votes table
-- RSVP responses (going/not going) from users to events

CREATE TABLE IF NOT EXISTS event_votes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    vote TEXT NOT NULL CHECK(vote IN ('going', 'not_going')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(event_id, user_id),
    FOREIGN KEY(event_id) REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
