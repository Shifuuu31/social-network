package handlers

import (
	"log"
	"net/http"
	"sync"

	"social-network/pkg/models"
	"social-network/pkg/tools"

	"github.com/gorilla/websocket"
)

type WSHub struct {
	Upgrader  websocket.Upgrader
	Clients   map[int]*websocket.Conn         // userID -> connection
	Groups    map[int]map[int]*websocket.Conn // groupID -> userID -> connection
	Broadcast chan any                        // broadcast channel for messages
	Lock      sync.RWMutex                    // protects Clients and Groups
}

type WSResponse struct {
	Status  string `json:"status"`  // "success" or "error"
	Message string `json:"message"` // additional text
}

func NewHub() *WSHub {
	return &WSHub{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins
		},
		Clients:   make(map[int]*websocket.Conn),
		Groups:    make(map[int]map[int]*websocket.Conn),
		Broadcast: make(chan any),
	}
}

func (rt *Root) NewWSHandler() *http.ServeMux {
	wsMux := http.NewServeMux()
	wsMux.HandleFunc("/connect", rt.Connect) // route for websocket connection
	return wsMux
}

func (rt *Root) Connect(w http.ResponseWriter, r *http.Request) {
	requesterID := rt.DL.GetRequesterID(w, r) // identify user

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "New websocket connection attempt",
		Metadata: map[string]any{
			"user_id": requesterID,
			"remote":  r.RemoteAddr,
		},
	})

	conn, err := rt.Hub.Upgrader.Upgrade(w, r, nil) // upgrade HTTP to websocket
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
	defer rt.Hub.RemoveClient(requesterID) // clean up on exit

	rt.Hub.AddClient(requesterID, conn) // add to clients map
	rt.DL.Logger.Log(models.LogEntry{
		Level:    "INFO",
		Message:  "Websocket client added",
		Metadata: map[string]any{"user_id": requesterID},
	})
WSLoop:
	for {
		var msg *models.Message
		err := conn.ReadJSON(&msg) // read incoming message JSON
		if err != nil {
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "INFO",
				Message: "Websocket read error or client disconnected",
				Metadata: map[string]any{
					"user_id": requesterID,
					"error":   err.Error(),
				},
			})
			conn.WriteJSON(WSResponse{
				Status:  "error",
				Message: "failed to send msg",
			})
			break WSLoop
		}

		msg.SenderID = requesterID
		// Insert message into DB, log errors if any
		if err := rt.DL.Messages.Insert(msg); err != nil {
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "ERROR",
				Message: "Failed to insert message",
				Metadata: map[string]any{
					"user_id": requesterID,
					"error":   err.Error(),
					"message": msg,
				},
			})
			conn.WriteJSON(WSResponse{
				Status:  "error",
				Message: "failed to send msg",
			})
			break WSLoop

		} else {
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "INFO",
				Message: "Message inserted",
				Metadata: map[string]any{
					"user_id": requesterID,
					"message": msg,
				},
			})
		}

		// Handle message by type
		switch msg.Type {
		case "private":
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "INFO",
				Message: "Sending private message",
				Metadata: map[string]any{
					"from": msg.SenderID,
					"to":   msg.ReceiverID,
				},
			})

			if err := rt.Hub.SendToUser(msg.ReceiverID, msg); err != nil {
				conn.WriteJSON(WSResponse{
					Status:  "error",
					Message: "failed to send msg",
				})
				break WSLoop

			} // send private msg

		case "group":
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "INFO",
				Message: "Broadcasting group message",
				Metadata: map[string]any{
					"group_id": msg.GroupID,
					"from":     msg.SenderID,
				},
			})
			if err := rt.Hub.BroadcastToGroup(msg.GroupID, msg.SenderID, msg); err != nil {
				conn.WriteJSON(WSResponse{
					Status:  "error",
					Message: "failed to send msg",
				})
				break WSLoop

			} // send to group

		case "join_group":
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "INFO",
				Message: "User joining group",
				Metadata: map[string]any{
					"user_id":  requesterID,
					"group_id": msg.GroupID,
				},
			})
			rt.Hub.JoinGroup(requesterID, msg.GroupID) // add user to group

		case "notification_subscribe":
			// Handle notification subscription - user wants to receive real-time notifications
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "INFO",
				Message: "User subscribed to notifications",
				Metadata: map[string]any{
					"user_id": requesterID,
				},
			})
			// Client is already in the Clients map, so they'll receive notifications
			conn.WriteJSON(WSResponse{
				Status:  "success",
				Message: "subscribed to notifications",
			})
			continue WSLoop // don't break, continue listening

		default:

			rt.DL.Logger.Log(models.LogEntry{
				Level:   "Error",
				Message: "invalid msg type",
				Metadata: map[string]any{
					"user_id":  requesterID,
					"group_id": msg.GroupID,
				},
			})
			conn.WriteJSON(WSResponse{
				Status:  "error",
				Message: "failed to send msg",
			})
			break WSLoop

		}
		conn.WriteJSON(WSResponse{
			Status:  "succes",
			Message: "msg sent",
		})
	}
}

