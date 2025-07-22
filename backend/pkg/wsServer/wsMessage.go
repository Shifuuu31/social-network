package ws

import (
	"encoding/json"
	"log"
	"social-network/pkg/models"
)

/* --------------------- actions for websocket messages --------------------- */
const NotificationAction = "notification"
const ChatAction = "chat"
const GroupAcceptAction = "groupAccept"

type WsMessage struct {
	UserID       string              `json:"uid"`
	Action       string              `json:"action"` //msg request action
	Notification models.Notification `json:"notification"`
	ChatMessage  models.ChatMessage  `json:"chatMessage"`
	Message      string              `json:"message"`
}

// encode method that can be called to create a json []byte object
func (message *WsMessage) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}
	return json
}
