package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"social-network/pkg/models"
	"social-network/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

func (handler *Handler) Signin(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	/* ---------------------------- read incoming data --------------------------- */
	// Try to decode the JSON request to a LoginUser
	var client models.User
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	/* --------------------------- validate user in db -------------------------- */
	// find user with email in db (need Password and user_id)
	dbUser, errDb := handler.Repos.UserRepo.FindUserByEmail(client.Email)
	if errDb != nil {
		utils.RespondWithError(w, "Wrong credentials", 200)
		return
	}
	// Compare passwords
	errPwd := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(client.Password))
	if errPwd != nil {
		utils.RespondWithError(w, "Wrong credentials", 200)
		return
	}
	/* ------------------- user valid - create/update session ------------------- */
	// get existing session from db
	_, errSession := handler.Repos.SessionRepo.GetByUser(dbUser.ID)
	newSession := utils.SessionStart(w, r, dbUser.ID)
	var errOnSave error
	// Update or create new row in db based on/ if session already exist
	if errSession != nil { // create new session
		errOnSave = handler.Repos.SessionRepo.Set(newSession)
	} else { // Update existing
		errOnSave = handler.Repos.SessionRepo.Update(newSession)
	}
	if errOnSave != nil {
		utils.RespondWithError(w, "Error on creating new session", 200)
		return
	}
	utils.RespondWithSuccess(w, "Login successful", 200)
}

// endpoint for checking if user session is already in progress
// responds with err if no session active
// responds with success if session valid
// updates session access time in db
func (handler *Handler) SessionActive(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// Get cookie value from request
	sessionId, errCookie := utils.GetCookie(r)
	if errCookie != nil {
		utils.RespondWithError(w, "Session not active", 200)
		return
	}
	// Get session based on session id
	session, errSession := handler.Repos.SessionRepo.Get(sessionId)
	if errSession != nil {
		utils.RespondWithError(w, "Session not active", 200)
		return
	}
	// check if session not expired
	sessionValid := utils.CheckSessionExpiration(session)
	if !sessionValid {
		// if not valid any more delete from db
		handler.Repos.SessionRepo.Delete(session)
		// Delete from client browser
		utils.DeleteCookie(w)
		utils.RespondWithError(w, "Session not active", 200)
		return
	} else {
		// Session stil valid -> prolong it by 30 min
		session.ExpirationTime = time.Now().Add(30 * time.Minute)
		handler.Repos.SessionRepo.Update(session)
		utils.RespondWithSuccess(w, "Session active", 200)
		return
	}
}
