package models

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Don't expose hash in JSON responses
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	AvatarURL    string    `json:"avatar_url"`
	Nickname     string    `json:"nickname"`
	AboutMe      string    `json:"about_me"`
	IsPublic     bool      `json:"is_public"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserModel struct {
	DB *sql.DB
}

// InsertUser inserts a new user into the database.
func (um *UserModel) InsertUser(user *User) error {
	query := `
		INSERT INTO users (
			email, password_hash, first_name, last_name, date_of_birth,
			avatar_url, nickname, about_me, is_public, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	res, err := um.DB.Exec(query,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.DateOfBirth,
		user.AvatarURL,
		user.Nickname,
		user.AboutMe,
		user.IsPublic,
	)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	lastID, err := res.LastInsertId()
	if err == nil {
		user.ID = int(lastID)
	}

	return nil
}

// UpdateUser updates an existing user's profile data.
func (um *UserModel) UpdateUser(user *User) error {
	query := `
		UPDATE users
		SET first_name = ?, last_name = ?, date_of_birth = ?, avatar_url = ?, nickname = ?,
			about_me = ?, is_public = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := um.DB.Exec(query,
		user.FirstName,
		user.LastName,
		user.DateOfBirth,
		user.AvatarURL,
		user.Nickname,
		user.AboutMe,
		user.IsPublic,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

// DeleteUser deletes a user by ID.
func (um *UserModel) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := um.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

// GetUserByID retrieves a user by ID.
func (um *UserModel) GetUserByID(id int) (*User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, date_of_birth,
		       avatar_url, nickname, about_me, is_public, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	var user User
	err := um.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.AvatarURL,
		&user.Nickname,
		&user.AboutMe,
		&user.IsPublic,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user with this id: %w", err) // Not found
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &user, nil
}
