CREATE TABLE IF NOT EXISTS notifications (
			"notif_id" VARCHAR(255) not null,
			"user_id" VARCHAR(255) not null,
			"type" VARCHAR(255) not null,
			"content" VARCHAR(255) not null,
			"sender" VARCHAR(255) not null
);