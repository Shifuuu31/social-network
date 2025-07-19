package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

// Debug logging utility
func debugLog(method, message string, data ...interface{}) {
	log.Printf("🔍 [MessageModel.%s] %s %v", method, message, data)
}

// Message represents a chat message in the database
type Message struct {
	ID         int       `json:"id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

// MessageWithUser includes sender/receiver user information
type MessageWithUser struct {
	Message
	SenderName   string `json:"sender_name"`
	SenderAvatar string `json:"sender_avatar"`
	ReceiverName string `json:"receiver_name"`
}

// MessageModel handles database operations for messages
type MessageModel struct {
	DB *sql.DB
}

// Validate checks if the message is valid
func (m *Message) Validate() error {
	debugLog("Validate", "Validating message: SenderID=%d, ReceiverID=%d, Content='%s'", m.SenderID, m.ReceiverID, m.Content)

	if m.SenderID <= 0 {
		debugLog("Validate", "Validation failed: sender_id must be a positive integer")
		return errors.New("sender_id must be a positive integer")
	}
	if m.ReceiverID <= 0 {
		debugLog("Validate", "Validation failed: receiver_id must be a positive integer")
		return errors.New("receiver_id must be a positive integer")
	}
	if m.SenderID == m.ReceiverID {
		debugLog("Validate", "Validation failed: cannot send message to yourself")
		return errors.New("cannot send message to yourself")
	}
	if m.Content == "" {
		debugLog("Validate", "Validation failed: message content cannot be empty")
		return errors.New("message content cannot be empty")
	}
	if len(m.Content) > 1000 {
		debugLog("Validate", "Validation failed: message content too long (max 1000 characters)")
		return errors.New("message content too long (max 1000 characters)")
	}

	debugLog("Validate", "Message validation passed")
	return nil
}

// CanSendMessage checks if sender can send message to receiver
// At least one user must be following the other
func (mm *MessageModel) CanSendMessage(senderID, receiverID int) (bool, error) {
	debugLog("CanSendMessage", "Checking if user %d can send message to user %d", senderID, receiverID)

	query := `
		SELECT COUNT(*) FROM follow_request 
		WHERE ((from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?))
		AND status = 'accepted'
	`
	debugLog("CanSendMessage", "Executing query: %s with params: [%d, %d, %d, %d]", query, senderID, receiverID, receiverID, senderID)

	var count int
	err := mm.DB.QueryRow(query, senderID, receiverID, receiverID, senderID).Scan(&count)
	if err != nil {
		debugLog("CanSendMessage", "Database error: %v", err)
		return false, fmt.Errorf("error checking follow relationship: %w", err)
	}

	debugLog("CanSendMessage", "Follow relationship count: %d", count)
	canSend := count > 0
	debugLog("CanSendMessage", "Can send message: %t", canSend)
	return canSend, nil
}

// CreateMessage saves a new message to the database
func (mm *MessageModel) CreateMessage(message *Message) error {
	debugLog("CreateMessage", "=== START: CreateMessage ===")
	debugLog("CreateMessage", "Message to create: %+v", message)

	// Validate message
	debugLog("CreateMessage", "Validating message...")
	if err := message.Validate(); err != nil {
		debugLog("CreateMessage", "Message validation failed: %v", err)
		return err
	}
	debugLog("CreateMessage", "Message validation passed")

	// Check if users can message each other
	debugLog("CreateMessage", "Checking if users can message each other...")
	canSend, err := mm.CanSendMessage(message.SenderID, message.ReceiverID)
	if err != nil {
		debugLog("CreateMessage", "Error checking if users can send messages: %v", err)
		return err
	}
	if !canSend {
		debugLog("CreateMessage", "Users cannot send messages to each other")
		return errors.New("cannot send message: no follow relationship exists")
	}
	debugLog("CreateMessage", "Users can send messages to each other")

	// Insert message into database
	query := `
		INSERT INTO messages (sender_id, receiver_id, content)
		VALUES (?, ?, ?)
	`
	debugLog("CreateMessage", "Executing INSERT query: %s with params: [%d, %d, '%s']", query, message.SenderID, message.ReceiverID, message.Content)

	result, err := mm.DB.Exec(query, message.SenderID, message.ReceiverID, message.Content)
	if err != nil {
		debugLog("CreateMessage", "Database INSERT failed: %v", err)
		return fmt.Errorf("error creating message: %w", err)
	}
	debugLog("CreateMessage", "Database INSERT successful")

	// Get the inserted message ID
	debugLog("CreateMessage", "Getting last insert ID...")
	id, err := result.LastInsertId()
	if err != nil {
		debugLog("CreateMessage", "Error getting last insert ID: %v", err)
		return fmt.Errorf("error getting last insert id: %w", err)
	}

	message.ID = int(id)
	debugLog("CreateMessage", "Message created successfully with ID: %d", message.ID)
	debugLog("CreateMessage", "=== END: CreateMessage ===")
	return nil
}

// GetConversation retrieves messages between two users
func (mm *MessageModel) GetConversation(user1ID, user2ID int, limit, offset int) ([]*MessageWithUser, error) {
	debugLog("GetConversation", "=== START: GetConversation ===")
	debugLog("GetConversation", "Getting conversation between users %d and %d, limit: %d, offset: %d", user1ID, user2ID, limit, offset)

	query := `
		SELECT m.id, m.sender_id, m.receiver_id, m.content, m.created_at,
		       COALESCE(s.nickname, s.first_name || ' ' || s.last_name) as sender_name,
		       s.avatar_url as sender_avatar,
		       COALESCE(r.nickname, r.first_name || ' ' || r.last_name) as receiver_name
		FROM messages m
		JOIN users s ON m.sender_id = s.id
		JOIN users r ON m.receiver_id = r.id
		WHERE (m.sender_id = ? AND m.receiver_id = ?) OR (m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY m.created_at DESC
		LIMIT ? OFFSET ?
	`
	debugLog("GetConversation", "Executing query with params: [%d, %d, %d, %d, %d, %d]", user1ID, user2ID, user2ID, user1ID, limit, offset)

	rows, err := mm.DB.Query(query, user1ID, user2ID, user2ID, user1ID, limit, offset)
	if err != nil {
		debugLog("GetConversation", "Database query failed: %v", err)
		return nil, fmt.Errorf("error fetching conversation: %w", err)
	}
	defer rows.Close()
	debugLog("GetConversation", "Database query successful")

	var messages []*MessageWithUser
	for rows.Next() {
		var msg MessageWithUser
		err := rows.Scan(
			&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.CreatedAt,
			&msg.SenderName, &msg.SenderAvatar, &msg.ReceiverName,
		)
		if err != nil {
			debugLog("GetConversation", "Error scanning message row: %v", err)
			return nil, fmt.Errorf("error scanning message: %w", err)
		}
		debugLog("GetConversation", "Scanned message: ID=%d, SenderID=%d, ReceiverID=%d, Content='%s'", msg.ID, msg.SenderID, msg.ReceiverID, msg.Content)
		messages = append(messages, &msg)
	}

	debugLog("GetConversation", "Retrieved %d messages", len(messages))
	debugLog("GetConversation", "=== END: GetConversation ===")
	return messages, nil
}

// GetRecentConversations gets recent conversations for a user
func (mm *MessageModel) GetRecentConversations(userID int, limit int) ([]*MessageWithUser, error) {
	query := `
		SELECT DISTINCT m.id, m.sender_id, m.receiver_id, m.content, m.created_at,
		       COALESCE(s.nickname, s.first_name || ' ' || s.last_name) as sender_name,
		       s.avatar_url as sender_avatar,
		       COALESCE(r.nickname, r.first_name || ' ' || r.last_name) as receiver_name
		FROM messages m
		JOIN users s ON m.sender_id = s.id
		JOIN users r ON m.receiver_id = r.id
		WHERE m.id IN (
			SELECT MAX(id) FROM messages 
			WHERE sender_id = ? OR receiver_id = ?
			GROUP BY 
				CASE 
					WHEN sender_id = ? THEN receiver_id 
					ELSE sender_id 
				END
		)
		ORDER BY m.created_at DESC
		LIMIT ?
	`

	rows, err := mm.DB.Query(query, userID, userID, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("error fetching recent conversations: %w", err)
	}
	defer rows.Close()

	var messages []*MessageWithUser
	for rows.Next() {
		var msg MessageWithUser
		err := rows.Scan(
			&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.CreatedAt,
			&msg.SenderName, &msg.SenderAvatar, &msg.ReceiverName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning message: %w", err)
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

// DeleteMessage deletes a message (only by sender)
func (mm *MessageModel) DeleteMessage(messageID, userID int) error {
	query := `DELETE FROM messages WHERE id = ? AND sender_id = ?`
	result, err := mm.DB.Exec(query, messageID, userID)
	if err != nil {
		return fmt.Errorf("error deleting message: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %w", err)
	}

	if affected == 0 {
		return errors.New("message not found or you don't have permission to delete it")
	}

	return nil
}
