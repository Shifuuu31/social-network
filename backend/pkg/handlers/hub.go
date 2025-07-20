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
	Clients   map[int][]*websocket.Conn         // userID -> connection
	Groups    map[int]map[int][]*websocket.Conn // groupID -> userID -> connection
	Broadcast chan any                          // broadcast channel for messages
	Mutex     sync.RWMutex                      // protects Clients and Groups
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
		Clients:   make(map[int][]*websocket.Conn),
		Groups:    make(map[int]map[int][]*websocket.Conn),
		Broadcast: make(chan any),
	}
}

func (rt *Root) NewWSHandler() *http.ServeMux {
	wsMux := http.NewServeMux()
	wsMux.HandleFunc("/connect", rt.Connect) // route for websocket connection
	return wsMux
}

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

	rt.Hub.AddClient(requesterID, conn)

	rt.DL.Logger.Log(models.LogEntry{
		Level:    "INFO",
		Message:  "WebSocket client added",
		Metadata: map[string]any{"user_id": requesterID},
	})

	go rt.Send(requesterID, conn)
}

// AddClient adds a new client connection
func (hub *WSHub) AddClient(userID int, conn *websocket.Conn) {
	hub.Mutex.RLock()

	hub.Clients[userID] = append(hub.Clients[userID], conn)
	log.Printf("Added client userID=%d", userID)

	if len(hub.Clients[userID]) == 1 {
		// TODO send user online
	}
	hub.Mutex.RUnlock()
}

func (hub *WSHub) RemoveClientConn(userID int, conn *websocket.Conn) {
	hub.Mutex.RLock()

	conns := hub.Clients[userID]
	for i, c := range conns {
		if c == conn {
			conn.Close()
			// Remove the connection from the slice
			hub.Clients[userID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}

	// Clean up empty slice
	if len(hub.Clients[userID]) == 0 {
		delete(hub.Clients, userID)
		log.Printf("All clients disconnected for userID=%d", userID)
	}
	hub.Mutex.RUnlock()
}

// JoinGroup adds a client to a group
func (hub *WSHub) JoinGroup(userID, groupID int) {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()

	if hub.Groups[groupID] == nil {
		hub.Groups[groupID] = make(map[int][]*websocket.Conn)
	}
	if conns, ok := hub.Clients[userID]; ok {
		hub.Groups[groupID][userID] = conns
		log.Printf("User %d joined group %d", userID, groupID)
	}
}

func (rt *Root) Send(userID int, conn *websocket.Conn) {
	defer rt.Hub.RemoveClientConn(userID, conn)

	for {
		var msg *models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "INFO",
				Message: "WebSocket read error or client disconnected",
				Metadata: map[string]any{
					"user_id": userID,
					"error":   err.Error(),
				},
			})
			conn.WriteJSON(WSResponse{
				Status:  "error",
				Message: "failed to send msg",
			})
			return
		}

		msg.SenderID = userID

		if err := rt.DL.Messages.Insert(msg); err != nil {
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "ERROR",
				Message: "Failed to insert message",
				Metadata: map[string]any{
					"user_id": userID,
					"error":   err.Error(),
					"message": msg,
				},
			})
			conn.WriteJSON(WSResponse{
				Status:  "error",
				Message: "failed to send msg",
			})
			return
		}

		rt.DL.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "Message inserted",
			Metadata: map[string]any{
				"user_id": userID,
				"message": msg,
			},
		})

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
				return
			}

		case "group":
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "INFO",
				Message: "Broadcasting group message",
				Metadata: map[string]any{
					"group_id": msg.GroupID,
					"from":     msg.SenderID,
				},
			})
			if err := rt.Hub.BroadcastToGroupMembers(msg.GroupID, msg.SenderID, msg); err != nil {
				conn.WriteJSON(WSResponse{
					Status:  "error",
					Message: "failed to send msg",
				})
				return
			}

		case "join_group":
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "INFO",
				Message: "User joining group",
				Metadata: map[string]any{
					"user_id":  userID,
					"group_id": msg.GroupID,
				},
			})
			rt.Hub.JoinGroup(userID, msg.GroupID)

		default:
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "ERROR",
				Message: "Invalid message type",
				Metadata: map[string]any{
					"user_id":  userID,
					"group_id": msg.GroupID,
				},
			})
			conn.WriteJSON(WSResponse{
				Status:  "error",
				Message: "invalid message type",
			})
			return
		}

		conn.WriteJSON(WSResponse{
			Status:  "success",
			Message: "msg sent",
		})
	}
}

// BroadcastToAll sends the given data to all WebSocket connections in the slice.
// It removes any closed connections from the slice if needed (you can customize that).
func (hub *WSHub) SendToUser(userID int, data any) error {
	hub.Mutex.RLock()
	conns := hub.Clients[userID]
	hub.Mutex.RUnlock()

	for _, conn := range conns {
		if err := conn.WriteJSON(data); err != nil {
			log.Printf("SendToUser: Failed to send to user %d: %v", userID, err)
			// Optionally close and remove the conn here
		}
	}
	return nil
}

func (rt *Root) NotifyUser(userID int, notif *models.Notification) {
	if err := rt.DL.Notifications.Upsert(notif); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to save notification",
			Metadata: map[string]any{
				"user_id": userID,
				"error":   err.Error(),
			},
		})
		return
	}

	rt.Hub.Mutex.RLock()
	conns := rt.Hub.Clients[userID]
	rt.Hub.Mutex.RUnlock()

	msg := models.Message{
		Type:  "notify",
		Notif: notif,
	}

	for _, conn := range conns {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("SendNotificationToUser error: userID=%d, err=%v", userID, err)
			conn.Close()
		}
	}
}

func (hub *WSHub) BroadcastToGroupMembers(groupID, senderID int, data any) error {
	hub.Mutex.RLock()
	defer hub.Mutex.RUnlock()

	groupConns, ok := hub.Groups[groupID]
	if !ok {
		log.Printf("BroadcastToGroup: group %d not found", groupID)
		return nil
	}

	for userID, conns := range groupConns {
		if userID == senderID {
			continue
		}
		for _, conn := range conns {
			if err := conn.WriteJSON(data); err != nil {
				log.Printf("BroadcastToGroup: failed to send to user %d in group %d: %v", userID, groupID, err)
			}
		}
	}

	return nil
}

func (rt *Root) NotifyGroupMembers(groupID int, notif *models.Notification) {
	if err := rt.DL.Notifications.Upsert(notif); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to save notification",
			Metadata: map[string]any{
				"user_id": groupID,
				"error":   err.Error(),
			},
		})
		return
	}

	rt.Hub.Mutex.RLock()
	groupUsers := rt.Hub.Groups[groupID]
	rt.Hub.Mutex.RUnlock()

	msg := models.Message{
		Type:  "notify",
		Notif: notif,
	}

	for userID, conns := range groupUsers {
		for _, conn := range conns {
			if err := conn.WriteJSON(msg); err != nil {
				log.Printf("SendNotificationToGroup error: userID=%d, err=%v", userID, err)
				conn.Close()
			}
		}
	}
}

func removeDeadConns(conns []*websocket.Conn) []*websocket.Conn {
	var alive []*websocket.Conn
	for _, c := range conns {
		if c != nil {
			if err := c.WriteMessage(websocket.PingMessage, nil); err == nil {
				alive = append(alive, c)
			} else {
				c.Close()
			}
		}
	}
	return alive
}
