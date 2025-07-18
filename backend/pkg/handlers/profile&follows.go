package handlers

import (
	"fmt"
	"net/http"
	"time"

	"social-network/pkg/models"
	"social-network/pkg/tools"
)

func (rt *Root) NewUsersHandler() (usersMux *http.ServeMux) {
	usersMux = http.NewServeMux()
	fmt.Println("ss")
	usersMux.HandleFunc("POST /profile/info", rt.ProfileInfo)
	usersMux.HandleFunc("POST /profile/me", rt.GetCurrentUser)
	usersMux.HandleFunc("POST /profile/activity", rt.ProfileActivity)
	usersMux.HandleFunc("POST /profile/followers", rt.ProfileFollowers)
	usersMux.HandleFunc("POST /profile/following", rt.ProfileFollowing)
	usersMux.HandleFunc("POST /profile/visibility", rt.UpdateProfileVisibility)
	usersMux.HandleFunc("POST /follow/follow-unfollow", rt.FollowUnfollow)
	usersMux.HandleFunc("POST /follow/accept-decline", rt.AcceptDeclineFollowRequest)
	usersMux.HandleFunc("POST /close-friends/add", rt.AddCloseFriendHandler)
	usersMux.HandleFunc("POST /close-friends/remove", rt.RemoveCloseFriendHandler)
	usersMux.HandleFunc("GET /close-friends/list", rt.ListCloseFriendsHandler)

	return usersMux
}

type ProfileInfoResponse struct {
	User         *models.User `json:"user"`
	FollowStatus string       `json:"follow_status"`
}

func NewUserDTO(u *models.User) *models.UserDTO {
	return &models.UserDTO{
		ID:          u.ID,
		Nickname:    u.Nickname,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		AboutMe:     u.AboutMe,
		DateOfBirth: u.DateOfBirth,
		AvatarURL:   u.AvatarURL,
		IsPublic:    u.IsPublic,
		CreatedAt:   u.CreatedAt,
	}
}

func (rt *Root) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println(15551)
	requesterID := rt.DL.GetRequesterID(w, r)
	user := &models.User{ID: requesterID}

	if err := rt.DL.Users.GetUserByID(user); err != nil {
		tools.RespondError(w, "User not found", http.StatusUnauthorized)
		return
	}

	dto := NewUserDTO(user)
	tools.EncodeJSON(w, http.StatusOK, dto)
}

func (rt *Root) ProfileAccess(w http.ResponseWriter, r *http.Request, targetUser *models.User) bool {
	fmt.Println("PrfAcc")
	requesterID := rt.DL.GetRequesterID(w, r)
	fmt.Println("requesterID", requesterID)
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
	if requesterID != targetUser.ID {
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
			fmt.Println(1, err)
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
			fmt.Println("ddd", 2)
			// tools.RespondError(w, "Private profile — follow to see more", http.StatusForbidden)
			return false
		}
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
	fmt.Println("S_prfl")

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
		println(err)
		tools.RespondError(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	// fmt.Println("VVV")

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
		fmt.Println("Profil access denied")

		user.DateOfBirth = time.Unix(0, 0)
		user.AboutMe = ""
	}

	requesterID := rt.DL.GetRequesterID(w, r)
	followRequest := &models.FollowRequest{FromUserID: requesterID, ToUserID: user.ID}
	// followStaus := "none"
	if err := rt.DL.Follows.GetFollowStatus(followRequest); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Follow status not found or error",
			Metadata: map[string]any{
				"from_user": requesterID,
				"to_user":   user.ID,
				"ip":        r.RemoteAddr,
				"path":      r.URL.Path,
				"error":     err.Error(),
			},
		})
		fmt.Println(1, err)
		tools.RespondError(w, "Private profile — follow to see more", http.StatusForbidden)
		return
	}

	fmt.Println("Follow", followRequest.Status)
	response := ProfileInfoResponse{
		User:         user,
		FollowStatus: followRequest.Status,
	}
	if err := tools.EncodeJSON(w, http.StatusOK, response); err != nil {
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
		fmt.Println("err", err)

		return
	}

	fmt.Println("user", user)
	fmt.Println("E_prfl")

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
		fmt.Println(3)
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
	fmt.Println("ProflieFollowers")
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
		fmt.Println(4)
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
	fmt.Println("follow")
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
		fmt.Println(5)
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
		fmt.Println("err", err)
		return
	}
	fmt.Println(user.IsPublic)
	user.IsPublic = !user.IsPublic
	fmt.Println(user.IsPublic)

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
		fmt.Println("err", err)
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(user.IsPublic)

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
	fmt.Println("Following")
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
			fmt.Println("'Err'", err)
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

// --- Close Friends Handlers ---

// POST /close-friends/add
func (rt *Root) AddCloseFriendHandler(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		FriendID int `json:"friend_id"`
	}{}
	if err := tools.DecodeJSON(r, &payload); err != nil {
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid payload"})
		return
	}
	userID := rt.DL.GetRequesterID(w, r)
	if userID <= 0 || payload.FriendID <= 0 {
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user or friend ID"})
		return
	}
	if err := rt.DL.Follows.AddCloseFriend(userID, payload.FriendID); err != nil {
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to add close friend"})
		return
	}
	tools.EncodeJSON(w, http.StatusOK, map[string]string{"message": "Close friend added"})
}

// POST /close-friends/remove
func (rt *Root) RemoveCloseFriendHandler(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		FriendID int `json:"friend_id"`
	}{}
	if err := tools.DecodeJSON(r, &payload); err != nil {
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid payload"})
		return
	}
	userID := rt.DL.GetRequesterID(w, r)
	if userID <= 0 || payload.FriendID <= 0 {
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user or friend ID"})
		return
	}
	if err := rt.DL.Follows.RemoveCloseFriend(userID, payload.FriendID); err != nil {
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to remove close friend"})
		return
	}
	tools.EncodeJSON(w, http.StatusOK, map[string]string{"message": "Close friend removed"})
}

// GET /close-friends/list
func (rt *Root) ListCloseFriendsHandler(w http.ResponseWriter, r *http.Request) {
	userID := rt.DL.GetRequesterID(w, r)
	if userID <= 0 {
		tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}
	friends, err := rt.DL.Follows.ListCloseFriends(userID)
	if err != nil {
		tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list close friends"})
		return
	}
	tools.EncodeJSON(w, http.StatusOK, friends)
}
