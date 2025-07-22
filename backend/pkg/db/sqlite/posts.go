package sqlite

import (
	"database/sql"

	"social-network/pkg/models"
)

type PostRepository struct {
	DB *sql.DB
}

// Returns all posts for user ->
// group posts if is a member
// all public posts
// Private posts if is a follower
// almost_private if has access
// all posts if user is an author
func (repo *PostRepository) GetAll(userID string) ([]models.Post, error) {
	var posts []models.Post
	rows, err := repo.DB.Query("SELECT post_id , created_by, content, image  FROM posts WHERE (group_id IS NULL  AND visibility = 'PUBLIC') OR (group_id IS NULL AND visibility = 'PRIVATE' AND (SELECT COUNT() FROM followers WHERE posts.created_by = followers.user_id AND followers.follower_id = '"+userID+"') = 1 ) OR (group_id IS NULL  AND visibility = 'ALMOST_PRIVATE' AND (SELECT COUNT() FROM almost_private WHERE almost_private.post_id = posts.post_id AND almost_private.user_id = '"+userID+"')=1) OR  (group_id IS NULL  AND created_by = '"+userID+"' )ORDER BY created_at DESC;", userID)
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		var post models.Post
		rows.Scan(&post.ID, &post.AuthorID, &post.Content, &post.ImagePath)
		posts = append(posts, post)
	}
	return posts, nil
}

func (repo *PostRepository) GetUserPosts(userID, currentUserID string) ([]models.Post, error) {
	// get all posts that do not belong to group
	// get group posts only if current user alsa a member or admin
	var posts []models.Post
	rows, err := repo.DB.Query("SELECT post_id , created_by, content, image  FROM posts WHERE (group_id IS NULL  AND visibility = 'PUBLIC' AND created_by = '" + userID + "') OR (group_id IS NULL AND visibility = 'PRIVATE' AND created_by = '" + userID + "' AND (SELECT COUNT() FROM followers WHERE posts.created_by = followers.user_id AND followers.follower_id = '" + currentUserID + "') = 1 ) OR (group_id IS NULL  AND visibility = 'ALMOST_PRIVATE' AND created_by = '" + userID + "' AND (SELECT COUNT() FROM almost_private WHERE almost_private.post_id = posts.post_id AND almost_private.user_id = '" + currentUserID + "')=1) OR (group_id IS NULL  AND created_by = '" + userID + "' AND '" + userID + "' = '" + currentUserID + "') ORDER BY created_at DESC;")
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		var post models.Post
		rows.Scan(&post.ID, &post.AuthorID, &post.Content, &post.ImagePath)
		posts = append(posts, post)
	}
	return posts, nil
}

func (repo *PostRepository) GetGroupPosts(groupID string) ([]models.Post, error) {
	var posts []models.Post
	rows, err := repo.DB.Query("SELECT post_id , created_by, content, image  FROM posts WHERE group_id = ? ORDER BY created_at DESC;", groupID)
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		var post models.Post
		rows.Scan(&post.ID, &post.AuthorID, &post.Content, &post.ImagePath)
		posts = append(posts, post)
	}
	return posts, nil
}

func (repo *PostRepository) New(post models.Post) error {
	stmt, err := repo.DB.Prepare("INSERT INTO posts (post_id, group_id, created_by, content,image,visibility) values (?,(NULLIF(?,'')),?,?,?,?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(post.ID, post.GroupID, post.AuthorID, post.Content, post.ImagePath, post.Visibility); err != nil {
		return err
	}
	return nil
}

func (repo *PostRepository) SaveAccess(postId, userId string) error {
	stmt, err := repo.DB.Prepare("INSERT INTO almost_private (post_id, user_id) values (?,?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(postId, userId); err != nil {
		return err
	}
	return nil
}
