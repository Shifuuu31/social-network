package models

type ChatMessage struct {
	ID         string `json:"id"`
	SenderId   string `json:"senderId"`
	ReceiverId string `json:"receiverId"`
	Type       string `json:"type"` //GROUP|PERSON
	Content    string `json:"content"`
	Sender User `json:"sender"`
}

type ChatStats struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	UnreadMsgCount int    `json:"unreadMsgCount"`
}

type MsgRepository interface {
	Save(ChatMessage) error
	//get all for specific chat
	// needs  RECEIVER and SENDER as input
	GetAll(ChatMessage) ([]ChatMessage, error)
	GetAllGroup(userId, groupId string) ([]ChatMessage, error)
	GetUnread(userId string) ([]ChatStats, error)
	GetUnreadGroup(userId string) ([]ChatStats, error)
	// mark as read
	MarkAsRead(ChatMessage) error
	MarkAsReadGroup(ChatMessage) error

	SaveGroupMsg(ChatMessage) error

	//returns list of user id's that hve chat history with provided user
	GetChatHistoryIds(userId string)(map[string]bool, error)
	// responds tru if both users have chat history
	HasHistory(senderId, receiverId string) (bool, error)
}
