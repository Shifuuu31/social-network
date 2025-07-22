CREATE TABLE IF NOT EXISTS "event" (
    "event_id" VARCHAR(255) not null,
    "group_id" VARCHAR(255) not null,
    "created_by" VARCHAR(255) not null,
    "created_at" DATETIME not null default CURRENT_TIMESTAMP,
    "title" VARCHAR(255) not null,
    "content" VARCHAR(255) not null,
    "date" DATETIME not null,
    primary key ("event_id")
);
