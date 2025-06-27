package handlers

import (
	"net/http"
	"time"

	"social-network/pkg/models"
	"social-network/pkg/tools"
)

func (rt *Root) NewUsersHandler() (usersMux *http.ServeMux) {
	usersMux = http.NewServeMux()

	usersMux.HandleFunc("POST /profile/info", rt.ProfileInfo)
	usersMux.HandleFunc("POST /profile/activity", rt.ProfileActivity)
	usersMux.HandleFunc("POST /profile/followers", rt.ProfileFollowers)
	usersMux.HandleFunc("POST /profile/following", rt.ProfileFollowing)
	usersMux.HandleFunc("POST /profile/visibility", rt.UpdateProfileVisibility)
	usersMux.HandleFunc("POST /follow/follow-unfollow", rt.FollowUnfollow)
	usersMux.HandleFunc("POST /follow/accept-decline", rt.AcceptDeclineFollowRequest)

	return usersMux
}

func (rt *Root) ProfileAccess(w http.ResponseWriter, r *http.Request, targetUser *models.User) bool {
	requesterID := rt.DL.GetRequesterID(w, r)

	// Fetch latest user info
	if err := rt.DL.Users.GetUserByID(targetUser); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to get user info for profile view",
			Metadata: map[string]any{
				"user_id": targetUser.ID,
				"ip":      r.RemoteAddr,
				"path":    r.URL.Path,
				"error":   err.Error(),
			},
		})
		tools.RespondError(w, "User not found", http.StatusNotFound)
		return false
	}

	targetUser.PasswordHash = nil

	if targetUser.IsPublic {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "Profile is public, access granted",
			Metadata: map[string]any{
				"user_id":   targetUser.ID,
				"requester": requesterID,
				"ip":        r.RemoteAddr,
				"path":      r.URL.Path,
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
			Metadata: map[string]any{
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
			Metadata: map[string]any{
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
		Metadata: map[string]any{
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
			Metadata: map[string]any{
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"error": err.Error(),
			},
		})
		tools.RespondError(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if !rt.ProfileAccess(w, r, user) {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Profile info access denied",
			Metadata: map[string]any{
				"user_id": user.ID,
				"ip":      r.RemoteAddr,
				"path":    r.URL.Path,
			},
		})

		user.DateOfBirth = time.Unix(0, 0)
		user.AboutMe = ""
	}

	if err := tools.EncodeJSON(w, http.StatusOK, user); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send profile info response",
			Metadata: map[string]any{
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
		Metadata: map[string]any{
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
			Metadata: map[string]any{
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"error": err.Error(),
			},
		})
		tools.RespondError(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if !rt.ProfileAccess(w, r, user) {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Profile activity access denied",
			Metadata: map[string]any{
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
		Metadata: map[string]any{
			"user_id": user.ID,
			"ip":      r.RemoteAddr,
			"path":    r.URL.Path,
		},
	})

	tools.RespondError(w, "Not implemented", http.StatusNotImplemented)
}

func (rt *Root) ProfileFollowers(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	if err := tools.DecodeJSON(r, &user); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to decode profile connections JSON",
			Metadata: map[string]any{
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"error": err.Error(),
			},
		})
		tools.RespondError(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if !rt.ProfileAccess(w, r, user) {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Profile connections access denied",
			Metadata: map[string]any{
				"user_id": user.ID,
				"ip":      r.RemoteAddr,
				"path":    r.URL.Path,
			},
		})
		tools.RespondError(w, "Forbidden", http.StatusForbidden)
		return
	}

	followers, err := rt.DL.Follows.GetFollows(user.ID, "followers")
	if err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to fetch profile connections",
			Metadata: map[string]any{
				"user_id": user.ID,
				"ip":      r.RemoteAddr,
				"path":    r.URL.Path,
				"error":   err,
			},
		})
		tools.RespondError(w, "Failed to fetch connections", http.StatusInternalServerError)
		return
	}

	if err := tools.EncodeJSON(w, http.StatusOK, followers); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send profile followers response",
			Metadata: map[string]any{
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
		Metadata: map[string]any{
			"user_id": user.ID,
			"ip":      r.RemoteAddr,
			"path":    r.URL.Path,
		},
	})
}

