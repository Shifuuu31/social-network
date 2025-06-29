package handlers

import (
	"net/http"

	"social-network/pkg/models"
	"social-network/pkg/tools"
)

func (rt *Root) NewUserHandler() (userMux *http.ServeMux) {
	userMux = http.NewServeMux()

	userMux.HandleFunc("POST /profile/info", rt.ProfileInfo)
	userMux.HandleFunc("POST /profile/activity", rt.ProfileActivity)
	userMux.HandleFunc("POST /profile/connections", rt.ProfileConnections)
	userMux.HandleFunc("POST /profile/visibility", rt.UpdateProfileVisibility)
	userMux.HandleFunc("POST /follow/follow-unfollow", rt.FollowUnfollow)
	userMux.HandleFunc("POST /follow/accept-decline", rt.AcceptDecline)

	return userMux
}

func (rt *Root) canViewProfile(w http.ResponseWriter, r *http.Request, targetUser *models.User) bool {
	requesterID, ok := r.Context().Value(models.UserIDKey).(int)
	if !ok {
		tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Unauthorized profile view attempt",
			Metadata: map[string]interface{}{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
			},
		})
		return false
	}

	// Fetch latest user info
	if err := rt.DL.Users.GetUserByID(targetUser); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to get user info for profile view",
			Metadata: map[string]interface{}{
				"user_id": targetUser.ID,
				"ip":      r.RemoteAddr,
				"path":    r.URL.Path,
				"error":   err.Error(),
			},
		})
		tools.RespondError(w, "User not found", http.StatusNotFound)
		return false
	}

	if targetUser.IsPublic {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "Profile is public, access granted",
			Metadata: map[string]interface{}{
				"user_id":    targetUser.ID,
				"requester":  requesterID,
				"ip":         r.RemoteAddr,
				"path":       r.URL.Path,
			},
		})
		return true
	}

	// Check follow status if private
	followRequest := &models.FollowRequest{FromUserID: requesterID, ToUserID: targetUser.ID}
	if err := rt.DL.Follows.GetFollowStatus(followRequest); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Follow status not found or error",
			Metadata: map[string]interface{}{
				"from_user": requesterID,
				"to_user":   targetUser.ID,
				"ip":        r.RemoteAddr,
				"path":      r.URL.Path,
				"error":     err.Error(),
			},
		})
		tools.RespondError(w, "Private profile — follow to see more", http.StatusForbidden)
		return false
	}
	if followRequest.Status != "accepted" {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "Access denied due to non-accepted follow request",
			Metadata: map[string]interface{}{
				"from_user": requesterID,
				"to_user":   targetUser.ID,
				"status":    followRequest.Status,
				"ip":        r.RemoteAddr,
				"path":      r.URL.Path,
			},
		})
		tools.RespondError(w, "Private profile — follow to see more", http.StatusForbidden)
		return false
	}

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Access granted to private profile via accepted follow",
		Metadata: map[string]interface{}{
			"from_user": requesterID,
			"to_user":   targetUser.ID,
			"ip":        r.RemoteAddr,
			"path":      r.URL.Path,
		},
	})

	return true
}

func (rt *Root) ProfileInfo(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	if err := tools.DecodeJSON(r, &user); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to decode profile info JSON",
			Metadata: map[string]interface{}{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"error":  err.Error(),
			},
		})
		tools.RespondError(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if !rt.canViewProfile(w, r, user) {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Profile info access denied",
			Metadata: map[string]interface{}{
				"user_id": user.ID,
				"ip":      r.RemoteAddr,
				"path":    r.URL.Path,
			},
		})
		tools.RespondError(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := tools.EncodeJSON(w, http.StatusOK, user); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send profile info response",
			Metadata: map[string]interface{}{
				"user_id": user.ID,
				"ip":      r.RemoteAddr,
				"path":    r.URL.Path,
				"error":   err.Error(),
			},
		})
		return
	}

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Profile info sent successfully",
		Metadata: map[string]interface{}{
			"user_id": user.ID,
			"ip":      r.RemoteAddr,
			"path":    r.URL.Path,
		},
	})
}

