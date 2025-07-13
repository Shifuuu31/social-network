package handlers

import (
	"fmt"
	"net/http"
	"strings"
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

	return authMux
}

func (rt *Root) SignUp(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	var avatarPath string
	// fmt.Println("s")
	fmt.Println("DEBUG Content-Type:", r.Header.Get("Content-Type"))

	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "multipart/form-data") {
		// Handle FormData (with avatar file)
		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB limit
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "ERROR",
				Message: "Failed to parse multipart form",
				Metadata: map[string]interface{}{
					"ip":   r.RemoteAddr,
					"path": r.URL.Path,
					"err":  err.Error(),
				},
			})
			tools.RespondError(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		// Parse user data from form fields
		user = &models.User{
			Email:     r.FormValue("email"),
			Password:  r.FormValue("password"),
			FirstName: r.FormValue("first_name"),
			LastName:  r.FormValue("last_name"),
			Gender:    r.FormValue("gender"),
			Nickname:  r.FormValue("nickname"),
			AboutMe:   r.FormValue("about_me"),
			IsPublic:  r.FormValue("is_public") == "true",
		}

		// Parse date of birth
		if dobStr := r.FormValue("date_of_birth"); dobStr != "" {
			if dob, err := time.Parse(time.RFC3339, dobStr); err == nil {
				user.DateOfBirth = dob
			}
		}

		// Handle avatar file upload
		if file, header, err := r.FormFile("avatar_file"); err == nil {
			defer file.Close()

			// Validate file size
			if header.Size > 5<<20 { // 5 MB limit
				tools.RespondError(w, "Avatar file too large (max 5MB)", http.StatusBadRequest)
				return
			}

			// Validate file type using existing tools function
			if !tools.IsAllowedFile(header.Filename, file) {
				tools.RespondError(w, "Avatar file format is not supported", http.StatusBadRequest)
				return
			}

			// Upload file using existing tools function
			uploadedPath, status := tools.UploadHandler(file, header)
			if status != 200 {
				tools.RespondError(w, "Failed to upload avatar", status)
				return
			}
			avatarPath = uploadedPath
		}
	} else {
		// Handle JSON request (no avatar)
		if err := tools.DecodeJSON(r, &user); err != nil {
			rt.DL.Logger.Log(models.LogEntry{
				Level:   "ERROR",
				Message: "Failed to decode signup JSON",
				Metadata: map[string]interface{}{
					"ip":   r.RemoteAddr,
					"path": r.URL.Path,
					"err":  err.Error(),
				},
			})
			tools.RespondError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// verify user input
	if err := user.Validate(); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "Invalid signup input",
			Metadata: map[string]interface{}{
				"email": user.Email,
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"err":   err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// hash user password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Password hashing failed",
			Metadata: map[string]interface{}{
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

	// Set avatar URL if avatar was uploaded
	if avatarPath != "" {
		// Convert uploads/filename.jpg to /images/filename.jpg format
		// to match the profile upload format
		filename := strings.TrimPrefix(avatarPath, "uploads/")
		user.AvatarURL = "/images/" + filename
	}

	// insert user into db
	if err := rt.DL.Users.InsertUser(user); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to insert user into DB",
			Metadata: map[string]interface{}{
				"email": user.Email,
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"err":   err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "New user registered successfully",
		Metadata: map[string]interface{}{
			"email": user.Email,
			"ip":    r.RemoteAddr,
			"path":  r.URL.Path,
		},
	})

	if err := tools.EncodeJSON(w, http.StatusCreated, nil); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send signup response",
			Metadata: map[string]interface{}{
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
			Metadata: map[string]interface{}{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"err":  err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := rt.DL.Users.ValidateCredentials(user); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "INFO",
			Message: "Invalid signin credentials",
			Metadata: map[string]interface{}{
				"email": user.Email,
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"err":   err.Error(),
			},
		})
		tools.RespondError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := rt.DL.Sessions.SetSession(w, user.ID); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to set session during signin",
			Metadata: map[string]interface{}{
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

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "User signed in successfully",
		Metadata: map[string]interface{}{
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
			Metadata: map[string]interface{}{
				"email": user.Email,
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"err":   err.Error(),
			},
		})
	}
}

func (rt *Root) SignOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	fmt.Println(cookie, err)
	if err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "WARN",
			Message: "Missing session cookie during signout",
			Metadata: map[string]interface{}{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"err":  err.Error(),
			},
		})
		tools.RespondError(w, "Session not found", http.StatusBadRequest)
		return
	}

	if err := rt.DL.Sessions.DeleteSessionByToken(cookie.Value); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to delete session from DB",
			Metadata: map[string]interface{}{
				"token": cookie.Value,
				"ip":    r.RemoteAddr,
				"path":  r.URL.Path,
				"err":   err.Error(),
			},
		})
		tools.RespondError(w, "Failed to sign out", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	})

	rt.DL.Logger.Log(models.LogEntry{
		Level:   "INFO",
		Message: "User signed out successfully",
		Metadata: map[string]interface{}{
			"ip":   r.RemoteAddr,
			"path": r.URL.Path,
		},
	})

	if err := tools.EncodeJSON(w, http.StatusOK, nil); err != nil {
		rt.DL.Logger.Log(models.LogEntry{
			Level:   "ERROR",
			Message: "Failed to send signout response",
			Metadata: map[string]interface{}{
				"ip":   r.RemoteAddr,
				"path": r.URL.Path,
				"err":  err.Error(),
			},
		})
	}
}
