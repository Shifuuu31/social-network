
CREATE TABLE IF NOT EXISTS logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,     -- or SERIAL in PostgreSQL
    level TEXT NOT NULL,                      -- e.g. 'INFO', 'ERROR', 'WARN', 'DEBUG'
    message TEXT NOT NULL,                    -- log message
    metadata TEXT,                            -- optional JSON metadata (e.g. IP, path, etc.)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
