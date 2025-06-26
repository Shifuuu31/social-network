package models

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Post struct {
	Replies        int    `json:"replies"`
	Owner          string `json:"owner"`
	Id             int    `json:"id"`
	OwnerId        int    `json:"owner_id"`
	GroupId        int    `json:"group_id"`
	Content        string `json:"content"`
	Image          string `json:"image"`
	Privacy        string `json:"privacy"` // [public', 'almost_private', 'private']
	CreatedAt      string `json:"created_at"`
	ChosenUsersIds []int  `json:"chosen_users_ids"`
}
type PostModel struct {
	DB *sql.DB
}

type PostFilter struct {
	Id    int    `json:"id"`
	Type  string `json:"type"`
	Start int    `json:"start"`
	NPost int    `json:"n_post"`
}

func (pm *PostModel) GetPosts(filter *PostFilter) (posts []Post, err error) {
	var query string
	var rows *sql.Rows
	fmt.Println("PostFilter:", filter)

	switch filter.Type {
	case "group":
		query = `
            SELECT posts.id, posts.user_id, users.nickname, posts.group_id, 
                   posts.content, posts.image, posts.privacy, posts.created_at,
                   COALESCE(comment_counts.reply_count, 0) as replies
            FROM posts
            JOIN users ON posts.user_id = users.id
            LEFT JOIN (
                SELECT post_id, COUNT(*) as reply_count
                FROM comments
                GROUP BY post_id
            ) comment_counts ON CAST(posts.id as TEXT) = comment_counts.post_id
            WHERE posts.group_id = ? AND posts.privacy = '' AND posts.id > ?
            ORDER BY posts.id ASC
            LIMIT ?`
		rows, err = pm.DB.Query(query, filter.Id, filter.Start, filter.NPost)

	case "privacy":
		query = `
            SELECT posts.id, posts.user_id, users.nickname, posts.group_id, 
                   posts.content, posts.image, posts.privacy, posts.created_at,
                   COALESCE(comment_counts.reply_count, 0) as replies
            FROM posts
            JOIN users ON posts.user_id = users.id
            LEFT JOIN follows 
                ON follows.followee_id = posts.user_id AND follows.follower_id = ? AND follows.status = 'accepted'
            LEFT JOIN post_privacy
                ON post_privacy.post_id = posts.id AND post_privacy.chosen_id = ?
            LEFT JOIN (
                SELECT post_id, COUNT(*) as reply_count
                FROM comments
                GROUP BY post_id
            ) comment_counts ON CAST(posts.id as TEXT) = comment_counts.post_id
            WHERE NOT (posts.group_id IS NOT NULL AND posts.privacy = '')
              AND (
                posts.privacy IN ('public')
                OR (posts.privacy = 'almost_private' AND follows.follower_id IS NOT NULL)
                OR (posts.privacy = 'private' AND post_privacy.chosen_id IS NOT NULL)
              )
              AND posts.id > ?
            ORDER BY posts.id DESC
            LIMIT ?`
		rows, err = pm.DB.Query(query, filter.Id, filter.Id, filter.Start, filter.NPost)

	case "single":
		query = `
            SELECT posts.id, posts.user_id, users.nickname, posts.group_id, 
                   posts.content, posts.image, posts.privacy, posts.created_at,
                   COALESCE(comment_counts.reply_count, 0) as replies
            FROM posts
            JOIN users ON posts.user_id = users.id
            LEFT JOIN (
                SELECT post_id, COUNT(*) as reply_count
                FROM comments
                GROUP BY post_id
            ) comment_counts ON CAST(posts.id as TEXT) = comment_counts.post_id
            WHERE posts.id = ?`
		rows, err = pm.DB.Query(query, filter.Id)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.Id,
			&post.OwnerId,
			&post.Owner,
			&post.GroupId,
			&post.Content,
			&post.Image,
			&post.Privacy,
			&post.CreatedAt,
			&post.Replies,
		)
		fmt.Println("Post:", post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func ValidatePost(post *Post) int {
	// Validate OwnerId
	if post.OwnerId <= 0 {
		return 400
	}

	if post.GroupId < 0 {
		return 400
	}

	// Validate Content
	post.Content = strings.TrimSpace(post.Content)
	if post.Content == "" || len(post.Content) > 1000 {
		return 400
	}

	// Validate Privacy
	validPrivacyLevels := map[string]bool{
		"public":         true,
		"almost_private": true,
		"private":        true,
	}
	if !validPrivacyLevels[post.Privacy] {
		return 400
	}

	// Validate ChosenUsersIds for private posts
	if post.Privacy == "private" {
		if len(post.ChosenUsersIds) == 0 {
			fmt.Println("ChosenUsersIds is empty for private post")
			return 400
		}
		// Check for duplicate user IDs and validate positive integers
		userIdMap := make(map[int]bool)
		for _, id := range post.ChosenUsersIds {
			if id <= 0 {
				return 400
			}
			if userIdMap[id] {
				return 400 // Duplicate user ID
			}
			userIdMap[id] = true
		}
	} else if len(post.ChosenUsersIds) > 0 {
		return 400 // chosen_users_ids should only be specified for private posts
	}

	return 200 // OK
}

func ParsePostFromForm(r *http.Request, post *Post) int {
	var err error

	// Parse owner_id
	if ownerIdStr := r.FormValue("owner_id"); ownerIdStr != "" {
		if post.OwnerId, err = strconv.Atoi(ownerIdStr); err != nil {
			return 400
		}
	}

	// Parse group_id
	if groupIdStr := r.FormValue("group_id"); groupIdStr != "" {
		if post.GroupId, err = strconv.Atoi(groupIdStr); err != nil {
			return 400
		}
	}

	// Parse content and privacy
	post.Content = r.FormValue("content")
	post.Privacy = r.FormValue("privacy")

	// Parse chosen_users_ids for private posts
	if chosenUsersStr := r.FormValue("chosen_users_ids"); chosenUsersStr != "" {
		userIds := strings.Split(chosenUsersStr, ",")
		post.ChosenUsersIds = make([]int, 0, len(userIds))

		for _, idStr := range userIds {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			if id, err := strconv.Atoi(idStr); err != nil {
				return 400
			} else {
				post.ChosenUsersIds = append(post.ChosenUsersIds, id)
			}
		}
	}

	return 200 // OK
}

func ValidatePostFilter(filter *PostFilter) error {
	var errs []string

	// Validate Type if provided
	if filter.Type != "" {
		validTypes := map[string]bool{
			"feed":   true,
			"user":   true,
			"group":  true,
			"public": true,
		}
		if !validTypes[filter.Type] {
			errs = append(errs, "type must be one of: feed, user, group, public")
		}
	}

	if filter.Start < 0 {
		errs = append(errs, "start cannot be negative")
	}

	if filter.NPost <= 0 {
		filter.NPost = 10
	}
	if filter.NPost > 100 {
		errs = append(errs, "n_post cannot exceed 100")
	}

	if filter.Id < 0 {
		errs = append(errs, "id cannot be negative")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}
