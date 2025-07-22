CREATE TABLE IF NOT EXISTS posts (
    "post_id" VARCHAR(255) not null,
    "group_id" varchar(255) null,
    "created_by" varchar(255) not null,
    "created_at" datetime not null default CURRENT_TIMESTAMP,
    "content" TEXT null,
    "image" varchar(255) null,
    "visibility" varchar(255) null default PUBLIC,
    primary key ("post_id")
);

