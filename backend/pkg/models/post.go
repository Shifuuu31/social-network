package models

import (
	"database/sql"
	"time"
)

// Post represents a post in the social network.
type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url"`
	Privacy   string    `json:"privacy"` // "public", "followers", or "selected"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostModel struct {
	DB *sql.DB
}

func (pm *PostModel) GetAllUserPosts(userId int) (posts []Post, err error) {

	// TODO to be implemented
	return posts, nil
}
