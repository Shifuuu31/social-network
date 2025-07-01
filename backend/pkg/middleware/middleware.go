package middleware

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"social-network/pkg/models"
	"social-network/pkg/tools"
)

type DataLayer struct {
	Users    *models.UserModel
	Sessions *models.SessionModel
	Posts    *models.PostModel
	Follows  *models.FollowRequestModel
	Groups   *models.GroupModel
	Members  *models.GroupMemberModel
	Events   *models.EventModel
	Votes    *models.EventVoteModel
	Logger   *models.LoggerModel
	// link to other models db connection
}

// Define paths to skip (exact match only)
var skipPaths = []string{
	"/auth/signup",
	"/signup",
	"/signin",

	"/auth/signin",
}

// Optional: define prefixes to skip (e.g., for /static/*)
var skipPrefixes = []string{
	"/public/",
}

func (dl *DataLayer) GetRequesterID(w http.ResponseWriter, r *http.Request) (requesterID int) {
	requesterID, ok := r.Context().Value(models.UserIDKey).(int)
	if !ok {
		fmt.Println("StatusUnauthorized, Context")
		tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
		dl.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Unauthorized",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
			},
		})
		return 0
	}
	fmt.Println("ReqsterID:", requesterID)
	return requesterID
}

// RequireAuth checks if a user is authenticated by session
func (dl *DataLayer) AccessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:5173" {
			fmt.Println("wwwwwaaaa")
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // ✅ THIS LINE IS MANDATORY
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		}
		// Check if request path is in skipPaths
		// Check if request path matches any skipPrefixes
		if slices.Contains(skipPaths, r.URL.Path) || tools.SliceHasPrefix(skipPrefixes, r.URL.Path) {
			fmt.Println("1")
			next.ServeHTTP(w, r)
			return
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		fmt.Println("didn't skip")
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
			fmt.Println("Cookie", cookie)
			fmt.Println("StatusUnauthorized Cookie", err)
			tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		fmt.Println("Cookie", cookie)
		fmt.Println("Coosssskie", cookie.Value)

		fmt.Println(111)
		// validate token and get user ID
		session, err := dl.Sessions.GetSessionByToken(cookie.Value)
		if err != nil {
			dl.Logger.Log(models.LogEntry{
				Level:    "WARN",
				Message:  "Missing osr Invalid session cookie",
				Metadata: map[string]any{},
			})
			fmt.Println("StatusUnauthorized Get", err)
			tools.RespondError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Println(222)
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
		fmt.Println("CORS middleware reached:", r.Method, r.URL.Path)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
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
