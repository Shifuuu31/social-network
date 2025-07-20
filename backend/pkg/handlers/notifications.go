package handlers

import (
	"net/http"
	"strconv"
	"time"

	"social-network/pkg/middleware"
	"social-network/pkg/models"
	"social-network/pkg/tools"
)

// NotificationService handles all notification operations
type NotificationService struct {
	DL  *middleware.DataLayer
	Hub *WSHub
}

// NotificationData represents additional data stored with notifications
type NotificationData struct {
	GroupID   int    `json:"group_id,omitempty"`
	EventID   int    `json:"event_id,omitempty"`
	SenderID  int    `json:"sender_id,omitempty"`
	ActionURL string `json:"action_url,omitempty"`
}

func (rt *Root) NewNotificationHandler() *http.ServeMux {
	notifMux := http.NewServeMux()

	// Notification routes
	notifMux.HandleFunc("/notifications", rt.GetNotifications)
	notifMux.HandleFunc("/notifications/mark-read", rt.MarkNotificationsRead)
	// notifMux.HandleFunc("POST /notifications/{id}/mark-read", rt.MarkNotificationRead)
	notifMux.HandleFunc("/notifications/{id}", rt.DeleteNotification)
	notifMux.HandleFunc("/notifications/unread-count", rt.GetUnreadCount)

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Notification routes registered",
		Metadata: map[string]any{
			"routes": "/notifications, /notifications/mark-read, /notifications/{id}/mark-read, /notifications/{id}, /notifications/unread-count",
		},
	})

	return notifMux
}

// CreateNotification creates a new notification and sends it via WebSocket if user is online
func (ns *NotificationService) CreateNotification(userID int, notifType, title, message string, data *NotificationData) error {
	notification := &models.Notification{
		UserID:  userID,
		Type:    notifType,
		Message: message,
		Seen:    false,
	}

	// Insert notification into database
	if err := ns.DL.Notifications.Insert(notification); err != nil {
		ns.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to create notification",
			Metadata: map[string]any{
				"user_id": userID,
				"type":    notifType,
				"error":   err.Error(),
			},
		})
		return err
	}

	// Send real-time notification via WebSocket if user is online
	ns.SendRealTimeNotification(userID, notification)

	ns.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Notification created successfully",
		Metadata: map[string]any{
			"user_id":         userID,
			"type":            notifType,
			"notification_id": notification.ID,
		},
	})

	return nil
}

// SendRealTimeNotification sends notification via WebSocket
func (ns *NotificationService) SendRealTimeNotification(userID int, notification *models.Notification) {
	if conns, ok := ns.Hub.Clients[userID]; ok {
		wsMessage := map[string]any{
			"type":         "notification",
			"notification": notification,
			"timestamp":    time.Now(),
			"unread_count": ns.GetUserUnreadCount(userID),
		}

		for _, conn := range conns {
			if err := conn.WriteJSON(wsMessage); err != nil {
				ns.DL.Logger.Log(models.LogEntry{
					Level:   "WARN",
					Message: "Failed to send real-time notification",
					Metadata: map[string]any{
						"user_id":         userID,
						"notification_id": notification.ID,
						"error":           err.Error(),
					},
				})
			}
		}
	}
}

// GetUserUnreadCount returns the number of unread notifications for a user
func (ns *NotificationService) GetUserUnreadCount(userID int) int {
	count, err := ns.DL.Notifications.GetUnreadCount(userID)
	if err != nil {
		ns.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to get unread notification count",
			Metadata: map[string]any{
				"user_id": userID,
				"error":   err.Error(),
			},
		})
		return 0
	}
	return count
}

// GetNotifications retrieves notifications for the authenticated user
func (rt *Root) GetNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		tools.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	requesterID := rt.DL.GetRequesterID(w, r)
	if requesterID <= 0 {
		tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse query parameters
	page := 1
	limit := 20
	unseenOnly := false

	if p := r.URL.Query().Get("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}
	if r.URL.Query().Get("unseen_only") == "true" {
		unseenOnly = true
	}

	offset := (page - 1) * limit

	var notifications []*models.Notification
	var err error

	if unseenOnly {
		notifications, err = rt.DL.Notifications.GetUserUnseenNotifications(requesterID, limit, offset)
	} else {
		notifications, err = rt.DL.Notifications.GetUserNotifications(requesterID, limit, offset)
	}
	if err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to retrieve notifications",
			Metadata: map[string]any{
				"user_id": requesterID,
				"error":   err.Error(),
			},
		})
		tools.RespondError(w, "Failed to retrieve notifications", http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"notifications": notifications,
		"page":          page,
		"limit":         limit,
		"total":         len(notifications),
	}

	if err := tools.EncodeJSON(w, http.StatusOK, response); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send notifications response",
			Metadata: map[string]any{
				"user_id": requesterID,
				"error":   err.Error(),
			},
		})
	}
}

