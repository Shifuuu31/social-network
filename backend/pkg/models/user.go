package models

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const UserIDKey = contextKey("userID")

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`

	DateOfBirth  time.Time `json:"date_of_birth"`
	AvatarURL    string    `json:"avatar_url"`
	Nickname     string    `json:"nickname"`
	AboutMe      string    `json:"about_me"`
	IsPublic     bool      `json:"is_public"`
	PasswordHash []byte    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

// UserDTO is a safe representation of a user to be sent to the client
type UserDTO struct {
	ID          int       `json:"id"`
	Nickname    string    `json:"nickname"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Gender      string    `json:"gender"`
	AboutMe     string    `json:"about_me"`
	DateOfBirth time.Time `json:"date_of_birth"`
	AvatarURL   string    `json:"avatar_url"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
}

// Validate checks all required fields and ensures correct formats
func (u *User) Validate() error {
	// Trim spaces
	u.Email = strings.TrimSpace(u.Email)
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Nickname = strings.TrimSpace(u.Nickname)
	u.AvatarURL = strings.TrimSpace(u.AvatarURL)
	u.AboutMe = strings.TrimSpace(u.AboutMe)

	// Email validation
	emailExp := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if u.Email == "" || !emailExp.MatchString(u.Email) {
		return errors.New("invalid email format")
	}

	// Password validation
	if len(u.Password) < 8 || len(u.Password) > 64 {
		return errors.New("password must be between 8/64 characters")
	}

	// Name validations
	if u.FirstName == "" {
		return errors.New("first name is required")
	}
	if u.LastName == "" {
		return errors.New("last name is required")
	}

	// Date of Birth should be valid and user should be at least 13
	now := time.Now()
	minDOB := now.AddDate(-13, 0, 0)
	if u.DateOfBirth.IsZero() || u.DateOfBirth.After(minDOB) {
		return errors.New("user must be at least 13 years old")
	}

	// TODO — add if https is prsent (allow empty, but if present, validate format)
	// if u.AvatarURL != "" {
	// 	urlRegex := regexp.MustCompile(`^https?://[^\s]+$`)
	// 	if !urlRegex.MatchString(u.AvatarURL) {
	// 		return errors.New("invalid avatar URL")
	// 	}
	// }

	if u.Nickname != "" {
		// Alphanumeric and underscore, 3–30 chars
		nickRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)
		if !nickRegex.MatchString(u.Nickname) {
			return errors.New("nickname must be alphanumeric/underscore and 3–30 chars")
		}
	}

	if len(u.AboutMe) > 500 {
		return errors.New("about me section is too long (max 500 characters)")
	}

	return nil
}

type UserModel struct {
	DB *sql.DB
}

func (um *UserModel) ValidateCredentials(user *User) error {
	if user.Email == "" && user.Nickname == "" {
		return errors.New("email or username must be provided")
	}

	var query string
	var arg string

	if user.Email != "" {
		query = `SELECT id, password_hash FROM users WHERE email = ? LIMIT 1`
		arg = user.Email
	} else {
		query = `SELECT id, password_hash FROM users WHERE nickname = ? LIMIT 1`
		arg = user.Nickname
	}

	if err := um.DB.QueryRow(query, arg).Scan(&user.ID, &user.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("User not found: %w", err)
		}
		return fmt.Errorf("query error: %w", err)
	}
	// Compare hashed password with input
	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(user.Password)); err != nil {
		return fmt.Errorf("Invalid password: %w", err)
	}

	return nil
}

// InsertUser inserts a new user into the database.
func (um *UserModel) InsertUser(user *User) error {
	query := `
		INSERT INTO users (
			email, password_hash, first_name, last_name, date_of_birth,
			avatar_url, nickname, about_me, is_public
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ? )
	`
	res, err := um.DB.Exec(query,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.AvatarURL,
		&user.Nickname,
		&user.AboutMe,
		&user.IsPublic,
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
func (um *UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET first_name = ?, last_name = ?, date_of_birth = ?, avatar_url = ?, nickname = ?,
			about_me = ?, is_public = ?
		WHERE id = ?
	`

	if _, err := um.DB.Exec(query,
		user.FirstName,
		user.LastName,
		user.DateOfBirth,
		user.AvatarURL,
		user.Nickname,
		user.AboutMe,
		user.IsPublic,
		user.ID,
	); err != nil {
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
func (um *UserModel) GetUserByID(user *User) error {
	query := `
		SELECT id, email, password_hash, first_name, last_name, date_of_birth,
		       avatar_url, nickname, about_me, is_public, created_at
		FROM users
		WHERE id = ?
	`

	err := um.DB.QueryRow(query, user.ID).Scan(
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
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no user with this id: %w", err) // Not found
		}
		return fmt.Errorf("get user by id: %w", err)
	}
	return nil
}
