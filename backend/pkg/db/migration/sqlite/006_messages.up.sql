CREATE TABLE IF NOT EXISTS messages (
			"message_id" VARCHAR(255) not null,
			"sender_id" VARCHAR(255) not null,
			"receiver_id" VARCHAR(255) not null,
			"type" VARCHAR(255) not null,
			"created_at" datetime not null default CURRENT_TIMESTAMP,
			"content" VARCHAR(255) not null,
			"is_read" INT default 0,
			primary key ("message_id")
		);