func (rt *Root) ProfileFollowing(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	if err := tools.DecodeJSON(r, &user); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to decode profile connections JSON",
			Metadata: map[string]any{
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"error": err.Error(),
			},
		})
		tools.RespondError(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if !rt.ProfileAccess(w, r, user) {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Profile connections access denied",
			Metadata: map[string]any{
				"user_id": user.ID,
				"ip":      r.RemoteAddr,
				"path":    r.URL.Path,
			},
		})
		tools.RespondError(w, "Forbidden", http.StatusForbidden)
		return
	}

	following, err := rt.DL.Follows.GetFollows(user.ID, "following")
	if err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to fetch profile connections",
			Metadata: map[string]any{
				"user_id": user.ID,
				"ip":      r.RemoteAddr,
				"path":    r.URL.Path,
				"error":   err,
			},
		})
		tools.RespondError(w, "Failed to fetch connections", http.StatusInternalServerError)
		return
	}

	if err := tools.EncodeJSON(w, http.StatusOK, following); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send profile following response",
			Metadata: map[string]any{
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
		Metadata: map[string]any{
			"user_id": user.ID,
			"ip":      r.RemoteAddr,
			"path":    r.URL.Path,
		},
	})
}

func (rt *Root) UpdateProfileVisibility(w http.ResponseWriter, r *http.Request) {
	requesterID := rt.DL.GetRequesterID(w, r)

	user := &models.User{ID: requesterID}
	if err := rt.DL.Users.GetUserByID(user); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to fetch user for visibility toggle",
			Metadata: map[string]any{
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

	if err := rt.DL.Users.Update(user); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to update user visibility",
			Metadata: map[string]any{
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
			Metadata: map[string]any{
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
		Metadata: map[string]any{
			"user_id":   user.ID,
			"is_public": user.IsPublic,
			"ip":        r.RemoteAddr,
			"path":      r.URL.Path,
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

	followRequest := &models.FollowRequest{FromUserID: requesterID, ToUserID: payload.TargetId}

	switch payload.Action {
	case "follow":
		if err := rt.DL.Follows.Insert(followRequest); err != nil {
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "ERROR",
				Message: "Failed to insert follow request",
				Metadata: map[string]any{
					"from":  requesterID,
					"to":    payload.TargetId,
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
			Metadata: map[string]any{
				"from": requesterID,
				"to":   payload.TargetId,
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
			},
		})
	case "unfollow":
		if err := rt.DL.Follows.Delete(followRequest); err != nil {
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "ERROR",
				Message: "Failed to unfollow user",
				Metadata: map[string]any{
					"from":  requesterID,
					"to":    payload.TargetId,
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
			Metadata: map[string]any{
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
			Metadata: map[string]any{
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
			Metadata: map[string]any{
				"payload": payload,
				"ip":      r.RemoteAddr,
				"path":    r.URL.Path,
				"error":   err.Error(),
			},
		})
	}
}

func (rt *Root) AcceptDeclineFollowRequest(w http.ResponseWriter, r *http.Request) {
	payload := &struct {
		TargetId int    `json:"target_id"`
		Action   string `json:"action"`
	}{}

	if err := tools.DecodeJSON(r, &payload); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to decode accept/decline JSON",
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

	if payload.Action != "accepted" && payload.Action != "declined" {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Invalid accept/decline action",
			Metadata: map[string]any{
				"action": payload.Action,
				"from":   payload.TargetId,
				"to":     requesterID,
				"ip":     r.RemoteAddr,
				"path":   r.URL.Path,
			},
		})
		tools.RespondError(w, "Invalid action", http.StatusBadRequest)
		return
	}

	followRequest := &models.FollowRequest{
		FromUserID: payload.TargetId,
		ToUserID:   requesterID,
		Status:     payload.Action,
	}

	if err := rt.DL.Follows.UpdateStatus(followRequest); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to update follow request status",
			Metadata: map[string]any{
				"from":   payload.TargetId,
				"to":     requesterID,
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
		Metadata: map[string]any{
			"from":   payload.TargetId,
			"to":     requesterID,
			"ip":     r.RemoteAddr,
			"path":   r.URL.Path,
			"action": payload.Action,
		},
	})

	if err := tools.EncodeJSON(w, http.StatusOK, nil); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send accept/decline response",
			Metadata: map[string]any{
				"payload": payload,
				"ip":      r.RemoteAddr,
				"path":    r.URL.Path,
				"error":   err.Error(),
			},
		})
	}
}
