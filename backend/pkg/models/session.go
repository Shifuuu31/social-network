package models

import "time"

type Session struct {
	ID             string
	UserID         string
	ExpirationTime time.Time
}
// repository represent functions that communicate with sessions table in db
type SessionRepository interface {
	// save new session to db
	Set(Session) error
	// Gets session from db based on session id
	Get(sID string) (Session, error)
	// Gets session from db based on user id
	GetByUser(userID string) (Session, error)
	// Update sessions expiration time
	Update(Session) error
	// Delete session 
	Delete(Session) error
}
