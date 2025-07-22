package models

type Notification struct {
	ID       string `json:"id"`
	TargetID string `json:"targetId"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	Sender   string `json:"sender"`

	//additional info for notification
	User  User  `json:"user"`
	Event Event `json:"event"`
	Group Group `json:"group"`
}

type NotifRepository interface {
	Save(Notification) error
	Delete(notificationId string) error
	DeleteByType(Notification)error
	CheckIfExists(Notification)(bool, error) // true if exists, false otherwise
	
	//get all pending requests to join group
	GetGroupRequests(groupId string) ([]Notification, error)
	// get specific user_id from request to join
	GetUserFromRequest(notificationId string) (string, error)
	// get group id from specific request
	GetGroupId(notificationId string) (string, error)
	// get all notifications for client
	GetAll(userId string) ([]Notification, error)
	// get Chat_request notifications based on receiver_id
	GetCahtNotifById(notificationId string) (Notification, error)
	// get content form chat_request notification
	GetContentFromChatRequest(senderId, receiverId string)(string, error)
	CheckIfChatRequestExists(senderId, receiverId string)(bool, error) // true if exists, false otherwise
}
