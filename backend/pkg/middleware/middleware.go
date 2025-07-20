package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"time"

	"social-network/pkg/models"
	"social-network/pkg/tools"
)

type DataLayer struct {
	Users         *models.UserModel
	Sessions      *models.SessionModel
	Posts         *models.PostModel
	Follows       *models.FollowRequestModel
	Groups        *models.GroupModel
	Members       *models.GroupMemberModel
	Events        *models.EventModel
	Votes         *models.EventVoteModel
	Messages      *models.MessageModel
	Notifications *models.NotificationModel
	Images        *models.ImageModel
	Logger        *models.LoggerModel
	// link to other models db connection
}

// Define paths to skip (exact match only)
var skipPaths = []string{
	"/auth/signup",
	"/auth/signin",
	"/connect", // WebSocket connection endpoint
	// Notification endpoints for testing without auth
	"/notifications/fetch",
	"/notifications/mark-seen",
	"/notifications/mark-all-seen",
	"/notifications/unread-count",
	// Group endpoints for testing notifications (correct paths)
	"/groups/group/browse",
	"/groups/group/new",
	"/groups/group/invite",
	"/groups/group/request",
	"/groups/group/accept-decline",
	"/groups/group/event/new",
	// Profile endpoints for testing notifications
	"/profile/follow",
	"/profile/accept-decline-follow",
}

// Optional: define prefixes to skip (e.g., for /static/*)
var skipPrefixes = []string{
	"/public/",
}

func (dl *DataLayer) GetRequesterID(w http.ResponseWriter, r *http.Request) (requesterID int) {
	// First, try to get the user ID from context (set by AccessMiddleware for skipped paths)
	if userID, ok := r.Context().Value(models.UserIDKey).(int); ok {
		dl.Logger.Log(models.LogEntry{
			Level:   "DEBUG",
			Message: "Using user ID from context",
			Metadata: map[string]any{
				"user_id": userID,
				"path":    r.URL.Path,
				"method":  r.Method,
			},
		})
		return userID
	}
	
	// FOR PRODUCTION: User is not authenticated
	dl.Logger.Log(models.LogEntry{
		Level:   "ERROR",
		Message: "Unauthorized: user not authenticated",
		Metadata: map[string]any{
			"ip":     r.RemoteAddr,
			"path":   r.URL.Path,
			"method": r.Method,
		},
	})
	tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
	return 0
}

// getTestUserID provides a basic user ID for testing
func (dl *DataLayer) getTestUserID(r *http.Request) int {
	// For demo/testing purposes, return a default user ID
	// In a real application, this would come from authentication
	return 1
}

