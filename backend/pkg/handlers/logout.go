package handlers

import (
	"fmt"
	"net/http"

	"social-network/pkg/models"
	"social-network/pkg/utils"
)

// handler for logout/ validate user by id and delete session
func (handler *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	userId := r.Context().Value(utils.UserKey).(string)
	// delete session
	session := models.Session{UserID: userId}
	errSession := handler.Repos.SessionRepo.Delete(session)
	if errSession != nil {
		fmt.Println("error on deleting session", errSession)
		return
	}
	// delete cookie
	utils.DeleteCookie(w)
	utils.RespondWithSuccess(w, "Logout successful", 200)
}
