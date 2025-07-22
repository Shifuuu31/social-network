package handlers

import (
	"fmt"
	"net/http"

	"social-network/pkg/models"
	"social-network/pkg/tools"
)

// NotificationActionPayload represents the payload for notification-based actions
type NotificationActionPayload struct {
	NotificationID int    `json:"notification_id"`
	Action         string `json:"action"` // "accept" or "decline"
	GroupID        int    `json:"group_id,omitempty"`
	UserID         int    `json:"user_id,omitempty"`
}

// AcceptDeclineFromNotification handles accept/decline actions directly from notifications
func (rt *Root) AcceptDeclineFromNotification(w http.ResponseWriter, r *http.Request) {
	var payload NotificationActionPayload
	if err := tools.DecodeJSON(r, &payload); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to decode notification action JSON",
			Metadata: map[string]any{
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"error": err.Error(),
			},
		})
		tools.RespondError(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	requesterID := rt.DL.GetRequesterID(w, r)
	if requesterID <= 0 {
		tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the notification to determine the type and extract relevant data
	notification := &models.Notification{ID: payload.NotificationID}
	if err := rt.DL.Notifications.GetByID(notification); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to get notification details",
			Metadata: map[string]any{
				"notification_id": payload.NotificationID,
				"user_id":         requesterID,
				"error":           err.Error(),
			},
		})
		tools.RespondError(w, "Notification not found", http.StatusNotFound)
		return
	}

	// Verify the notification belongs to the requester
	if notification.UserID != requesterID {
		tools.RespondError(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Handle based on notification type
	switch notification.Type {
	case "group_invite":
		err := rt.handleGroupInviteAction(payload, requesterID)
		if err != nil {
			tools.RespondError(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "follow_request":
		err := rt.handleFollowRequestAction(payload, requesterID)
		if err != nil {
			tools.RespondError(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "group_request":
		// For group join requests, verify requester is group creator
		err := rt.handleGroupJoinRequestAction(payload, requesterID)
		if err != nil {
			tools.RespondError(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "group_join_request":
		// For group join requests, verify requester is group creator
		err := rt.handleGroupJoinRequestAction(payload, requesterID)
		if err != nil {
			tools.RespondError(w, err.Error(), http.StatusInternalServerError)
			return
		}

	default:
		tools.RespondError(w, "Invalid notification type for action", http.StatusBadRequest)
		return
	}

	// Mark notification as read
	if err := rt.DL.Notifications.MarkAsSeen(payload.NotificationID, requesterID); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Failed to mark notification as seen after action",
			Metadata: map[string]any{
				"notification_id": payload.NotificationID,
				"user_id":         requesterID,
			},
		})
	}

	// Send success response
	tools.EncodeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"action":  payload.Action,
		"message": fmt.Sprintf("Successfully %sed", payload.Action),
	})

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Notification action completed successfully",
		Metadata: map[string]any{
			"notification_id": payload.NotificationID,
			"action":          payload.Action,
			"type":            notification.Type,
			"user_id":         requesterID,
		},
	})
}

// handleGroupInviteAction handles accepting/declining group invitations
func (rt *Root) handleGroupInviteAction(payload NotificationActionPayload, userID int) error {
	// Extract group ID from payload or parse from notification message
	if payload.GroupID <= 0 {
		return fmt.Errorf("group ID is required for group invite actions")
	}

	member := &models.GroupMember{
		GroupID:    payload.GroupID,
		UserID:     userID,
		PrevStatus: "invited",
	}

	if payload.Action == "accept" {
		member.Status = "member"
		// Use existing AcceptDeclineGroup logic
		if err := rt.DL.Members.Upsert(member); err != nil {
			return fmt.Errorf("failed to accept group invitation: %w", err)
		}

		// Join group chat if user is connected
		rt.Hub.InitializeGroupChat(payload.GroupID)
		if conn, hasConnection := rt.Hub.Clients[userID]; hasConnection {
			rt.Hub.JoinGroup(userID, payload.GroupID)
			if conn != nil {
				welcomeMsg := &models.Notification{
					Type: "group_joined",
					Message: &models.Message{
						GroupID: payload.GroupID,
						Content: "Your join request was accepted! Welcome to the group!",
					},
				}
				rt.SendNotificationToUser(member.UserID, welcomeMsg)
			}
		}
	} else if payload.Action == "decline" {
		// Delete the invitation record
		if err := rt.DL.Members.Delete(member); err != nil {
			return fmt.Errorf("failed to decline group invitation: %w", err)
		}
	} else {
		return fmt.Errorf("invalid action: %s", payload.Action)
	}

	return nil
}

// handleFollowRequestAction handles accepting/declining follow requests
func (rt *Root) handleFollowRequestAction(payload NotificationActionPayload, userID int) error {
	if payload.UserID <= 0 {
		return fmt.Errorf("user ID is required for follow request actions")
	}

	followRequest := &models.FollowRequest{
		FromUserID: payload.UserID, // The user who sent the follow request
		ToUserID:   userID,         // The current user (notification recipient)
	}

	if payload.Action == "accept" {
		followRequest.Status = "accepted"
	} else if payload.Action == "decline" {
		followRequest.Status = "declined"
	} else {
		return fmt.Errorf("invalid action: %s", payload.Action)
	}

	if err := rt.DL.Follows.UpdateStatus(followRequest); err != nil {
		return fmt.Errorf("failed to %s follow request: %w", payload.Action, err)
	}

	return nil
}

// handleGroupJoinRequestAction handles accepting/declining group join requests
func (rt *Root) handleGroupJoinRequestAction(payload NotificationActionPayload, creatorID int) error {
	if payload.GroupID <= 0 || payload.UserID <= 0 {
		return fmt.Errorf("both group ID and user ID are required for group join request actions")
	}

	// Verify requester is group creator
	if err := rt.DL.Groups.IsUserCreator(payload.GroupID, creatorID); err != nil {
		return fmt.Errorf("unauthorized: not group creator")
	}

	member := &models.GroupMember{
		GroupID:    payload.GroupID,
		UserID:     payload.UserID,
		PrevStatus: "requested",
	}

	if payload.Action == "accept" {
		member.Status = "member"
		if err := rt.DL.Members.Upsert(member); err != nil {
			return fmt.Errorf("failed to accept join request: %w", err)
		}

		// Join group chat if user is connected
		rt.Hub.InitializeGroupChat(payload.GroupID)
		if conn, hasConnection := rt.Hub.Clients[payload.UserID]; hasConnection {
			rt.Hub.JoinGroup(payload.UserID, payload.GroupID)
			if conn != nil {
				welcomeMsg := &models.Notification{
					Type: "group_joined",
					Message: &models.Message{
						GroupID: payload.GroupID,
						Content: "Your join request was accepted! Welcome to the group!",
					},
				}
				rt.SendNotificationToUser(member.UserID, welcomeMsg)
			}
		}
	} else if payload.Action == "decline" {
		// Delete the join request record
		if err := rt.DL.Members.Delete(member); err != nil {
			return fmt.Errorf("failed to decline join request: %w", err)
		}
	} else {
		return fmt.Errorf("invalid action: %s", payload.Action)
	}

	return nil
}
