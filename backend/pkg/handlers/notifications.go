package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"social-network/pkg/models"
	"social-network/pkg/tools"
)

// NotificationResponse wraps notification data for API responses
type NotificationResponse struct {
	Notifications []*models.Notification `json:"notifications"`
	UnreadCount   int                    `json:"unread_count"`
	Total         int                    `json:"total"`
}

// WebSocket notification structure for real-time updates
type WSNotification struct {
	Type         string                `json:"type"`         // "notification"
	Action       string                `json:"action"`       // "new", "update", "delete"
	Notification *models.Notification  `json:"notification"`
	UnreadCount  int                   `json:"unread_count"`
}

func (rt *Root) NewNotificationsHandler() (notificationsMux *http.ServeMux) {
	notificationsMux = http.NewServeMux()

	notificationsMux.HandleFunc("POST /fetch", rt.GetNotifications)
	notificationsMux.HandleFunc("POST /mark-seen", rt.MarkNotificationSeen)
	notificationsMux.HandleFunc("POST /mark-all-seen", rt.MarkAllNotificationsSeen)
	notificationsMux.HandleFunc("GET /unread-count", rt.GetUnreadCount)
	notificationsMux.HandleFunc("DELETE /{id}", rt.DeleteNotification)

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Notification routes registered",
		Metadata: map[string]any{
			"routes": "/fetch, /mark-seen, /mark-all-seen, /unread-count, /{id}",
		},
	})

	return notificationsMux
}

// GetNotifications retrieves notifications for the authenticated user
func (rt *Root) GetNotifications(w http.ResponseWriter, r *http.Request) {
	requesterID := rt.DL.GetRequesterID(w, r)
	if requesterID <= 0 {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Unauthorized: requester ID not found",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
			},
		})
		tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload *models.NotificationPayload
	if err := tools.DecodeJSON(r, &payload); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to decode notification payload JSON",
			Metadata: map[string]any{
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"error": err.Error(),
			},
		})
		tools.RespondError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Set user ID from authenticated user
	payload.UserID = requesterID

	// Set defaults if not provided
	if payload.NumOfItems <= 0 {
		payload.NumOfItems = 20
	}
	if payload.Start < 0 {
		payload.Start = 0
	}
	if payload.Type == "" {
		payload.Type = "all"
	}

	notifications, err := rt.DL.Notifications.GetByUser(payload)
	if err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to fetch notifications",
			Metadata: map[string]any{
				"user_id": requesterID,
				"error":   err.Error(),
			},
		})
		tools.RespondError(w, "Failed to fetch notifications", http.StatusInternalServerError)
		return
	}

	// Get unread count
	unreadCount, err := rt.DL.Notifications.CountUnseen(requesterID)
	if err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Failed to get unread count",
			Metadata: map[string]any{
				"user_id": requesterID,
				"error":   err.Error(),
			},
		})
		unreadCount = 0
	}

	response := NotificationResponse{
		Notifications: notifications,
		UnreadCount:   unreadCount,
		Total:         len(notifications),
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

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Notifications fetched successfully",
		Metadata: map[string]any{
			"user_id": requesterID,
			"count":   len(notifications),
			"unread":  unreadCount,
		},
	})
}

// MarkNotificationSeen marks a specific notification as seen
func (rt *Root) MarkNotificationSeen(w http.ResponseWriter, r *http.Request) {
	requesterID := rt.DL.GetRequesterID(w, r)
	if requesterID <= 0 {
		tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload struct {
		NotificationID int `json:"notification_id"`
	}
	
	if err := tools.DecodeJSON(r, &payload); err != nil {
		tools.RespondError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if err := rt.DL.Notifications.MarkAsSeen(payload.NotificationID, requesterID); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to mark notification as seen",
			Metadata: map[string]any{
				"user_id":         requesterID,
				"notification_id": payload.NotificationID,
				"error":           err.Error(),
			},
		})
		tools.RespondError(w, "Failed to mark notification as seen", http.StatusInternalServerError)
		return
	}

	// Get updated unread count
	unreadCount, _ := rt.DL.Notifications.CountUnseen(requesterID)

	// Send real-time update via WebSocket
	rt.SendNotificationUpdate(requesterID, "update", nil, unreadCount)

	tools.EncodeJSON(w, http.StatusOK, map[string]any{
		"success":      true,
		"unread_count": unreadCount,
	})
}

