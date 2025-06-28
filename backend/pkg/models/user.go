package models

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const UserIDKey = contextKey("userID")

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	AvatarPath   string    `json:"avatar_path"`
	Nickname     string    `json:"nickname"`
	AboutMe      string    `json:"about_me"`
	IsPublic     bool      `json:"is_public"`
	PasswordHash []byte    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

// Validate checks all required fields and ensures correct formats
func Validate(r *http.Request) (User, error) {
	// Parse form data (works for both application/x-www-form-urlencoded and multipart/form-data)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return User{}, fmt.Errorf("failed to parse form: %w", err)
	}

	user := User{
		Email:      strings.TrimSpace(r.FormValue("email")),
		Password:   r.FormValue("password"),
		FirstName:  strings.TrimSpace(r.FormValue("first_name")),
		LastName:   strings.TrimSpace(r.FormValue("last_name")),
		AvatarPath: strings.TrimSpace(r.FormValue("avatar_path")),
		Nickname:   strings.TrimSpace(r.FormValue("nickname")),
		AboutMe:    strings.TrimSpace(r.FormValue("about_me")),
	}

	fmt.Println(user)
	emailExp := regexp.MustCompile(`^[a-zA-Z0-9._%%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if user.Email == "" || !emailExp.MatchString(user.Email) {
		return User{}, errors.New("invalid email format")
	}

	if len(user.Password) < 8 || len(user.Password) > 64 {
		return User{}, errors.New("password must be between 8 and 64 characters")
	}

	if user.FirstName == "" {
		return User{}, errors.New("first name is required")
	}

	if user.LastName == "" {
		return User{}, errors.New("last name is required")
	}

	dobStr := r.FormValue("date_of_birth")
	if dobStr == "" {
		return User{}, errors.New("date of birth is required")
	}
	dob, err := time.Parse("2006-01-02", dobStr) // e.g., "1990-04-23"
	if err != nil {
		return User{}, errors.New("invalid date of birth format (expected YYYY-MM-DD)")
	}
	minDOB := time.Now().AddDate(-13, 0, 0)
	if dob.After(minDOB) {
		return User{}, errors.New("user must be at least 13 years old")
	}
	user.DateOfBirth = dob

	// Optional: validate URL format
	// if user.AvatarPath != "" {
	// 	urlRegex := regexp.MustCompile(`^https?://[^\s]+$`)
	// 	if !urlRegex.MatchString(user.AvatarPath) {
	// 		return User{},errors.New("invalid avatar URL")
	// 	}
	// }

	if user.Nickname != "" {
		nickRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)
		if !nickRegex.MatchString(user.Nickname) {
			return User{}, errors.New("nickname must be alphanumeric/underscore and 3–30 characters")
		}
	}
	if len(user.AboutMe) > 500 {
		return User{}, errors.New("about me section is too long (max 500 characters)")
	}

	return user, nil
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

// Insert inserts a new user into the database.
func (um *UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (
			email, password_hash, first_name, last_name, date_of_birth,
			avatar_path, nickname, about_me, is_public
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ? )
	`
	res, err := um.DB.Exec(query,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.AvatarPath,
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
		SET first_name = ?, last_name = ?, date_of_birth = ?, avatar_path = ?, nickname = ?,
			about_me = ?, is_public = ? = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	if _, err := um.DB.Exec(query,
		user.FirstName,
		user.LastName,
		user.DateOfBirth,
		user.AvatarPath,
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
func (um *UserModel) Delete(id int) error {
	query := `DELETE FROM users WHERE id = ?`

	if _, err := um.DB.Exec(query, id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

// GetUserByID retrieves a user by ID.
func (um *UserModel) GetUserByID(user *User) error {
	query := `
		SELECT id, email, password_hash, first_name, last_name, date_of_birth,
		       avatar_path, nickname, about_me, is_public, created_at
		FROM users
		WHERE id = ?
	`

	if err := um.DB.QueryRow(query, user.ID).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.AvatarPath,
		&user.Nickname,
		&user.AboutMe,
		&user.IsPublic,
		&user.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no user with this id: %w", err) // Not found
		}
		return fmt.Errorf("get user by id: %w", err)
	}
	return nil
}