// RequireAuth checks if a user is authenticated by session
func (dl *DataLayer) AccessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request path is in skipPaths or matches any skipPrefixes
		if slices.Contains(skipPaths, r.URL.Path) || tools.SliceHasPrefix(skipPrefixes, r.URL.Path) {
			// Add test user context for skipped paths
			testUserID := dl.getTestUserID(r)
			ctx := context.WithValue(r.Context(), models.UserIDKey, testUserID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		cookie, err := r.Cookie("session_token")
		if err != nil || cookie.Value == "" {
			dl.Logger.Log(models.LogEntry{
				Level:   "WARN",
				Message: "Missing or empty session cookie",
				Metadata: map[string]any{
					"method": r.Method,
					"path":   r.URL.Path,
				},
			})
			tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// validate token and get user ID
		session, err := dl.Sessions.GetSessionByToken(cookie.Value)
		if err != nil {
			dl.Logger.Log(models.LogEntry{
				Level:    "WARN",
				Message:  "Missing osr Invalid session cookie",
				Metadata: map[string]any{},
			})
			tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// and get user id by session token
		requesterID := session.UserID
		dl.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "Authorized access granted",
			Metadata: map[string]any{
				"requesterID": requesterID,
				"token":       cookie.Value,
				"path":        r.URL.Path,
			},
		})
		ctx := context.WithValue(r.Context(), models.UserIDKey, requesterID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// func (dl *DataLayer) GroupAccessMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var member *models.GroupMember
// 	if err := tools.DecodeJSON(r, &member); err != nil {
// 		dl.Logger.Log(models.LogEntry{
// 			Level:   "ERROR",
// 			Message: "Failed to decode group invite info JSON",
// 			Metadata: map[string]any{
// 				"ip":    r.RemoteAddr,
// 				"path":  r.URL.Path,
// 				"error": err.Error(),
// 			},
// 		})
// 		tools.RespondError(w, "Invalid payload", http.StatusBadRequest)
// 		return
// 	}
// 		// decode paylod
// 		if err := dl.Members.IsUserGroupMember(member); err != nil {
// 		dl.Logger.Log(models.LogEntry{
// 				Level:    "WARN",
// 				Message:  "Forbidden",
// 				Metadata: map[string]any{},
// 			})
// 			tools.RespondError(w, "Forbidden", http.StatusForbidden)
// 			return
// 		}
// 	})
// }

func (dl *DataLayer) GroupAccessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var groupID int
		var err error

		// Get user ID from authenticated context
		userID, ok := r.Context().Value(models.UserIDKey).(int)
		if !ok || userID <= 0 {
			dl.Logger.Log(models.LogEntry{
				Level:   "WARN",
				Message: "Missing or invalid user ID in context",
				Metadata: map[string]any{
					"path":  r.URL.Path,
					"blasa": "groupMiddleWare",
				},
			})
			tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		switch r.Method {
		case http.MethodGet:
			groupIDStr := r.PathValue("id")
			groupID, err = strconv.Atoi(groupIDStr)
			if err != nil {
				tools.RespondError(w, "Invalid group ID", http.StatusBadRequest)
				return
			}

		case http.MethodPost:
			var payload struct {
				GroupID int `json:"group_id"`
			}

			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				tools.RespondError(w, "Failed to read body", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			if len(bodyBytes) > 0 {
				if err := json.Unmarshal(bodyBytes, &payload); err != nil {
					dl.Logger.Log(models.LogEntry{
						Level:   "WARN",
						Message: "Failed to decode group payload",
						Metadata: map[string]any{
							"path":  r.URL.Path,
							"error": err.Error(),
						},
					})
					tools.RespondError(w, "Invalid payload", http.StatusBadRequest)
					return
				}
				groupID = payload.GroupID
			}

			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		default:
			tools.RespondError(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}

		if err := dl.Members.IsUserGroupMember(&models.GroupMember{
			GroupID: groupID,
			UserID:  userID,
		}); err != nil {
			dl.Logger.Log(models.LogEntry{
				Level:   "WARN",
				Message: "User not in group",
				Metadata: map[string]any{
					"group_id": groupID,
					"user_id":  userID,
					"path":     r.URL.Path,
				},
			})
			tools.RespondError(w, "Forbidden", http.StatusForbidden)
			return
		}
		fmt.Println("User is in group, proceeding with request", groupID, userID)

		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware sets CORS headers
func (dl *DataLayer) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5174")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if r.Method == "OPTIONS" {
			dl.Logger.Log(models.LogEntry{
				Level:   "INFO",
				Message: "CORS preflight handled",
				Metadata: map[string]any{
					"method": r.Method,
					"path":   r.URL.Path,
				},
			})
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RecoverMiddleware catches panics to avoid server crash
func (dl *DataLayer) RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				dl.Logger.Log(models.LogEntry{
					Level:   "ERROR",
					Message: "Panic recovered",
					Metadata: map[string]any{
						"error": err,
						"path":  r.URL.Path,
					},
				})
				tools.RespondError(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// TimeoutMiddleware applies a timeout to requests, but excludes WebSocket connections
func (dl *DataLayer) TimeoutMiddleware(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		dl.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "Applied timeout middleware",
			Metadata: map[string]any{
				"timeout": duration.String(),
			},
		})
		
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip timeout for WebSocket connections
			if r.URL.Path == "/connect" || r.Header.Get("Upgrade") == "websocket" {
				next.ServeHTTP(w, r)
				return
			}
			
			// Apply timeout for regular HTTP requests
			http.TimeoutHandler(next, duration, "Request timed out").ServeHTTP(w, r)
		})
	}
}

// ChainMiddlewares wraps all middleware in correct order
func (dl *DataLayer) GlobalMiddleware(handler http.Handler) http.Handler {
	return dl.RecoverMiddleware(
		dl.AccessMiddleware(
			handler,
		),
	)
}