// MarkAllNotificationsSeen marks all notifications as seen for the user
func (rt *Root) MarkAllNotificationsSeen(w http.ResponseWriter, r *http.Request) {
	requesterID := rt.DL.GetRequesterID(w, r)
	if requesterID <= 0 {
		tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := rt.DL.Notifications.MarkAllAsSeen(requesterID); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to mark all notifications as seen",
			Metadata: map[string]any{
				"user_id": requesterID,
				"error":   err.Error(),
			},
		})
		tools.RespondError(w, "Failed to mark notifications as seen", http.StatusInternalServerError)
		return
	}

	// Send real-time update via WebSocket
	rt.SendNotificationUpdate(requesterID, "update", nil, 0)

	tools.EncodeJSON(w, http.StatusOK, map[string]any{
		"success":      true,
		"unread_count": 0,
	})

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "All notifications marked as seen",
		Metadata: map[string]any{
			"user_id": requesterID,
		},
	})
}

// GetUnreadCount returns the count of unread notifications
func (rt *Root) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	requesterID := rt.DL.GetRequesterID(w, r)
	if requesterID <= 0 {
		tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	count, err := rt.DL.Notifications.CountUnseen(requesterID)
	if err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to get unread count",
			Metadata: map[string]any{
				"user_id": requesterID,
				"error":   err.Error(),
			},
		})
		tools.RespondError(w, "Failed to get unread count", http.StatusInternalServerError)
		return
	}

	tools.EncodeJSON(w, http.StatusOK, map[string]any{
		"unread_count": count,
	})
}

// DeleteNotification deletes a specific notification
func (rt *Root) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	requesterID := rt.DL.GetRequesterID(w, r)
	if requesterID <= 0 {
		tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	notificationIDStr := r.PathValue("id")
	notificationID, err := strconv.Atoi(notificationIDStr)
	if err != nil {
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

	// Get updated unread count
	unreadCount, _ := rt.DL.Notifications.CountUnseen(requesterID)

	// Send real-time update via WebSocket
	rt.SendNotificationUpdate(requesterID, "delete", nil, unreadCount)

	tools.EncodeJSON(w, http.StatusOK, map[string]any{
		"success":      true,
		"unread_count": unreadCount,
	})

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Notification deleted successfully",
		Metadata: map[string]any{
			"user_id":         requesterID,
			"notification_id": notificationID,
		},
	})
}

// SendNotificationUpdate sends real-time notification updates via WebSocket
func (rt *Root) SendNotificationUpdate(userID int, action string, notification *models.Notification, unreadCount int) {
	wsNotification := WSNotification{
		Type:         "notification",
		Action:       action,
		Notification: notification,
		UnreadCount:  unreadCount,
	}

	if err := rt.Hub.SendToUser(userID, wsNotification); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Failed to send notification update via WebSocket",
			Metadata: map[string]any{
				"user_id": userID,
				"action":  action,
				"error":   err.Error(),
			},
		})
	}
}

// CreateAndSendNotification creates a notification and sends it via WebSocket
func (rt *Root) CreateAndSendNotification(notification *models.Notification) error {
	// Insert notification into database
	if err := rt.DL.Notifications.Insert(notification); err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	// Get updated unread count
	unreadCount, err := rt.DL.Notifications.CountUnseen(notification.UserID)
	if err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Failed to get unread count for new notification",
			Metadata: map[string]any{
				"user_id": notification.UserID,
				"error":   err.Error(),
			},
		})
		unreadCount = 1 // fallback
	}

	// Send real-time notification via WebSocket
	rt.SendNotificationUpdate(notification.UserID, "new", notification, unreadCount)

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Notification created and sent",
		Metadata: map[string]any{
			"user_id":      notification.UserID,
			"type":         notification.Type,
			"unread_count": unreadCount,
		},
	})

	return nil
}
