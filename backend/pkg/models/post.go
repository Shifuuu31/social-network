package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

// Post represents a post in the social network.
type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	GroupID   int       `json:"group_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url"`
	Privacy   string    `json:"privacy"` // "public", "followers", "selected" or "group"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostModel struct {
	DB *sql.DB
}

func ValidatePost(post *Post) error {
	if post == nil {
		return errors.New("post is nil")
	}

	post.Title = strings.TrimSpace(post.Title)
	if post.Title == "" {
		return errors.New("post.Title cannot be empty or whitespace")
	}
	if len(post.Title) > 70 {
		return errors.New("post.Title must be at most 70 characters long")
	}

	post.Content = strings.TrimSpace(post.Content)
	if post.Content == "" {
		return errors.New("post.Content cannot be empty or whitespace")
	}
	if len(post.Content) > 1000 {
		return errors.New("post.Content must be at most 1000 characters long")
	}

	return nil
}

func (pm *PostModel) InsertPost(post Post) error {
	postQuery := `
		INSERT INTO posts (user_id, title, content, privacy)
		VALUES (?, ?, ?)
	`

	_, err := pm.DB.Exec(postQuery, post.UserID, post.Title, post.Content, post.Privacy)
	if err != nil {
		return err
	}

	// postID, err := res.LastInsertId()
	// if err != nil {
	// 	return err
	// }

	return nil
}

// TODO to be implemented
func (pm *PostModel) GetAllUserPosts(userId int) (posts []Post, err error) {
	return posts, nil
}
