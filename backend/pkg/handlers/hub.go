package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"social-network/pkg/models"
	"social-network/pkg/tools"

	"github.com/gorilla/websocket"
)

// **WEBSOCKET HUB - MANAGES ALL CLIENT CONNECTIONS AND GROUPS**
type WSHub struct {
	Upgrader websocket.Upgrader
	Clients  map[int][]*websocket.Conn         // userID -> connection
	Groups   map[int]map[int][]*websocket.Conn // groupID -> userID -> connection
	Mutex    sync.RWMutex                      // protects Clients and Groups
}

// **WEBSOCKET RESPONSE STRUCTURE - USED FOR SUCCESS/ERROR RESPONSES**
type WSResponse struct {
	Status  string `json:"status"`  // "success" or "error"
	Message string `json:"message"` // additional text
}

// **CREATES NEW WEBSOCKET HUB INSTANCE**
func NewHub() *WSHub {
	return &WSHub{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins
		},
		Clients: make(map[int][]*websocket.Conn),
		Groups:  make(map[int]map[int][]*websocket.Conn),
	}
}

// **SETS UP WEBSOCKET ROUTES**
func (rt *Root) NewWSHandler() *http.ServeMux {
	wsMux := http.NewServeMux()
	wsMux.HandleFunc("/connect", rt.Connect) // route for websocket connection
	return wsMux
}

// **HANDLES INITIAL WEBSOCKET CONNECTION UPGRADE**
func (rt *Root) Connect(w http.ResponseWriter, r *http.Request) {
	requesterID := rt.DL.GetRequesterID(w, r)
	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "New websocket connection attempt",
		Metadata: map[string]any{
			"user_id": requesterID,
			"remote":  r.RemoteAddr,
		},
	})

	// **UPGRADE HTTP CONNECTION TO WEBSOCKET**
	conn, err := rt.Hub.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to upgrade websocket connection",
			Metadata: map[string]any{
				"user_id": requesterID,
				"error":   err.Error(),
			},
		})
		tools.RespondError(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	// **ADD CLIENT TO HUB AND START HANDLING**
	rt.Hub.AddClient(requesterID, conn)
	rt.DL.Logger.Log(models.LogEntry{
		Level:    "INFO",
		Message:  "WebSocket client added",
		Metadata: map[string]any{"user_id": requesterID},
	})

	go rt.HandleClient(requesterID, conn)
}

// **MAIN CLIENT HANDLER - PROCESSES INCOMING NOTIFICATIONS**
func (rt *Root) HandleClient(userID int, conn *websocket.Conn) {
	defer rt.Hub.RemoveClientConn(userID, conn)

	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				conn.Close()
				return
			}
		}
	}()

	for {
		// **READ NOTIFICATION FROM CLIENT (NOTIFICATION-FIRST ARCHITECTURE)**
		var notif *models.Notification
		err := conn.ReadJSON(&notif)
		if err != nil {
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "INFO",
				Message: "WebSocket read error or client disconnected",
				Metadata: map[string]any{
					"user_id": userID,
					"error":   err.Error(),
				},
			})
			// **SEND ERROR RESPONSE ON READ FAILURE**
			WriteWSError(conn, "failed to read notification")
			return
		}

		// **ROUTE NOTIFICATION BY TYPE**
		switch notif.Type {
		case "message":
			if err := rt.handleMessageNotification(userID, notif); err != nil {
				// **ERROR HANDLING - SEND WSResponse ON FAILURE**
				WriteWSError(conn, "failed to process message")
				return
			}
		case "join_group":
			if err := rt.handleJoinGroupNotification(userID, notif); err != nil {
				// **ERROR HANDLING - SEND WSResponse ON FAILURE**

				WriteWSError(conn, "failed to join group")

				return
			}
		default:
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "ERROR",
				Message: "Invalid notification type",
				Metadata: map[string]any{
					"user_id": userID,
					"type":    notif.Type,
				},
			})
			// **ERROR HANDLING - INVALID TYPE**
			WriteWSError(conn, "invalid notification type")
			continue
		}

		// **SUCCESS RESPONSE - NOTIFICATION PROCESSED**
		WriteWSSuccess(conn, "notification processed")
	}
}

