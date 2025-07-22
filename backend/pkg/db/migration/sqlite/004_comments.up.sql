CREATE TABLE IF NOT EXISTS comments (
    "comment_id" VARCHAR(255) not null,
    "post_id" VARCHAR(255) not null,
    "created_at" datetime not null default CURRENT_TIMESTAMP,
    "created_by" varchar(255) not null,
    "content" text null,
    "image" varchar(255) null,
    primary key ("comment_id")
);

