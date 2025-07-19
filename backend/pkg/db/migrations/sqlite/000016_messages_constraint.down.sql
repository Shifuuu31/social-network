-- Migration: Remove constraint for message type
-- This removes the constraint that ensures either receiver_id or group_id exists

-- First, create a backup of existing data
CREATE TABLE messages_backup AS SELECT * FROM messages;

-- Drop the existing table
DROP TABLE messages;

-- Recreate the table without the constraint
CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER NOT NULL,
    receiver_id INTEGER,
    group_id INTEGER,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(sender_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(receiver_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(group_id) REFERENCES groups(id) ON DELETE CASCADE
);

-- Restore the data
INSERT INTO messages SELECT * FROM messages_backup;

-- Drop the backup table
DROP TABLE messages_backup; 