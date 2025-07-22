package handlers

import (
	"net/http"
	"strings"

	"social-network/pkg/models"
	"social-network/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

// Register user endpont -> validate inputs / save in db / start session
func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)

	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on form submittion", 400)
		return
	}
	err := r.ParseMultipartForm(3145728) // 3MB
	if err != nil {
		utils.RespondWithError(w, "Error in form validation", 400)
		return
	}

	// Create new user instance
	newUser := models.User{
		Email:       strings.ToLower(r.PostFormValue("email")),
		FirstName:   r.PostFormValue("firstname"),
		LastName:    r.PostFormValue("lastname"),
		Password:    r.PostFormValue("password"),
		Nickname:    r.PostFormValue("nickname"),
		About:       r.PostFormValue("aboutme"),
		DateOfBirth: r.PostFormValue("dateofbirth"),
	}
	// Validate all user fields
	errValid := utils.ValidateNewUser(newUser)
	if errValid != nil {
		utils.RespondWithError(w, "Error in validation", 400)
		return
	}

	// Check if email alredy taken
	if emailUnique, _ := handler.Repos.UserRepo.EmailNotTaken(newUser.Email); !emailUnique {
		utils.RespondWithError(w, "Email already taken", 409)
		return
	}
	// Hash password
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	newUser.Password = string(hashedPwd)
	// create user id
	userID := utils.UniqueId()
	newUser.ID = userID
	// check if avatar added / save in filesystem
	newUser.ImagePath = utils.SaveAvatar(r)
	// Save user in db
	errSave := handler.Repos.UserRepo.Add(newUser)
	if errSave != nil {
		utils.RespondWithError(w, "Couldn't save new user", 500)
		return
	}
	// Start new session for user (Including cookies)
	/* ---------- code commented out if from register redirect to login --------- */
	/*
		session := utils.SessionStart(w, r, userID)
		// Save session in database
		errSession := handler.Repos.SessionRepo.Set(session)
		if errSession != nil {
			utils.RespondWithError(w, "Error on creating new session", 500)
		}
	*/
}