func (rt *Root) ProfileActivity(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	if err := tools.DecodeJSON(r, &user); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to decode profile activity JSON",
			Metadata: map[string]interface{}{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"error":  err.Error(),
			},
		})
		tools.RespondError(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	
	if !rt.canViewProfile(w, r, user) {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Profile activity access denied",
			Metadata: map[string]interface{}{
				"user_id": user.ID,
				"ip":      r.RemoteAddr,
				"path":    r.URL.Path,
			},
		})
		tools.RespondError(w, "Forbidden", http.StatusForbidden)
		return
	}

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Profile activity request received but not implemented",
		Metadata: map[string]interface{}{
			"user_id": user.ID,
			"ip":      r.RemoteAddr,
			"path":    r.URL.Path,
		},
	})

	tools.RespondError(w, "Not implemented", http.StatusNotImplemented)
}

func (rt *Root) ProfileConnections(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	if err := tools.DecodeJSON(r, &user); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to decode profile connections JSON",
			Metadata: map[string]interface{}{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"error":  err.Error(),
			},
		})
		tools.RespondError(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if !rt.canViewProfile(w, r, user) {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Profile connections access denied",
			Metadata: map[string]interface{}{
				"user_id": user.ID,
				"ip":      r.RemoteAddr,
				"path":    r.URL.Path,
			},
		})
		tools.RespondError(w, "Forbidden", http.StatusForbidden)
		return
	}

	followers, err1 := rt.DL.Follows.GetFollows(user.ID, "followers")
	following, err2 := rt.DL.Follows.GetFollows(user.ID, "following")

	if err1 != nil || err2 != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to fetch profile connections",
			Metadata: map[string]interface{}{
				"user_id": user.ID,
				"ip":      r.RemoteAddr,
				"path":    r.URL.Path,
				"error1":  err1,
				"error2":  err2,
			},
		})
		tools.RespondError(w, "Failed to fetch connections", http.StatusInternalServerError)
		return
	}

	follows := &struct {
		Followers []models.User `json:"followers"`
		Following []models.User `json:"following"`
	}{
		Followers: followers,
		Following: following,
	}

	if err := tools.EncodeJSON(w, http.StatusOK, follows); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send profile connections response",
			Metadata: map[string]interface{}{
				"user_id": user.ID,
				"ip":      r.RemoteAddr,
				"path":    r.URL.Path,
				"error":   err.Error(),
			},
		})
		return
	}

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Profile connections sent successfully",
		Metadata: map[string]interface{}{
			"user_id": user.ID,
			"ip":      r.RemoteAddr,
			"path":    r.URL.Path,
		},
	})
}

func (rt *Root) UpdateProfileVisibility(w http.ResponseWriter, r *http.Request) {
	requesterID, ok := r.Context().Value(models.UserIDKey).(int)
	if !ok {
		tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Unauthorized visibility toggle attempt",
			Metadata: map[string]interface{}{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
			},
		})
		return
	}
	
	var user = &models.User{ID: requesterID}
	if err := rt.DL.Users.GetUserByID(user); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to fetch user for visibility toggle",
			Metadata: map[string]interface{}{
				"user_id": requesterID,
				"ip":      r.RemoteAddr,
				"path":    r.URL.Path,
				"error":   err.Error(),
			},
		})
		tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	user.IsPublic = !user.IsPublic
	
	if err := rt.DL.Users.UpdateUser(user); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to update user visibility",
			Metadata: map[string]interface{}{
				"user_id": user.ID,
				"ip":      r.RemoteAddr,
				"path":    r.URL.Path,
				"error":   err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tools.EncodeJSON(w, http.StatusOK, user); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send profile visibility response",
			Metadata: map[string]interface{}{
				"user_id": user.ID,
				"ip":      r.RemoteAddr,
				"path":    r.URL.Path,
				"error":   err.Error(),
			},
		})
		return
	}

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "User visibility toggled successfully",
		Metadata: map[string]interface{}{
			"user_id": user.ID,
			"is_public": user.IsPublic,
			"ip":      r.RemoteAddr,
			"path":    r.URL.Path,
		},
	})
}