// **ADDS NEW CLIENT CONNECTION TO HUB**
func (hub *WSHub) AddClient(userID int, conn *websocket.Conn) {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	hub.Clients[userID] = append(hub.Clients[userID], conn)
	log.Printf("Added client userID=%d", userID)

	// **HANDLE FIRST CONNECTION - USER COMES ONLINE**
	if len(hub.Clients[userID]) == 1 {
		// TODO send user online notification
	}
}

// **REMOVES CLIENT CONNECTION FROM HUB**
func (hub *WSHub) RemoveClientConn(userID int, conn *websocket.Conn) {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	conns := hub.Clients[userID]
	for i, c := range conns {
		if c == conn {
			conn.Close()
			// **REMOVE CONNECTION FROM SLICE**
			hub.Clients[userID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}

	// **CLEAN UP EMPTY SLICE - USER GOES OFFLINE**
	if len(hub.Clients[userID]) == 0 {
		delete(hub.Clients, userID)
		log.Printf("All clients disconnected for userID=%d", userID)
	}
}

// **ADDS USER TO GROUP FOR GROUP MESSAGING**
func (hub *WSHub) JoinGroup(userID, groupID int) {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	// **INITIALIZE GROUP IF NOT EXISTS**
	if hub.Groups[groupID] == nil {
		hub.Groups[groupID] = make(map[int][]*websocket.Conn)
	}
	// **ADD USER'S CONNECTIONS TO GROUP**
	if conns, ok := hub.Clients[userID]; ok {
		hub.Groups[groupID][userID] = conns
		log.Printf("User %d joined group %d", userID, groupID)
	}
}

// **HANDLES MESSAGE-TYPE NOTIFICATIONS - PROCESSES AND STORES MESSAGES**
func (rt *Root) handleMessageNotification(userID int, notif *models.Notification) error {
	// **VALIDATE MESSAGE EXISTS IN NOTIFICATION**
	if notif.Message == nil {
		return fmt.Errorf("message is required for message notification")
	}

	// **SET MESSAGE METADATA**
	notif.Message.SenderID = userID
	notif.Message.CreatedAt = time.Now()

	// **INSERT MESSAGE INTO DATABASE FIRST**
	if err := rt.DL.Messages.Insert(notif.Message); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to insert message",
			Metadata: map[string]any{
				"user_id": userID,
				"error":   err.Error(),
				"message": notif.Message,
			},
		})
		return err
	}

	// **UPDATE NOTIFICATION WITH MESSAGE ID AND TIMESTAMP**
	notif.MessageID = notif.Message.ID
	notif.CreatedAt = int(time.Now().Unix())

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Message inserted",
		Metadata: map[string]any{
			"user_id":    userID,
			"message_id": notif.Message.ID,
		},
	})

	// **ROUTE NOTIFICATION BASED ON MESSAGE TYPE**
	switch notif.Message.Type {
	case "private":
		return rt.sendPrivateNotification(notif)
	case "group":
		return rt.sendGroupNotification(notif)
	default:
		return fmt.Errorf("invalid message type: %s", notif.Message.Type)
	}
}

// **HANDLES JOIN GROUP NOTIFICATIONS**
func (rt *Root) handleJoinGroupNotification(userID int, notif *models.Notification) error {
	// **VALIDATE GROUP ID EXISTS**
	if notif.Message == nil || notif.Message.GroupID == 0 {
		return fmt.Errorf("group_id is required for join_group notification")
	}

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "User joining group",
		Metadata: map[string]any{
			"user_id":  userID,
			"group_id": notif.Message.GroupID,
		},
	})

	// **ADD USER TO GROUP IN HUB**
	rt.Hub.JoinGroup(userID, notif.Message.GroupID)
	return nil
}

// **SENDS PRIVATE MESSAGE NOTIFICATION TO SPECIFIC USER**
func (rt *Root) sendPrivateNotification(notif *models.Notification) error {
	// **VALIDATE RECEIVER ID EXISTS**
	if notif.Message == nil || notif.Message.ReceiverID == 0 {
		return fmt.Errorf("receiver_id is required for private message")
	}

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Sending private notification",
		Metadata: map[string]any{
			"from": notif.Message.SenderID,
			"to":   notif.Message.ReceiverID,
		},
	})

	// **CREATE NOTIFICATION FOR RECEIVER**
	receiverNotif := &models.Notification{
		UserID:     notif.Message.ReceiverID,
		Type:       "message_received",
		SubMessage: "You have a new message",
		MessageID:  notif.Message.ID,
		Message:    notif.Message,
		Seen:       false,
		CreatedAt:  int(time.Now().Unix()),
	}

	return rt.SendNotificationToUser(notif.Message.ReceiverID, receiverNotif)
}

