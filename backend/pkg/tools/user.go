package tools

import (
	"net/http"
)

type contextKey string

const UserIDKey = contextKey("userID")

func GetUserID(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	return userID, ok
}
