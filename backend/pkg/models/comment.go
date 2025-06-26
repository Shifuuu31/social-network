package models

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Comment struct {
	Id        int    `json:"id"`
	Post_id   int    `json:"post_id"`
	OwnerId   int    `json:"owner_id"`
	Content   string `json:"content"`
	Image     string `json:"image"`
	CreatedAt string `json:"created_at"`
}

type CommentModel struct {
	DB *sql.DB
}

type CommentFilter struct {
	Start    int `json:"start"`
	Ncomment int `json:"n_comment"`
}

func (app *CommentModel) GetComments(postId int) ([]Comment, error) {
	comments := []Comment{}
	comment := Comment{}
	rows, err := app.DB.Query("SELECT * FROM comments where post_id = ?", postId)
	//  rows.Scan(&comment.Id, &comment.OwnerId, &comment.Post_id, &comment.Content, &comment.Image, &comment.Created_at)
	if err != nil {
		log.Printf("Error fetching comments for post %d: %v", postId, err)
		return nil, fmt.Errorf("failed to fetch comments for post %d: %w", postId, err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {

		if err := rows.Scan(&comment.Id, &comment.OwnerId, &comment.Post_id, &comment.Content, &comment.Image, &comment.CreatedAt); err != nil {
			log.Printf("Error scanning comment: %v", err)
			continue
		}
		i++
		comments = append(comments, comment)
		// fmt.Println(comment)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over comments: %v", err)
		return nil, fmt.Errorf("error iterating over comments: %w", err)
	}
	return comments, nil

}

func ValidateComment(db *sql.DB, comment *Comment) (bool, int) {

	if comment.OwnerId <= 0 || comment.Post_id <= 0 {
		return false, 400
	}

	comment.Content = strings.TrimSpace(comment.Content)
	if len(comment.Content) > 300 || comment.Content == "" {
		return false, 400
	}

	var postOwnerID int
	var privacy string
	err := db.QueryRow("SELECT user_id, privacy FROM posts WHERE id = ?", comment.Post_id).Scan(&postOwnerID, &privacy)
	if err == sql.ErrNoRows {
		return false, 404
	} else if err != nil {
		fmt.Println("Error querying post:", err)
		return false, 500
	}

	var userExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", comment.OwnerId).Scan(&userExists)
	if err != nil {
		fmt.Println("Error checking user existence:", err)
		return false, 500
	}
	if !userExists {
		return false, 404
	}

	switch privacy {
	case "public":
		return true, 200

	case "almost_private":
		var isFollowing bool
		err = db.QueryRow(`
            SELECT EXISTS(
                SELECT 1 FROM followers
                WHERE follower_id = ? AND followee_id = ? AND status = 'accepted'
            )
        `, comment.OwnerId, postOwnerID).Scan(&isFollowing)
		if err != nil {
			return false, 500
		}
		return isFollowing, 200

	case "private":
		var isAllowed bool
		err = db.QueryRow(`
            SELECT EXISTS(
                SELECT 1 FROM post_privacy
                WHERE post_id = ? AND chosen_id = ?
            )
        `, comment.Post_id, comment.OwnerId).Scan(&isAllowed)
		if err != nil {
			return false, 500
		}
		return isAllowed, 200

	default:
		return false, 400 // Unknown privacy type
	}
}

func ParseCommentFromForm(r *http.Request, comment *Comment) int {
	var err error

	if ownerIdStr := r.FormValue("owner_id"); ownerIdStr != "" {
		if comment.OwnerId, err = strconv.Atoi(ownerIdStr); err != nil {
			return 400
		}
	}

	if PostIdStr := r.FormValue("group_id"); PostIdStr != "" {
		if comment.Post_id, err = strconv.Atoi(PostIdStr); err != nil && comment.Post_id <= 0 {
			return 400
		}
	}

	comment.Content = r.FormValue("content")

	return 200 // OK
}
