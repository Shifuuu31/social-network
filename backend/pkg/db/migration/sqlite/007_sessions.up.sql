CREATE TABLE IF NOT EXISTS "sessions" (
    "session_id" VARCHAR(255) NOT NULL PRIMARY KEY,
    "user_id" VARCHAR(255) NOT NULL,
    "expiration_time" DATETIME NOT NULL
);
