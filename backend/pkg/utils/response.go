package utils

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/models"
)

type ResponseMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"` // message itself
}
type PostMessage struct {
	Type  string        `json:"type"`
	Posts []models.Post `json:"posts"`
}

type UserMessage struct {
	Type  string        `json:"type"`
	Users []models.User `json:"users"`
}

type GroupMessage struct {
	Type   string         `json:"type"`
	Groups []models.Group `json:"groups"`
}

type EventMessage struct {
	Type   string         `json:"type"`
	Events []models.Event `json:"events"`
}

type NotifMessage struct {
	Type          string                `json:"type"`
	Notifications []models.Notification `json:"notifications"`
}

type ChatMsgMessage struct {
	Type     string               `json:"type"`
	Messages []models.ChatMessage `json:"chatMessage"`
}

type ChatStatMessage struct {
	Type      string             `json:"type"`
	ChatStats []models.ChatStats `json:"chatStats"`
}

// Error takes writer, message, status code and additional error property
// Sets status code in header and encode resp in json
func RespondWithError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	err := ResponseMessage{Message: message, Type: "Error"}
	jsonResp, _ := json.Marshal(err)
	w.Write(jsonResp)
}

// responds with success message
func RespondWithSuccess(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	err := ResponseMessage{Message: message, Type: "Success"}
	jsonResp, _ := json.Marshal(err)
	w.Write(jsonResp)
}

// responds with success group
func RespondWithGroups(w http.ResponseWriter, groups []models.Group, code int) {
	w.WriteHeader(code)
	err := GroupMessage{Groups: groups, Type: "Success"}
	jsonResp, _ := json.Marshal(err)
	w.Write(jsonResp)
}

// responds with success group
func RespondWithUsers(w http.ResponseWriter, users []models.User, code int) {
	w.WriteHeader(code)
	err := UserMessage{Users: users, Type: "Success"}
	jsonResp, _ := json.Marshal(err)
	w.Write(jsonResp)
}

// responds with success group
func RespondWithPosts(w http.ResponseWriter, posts []models.Post, code int) {
	w.WriteHeader(code)
	err := PostMessage{Posts: posts, Type: "Success"}
	jsonResp, _ := json.Marshal(err)
	w.Write(jsonResp)
}

// responds with success events
func RespondWithEvents(w http.ResponseWriter, events []models.Event, code int) {
	w.WriteHeader(code)
	err := EventMessage{Events: events, Type: "Success"}
	jsonResp, _ := json.Marshal(err)
	w.Write(jsonResp)
}

// responds with success notifs
func RespondWithNotifications(w http.ResponseWriter, notifs []models.Notification, code int) {
	w.WriteHeader(code)
	err := NotifMessage{Notifications: notifs, Type: "Success"}
	jsonResp, _ := json.Marshal(err)
	w.Write(jsonResp)
}

// responds with success chat msg
func RespondWithMessages(w http.ResponseWriter, msgs []models.ChatMessage, code int) {
	w.WriteHeader(code)
	err := ChatMsgMessage{Messages: msgs, Type: "Success"}
	jsonResp, _ := json.Marshal(err)
	w.Write(jsonResp)
}

// responds with success chat stats
func RespondWithChatStats(w http.ResponseWriter, msgs []models.ChatStats, code int) {
	w.WriteHeader(code)
	err := ChatStatMessage{ChatStats: msgs, Type: "Success"}
	jsonResp, _ := json.Marshal(err)
	w.Write(jsonResp)
}
