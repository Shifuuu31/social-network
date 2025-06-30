package handlers

import (
	"fmt"
	"net/http"
	"time"

	"social-network/pkg/models"
	"social-network/pkg/tools"

	"golang.org/x/crypto/bcrypt"
)

func (rt *Root) NewAuthHandler() (authMux *http.ServeMux) {
	authMux = http.NewServeMux()

	authMux.HandleFunc("POST /signup", rt.SignUp)
	authMux.HandleFunc("POST /signin", rt.SignIn)
	authMux.HandleFunc("DELETE /signout", rt.SignOut)

	// Log handler setup
	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "Auth routes registered",
		Metadata: map[string]any{
			"path": "/signup, /signin, /signout",
		},
	})

	return authMux
}

func (rt *Root) SignUp(w http.ResponseWriter, r *http.Request) {
	// verify user input
	user, err := models.Validate(r)
	if err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "Invalid signup input",
			Metadata: map[string]any{
				"email": user.Email,
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"err":   err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Signup input validated"})

	// hash user password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Password hashing failed",
			Metadata: map[string]any{
				"email": user.Email,
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"err":   err.Error(),
			},
		})
		tools.RespondError(w, "Oops, try again later.", http.StatusInternalServerError)
		return
	}
	user.PasswordHash = hash
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Password hashed successfully"})

	// insert user into db
	if err := rt.DL.Users.Insert(&user); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to insert user into DB",
			Metadata: map[string]any{
				"email": user.Email,
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"err":   err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "User inserted into DB"})

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "New user registered successfully",
		Metadata: map[string]any{
			"email": user.Email,
			"ip":    r.RemoteAddr,
			"path":  r.URL.Path,
		},
	})

	if err := tools.EncodeJSON(w, http.StatusCreated, nil); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send signup response",
			Metadata: map[string]any{
				"email": user.Email,
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"err":   err.Error(),
			},
		})
	}
}

func (rt *Root) SignIn(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	if err := tools.DecodeJSON(r, &user); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to decode signin JSON",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"err":  err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Signin JSON decoded"})

	if err := rt.DL.Users.ValidateCredentials(user); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "Invalid signin credentials",
			Metadata: map[string]any{
				"email": user.Email,
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"err":   err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusUnauthorized)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Signin credentials validated"})

	if err := rt.DL.Sessions.SetSession(w, user.ID); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to set session during signin",
			Metadata: map[string]any{
				"user_id": user.ID,
				"email":   user.Email,
				"ip":      r.RemoteAddr,
				"path":    r.URL.Path,
				"err":     err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Session set for user"})

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "User signed in successfully",
		Metadata: map[string]any{
			"user_id": user.ID,
			"email":   user.Email,
			"ip":      r.RemoteAddr,
			"path":    r.URL.Path,
		},
	})
	if err := tools.EncodeJSON(w, http.StatusOK, nil); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send signin response",
			Metadata: map[string]any{
				"email": user.Email,
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"err":   err.Error(),
			},
		})
	}
}

func (rt *Root) SignOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	fmt.Println(cookie, err) // existing debug print left as is

	if err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Missing session cookie during signout",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"err":  err.Error(),
			},
		})
		tools.RespondError(w, "Session not found", http.StatusBadRequest)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Session cookie found for signout"})

	if err := rt.DL.Sessions.DeleteSessionByToken(cookie.Value); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to delete session from DB",
			Metadata: map[string]any{
				"token": cookie.Value,
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"err":   err.Error(),
			},
		})
		tools.RespondError(w, "Failed to sign out", http.StatusInternalServerError)
		return
	}
	rt.DL.Logger.Log(models.LogEntry{Level: "DEBUG", Message: "Session deleted from DB"})

	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	})

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "User signed out successfully",
		Metadata: map[string]any{
			"ip":   r.RemoteAddr,
			"path": r.URL.Path,
		},
	})

	if err := tools.EncodeJSON(w, http.StatusOK, nil); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send signout response",
			Metadata: map[string]any{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"err":  err.Error(),
			},
		})
	}
}
