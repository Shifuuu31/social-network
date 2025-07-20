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

func (m *MessageModel) Insert(msg *Message) error {
	query := `
		INSERT INTO messages (sender_id, receiver_id, group_id, content, type)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id, created_at
	`
	return m.DB.QueryRow(query, msg.SenderID, msg.ReceiverID, msg.GroupID, msg.Content, msg.Type).
		Scan(&msg.ID, &msg.CreatedAt)
}