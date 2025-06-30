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

func (post *Post) Validate() error {
	if post.OwnerId <= 0 {
		return errors.New("owner_id must be a positive integer")
	}

	if post.GroupId < 0 {
		return errors.New("group_id cannot be negative")
	}

	post.Content = strings.TrimSpace(post.Content)
	if post.Content == "" {
		return errors.New("content cannot be empty")
	}
	if len(post.Content) > 1000 {
		return errors.New("content exceeds 1000 character limit")
	}

	// Validate privacy value
	validPrivacyLevels := map[string]bool{
		"public":   true,
		"selected": true,
		"private":  true,
		"group":    true,
	}
	if !validPrivacyLevels[post.Privacy] {
		return errors.New("invalid privacy value; must be 'public', 'selected', 'group', or 'private'")
	}

	// Validate ChosenUsersIds
	if post.Privacy == "private" {
		if len(post.ChosenUsersIds) == 0 {
			return errors.New("chosen_users_ids required for private posts")
		}
		// Check for duplicates and valid IDs
		seen := make(map[int]bool)
		for _, id := range post.ChosenUsersIds {
			if id <= 0 {
				return errors.New("all chosen_users_ids must be positive integers")
			}
			if seen[id] {
				return errors.New("duplicate user ID found in chosen_users_ids")
			}
			seen[id] = true
		}
	} else {
		if len(post.ChosenUsersIds) > 0 {
			return errors.New("chosen_users_ids should only be set for private posts")
		}
	}

	return nil
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

func (pfl *PostFilter) Validate() error {
	pfl.Type = strings.TrimSpace(pfl.Type)

	if pfl.NPost <= 0 {
		pfl.NPost = 10
	}

	if pfl.Id < 0 {
		return errors.New("id cannot be negative")
	}

	if pfl.Start < 0 {
		return errors.New("start cannot be negative")
	}

	if pfl.Type != "" {
		validTypes := map[string]bool{
			"feed":   true,
			"user":   true,
			"group":  true,
			"public": true,
		}
		if !validTypes[pfl.Type] {
			return errors.New("type must be one of: feed, user, group, public")
		}
	}

	if pfl.NPost > 100 {
		return errors.New("n_post cannot exceed 100")
	}

	return nil
}

func (pm *PostModel) GetPosts(filter *PostFilter) (posts []Post, err error) {
	var query string
	var rows *sql.Rows
	var args []any
	
	fmt.Printf("Filter: Type=%s, Id=%v, Start=%v, NPost=%v\n", filter.Type, filter.Id, filter.Start, filter.NPost)

	switch filter.Type {
	case "group":
		query = `
		SELECT posts.id, posts.user_id, users.nickname, posts.group_id, 
			   posts.content, posts.image_path as image, posts.privacy, posts.created_at,
			   COALESCE(comment_counts.reply_count, 0) as replies
		FROM posts
		JOIN users ON posts.user_id = users.id
		LEFT JOIN (
			SELECT post_id, COUNT(*) as reply_count
			FROM comments
			GROUP BY post_id
		) comment_counts ON CAST(posts.id AS TEXT) = comment_counts.post_id
		WHERE posts.group_id = ? AND posts.privacy = 'group' AND posts.id > ?
		ORDER BY posts.id ASC
		LIMIT ?`
		args = append(args, filter.Id, filter.Start, filter.NPost)
		
	case "feed":
		query = `
		SELECT posts.id, posts.user_id, users.nickname, posts.group_id, 
			   posts.content, posts.image_path as image, posts.privacy, posts.created_at,
			   COALESCE(comment_counts.reply_count, 0) as replies
			   FROM posts
			   JOIN users ON posts.user_id = users.id
			   LEFT JOIN follow_request 
			   ON follow_request.to_user_id = posts.user_id 
			   AND follow_request.from_user_id = ? 
			   AND follow_request.status = 'accepted'
			   LEFT JOIN post_privacy_selected
			   ON post_privacy_selected.post_id = posts.id 
			   AND post_privacy_selected.user_id = ?
			   LEFT JOIN (
				SELECT post_id, COUNT(*) as reply_count
				FROM comments
			GROUP BY post_id
			) comment_counts ON posts.id = comment_counts.post_id
			WHERE NOT (posts.group_id IS NOT NULL AND posts.privacy = 'group')
			AND (
				posts.privacy = 'public'
				OR (posts.privacy = 'followers' AND follow_request.id IS NOT NULL)
				OR (posts.privacy = 'selected' AND post_privacy_selected.id IS NOT NULL)
				)
				AND posts.id > ?
				ORDER BY posts.id DESC
				LIMIT ?`
				args = append(args, filter.Id, filter.Id, filter.Start, filter.NPost)
				
	case "user":
		query = `
		SELECT posts.id, posts.user_id, users.nickname, posts.group_id, 
			   posts.content, posts.image_path as image, posts.privacy, posts.created_at,
			   COALESCE(comment_counts.reply_count, 0) as replies
		FROM posts
		JOIN users ON posts.user_id = users.id
		LEFT JOIN (
			SELECT post_id, COUNT(*) as reply_count
			FROM comments
			GROUP BY post_id
		) comment_counts ON posts.id = comment_counts.post_id
		WHERE posts.user_id = ?`
		args = append(args, filter.Id)
	}

	rows, err = pm.DB.Query(query, args...)
	if err != nil {
		fmt.Printf("Query error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	postCount := 0
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
		if err != nil {
			fmt.Printf("Scan error: %v\n", err)
			return nil, err
		}
		posts = append(posts, post)
		postCount++
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Rows error: %v\n", err)
		return nil, err
	}

	return posts, nil
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
