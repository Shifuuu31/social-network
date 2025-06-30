package models

import (
	"database/sql"
	"time"
)

type Message struct {
	ID         int       `json:"id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	GroupID    int       `json:"group_id"`
	Content    string    `json:"content"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"created_at"`
}

type MessageModel struct {
	DB *sql.DB
}


// Insert inserts a message into the DB and updates the message with generated ID and timestamp
func (m *MessageModel) Insert(msg *Message) error {
	query := `
		INSERT INTO messages (sender_id, receiver_id, group_id, content)
		VALUES (?, NULLIF(?, 0), NULLIF(?, 0), ?)
		RETURNING id, created_at
	`

	// Use QueryRow to get generated id and created_at
	return m.DB.QueryRow(query, msg.SenderID, msg.ReceiverID, msg.GroupID, msg.Content).
		Scan(&msg.ID, &msg.CreatedAt)
}