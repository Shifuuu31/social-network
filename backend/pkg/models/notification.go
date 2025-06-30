package models

import (
	"database/sql"
	"fmt"
)

type Notification struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Type      string `json:"type"`
	Message   string `json:"message"` 
	Seen      bool   `json:"seen"`
	CreatedAt int    `json:"created_at"`
}

type NotificationModel struct {
	DB *sql.DB
}

// Upsert inserts or replaces a notification
func (nm *NotificationModel) Upsert(notification *Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, type, message, seen, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			user_id = excluded.user_id,
			type = excluded.type,
			message = excluded.message,
			seen = excluded.seen,
			created_at = excluded.created_at
	`

	if _, err := nm.DB.Exec(query, notification.ID, notification.UserID, notification.Type, notification.Message, notification.Seen, notification.CreatedAt); err != nil {
		return fmt.Errorf("upsert notification: %w", err)
	}
	return nil
}

// Delete deletes a specific notification
func (nm *NotificationModel) Delete(notificationID int) error {
	if _, err := nm.DB.Exec(`DELETE FROM notifications WHERE id = ?`, notificationID); err != nil {
		return fmt.Errorf("delete notification: %w", err)
	}
	return nil
}

// TODO
// func (nm *Notifica.

// MarkAllAsSeen marks all user notifications as seen
func (nm *NotificationModel) MarkAllAsSeen(userID int) error {
	if _, err := nm.DB.Exec(`UPDATE notifications SET seen = 1 WHERE user_id = ? AND seen = 0`, userID); err != nil {
		return fmt.Errorf("mark all as seen: %w", err)
	}
	return nil
}

// CountUnseen returns the number of unseen notifications for a user
func (nm *NotificationModel) CountUnseen(userID int) (count int, err error) {
	if err = nm.DB.QueryRow(`SELECT COUNT(*) FROM notifications WHERE user_id = ? AND seen = 0`, userID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count unseen: %w", err)
	}
	return count, nil
}

type NotificationPayload struct {
	UserID     string `json:"user_id"`
	Start      int    `json:"star"`
	NumOfItems int    `json:"n_items"`
	Type       string `json:"type"` // all, seen, unseen
}

// GetByUser retrieves notifications for a user with optional filters
func (nm *NotificationModel) GetByUser(payload *NotificationPayload) ([]*Notification, error) {
	query := `SELECT id, user_id, type, message, seen, created_at FROM notifications WHERE user_id = ?`
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
		n := &Notification{}
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Message, &n.Seen, &n.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan notification: %w", err)
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}
