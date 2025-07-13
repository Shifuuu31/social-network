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
	Image_url      string `json:"image_url"`
	AvatarURL      string `json:"avatar_url"` // User's profile image
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

	if pfl.Start < 0 {
		return errors.New("start cannot be negative")
	}

	if pfl.Type != "" {
		validTypes := map[string]bool{
			"feed":      true,
			"user":      true,
			"group":     true,
			"public":    true,
			"followers": true,
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

	fmt.Println("filter.Type: momo", filter.Type)
	fmt.Println("filter.Id: momo", filter.Id)

	// UNIFIED QUERY: Handle all privacy types in one query
	// This query will show posts based on each post's individual privacy setting
	query = `
		SELECT 
			posts.id, 
			posts.user_id, 
			users.nickname, 
			posts.group_id,
			posts.content,
			posts.image_url,
			posts.privacy, 
			posts.created_at,
			COALESCE(comment_counts.reply_count, 0) as replies,
			users.avatar_url
		FROM posts
		JOIN users ON posts.user_id = users.id
		LEFT JOIN (
			SELECT post_id, COUNT(*) as reply_count
			FROM comments
			GROUP BY post_id
		) comment_counts ON posts.id = comment_counts.post_id
		LEFT JOIN follow_request 
			ON follow_request.from_user_id = ? 
			AND follow_request.to_user_id = posts.user_id 
			AND follow_request.status = 'accepted'
		LEFT JOIN post_privacy_selected
			ON post_privacy_selected.post_id = posts.id 
			AND post_privacy_selected.user_id = ?
		WHERE
			-- Show user's own posts (regardless of privacy)
			posts.user_id = ?
			-- OR show public posts (visible to everyone)
			OR posts.privacy = 'public'
			-- OR show followers-only posts if user is a follower
			OR (posts.privacy = 'followers' AND follow_request.from_user_id IS NOT NULL)
			-- OR show private posts if user is specifically selected
			OR (posts.privacy = 'private' AND post_privacy_selected.user_id IS NOT NULL)
		ORDER BY posts.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err = pm.DB.Query(query, filter.Id, filter.Id, filter.Id, filter.NPost, filter.Start)

	if err != nil {
		fmt.Printf("ERROR: Query execution failed: %v\n", err)
		fmt.Printf("Query was: %s\n", query)
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
			&post.Image_url,
			&post.Privacy,
			&post.CreatedAt,
			&post.Replies,
			&post.AvatarURL,
		)

		if err != nil {
			fmt.Printf("ERROR: Failed to scan row %d: %v\n", postCount, err)
			return nil, err
		}

		postCount++
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("ERROR: Row iteration error: %v\n", err)
		return nil, err
	}

	fmt.Printf("DEBUG: First post avatar_url: %v\n", posts[0].AvatarURL)
	return posts, nil
}

func ParsePostFromForm(r *http.Request, post *Post) int {
	var err error

	// Parse owner_id (required)
	ownerIdStr := r.FormValue("owner_id")
	if ownerIdStr == "" {
		// Try alternative field name
		ownerIdStr = r.FormValue("ownerId")
	}
	if ownerIdStr == "" {
		return 400 // owner_id is required
	}
	if post.OwnerId, err = strconv.Atoi(ownerIdStr); err != nil {
		return 400
	}

	// Parse group_id (optional)
	groupIdStr := r.FormValue("group_id")
	if groupIdStr == "" {
		// Try alternative field name
		groupIdStr = r.FormValue("groupId")
	}
	if groupIdStr != "" {
		if post.GroupId, err = strconv.Atoi(groupIdStr); err != nil {
			return 400
		}
	} else {
		post.GroupId = 0 // Default to 0 (no group)
	}

	// Parse content and privacy
	post.Content = r.FormValue("content")
	post.Privacy = r.FormValue("privacy")

	// Set default privacy if not provided
	if post.Privacy == "" {
		post.Privacy = "public"
	}

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