// MarkNotificationsRead marks multiple notifications as read
func (rt *Root) MarkNotificationsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tools.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	requesterID := rt.DL.GetRequesterID(w, r)
	if requesterID <= 0 {
		tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload struct {
		NotificationIDs []int `json:"notification_ids"`
		MarkAll         bool  `json:"mark_all"`
	}

	if err := tools.DecodeJSON(r, &payload); err != nil {
		tools.RespondError(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	var err error
	if payload.MarkAll {
		err = rt.DL.Notifications.MarkAllAsRead(requesterID)
	} else {
		err = rt.DL.Notifications.MarkAsRead(requesterID, payload.NotificationIDs)
	}

	if err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to mark notifications as read",
			Metadata: map[string]any{
				"user_id": requesterID,
				"error":   err.Error(),
			},
		})
		tools.RespondError(w, "Failed to mark notifications as read", http.StatusInternalServerError)
		return
	}

	// Send updated unread count via WebSocket
	ns := &NotificationService{DL: rt.DL, Hub: rt.Hub}
	if conns, ok := rt.Hub.Clients[requesterID]; ok {
		wsMessage := map[string]any{
			"type":         "notification_count_updated",
			"unread_count": ns.GetUserUnreadCount(requesterID),
			"timestamp":    time.Now(),
		}

		for _, conn := range conns {
			conn.WriteJSON(wsMessage)
		}
	}

	tools.EncodeJSON(w, http.StatusOK, map[string]string{"message": "Notifications marked as read"})
}

// MarkNotificationRead marks a single notification as read
// func (rt *Root) MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
// 	requesterID := rt.DL.GetRequesterID(w, r)
// 	if requesterID <= 0 {
// 		tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	notificationID, err := strconv.Atoi(r.PathValue("id"))
// 	if err != nil || notificationID <= 0 {
// 		tools.RespondError(w, "Invalid notification ID", http.StatusBadRequest)
// 		return
// 	}

// 	if err := rt.DL.Notifications.MarkAsRead(requesterID, []int{notificationID}); err != nil {
// 		rt.DL.Logger.Log(models.LogEntry{
// 			Level:   "ERROR",
// 			Message: "Failed to mark notification as read",
// 			Metadata: map[string]any{
// 				"user_id":         requesterID,
// 				"notification_id": notificationID,
// 				"error":           err.Error(),
// 			},
// 		})
// 		tools.RespondError(w, "Failed to mark notification as read", http.StatusInternalServerError)
// 		return
// 	}

// 	tools.EncodeJSON(w, http.StatusOK, map[string]string{"message": "Notification marked as read"})
// }

// DeleteNotification deletes a notification
func (rt *Root) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	requesterID := rt.DL.GetRequesterID(w, r)
	if requesterID <= 0 {
		tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	notificationID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || notificationID <= 0 {
		tools.RespondError(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}
	if err := rt.DL.Notifications.Delete(notificationID); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to delete notification",
			Metadata: map[string]any{
				"user_id":         requesterID,
				"notification_id": notificationID,
				"error":           err.Error(),
			},
		})
		tools.RespondError(w, "Failed to delete notification", http.StatusInternalServerError)
		return
	}

	tools.EncodeJSON(w, http.StatusOK, map[string]string{"message": "Notification deleted"})
}

// GetUnreadCount returns the number of unread notifications
func (rt *Root) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		tools.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	requesterID := rt.DL.GetRequesterID(w, r)
	if requesterID <= 0 {
		tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ns := &NotificationService{DL: rt.DL, Hub: rt.Hub}
	count := ns.GetUserUnreadCount(requesterID)

	response := map[string]int{"unread_count": count}
	tools.EncodeJSON(w, http.StatusOK, response)
}
