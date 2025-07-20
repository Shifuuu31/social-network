package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

type Message struct {
	ID         int       `json:"id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	GroupID    int       `json:"group_id"`
	Content    string    `json:"content"`
	Type       string    `json:"type"` // 'private', 'group', 'join_group'
	CreatedAt  time.Time `json:"created_at"`
}

type MessageModel struct {
	DB *sql.DB
}

// Validate validates message data before insertion
func (m *Message) Validate() error {
	// Check content
	if strings.TrimSpace(m.Content) == "" {
		return errors.New("message content cannot be empty")
	}
	
	if len(m.Content) > 1000 {
		return errors.New("message content too long (max 1000 characters)")
	}

	// Check sender ID
	if m.SenderID <= 0 {
		return errors.New("invalid sender ID")
	}

	// Validate message type
	validTypes := []string{"private", "group", "join_group"}
	isValidType := false
	for _, validType := range validTypes {
		if m.Type == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return errors.New("invalid message type")
	}

	// Type-specific validation
	switch m.Type {
	case "private":
		if m.ReceiverID <= 0 {
			return errors.New("private message requires valid receiver ID")
		}
		if m.SenderID == m.ReceiverID {
			return errors.New("cannot send message to yourself")
		}
	case "group", "join_group":
		if m.GroupID <= 0 {
			return errors.New("group message requires valid group ID")
		}
	}

	return nil
}

// Insert inserts a message into the DB and updates the message with generated ID and timestamp
func (m *MessageModel) Insert(msg *Message) error {
	// Validate message before insertion
	if err := msg.Validate(); err != nil {
		return err
	}

	query := `
		INSERT INTO messages (sender_id, receiver_id, group_id, content, type)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id, created_at
	`

	// Use QueryRow to get generated id and created_at
	return m.DB.QueryRow(query, msg.SenderID, msg.ReceiverID, msg.GroupID, msg.Content, msg.Type).
		Scan(&msg.ID, &msg.CreatedAt)
}

// GetMessagesByConversation retrieves messages between two users (private chat)
func (m *MessageModel) GetMessagesByConversation(userID1, userID2 int, limit int) ([]*Message, error) {
	if userID1 <= 0 || userID2 <= 0 {
		return nil, errors.New("invalid user IDs")
	}

	query := `
		SELECT id, sender_id, receiver_id, group_id, content, type, created_at
		FROM messages 
		WHERE type = 'private' 
		AND ((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?))
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := m.DB.Query(query, userID1, userID2, userID2, userID1, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		msg := &Message{}
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.GroupID, 
			&msg.Content, &msg.Type, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

// GetMessagesByGroup retrieves messages for a specific group
func (m *MessageModel) GetMessagesByGroup(groupID int, limit int) ([]*Message, error) {
	if groupID <= 0 {
		return nil, errors.New("invalid group ID")
	}

	query := `
		SELECT id, sender_id, receiver_id, group_id, content, type, created_at
		FROM messages 
		WHERE group_id = ? AND (type = 'group' OR type = 'join_group')
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := m.DB.Query(query, groupID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		msg := &Message{}
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.GroupID, 
			&msg.Content, &msg.Type, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

// Delete removes a message (only sender can delete their own messages)
func (m *MessageModel) Delete(messageID, userID int) error {
	if messageID <= 0 || userID <= 0 {
		return errors.New("invalid message or user ID")
	}

	query := `DELETE FROM messages WHERE id = ? AND sender_id = ?`
	result, err := m.DB.Exec(query, messageID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("message not found or unauthorized to delete")
	}

	return nil
}
