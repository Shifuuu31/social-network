package sqlite

import (
	"database/sql"

	"social-network/pkg/models"
)

type MsgRepository struct {
	DB *sql.DB
}

func (repo *MsgRepository) Save(msg models.ChatMessage) error {
	stmt, err := repo.DB.Prepare("INSERT INTO messages (message_id, sender_id, receiver_id, type, content) values (?,?,?,?,?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(msg.ID, msg.SenderId, msg.ReceiverId, msg.Type, msg.Content); err != nil {
		return err
	}
	return nil
}

func (repo *MsgRepository) SaveGroupMsg(msg models.ChatMessage) error {
	stmt, err := repo.DB.Prepare("INSERT INTO group_messages (message_id, receiver_id) values (?,?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(msg.ID, msg.ReceiverId); err != nil {
		return err
	}
	return nil
}

// needs RECEIVER and SENDER as input
func (repo *MsgRepository) GetAll(msgIn models.ChatMessage) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	rows, err := repo.DB.Query("SELECT message_id,sender_id, receiver_id, type, content FROM messages WHERE (receiver_id = ? AND sender_id = ? )OR (receiver_id = ? AND sender_id = ? ) ORDER BY created_at ASC;", msgIn.ReceiverId, msgIn.SenderId, msgIn.SenderId, msgIn.ReceiverId)
	if err != nil {
		return messages, err
	}
	for rows.Next() {
		var msg models.ChatMessage
		rows.Scan(&msg.ID, &msg.SenderId, &msg.ReceiverId, &msg.Type, &msg.Content)
		messages = append(messages, msg)
	}
	return messages, nil
}

func (repo *MsgRepository) GetAllGroup(userId, groupId string) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	rows, err := repo.DB.Query("SELECT message_id,sender_id, receiver_id, type, content FROM messages WHERE (sender_id = ? AND receiver_id = ? ) OR (receiver_id = ? AND ((SELECT COUNT() FROM groups WHERE group_id = ? AND administrator = ?) = 1 OR (SELECT COUNT() FROM group_users WHERE group_id =? AND user_id =?) = 1) ) ORDER BY created_at ASC;", userId, groupId, groupId, groupId, userId, groupId, userId)
	if err != nil {
		return messages, err
	}
	for rows.Next() {
		var msg models.ChatMessage
		rows.Scan(&msg.ID, &msg.SenderId, &msg.ReceiverId, &msg.Type, &msg.Content)
		messages = append(messages, msg)
	}
	return messages, nil
}

func (repo *MsgRepository) MarkAsRead(msg models.ChatMessage) error {
	_, err := repo.DB.Exec("UPDATE messages SET is_read = ? WHERE message_id=? AND receiver_id =?", 1, msg.ID, msg.ReceiverId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *MsgRepository) MarkAsReadGroup(msg models.ChatMessage) error {
	_, err := repo.DB.Exec("UPDATE group_messages SET is_read = ? WHERE message_id = ? AND  receiver_id = ?", 1, msg.ID, msg.ReceiverId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *MsgRepository) GetUnread(userId string) ([]models.ChatStats, error) {
	var messages []models.ChatStats
	rows, err := repo.DB.Query("SELECT sender_id, type, COUNT(*) FROM messages WHERE receiver_id = ? AND  is_read = 0 GROUP BY sender_id;", userId)
	if err != nil {
		return messages, err
	}
	for rows.Next() {
		var msg models.ChatStats
		rows.Scan(&msg.ID, &msg.Type, &msg.UnreadMsgCount)
		messages = append(messages, msg)
	}
	return messages, nil
}

func (repo *MsgRepository) GetUnreadGroup(userId string) ([]models.ChatStats, error) {
	var messages []models.ChatStats
	rows, err := repo.DB.Query("SELECT receiver_id, type, COUNT(*) FROM messages WHERE type = 'GROUP'AND ((SELECT administrator FROM groups WHERE group_id = messages.receiver_id) = ? OR (SELECT COUNT(*) FROM group_users WHERE group_id = messages.receiver_id AND user_id = ?) = 1) AND (SELECT is_read FROM group_messages WHERE message_id = messages.message_id AND receiver_id = ?) = 0 GROUP BY receiver_id;", userId, userId, userId)
	/*
		SELECT receiver_id, type, COUNT(*) FROM messages WHERE type = 'GROUP' AND
			// is user group admin ?																	-- is group member? --
		((SELECT administrator FROM groups WHERE group_id = messages.receiver_id) = ? OR (SELECT COUNT(*) FROM group_users WHERE group_id = messages.receiver_id AND user_id = ?) = 1)
			AND (SELECT is_read FROM group_messages WHERE message_id = messages.message_id AND receiver_id = ?) = 0 GROUP BY receiver_id;
	*/if err != nil {
		return messages, err
	}
	for rows.Next() {
		var msg models.ChatStats
		rows.Scan(&msg.ID, &msg.Type, &msg.UnreadMsgCount)
		messages = append(messages, msg)
	}
	return messages, nil
}

func (repo *MsgRepository) GetChatHistoryIds(userId string) (map[string]bool, error) {
	idmap := make(map[string]bool)
	// select ids if current is receiver
	rowsReceiver, err := repo.DB.Query("SELECT sender_id FROM messages WHERE receiver_id = ? AND type = 'PERSON';", userId)
	if err != nil {
		return idmap, err
	}
	for rowsReceiver.Next() {
		var id string
		rowsReceiver.Scan(&id)
		idmap[id] = true
	}
	// select ids if current is sender
	rowsSender, err := repo.DB.Query("SELECT receiver_id FROM messages WHERE sender_id = ? AND type = 'PERSON';", userId)
	if err != nil {
		return idmap, err
	}
	for rowsSender.Next() {
		var id string
		rowsSender.Scan(&id)
		idmap[id] = true
	}
	return idmap, nil
}

func (repo *MsgRepository) HasHistory(senderId, receiverId string) (bool, error) {
	row := repo.DB.QueryRow("SELECT COUNT() FROM messages WHERE sender_id = ? AND receiver_id = ? OR sender_id = ? AND receiver_id = ?;", senderId, receiverId, receiverId, senderId)
	var result int
	if err := row.Scan(&result); err != nil {
		return false, err
	}
	if result == 0 {
		return false, nil
	}
	return true, nil
}
