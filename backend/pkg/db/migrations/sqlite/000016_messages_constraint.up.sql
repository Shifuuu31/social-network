-- Migration: Recreate messages table with constraint to ensure either receiver_id or group_id exists
-- This ensures messages are either direct messages (receiver_id) or group messages (group_id)

-- First, create a backup of existing data
CREATE TABLE messages_backup AS SELECT * FROM messages;

-- Clean the data to ensure constraint compliance
-- If both receiver_id and group_id are NULL, set receiver_id to a default value (or delete the message)
-- If both are NOT NULL, prioritize receiver_id and set group_id to NULL
UPDATE messages_backup 
SET group_id = NULL 
WHERE receiver_id IS NOT NULL AND group_id IS NOT NULL;

-- For messages with both NULL, we'll set receiver_id to 1 (assuming user 1 exists)
-- If this is not appropriate, you might want to delete these messages instead
UPDATE messages_backup 
SET receiver_id = 1 
WHERE receiver_id IS NULL AND group_id IS NULL;

-- Drop the existing table
DROP TABLE messages;

-- Recreate the table with the constraint
CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER NOT NULL,
    receiver_id INTEGER,
    group_id INTEGER,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(sender_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(receiver_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(group_id) REFERENCES groups(id) ON DELETE CASCADE,
    CONSTRAINT check_message_type CHECK (
        (receiver_id IS NOT NULL AND group_id IS NULL) OR 
        (receiver_id IS NULL AND group_id IS NOT NULL)
    )
);

-- Restore the cleaned data
INSERT INTO messages SELECT * FROM messages_backup;

-- Drop the backup table
DROP TABLE messages_backup; 