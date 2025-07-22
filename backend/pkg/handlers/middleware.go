package handlers

import (
	"context"
	"net/http"
	"time"

	"social-network/pkg/utils"
)

// basic authentification/ check if user logged in
// If not logged in not logged in return
// if logged in continue to handler with user id added to context
// also update expiration time in database
func (handler *Handler) Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w = utils.ConfigHeader(w)
		// Get cookie value from request
		sessionId, errCookie := utils.GetCookie(r)
		if errCookie != nil {
			utils.RespondWithError(w, "Error on getting cookie", 200)
			return
		}
		// Get session based on session id
		session, errSession := handler.Repos.SessionRepo.Get(sessionId)
		if errSession != nil {
			utils.RespondWithError(w, "Error on getting session", 200)
			return
		}
		// check if session not expired
		sessionValid := utils.CheckSessionExpiration(session)
		if !sessionValid {
			// if not valid any more delete from db
			handler.Repos.SessionRepo.Delete(session)
			// Delete from client browser
			utils.DeleteCookie(w)
			utils.RespondWithError(w, "Session is not valid", 200)
			return
		} else {
			// Session stil valid -> prolong it by 30 min
			session.ExpirationTime = time.Now().Add(30 * time.Minute)
			handler.Repos.SessionRepo.Update(session)
		}
		// Auth successful, continue with adding User_id to request context
		ctx := context.WithValue(r.Context(), utils.UserKey, session.UserID)
		next(w, r.WithContext(ctx))
	})
}
