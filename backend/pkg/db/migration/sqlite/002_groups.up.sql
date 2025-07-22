
CREATE TABLE IF NOT EXISTS groups (
    "group_id" VARCHAR(255) not null,
    "administrator" VARCHAR(255) not null,
    "name" VARCHAR(255) not null,
    "description" VARCHAR(255) null,
    primary key ("group_id")
);
