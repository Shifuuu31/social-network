-- Migration: Create notifications table
-- Notifications for user alerts (follow requests, group invites, etc.)

CREATE TABLE IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    type TEXT NOT NULL,
    sub_message TEXT NOT NULL,
    message_id INTEGER, -- ✅ New field: reference to a message
    seen BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(message_id) REFERENCES messages(id) ON DELETE SET NULL
);
