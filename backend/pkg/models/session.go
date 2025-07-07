package models

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"session_token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type SessionModel struct {
	DB *sql.DB
}

func (sm *SessionModel) SetSession(w http.ResponseWriter, userId int) error {
	newSession := &Session{
		UserID:    userId,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := sm.UpsertSession(newSession); err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:     "session_token",
		Value:    newSession.Token,
		Path:     "/",
		// HttpOnly: true,
		// SameSite: http.SameSiteLaxMode,
		Expires:  newSession.ExpiresAt,
	}
	http.SetCookie(w, &cookie)
	return nil
}


// InsertSession adds a new session to the database.
func (sm *SessionModel) UpsertSession(session *Session) error {
	query := `
		INSERT INTO sessions (user_id, session_token, expires_at)
		VALUES (?, ?, ?)
		ON CONFLICT(session_token) DO UPDATE SET
			user_id = excluded.user_id,
			expires_at = excluded.expires_at
	`
	result, err := sm.DB.Exec(query, session.UserID, session.Token, session.ExpiresAt)
	if err != nil {
		return fmt.Errorf("upsert session: %w", err)
	}

	id, err := result.LastInsertId()
	if err == nil {
		session.ID = int(id)
	}

	return nil
}

// GetSessionByToken finds a session by its token.
func (sm *SessionModel) GetSessionByToken(token string) (*Session, error) {
	query := `
		SELECT id, user_id, session_token, expires_at, created_at
		FROM sessions
		WHERE session_token = ?
	`

	var session Session
	err := sm.DB.QueryRow(query, token).Scan(
		&session.ID,
		&session.UserID,
		&session.Token,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("get session by token: %w", err)
	}

	return &session, nil
}

// DeleteSessionByToken removes a session by its token.
func (sm *SessionModel) DeleteSessionByToken(token string) error {
	query := `DELETE FROM sessions WHERE session_token = ?`

	_, err := sm.DB.Exec(query, token)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

// DeleteExpiredSessions removes all expired sessions.
func (sm *SessionModel) DeleteExpiredSessions() error {
	query := `DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP`

	_, err := sm.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("delete expired sessions: %w", err)
	}
	return nil
}
