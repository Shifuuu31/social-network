package sqlite

import (
	"database/sql"

	"social-network/pkg/models"
)

type CommentRepository struct {
	DB *sql.DB
}

func (repo *CommentRepository) Get(postID string) ([]models.Comment, error) {
	var comments []models.Comment
	rows, err := repo.DB.Query("SELECT comment_id, created_by, content, image FROM comments WHERE post_id = ? ORDER BY created_at DESC;", postID)
	if err != nil {
		return comments, err
	}
	for rows.Next() {
		var comment models.Comment
		rows.Scan(&comment.ID, &comment.AuthorID, &comment.Content, &comment.ImagePath)
		comments = append(comments, comment)
	}
	return comments, nil
}

func (repo *CommentRepository) New(comment models.Comment) error {
	stmt, err := repo.DB.Prepare("INSERT INTO comments (comment_id, post_id, created_by, content,image) values (?,?,?,?,?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(comment.ID, comment.PostID, comment.AuthorID, comment.Content, comment.ImagePath); err != nil {
		return err
	}
	return nil
}
