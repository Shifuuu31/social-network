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
}

// Optional: define prefixes to skip (e.g., for /static/*)
var skipPrefixes = []string{
	"/public/",
}

func (dl *DataLayer) GetRequesterID(w http.ResponseWriter, r *http.Request) (requesterID int) {
	// FOR TESTING WITHOUT AUTHENTICATION - RETURN HARDCODED USER IDS
	// Comment out this section and uncomment below for production
	
	// Define test user IDs based on path patterns for different test scenarios
	path := r.URL.Path
	
	// Default test user ID
	testUserID := 1
	
	// Use different user IDs for different endpoints to simulate different users
	if path == "/notifications/fetch" || 
	   path == "/notifications/mark-seen" || 
	   path == "/notifications/mark-all-seen" || 
	   path == "/notifications/unread-count" {
		// For notification endpoints, use user ID 2 to see notifications sent to user 2
		testUserID = 2
	}
	
	// For group operations, use user ID 1 (group creator in most cases)
	if slices.Contains([]string{"POST", "PUT", "DELETE"}, r.Method) {
		if path == "/groups/accept-decline" || path == "/groups/invite" {
			testUserID = 1 // Group creator
		}
		if path == "/groups/request-join" {
			testUserID = 3 // Different user requesting to join
		}
	}
	
	// For follow operations, use user ID 1
	if path == "/profile/follow" || path == "/profile/accept-decline-follow" {
		testUserID = 1
	}
	
	dl.Logger.Log(models.LogEntry{
		Level:   "DEBUG",
		Message: "Using hardcoded test user ID",
		Metadata: map[string]any{
			"test_user_id": testUserID,
			"path":         path,
			"method":       r.Method,
		},
	})
	
	return testUserID
	
	// PRODUCTION CODE (commented out for testing)
	// requesterID, ok := r.Context().Value(models.UserIDKey).(int)
	// if !ok {
	// 	dl.Logger.Log(models.LogEntry{
	// 		Level:   "ERROR",
	// 		Message: "Unauthorized: requester ID not found",
	// 		Metadata: map[string]any{
	// 			"ip":   r.RemoteAddr,
	// 			"path": r.URL.Path,
	// 		},
	// 	})
	// 	tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
	// 	return 0
	// }
	// return requesterID
}

// getTestUserID returns a test user ID based on the request path for testing purposes
func (dl *DataLayer) getTestUserID(r *http.Request) int {
	path := r.URL.Path
	
	// Default test user ID
	testUserID := 1
	
	// Use different user IDs for different endpoints to simulate different users
	if path == "/notifications/fetch" || 
	   path == "/notifications/mark-seen" || 
	   path == "/notifications/mark-all-seen" || 
	   path == "/notifications/unread-count" {
		// For notification endpoints, use user ID 2 to see notifications sent to user 2
		testUserID = 2
	}
	
	// For group operations, use user ID 1 (group creator in most cases)
	if slices.Contains([]string{"POST", "PUT", "DELETE"}, r.Method) {
		if path == "/groups/accept-decline" || path == "/groups/invite" {
			testUserID = 1 // Group creator
		}
		if path == "/groups/request-join" {
			testUserID = 3 // Different user requesting to join
		}
	}
	
	// For follow operations, use user ID 1
	if path == "/profile/follow" || path == "/profile/accept-decline-follow" {
		testUserID = 1
	}
	
	dl.Logger.Log(models.LogEntry{
		Level:   "DEBUG",
		Message: "Using test user ID for skipped path",
		Metadata: map[string]any{
			"test_user_id": testUserID,
			"path":         path,
			"method":       r.Method,
		},
	})
	
	return testUserID
}

// RequireAuth checks if a user is authenticated by session
func (dl *DataLayer) AccessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request path is in skipPaths or matches any skipPrefixes
		if slices.Contains(skipPaths, r.URL.Path) || tools.SliceHasPrefix(skipPrefixes, r.URL.Path) {
			// For testing: Add a fake user context for skipped paths (notification endpoints)
			if slices.Contains(skipPaths, r.URL.Path) {
				// Add test user context for notification endpoints
				testUserID := dl.getTestUserID(r)
				ctx := context.WithValue(r.Context(), models.UserIDKey, testUserID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			next.ServeHTTP(w, r)
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

		// userID, ok := r.Context().Value(models.UserIDKey).(int) //TODO: Get user ID from context
		userID, ok := 1, true
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

// TimeoutMiddleware applies a timeout to requests
func (dl *DataLayer) TimeoutMiddleware(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		dl.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "Applied timeout middleware",
			Metadata: map[string]any{
				"timeout": duration.String(),
			},
		})
		return http.TimeoutHandler(next, duration, "Request timed out")
	}
}

// ChainMiddlewares wraps all middleware in correct order
func (dl *DataLayer) GlobalMiddleware(handler http.Handler) http.Handler {
	timeout := 10 * time.Second

	return dl.AccessMiddleware(dl.RecoverMiddleware(
		dl.TimeoutMiddleware(timeout)(
			// dl.CORSMiddleware(
			handler,
			// ),
		),
	),
	)
}
