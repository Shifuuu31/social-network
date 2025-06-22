package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Session struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	SessionToken string    `json:"session_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type SessionModel struct {
	DB *sql.DB
}

// InsertSession adds a new session to the database.
func (sm *SessionModel) InsertSession(session *Session) error {
	query := `
		INSERT INTO sessions (user_id, session_token, expires_at)
		VALUES (?, ?, ?)
	`
	result, err := sm.DB.Exec(query, session.UserID, session.SessionToken, session.ExpiresAt)
	if err != nil {
		return fmt.Errorf("insert session: %w", err)
	}

	// Save the new session ID
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
		&session.SessionToken,
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
