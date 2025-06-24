package middleware

import (
	"context"
	"net/http"
	"time"

	"social-network/pkg/models"
	"social-network/pkg/tools"
)

type DataLayer struct {
	Users *models.UserModel
	Sessions *models.SessionModel
	Logger   *models.LoggerModel
	// link to other models db connection
}

// RequireAuth checks if a user is authenticated by session
func (dl *DataLayer) AccessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil || cookie.Value == "" {
			dl.Logger.Log(models.LogEntry{
				Level:   "WARN",
				Message: "Missing or empty session cookie",
				Metadata: map[string]interface{}{
					"method": r.Method,
					"path":   r.URL.Path,
				},
			})
			tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// TODO: validate token and get user ID
		session, err := dl.Sessions.GetSessionByToken(cookie.Value)
		if err != nil {
			dl.Logger.Log(models.LogEntry{
				Level:    "WARN",
				Message:  "Missing osr Invalid session cookie",
				Metadata: map[string]interface{}{},
			})
			tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// TODO: and get user id by session token
		userID := session.UserID
		dl.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "Authorized access granted",
			Metadata: map[string]interface{}{
				"userID": userID,
				"token":  cookie.Value,
				"path":   r.URL.Path,
			},
		})
		ctx := context.WithValue(r.Context(), tools.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CORSMiddleware sets CORS headers
func (dl *DataLayer) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if r.Method == "OPTIONS" {
			dl.Logger.Log(models.LogEntry{
				Level:   "INFO",
				Message: "CORS preflight handled",
				Metadata: map[string]interface{}{
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
					Metadata: map[string]interface{}{
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
			Metadata: map[string]interface{}{
				"timeout": duration.String(),
			},
		})
		return http.TimeoutHandler(next, duration, "Request timed out")
	}
}

// ChainMiddlewares wraps all middleware in correct order
func (dl *DataLayer) GlobalMiddleware(handler http.Handler, requireAuth bool) http.Handler {
	timeout := 10 * time.Second
	if requireAuth {
		handler = dl.AccessMiddleware(handler)
	}
	return dl.RecoverMiddleware(
		dl.TimeoutMiddleware(timeout)(
			// dl.CORSMiddleware(
				handler,
			// ),
		),
	)
}
