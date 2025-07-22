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
	ImageUUID      string `json:"image_uuid"`
	Privacy        string `json:"privacy"` // [public', 'almost_private', 'private']
	CreatedAt      string `json:"created_at"`
	ChosenUsersIds []int  `json:"chosen_users_ids"`
	Image_url      string `json:"image_url"`
	// Title string `json:"title"`
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
		"public":         true,
		"almost_private": true,
		"private":        true,
		"group":          true,
	}
	if !validPrivacyLevels[post.Privacy] {
		return errors.New("invalid privacy value; must be 'public', 'almost_private', or 'private'")
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

	// if pfl.Start < 0 {
	// 	return errors.New("start cannot be negative")
	// }

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

	// Debug: Print the filter details
	fmt.Printf("=== DEBUG GetPosts ===\n")
	fmt.Printf("Filter Type: %s\n", filter.Type)
	fmt.Printf("Filter ID: %d\n", filter.Id)
	fmt.Printf("Filter Start: %d\n", filter.Start)
	fmt.Printf("Filter NPost: %d\n", filter.NPost)
	if filter.Start == -1 {
		switch filter.Type {
		case "group":
			// For groups, get the max ID to start from the newest post
			err = pm.DB.QueryRow(`SELECT IFNULL(MAX(id), 0) FROM posts WHERE group_id = ?`, filter.Id).Scan(&filter.Start)
			if err != nil {
				fmt.Printf("ERROR: Failed to get max ID for group %d: %v\n", filter.Id, err)
				return nil, err
			}
		case "privacy", "single":
			err = pm.DB.QueryRow(`SELECT IFNULL(MAX(id), 0) FROM posts`).Scan(&filter.Start)
		case "public":
			filter.Start = 0
		}
	}
	fmt.Println("Filter Start after adjustment:", filter.Start)

	switch filter.Type {

	case "group":
		query = `
            SELECT posts.id, posts.user_id, 
                   COALESCE(NULLIF(users.nickname, ''), users.first_name || ' ' || users.last_name) as owner, 
                   posts.group_id, 
                   posts.content, COALESCE(posts.image_uuid, '') as image_uuid, posts.privacy, posts.created_at,
                   COALESCE(comment_counts.reply_count, 0) as replies
            FROM posts
            JOIN users ON posts.user_id = users.id
            LEFT JOIN (
                SELECT post_id, COUNT(*) as reply_count
                FROM comments
                GROUP BY post_id
            ) comment_counts ON CAST(posts.id as TEXT) = comment_counts.post_id
            WHERE posts.group_id = ? AND posts.privacy = 'group' AND posts.id <= ?
            ORDER BY posts.id DESC
            LIMIT ?`

		fmt.Printf("Executing GROUP query with params: groupId=%d, start=%d, limit=%d\n",
			filter.Id, filter.Start, filter.NPost)
		rows, err = pm.DB.Query(query, filter.Id, filter.Start, filter.NPost)

	case "privacy":
		query = `
            SELECT posts.id, posts.user_id, users.nickname, posts.group_id, 
                   posts.content, posts.image_uuid, posts.privacy, posts.created_at,
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

		fmt.Printf("Executing PRIVACY query with params: userId=%d, userId=%d, start=%d, limit=%d\n",
			filter.Id, filter.Id, filter.Start, filter.NPost)
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

		fmt.Printf("Executing SINGLE query with params: postId=%d\n", filter.Id)
		rows, err = pm.DB.Query(query, filter.Id)

	case "public":
		// FIXED: Use OFFSET instead of WHERE posts.id > ?
		query = `
			SELECT 
				posts.id, 
				posts.user_id, 
				users.nickname, 
				posts.group_id,
				posts.content,
 				posts.privacy, 
				posts.created_at,
				COALESCE(comment_counts.reply_count, 0) as replies
			FROM posts
			JOIN users ON posts.user_id = users.id
			LEFT JOIN (
				SELECT post_id, COUNT(*) as reply_count
				FROM comments
				GROUP BY post_id
			) comment_counts ON posts.id = comment_counts.post_id
			WHERE posts.privacy = 'public'
			ORDER BY posts.created_at DESC
			LIMIT ? OFFSET ?`

		fmt.Printf("Executing PUBLIC query with params: limit=%d, offset=%d\n",
			filter.NPost, filter.Start)
		fmt.Printf("Full query: %s\n", query)
		rows, err = pm.DB.Query(query, filter.NPost, filter.Start)

		// Additional debug: Check row count before processing
		if err == nil {
			fmt.Printf("Query executed without error, starting to process rows...\n")
		}

	default:
		fmt.Printf("ERROR: Unknown filter type: %s\n", filter.Type)
		return nil, fmt.Errorf("unknown filter type: %s", filter.Type)
	}

	// Debug: Check for query errors
	if err != nil {
		fmt.Printf("ERROR: Query execution failed: %v\n", err)
		fmt.Printf("Query was: %s\n", query)
		return nil, err
	}

	fmt.Printf("Query executed successfully\n")
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
			&post.ImageUUID,
			&post.Privacy,
			&post.CreatedAt,
			&post.Replies,
		)
		if err != nil {
			fmt.Printf("ERROR: Failed to scan row %d: %v\n", postCount, err)
			return nil, err
		}

		postCount++
		fmt.Printf("Post %d: ID=%d, Owner=%s, Content=%.50s...\n",
			postCount, post.Id, post.Owner, post.Content)

		posts = append(posts, post)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		fmt.Printf("ERROR: Row iteration error: %v\n", err)
		return nil, err
	}

	fmt.Printf("=== DEBUG COMPLETE ===\n")
	fmt.Printf("Total posts retrieved: %d\n", len(posts))
	fmt.Printf("=======================\n")

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
