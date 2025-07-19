package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"social-network/pkg/middleware"
	"social-network/pkg/models"
	"social-network/pkg/tools"
	"social-network/pkg/websocket"
)

// Debug logging utility
func debugLog(method, message string, data ...interface{}) {
	log.Printf("🔍 [ChatHandler.%s] %s %v", method, message, data)
}

// ChatHandler handles chat-related HTTP requests
type ChatHandler struct {
	DL  *middleware.DataLayer
	Hub *websocket.Hub
}

// NewChatHandler creates a new chat handler
func NewChatHandler(dl *middleware.DataLayer, hub *websocket.Hub) *ChatHandler {
	debugLog("NewChatHandler", "Creating new chat handler")
	return &ChatHandler{
		DL:  dl,
		Hub: hub,
	}
}

// SendMessageRequest represents a request to send a message
type SendMessageRequest struct {
	ReceiverID int    `json:"receiver_id"`
	Content    string `json:"content"`
}

// SendMessage handles sending a new message
func (ch *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	debugLog("SendMessage", "=== START: SendMessage handler ===")

	if r.Method != http.MethodPost {
		debugLog("SendMessage", "Method not allowed: %s", r.Method)
		tools.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	requesterID := ch.DL.GetRequesterID(w, r)
	debugLog("SendMessage", "Requester ID: %d", requesterID)
	if requesterID == 0 {
		debugLog("SendMessage", "No requester ID found, returning")
		return
	}

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		debugLog("SendMessage", "Failed to decode request body: %v", err)
		tools.RespondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	debugLog("SendMessage", "Request decoded successfully: ReceiverID=%d, Content='%s'", req.ReceiverID, req.Content)

	message := &models.Message{
		SenderID:   requesterID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
		CreatedAt:  time.Now(),
	}
	debugLog("SendMessage", "Created message object: %+v", message)

	// Check if users can send messages to each other
	debugLog("SendMessage", "Checking if users can send messages...")
	canSend, err := ch.DL.Messages.CanSendMessage(requesterID, req.ReceiverID)
	if err != nil {
		debugLog("SendMessage", "Error checking if users can send messages: %v", err)
		tools.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}
	debugLog("SendMessage", "Can send message: %t", canSend)

	if !canSend {
		debugLog("SendMessage", "Users cannot send messages to each other")
		tools.RespondError(w, "Cannot send message: users must follow each other", http.StatusForbidden)
		return
	}

	// Save message to database
	debugLog("SendMessage", "Attempting to save message to database...")
	if err := ch.DL.Messages.CreateMessage(message); err != nil {
		debugLog("SendMessage", "Failed to save message to database: %v", err)
		tools.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}
	debugLog("SendMessage", "Message saved successfully! Message ID: %d", message.ID)

	// Send via WebSocket if receiver is online
	debugLog("SendMessage", "Broadcasting message via WebSocket...")
	wsMessage := &websocket.Message{
		Type:       "message",
		SenderID:   requesterID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
		Timestamp:  time.Now().Format(time.RFC3339),
	}
	ch.Hub.Broadcast <- wsMessage
	debugLog("SendMessage", "WebSocket message broadcasted")

	// Return success response
	response := map[string]interface{}{
		"message":    "Message sent successfully",
		"message_id": message.ID,
	}
	debugLog("SendMessage", "Sending success response: %+v", response)
	tools.EncodeJSON(w, http.StatusCreated, response)
	debugLog("SendMessage", "=== END: SendMessage handler ===")
}

// GetConversationRequest represents a request to get conversation
type GetConversationRequest struct {
	OtherUserID int `json:"other_user_id"`
	Limit       int `json:"limit"`
	Offset      int `json:"offset"`
}

// GetConversation retrieves messages between two users
func (ch *ChatHandler) GetConversation(w http.ResponseWriter, r *http.Request) {
	debugLog("GetConversation", "=== START: GetConversation handler ===")

	if r.Method != http.MethodPost {
		debugLog("GetConversation", "Method not allowed: %s", r.Method)
		tools.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	requesterID := ch.DL.GetRequesterID(w, r)
	debugLog("GetConversation", "Requester ID: %d", requesterID)
	if requesterID == 0 {
		debugLog("GetConversation", "No requester ID found, returning")
		return
	}

	var req GetConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		debugLog("GetConversation", "Failed to decode request body: %v", err)
		tools.RespondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	debugLog("GetConversation", "Request decoded successfully: OtherUserID=%d, Limit=%d, Offset=%d", req.OtherUserID, req.Limit, req.Offset)

	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 50
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	debugLog("GetConversation", "Calling Messages.GetConversation with: user1=%d, user2=%d, limit=%d, offset=%d", requesterID, req.OtherUserID, req.Limit, req.Offset)
	messages, err := ch.DL.Messages.GetConversation(requesterID, req.OtherUserID, req.Limit, req.Offset)
	if err != nil {
		debugLog("GetConversation", "Error getting conversation: %v", err)
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	debugLog("GetConversation", "Successfully retrieved %d messages", len(messages))
	response := map[string]interface{}{
		"messages": messages,
		"count":    len(messages),
	}
	debugLog("GetConversation", "Sending response: %+v", response)
	tools.EncodeJSON(w, http.StatusOK, response)
	debugLog("GetConversation", "=== END: GetConversation handler ===")
}

// GetRecentConversations retrieves recent conversations for the user
func (ch *ChatHandler) GetRecentConversations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		tools.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	requesterID := ch.DL.GetRequesterID(w, r)
	if requesterID == 0 {
		return
	}

	// Get limit from query params
	limitStr := r.URL.Query().Get("limit")
	limit := 20 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}
	fmt.Println("requesterID", requesterID)
	messages, err := ch.DL.Messages.GetRecentConversations(requesterID, limit)
	if err != nil {
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("messages", messages)

	tools.EncodeJSON(w, http.StatusOK, map[string]interface{}{
		"conversations": messages,
		"count":         len(messages),
	})
}

// DeleteMessage deletes a message
func (ch *ChatHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		tools.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	requesterID := ch.DL.GetRequesterID(w, r)
	if requesterID == 0 {
		return
	}

	// Get message ID from URL
	messageIDStr := r.URL.Query().Get("message_id")
	if messageIDStr == "" {
		tools.RespondError(w, "Message ID is required", http.StatusBadRequest)
		return
	}

	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		tools.RespondError(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	if err := ch.DL.Messages.DeleteMessage(messageID, requesterID); err != nil {
		tools.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	tools.EncodeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Message deleted successfully",
	})
}

// HandleWebSocket handles WebSocket upgrade and connection
func (ch *ChatHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Try to get user ID from session cookie first (normal HTTP auth)
	requesterID := ch.DL.GetRequesterID(w, r)

	// If no session cookie, try token from query parameter (for WebSocket)
	if requesterID == 0 {
		token := r.URL.Query().Get("token")
		if token != "" {
			session, err := ch.DL.Sessions.GetSessionByToken(token)
			if err == nil && session != nil {
				requesterID = session.UserID
			}
		}
	}

	if requesterID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	websocket.ServeWs(ch.Hub, w, r, requesterID)
}

// SetupChatRoutes sets up chat-related routes
func (ch *ChatHandler) SetupChatRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	fmt.Println("SetupChatRoutes555555")

	mux.HandleFunc("/send", ch.SendMessage)
	mux.HandleFunc("/conversation", ch.GetConversation)
	mux.HandleFunc("/recent", ch.GetRecentConversations)
	mux.HandleFunc("/delete", ch.DeleteMessage)
	mux.HandleFunc("/ws", ch.HandleWebSocket)

	return mux
}