// **SENDS GROUP MESSAGE NOTIFICATION TO ALL GROUP MEMBERS**
func (rt *Root) sendGroupNotification(notif *models.Notification) error {
	// **VALIDATE GROUP ID EXISTS**
	if notif.Message == nil || notif.Message.GroupID == 0 {
		return fmt.Errorf("group_id is required for group message")
	}

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Broadcasting group notification",
		Metadata: map[string]any{
			"group_id": notif.Message.GroupID,
			"from":     notif.Message.SenderID,
		},
	})

	// **CREATE NOTIFICATION FOR GROUP MEMBERS**
	groupNotif := &models.Notification{
		Type:       "group_message_received",
		SubMessage: "New message in group",
		MessageID:  notif.Message.ID,
		Message:    notif.Message,
		Seen:       false,
		CreatedAt:  int(time.Now().Unix()),
	}

	return rt.BroadcastNotificationToGroup(notif.Message.GroupID, notif.Message.SenderID, groupNotif)
}

// SendNotificationToUser sends a notification to a specific user
func (rt *Root) SendNotificationToUser(userID int, notif *models.Notification) error {
	// Set the notification's user ID
	notif.UserID = userID

	// Save notification to database
	if err := rt.DL.Notifications.Upsert(notif); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to save notification",
			Metadata: map[string]any{
				"user_id": userID,
				"error":   err.Error(),
			},
		})
		return err
	}

	// Send to WebSocket connections
	rt.Hub.Mutex.Lock()
	conns := removeDeadConns(rt.Hub.Clients[userID])
	rt.Hub.Clients[userID] = conns
	rt.Hub.Mutex.Unlock()

	for _, conn := range conns {
		if err := conn.WriteJSON(notif); err != nil {
			log.Printf("SendNotificationToUser: Failed to send to user %d: %v", userID, err)
			// Optionally close and remove the conn here
		}
	}

	return nil
}

// BroadcastNotificationToGroup sends notification to all group members except sender
func (rt *Root) BroadcastNotificationToGroup(groupID, senderID int, notif *models.Notification) error {
	rt.Hub.Mutex.RLock()
	groupConns, ok := rt.Hub.Groups[groupID]
	rt.Hub.Mutex.RUnlock()

	if !ok {
		log.Printf("BroadcastNotificationToGroup: group %d not found", groupID)
		return nil
	}

	for userID, conns := range groupConns {
		conns := removeDeadConns(conns)
		rt.Hub.Mutex.Lock()
		rt.Hub.Groups[groupID][userID] = conns
		rt.Hub.Mutex.Unlock()

		if userID == senderID {
			continue // Skip sender
		}

		// Set the notification's user ID for each recipient
		userNotif := *notif // Copy the notification
		userNotif.UserID = userID

		// Save notification to database for each user
		if err := rt.DL.Notifications.Upsert(&userNotif); err != nil {
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "ERROR",
				Message: "Failed to save group notification",
				Metadata: map[string]any{
					"user_id":  userID,
					"group_id": groupID,
					"error":    err.Error(),
				},
			})
			continue
		}

		// Send to WebSocket connections
		for _, conn := range conns {
			if err := conn.WriteJSON(&userNotif); err != nil {
				log.Printf("BroadcastNotificationToGroup: failed to send to user %d in group %d: %v", userID, groupID, err)
			}
		}
	}

	return nil
}

func removeDeadConns(conns []*websocket.Conn) []*websocket.Conn {
	var alive []*websocket.Conn
	for _, c := range conns {
		c.SetWriteDeadline(time.Now().Add(1 * time.Second))
		if err := c.WriteControl(websocket.PingMessage, nil, time.Now().Add(1*time.Second)); err == nil {
			alive = append(alive, c)
		} else {
			c.Close()
		}
	}
	return alive
}

func WriteWSSuccess(conn *websocket.Conn, msg string) {
	conn.WriteJSON(WSResponse{Status: "success", Message: msg})
}

func WriteWSError(conn *websocket.Conn, msg string) {
	conn.WriteJSON(WSResponse{Status: "error", Message: msg})
}