func (rt *Root) FollowUnfollow(w http.ResponseWriter, r *http.Request) {
	payload := &struct {
		TargetId int    `json:"target_id"`
		Action   string `json:"action"`
	}{}

	if err := tools.DecodeJSON(r, &payload); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to decode follow/unfollow JSON",
			Metadata: map[string]interface{}{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"error":  err.Error(),
			},
		})
		tools.RespondError(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	requesterID, ok := r.Context().Value(models.UserIDKey).(int)
	if !ok {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Unauthorized follow/unfollow request",
			Metadata: map[string]interface{}{
				"target_id": payload.TargetId,
				"ip":        r.RemoteAddr,
				"path":      r.URL.Path,
			},
		})
		tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	followRequest := &models.FollowRequest{FromUserID: requesterID, ToUserID: payload.TargetId}

	switch payload.Action {
	case "follow":
		if err := rt.DL.Follows.InsertFollowRequest(followRequest); err != nil {
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "ERROR",
				Message: "Failed to insert follow request",
				Metadata: map[string]interface{}{
					"from": requesterID,
					"to":   payload.TargetId,
					"error": err.Error(),
					"ip":    r.RemoteAddr,
					"path":  r.URL.Path,
				},
			})
			tools.RespondError(w, "Failed to send follow request", http.StatusInternalServerError)
			return
		}
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "Follow request sent",
			Metadata: map[string]interface{}{
				"from": requesterID,
				"to":   payload.TargetId,
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
			},
		})
	case "unfollow":
		if err := rt.DL.Follows.UnfollowUser(followRequest); err != nil {
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "ERROR",
				Message: "Failed to unfollow user",
				Metadata: map[string]interface{}{
					"from": requesterID,
					"to":   payload.TargetId,
					"error": err.Error(),
					"ip":    r.RemoteAddr,
					"path":  r.URL.Path,
				},
			})
			tools.RespondError(w, "Failed to unfollow user", http.StatusInternalServerError)
			return
		}
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "User unfollowed",
			Metadata: map[string]interface{}{
				"from": requesterID,
				"to":   payload.TargetId,
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
			},
		})
	default:
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Invalid follow/unfollow action",
			Metadata: map[string]interface{}{
				"action": payload.Action,
				"from":   requesterID,
				"to":     payload.TargetId,
				"ip":     r.RemoteAddr,
				"path":   r.URL.Path,
			},
		})
		tools.RespondError(w, "Invalid action", http.StatusBadRequest)
		return
	}

	if err := tools.EncodeJSON(w, http.StatusOK, nil); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send follow/unfollow response",
			Metadata: map[string]interface{}{
				"payload": payload,
				"ip":      r.RemoteAddr,
				"path":    r.URL.Path,
				"error":   err.Error(),
			},
		})
	}
}

func (rt *Root) AcceptDecline(w http.ResponseWriter, r *http.Request) {
	payload := &struct {
		TargetId int    `json:"target_id"`
		Action   string `json:"action"`
	}{}

	if err := tools.DecodeJSON(r, &payload); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to decode accept/decline JSON",
			Metadata: map[string]interface{}{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"error":  err.Error(),
			},
		})
		tools.RespondError(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(models.UserIDKey).(int)
	if !ok {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Unauthorized accept/decline attempt",
			Metadata: map[string]interface{}{
				"target_id": payload.TargetId,
				"ip":        r.RemoteAddr,
				"path":      r.URL.Path,
			},
		})
		tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if payload.Action != "accepted" && payload.Action != "declined" {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Invalid accept/decline action",
			Metadata: map[string]interface{}{
				"action": payload.Action,
				"from":   payload.TargetId,
				"to":     userID,
				"ip":     r.RemoteAddr,
				"path":   r.URL.Path,
			},
		})
		tools.RespondError(w, "Invalid action", http.StatusBadRequest)
		return
	}

	followRequest := &models.FollowRequest{
		FromUserID: payload.TargetId,
		ToUserID:   userID,
		Status:     payload.Action,
	}

	if err := rt.DL.Follows.UpdateFollowRequest(followRequest); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to update follow request status",
			Metadata: map[string]interface{}{
				"from":   payload.TargetId,
				"to":     userID,
				"status": payload.Action,
				"error":  err.Error(),
				"ip":     r.RemoteAddr,
				"path":   r.URL.Path,
			},
		})
		tools.RespondError(w, "Failed to update follow request", http.StatusInternalServerError)
		return
	}

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Follow request " + payload.Action,
		Metadata: map[string]interface{}{
			"from":   payload.TargetId,
			"to":     userID,
			"ip":     r.RemoteAddr,
			"path":   r.URL.Path,
			"action": payload.Action,
		},
	})

	if err := tools.EncodeJSON(w, http.StatusOK, nil); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send accept/decline response",
			Metadata: map[string]interface{}{
				"payload": payload,
				"ip":      r.RemoteAddr,
				"path":    r.URL.Path,
				"error":   err.Error(),
			},
		})
	}
}
