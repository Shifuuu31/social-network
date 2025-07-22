CREATE TABLE IF NOT EXISTS users (
    "user_id" varchar(255) not null,
    "created_at" datetime not null default CURRENT_TIMESTAMP,
    "email" varchar(255) not null,
    "first_name" varchar(255) not null,
    "last_name" varchar(255) not null,
    "nickname" varchar(255) null,
    "birthday" datetime not null,
    "image" varchar(255) null,
    "about" TEXT null,
    "status" varchar(255) not null default PUBLIC,
    "password" varchar(255) not null,
    primary key ("user_id")
);
