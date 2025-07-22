package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Notification struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
	// GroupID    int      `json:"group_id"`
	Type       string    `json:"type"`
	SubMessage string    `json:"sub_message"`
	Message    *Message  `json:"message,omitempty"`
	Seen       bool      `json:"seen"`
	CreatedAt  time.Time `json:"created_at"`
}

func (n *Notification) Validate() error {
	// Validate message type
	validTypes := []string{"private", "group", "join_group"}
	isValidType := false
	for _, validType := range validTypes {
		if n.Type == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return errors.New("invalid message type")
	}

	// Type-specific validation
	switch n.Type {
	case "private":
		if n.Message.ReceiverID <= 0 {
			return errors.New("private message requires valid receiver ID")
		}
		if n.Message.SenderID == n.Message.ReceiverID {
			return errors.New("cannot send message to yourself")
		}
	case "group", "join_group":
		if n.Message.GroupID <= 0 {
			return errors.New("group message requires valid group ID")
		}
	}

	return nil
}

type NotificationModel struct {
	DB *sql.DB
}

func (nm *NotificationModel) Upsert(notif *Notification) error {
	query := `
	INSERT INTO notifications (id, user_id, type, sub_message, message_id, seen, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		user_id = excluded.user_id,
		type = excluded.type,
		sub_message = excluded.sub_message,
		message_id = excluded.message_id,
		seen = excluded.seen,
		created_at = excluded.created_at
	RETURNING id, user_id, type, sub_message, message_id, seen, created_at
	`
	row := nm.DB.QueryRow(query,
		notif.ID,
		notif.UserID,
		notif.Type,
		notif.SubMessage,
		notif.Message.ID,
		notif.Seen,
		notif.CreatedAt,
	)

	if err := row.Scan(&notif.ID, &notif.UserID, &notif.Type, &notif.SubMessage, &notif.Message.ID, &notif.Seen, &notif.CreatedAt); err != nil {
		return fmt.Errorf("upsert scan: %w", err)
	}

	if notif.Message.ID > 0 {
		msg := &Message{}
		err := nm.DB.QueryRow(`SELECT id, sender_id, receiver_id, group_id, content, created_at FROM messages WHERE id = ?`, notif.Message.ID).
			Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.GroupID, &msg.Content, &msg.CreatedAt)
		if err == nil {
			notif.Message = msg
		}
	}

	return nil
}

// Insert inserts a new notification into the database
func (nm *NotificationModel) Insert(notification *Notification) error {
	query := `
		INSERT INTO notifications (user_id, type, message, seen)
		VALUES (?, ?, ?, ?)
	`

	res, err := nm.DB.Exec(query, notification.UserID, notification.Type, notification.Message, notification.Seen)
	if err != nil {
		return fmt.Errorf("insert notification: %w", err)
	}

	lastID, err := res.LastInsertId()
	if err == nil {
		notification.ID = int(lastID)
	}

	return nil
}

func (nm *NotificationModel) Delete(notificationID int) error {
	if _, err := nm.DB.Exec(`DELETE FROM notifications WHERE id = ?`, notificationID); err != nil {
		return fmt.Errorf("delete notification: %w", err)
	}
	return nil
}

// MarkAsSeen marks a specific notification as seen for a user
func (nm *NotificationModel) MarkAsSeen(notificationID, userID int) error {
	if _, err := nm.DB.Exec(`UPDATE notifications SET seen = 1 WHERE id = ? AND user_id = ?`, notificationID, userID); err != nil {
		return fmt.Errorf("mark notification as seen: %w", err)
	}
	return nil
}

// MarkAllAsSeen marks all user notifications as seen
func (nm *NotificationModel) MarkAllAsSeen(userID int) error {
	if _, err := nm.DB.Exec(`UPDATE notifications SET seen = 1 WHERE user_id = ? AND seen = 0`, userID); err != nil {
		return fmt.Errorf("mark all as seen: %w", err)
	}
	return nil
}

func (nm *NotificationModel) CountUnseen(userID int) (count int, err error) {
	if err = nm.DB.QueryRow(`SELECT COUNT(*) FROM notifications WHERE user_id = ? AND seen = 0`, userID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count unseen: %w", err)
	}
	return count, nil
}

type NotificationPayload struct {
	UserID     int    `json:"user_id"`
	Start      int    `json:"start"`
	NumOfItems int    `json:"n_items"`
	Type       string `json:"type"`
}

func (nm *NotificationModel) GetByUser(payload *NotificationPayload) ([]*Notification, error) {
	query := `SELECT id, user_id, type, sub_message, message_id, seen, created_at FROM notifications WHERE user_id = ?`
	args := []any{payload.UserID}

	switch payload.Type {
	case "seen":
		query += ` AND seen = 1`
	case "unseen":
		query += ` AND seen = 0`
	}

	query += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, payload.NumOfItems, payload.Start)

	rows, err := nm.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("get notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*Notification
	for rows.Next() {
		notif := &Notification{}
		if err := rows.Scan(&notif.ID, &notif.UserID, &notif.Type, &notif.SubMessage, &notif.Message.ID, &notif.Seen, &notif.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan notification: %w", err)
		}
		notifications = append(notifications, notif)
	}
	return notifications, nil
}

// GetByID retrieves a notification by its ID
func (nm *NotificationModel) GetByID(notification *Notification) error {
	query := `SELECT id, user_id, type, message, seen, created_at FROM notifications WHERE id = ?`

	err := nm.DB.QueryRow(query, notification.ID).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Type,
		&notification.Message,
		&notification.Seen,
		&notification.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("get notification by id: %w", err)
	}

	return nil
}
