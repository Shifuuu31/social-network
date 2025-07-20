package models

import (
	"database/sql"
	"fmt"
)

type Notification struct {
	ID         int      `json:"id"`
	UserID     int      `json:"user_id"`
	Type       string   `json:"type"`
	SubMessage string   `json:"sub_message"`
	MessageID  int      `json:"message_id"`
	Message    *Message `json:"message,omitempty"`
	Seen       bool     `json:"seen"`
	CreatedAt  int      `json:"created_at"`
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
		notif.MessageID,
		notif.Seen,
		notif.CreatedAt,
	)

	if err := row.Scan(&notif.ID, &notif.UserID, &notif.Type, &notif.SubMessage, &notif.MessageID, &notif.Seen, &notif.CreatedAt); err != nil {
		return fmt.Errorf("upsert scan: %w", err)
	}

	if notif.MessageID > 0 {
		msg := &Message{}
		err := nm.DB.QueryRow(`SELECT id, sender_id, receiver_id, group_id, content, type, created_at FROM messages WHERE id = ?`, notif.MessageID).
			Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.GroupID, &msg.Content, &msg.Type, &msg.CreatedAt)
		if err == nil {
			notif.Message = msg
		}
	}

	return nil
}

func (nm *NotificationModel) Delete(notificationID int) error {
	if _, err := nm.DB.Exec(`DELETE FROM notifications WHERE id = ?`, notificationID); err != nil {
		return fmt.Errorf("delete notification: %w", err)
	}
	return nil
}

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
	UserID     string `json:"user_id"`
	Start      int    `json:"star"`
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
		if err := rows.Scan(&notif.ID, &notif.UserID, &notif.Type, &notif.SubMessage, &notif.MessageID, &notif.Seen, &notif.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan notification: %w", err)
		}
		notifications = append(notifications, notif)
	}
	return notifications, nil
}
