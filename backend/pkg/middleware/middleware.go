package middleware

import (
	"context"
	"net/http"
	"slices"
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
}

// Optional: define prefixes to skip (e.g., for /static/*)
var skipPrefixes = []string{
	"/public/",
}

func (dl *DataLayer) GetRequesterID(w http.ResponseWriter, r *http.Request) (requesterID int) {
	requesterID, ok := r.Context().Value(models.UserIDKey).(int)
	if !ok {
		dl.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Unauthorized: requester ID not found",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
			},
		})
		tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
		return 0
	}
	return requesterID
}

// RequireAuth checks if a user is authenticated by session
func (dl *DataLayer) AccessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request path is in skipPaths
		// Check if request path matches any skipPrefixes
		if slices.Contains(skipPaths, r.URL.Path) && tools.SliceHasPrefix(skipPrefixes, r.URL.Path) {
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
