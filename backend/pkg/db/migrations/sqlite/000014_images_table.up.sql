-- Migration: Create images table
-- Stores metadata for all uploaded images (posts, profiles, comments)

CREATE TABLE IF NOT EXISTS images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    filename VARCHAR(255) NOT NULL,           -- Generated secure filename
    original_name VARCHAR(255),               -- Original uploaded filename
    mime_type VARCHAR(50) NOT NULL,          -- MIME type (image/jpeg, image/png, etc.)
    file_size INTEGER NOT NULL,              -- File size in bytes
    upload_path VARCHAR(500) NOT NULL,       -- Relative path (posts/filename.jpg)
    uploaded_by INTEGER NOT NULL,            -- User who uploaded the image
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(uploaded_by) REFERENCES users(id) ON DELETE CASCADE
);

-- Index for faster lookups
CREATE INDEX IF NOT EXISTS idx_images_uploaded_by ON images(uploaded_by);
CREATE INDEX IF NOT EXISTS idx_images_upload_path ON images(upload_path); 