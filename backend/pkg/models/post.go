package models

import (
	"database/sql"
	"time"
)

// Post represents a post in the social network.
type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	GroupID    int       `json:"group_id"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_path"`
	Privacy   string    `json:"privacy"` // "public", "followers", "selected" or "group"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostModel struct {
	DB *sql.DB
}

// TODO to be implemented
func (pm *PostModel) GetAllUserPosts(userId int) (posts []Post, err error) {

	return posts, nil
}