// AddClient adds a new client connection
func (hub *WSHub) AddClient(userID int, conn *websocket.Conn) {
	hub.Lock.Lock()
	defer hub.Lock.Unlock()

	hub.Clients[userID] = conn
	log.Printf("Added client userID=%d", userID)
}

// RemoveClient closes and removes client connection
func (hub *WSHub) RemoveClient(userID int) {
	hub.Lock.Lock()
	defer hub.Lock.Unlock()

	if conn, ok := hub.Clients[userID]; ok {
		conn.Close()
		delete(hub.Clients, userID)
		log.Printf("Removed client userID=%d", userID)
	}
}

// JoinGroup adds a client to a group
func (hub *WSHub) JoinGroup(requesterID, groupID int) {
	hub.Lock.Lock()
	defer hub.Lock.Unlock()

	if hub.Groups[groupID] == nil {
		hub.Groups[groupID] = make(map[int]*websocket.Conn)
	}
	if conn, ok := hub.Clients[requesterID]; ok {
		hub.Groups[groupID][requesterID] = conn
		log.Printf("User %d joined group %d", requesterID, groupID)
	}
}

// InitializeGroupChat creates a new group chat room in the hub
func (hub *WSHub) InitializeGroupChat(groupID int) {
	hub.Lock.Lock()
	defer hub.Lock.Unlock()

	if hub.Groups[groupID] == nil {
		hub.Groups[groupID] = make(map[int]*websocket.Conn)
		log.Printf("Initialized group chat for group %d", groupID)
	} else {
		log.Printf("Group chat for group %d already exists", groupID)
	}
}

// SendToUser sends data to a single user
func (hub *WSHub) SendToUser(receiverID int, data any) error {
	hub.Lock.RLock()
	conn, ok := hub.Clients[receiverID]
	hub.Lock.RUnlock()

	if ok {
		err := conn.WriteJSON(data)
		if err != nil {
			log.Printf("Error sending to user %d: %v", receiverID, err)
			return err
		}
	}
	log.Printf("SendToUser: user %d connection not found", receiverID)
	return nil
}

// BroadcastToGroup sends data to all users in a group
func (hub *WSHub) BroadcastToGroup(groupID, senderID int, data any) error {
	hub.Lock.RLock()
	defer hub.Lock.RUnlock()

	conns, ok := hub.Groups[groupID]
	if !ok {
		log.Printf("BroadcastToGroup: group %d not found", groupID)
		return nil
	}

	for userID, conn := range conns {
		if userID != senderID {
			err := conn.WriteJSON(data)
			if err != nil {
				log.Printf("BroadcastToGroup: failed to send to user %d in group %d: %v", userID, groupID, err)
				return err
			}
		}
	}

	return nil
}
