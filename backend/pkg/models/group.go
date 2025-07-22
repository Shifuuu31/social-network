package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Group struct {
	ID          int       `json:"id"`
	CreatorID   int       `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImgUUID     string    `json:"image_uuid"`
	IsMember    string    `json:"is_member"` // Indicates if the user is a member of the group
	MemberCount int       `json:"member_count"`
	CreatedAt   time.Time `json:"created_at"`
}

func (g *Group) Validate() error {
	if g.CreatorID <= 0 {
		return errors.New("creator_id must be a positive integer")
	}

	g.Title = strings.TrimSpace(g.Title)
	if g.Title == "" {
		return errors.New("title cannot be empty")
	}
	if len(g.Title) > 100 {
		return errors.New("title cannot exceed 100 characters")
	}

	g.Description = strings.TrimSpace(g.Description)
	if g.Description == "" {
		return errors.New("description cannot be empty")
	}
	if len(g.Description) > 500 {
		return errors.New("description cannot exceed 500 characters")
	}

	return nil
}

type GroupModel struct {
	DB *sql.DB
}

// InsertGroup inserts a new group into the database.
func (gm *GroupModel) Insert(group *Group) error {
	query := `
		INSERT INTO groups (
			creator_id, title, description, image_uuid
		) VALUES (?, ?, ?, ?)
	`

	res, err := gm.DB.Exec(query,
		group.CreatorID,
		group.Title,
		group.Description,
		group.ImgUUID, // <-- add this
	)
	if err != nil {
		return fmt.Errorf("insert group: %w", err)
	}

	lastID, err := res.LastInsertId()
	if err == nil {
		group.ID = int(lastID)
	}

	return nil
}

func (gm *GroupModel) Delete(group *Group) error {
	res, err := gm.DB.Exec(`DELETE FROM groups WHERE id = ? AND creator_id = ?`, group.ID, group.CreatorID)
	if err != nil {
		return fmt.Errorf("delete group: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking delete result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no group deleted: not found or unauthorized")
	}
	return nil
}

type SearchPayload struct {
	Keyword    string `json:"keyword"`
	Start      int    `json:"start"`
	NumOfItems int    `json:"n_items"`
}

func (gm *GroupModel) SearchGroups(search *SearchPayload) ([]*Group, error) {
	if search.Start == -1 {
		err := gm.DB.QueryRow(`SELECT MAX(id) FROM groups WHERE title LIKE ? OR description LIKE ?`, "%"+search.Keyword+"%", "%"+search.Keyword+"%").Scan(&search.Start)
		if err != nil {
			return nil, fmt.Errorf("get max group id: %w", err)
		}
	}

	query := `
		SELECT g.id, g.creator_id, g.title, g.description, g.image_uuid, g.created_at,
       COUNT(m.id) AS member_count
		FROM groups g
		LEFT JOIN group_members m ON g.id = m.group_id AND m.status = 'member'
		WHERE (g.title LIKE ? OR g.description LIKE ?) AND g.id <= ?
		GROUP BY g.id
		ORDER BY g.id DESC
		LIMIT ?
	`

	pattern := "%" + search.Keyword + "%"
	rows, err := gm.DB.Query(query, pattern, pattern, search.Start, search.NumOfItems)
	if err != nil {
		return nil, fmt.Errorf("search groups: %w", err)
	}
	defer rows.Close()

	var groups []*Group
	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.ID, &g.CreatorID, &g.Title, &g.Description, &g.ImgUUID, &g.CreatedAt, &g.MemberCount); err != nil {
			return nil, fmt.Errorf("scan group: %w", err)
		}
		groups = append(groups, &g)
	}

	return groups, nil
}

func (gm *GroupModel) GetGroupByID(group *Group) error {
	query := `
		SELECT g.id, g.creator_id, g.title, g.description,  g.image_uuid, g.created_at,
		       COUNT(m.id) AS member_count
		FROM groups g
		LEFT JOIN group_members m ON g.id = m.group_id AND m.status = 'member'
		WHERE g.id = ?
		GROUP BY g.id
	`

	if err := gm.DB.QueryRow(query, group.ID).Scan(
		&group.ID, &group.CreatorID, &group.Title, &group.Description, &group.ImgUUID, &group.CreatedAt, &group.MemberCount,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("group not found or unauthorized")
		}
		return fmt.Errorf("get group by id: %w", err)
	}
	return nil
}

// IsUserCreator checks if the user is the creator of the given group.
func (gm *GroupModel) IsUserCreator(groupID, userID int) error {
	query := `
		SELECT COUNT(*)
		FROM groups
		WHERE id = ? AND creator_id = ?
	`
	var count int

	if err := gm.DB.QueryRow(query, groupID, userID).Scan(&count); err != nil || count <= 0 {
		return fmt.Errorf("check user is creator: %w", err)
	}

	return nil
}

type GroupsPayload struct {
	UserID     string `json:"user_id"`
	Start      int    `json:"start"`
	NumOfItems int    `json:"n_items"`
	Type       string `json:"type"`
	Search     string `json:"search"`
}

func (gm *GroupModel) GetGroups(filter *GroupsPayload) ([]*Group, error) {
	var (
		query  string
		args   []any
		userID int
		err    error
	)

	// Prepare search condition
	searchCondition := ""
	if filter.Search != "" {
		searchCondition = "AND (g.title LIKE ? OR g.description LIKE ?)"
	}

	switch filter.Type {
	case "user":
		// "My Groups" - groups where user has any interaction (member, requested, invited, or creator)
		var maxID sql.NullInt64
		maxIDQuery := `
			SELECT MAX(g.id) 
			FROM groups g 
			LEFT JOIN group_members gm ON g.id = gm.group_id 
			WHERE (g.creator_id = ? OR gm.user_id = ?)
		`
		maxIDArgs := []any{userID, userID}

		if filter.Search != "" {
			maxIDQuery += " AND (g.title LIKE ? OR g.description LIKE ?)"
			searchPattern := "%" + filter.Search + "%"
			maxIDArgs = append(maxIDArgs, searchPattern, searchPattern)
		}

		if err := gm.DB.QueryRow(maxIDQuery, maxIDArgs...).Scan(&maxID); err != nil {
			return []*Group{}, fmt.Errorf("get max group id: %w", err)
		}
		if maxID.Valid {
			filter.Start = int(maxID.Int64)
		} else {
			filter.Start = 0
		}

		query = `
			SELECT g.id, g.creator_id, g.title, g.description, g.image_uuid, g.created_at,
				COUNT(DISTINCT CASE WHEN m.status = 'member' THEN m.id END) AS member_count
			FROM groups g
			LEFT JOIN group_members m ON g.id = m.group_id AND m.status = 'member'
			LEFT JOIN group_members gm ON g.id = gm.group_id
			WHERE (g.creator_id = ? OR gm.user_id = ?) AND g.id <= ?
		` + searchCondition + `
			GROUP BY g.id
			ORDER BY g.id DESC
			LIMIT ?
		`
		args = []any{userID, userID, filter.Start}
		if filter.Search != "" {
			searchPattern := "%" + filter.Search + "%"
			args = append(args, searchPattern, searchPattern)
		}
		args = append(args, filter.NumOfItems)

	case "all":
		// "Explore" - groups where user has no interaction (not creator, no record in group_members)
		var maxID sql.NullInt64
		maxIDQuery := `
			SELECT MAX(id) FROM groups 
			WHERE creator_id != ? AND id NOT IN (
				SELECT group_id FROM group_members WHERE user_id = ?
			)
		`
		maxIDArgs := []any{userID, userID}

		if filter.Search != "" {
			maxIDQuery += " AND (title LIKE ? OR description LIKE ?)"
			searchPattern := "%" + filter.Search + "%"
			maxIDArgs = append(maxIDArgs, searchPattern, searchPattern)
		}

		if err := gm.DB.QueryRow(maxIDQuery, maxIDArgs...).Scan(&maxID); err != nil {
			fmt.Println("Error getting max group id:", err)
			return []*Group{}, fmt.Errorf("get max group id: %w", err)
		}
		if maxID.Valid {
			filter.Start = int(maxID.Int64)
		} else {
			filter.Start = 0
		}

		query = `
			SELECT 
				g.id, g.creator_id, g.title, g.description, g.image_uuid, g.created_at,
				COUNT(m.id) AS member_count
			FROM groups g
			LEFT JOIN group_members m 
				ON g.id = m.group_id AND m.status = 'member'
			WHERE g.creator_id != ? AND g.id NOT IN (
				SELECT group_id FROM group_members WHERE user_id = ?
			) AND g.id <= ?
		` + searchCondition + `
			GROUP BY g.id
			ORDER BY g.id DESC
			LIMIT ?
		`
		args = []any{userID, userID, filter.Start}
		if filter.Search != "" {
			searchPattern := "%" + filter.Search + "%"
			args = append(args, searchPattern, searchPattern)
		}
		args = append(args, filter.NumOfItems)

	// case "not_joined":
	// 	var maxID sql.NullInt64
	// 	maxIDQuery := `
	// 		SELECT MAX(id) FROM groups
	// 		WHERE id NOT IN (
	// 			SELECT group_id FROM group_members WHERE user_id = ?
	// 		)
	// 	`
	// 	maxIDArgs := []any{userID}

	// 	if filter.Search != "" {
	// 		maxIDQuery += " AND (title LIKE ? OR description LIKE ?)"
	// 		searchPattern := "%" + filter.Search + "%"
	// 		maxIDArgs = append(maxIDArgs, searchPattern, searchPattern)
	// 	}

	// 	if err := gm.DB.QueryRow(maxIDQuery, maxIDArgs...).Scan(&maxID); err != nil {
	// 		fmt.Println("Error getting max group id:", err)
	// 		return []*Group{}, fmt.Errorf("get max group id: %w", err)
	// 	}
	// 	if maxID.Valid {
	// 		filter.Start = int(maxID.Int64)
	// 	} else {
	// 		filter.Start = 0
	// 	}

	// 	query = `
	// 		SELECT
	// 			g.id, g.creator_id, g.title, g.description, g.image_uuid, g.created_at,
	// 			COUNT(m.id) AS member_count
	// 		FROM groups g
	// 		LEFT JOIN group_members m
	// 			ON g.id = m.group_id AND m.status = 'member'
	// 		WHERE g.id NOT IN (
	// 			SELECT group_id FROM group_members WHERE user_id = ?
	// 		) AND g.id <= ?
	// 	` + searchCondition + `
	// 		GROUP BY g.id
	// 		ORDER BY g.id DESC
	// 		LIMIT ?
	// 	`
	// 	args = []any{userID, filter.Start}
	// 	if filter.Search != "" {
	// 		searchPattern := "%" + filter.Search + "%"
	// 		args = append(args, searchPattern, searchPattern)
	// 	}
	// 	args = append(args, filter.NumOfItems)

	default:
		return []*Group{}, fmt.Errorf("invalid group type: must be 'user', 'all', or 'not_joined', got: %s", filter.Type)
	}

	rows, err := gm.DB.Query(query, args...)
	if err != nil {
		return []*Group{}, fmt.Errorf("get groups (%s): %w", filter.Type, err)
	}
	defer rows.Close()

	var groups []*Group
	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.ID, &g.CreatorID, &g.Title, &g.Description, &g.ImgUUID, &g.CreatedAt, &g.MemberCount); err != nil {
			return nil, fmt.Errorf("scan group: %w", err)
		}

		if filter.Type == "user" {
			var isMember string
			err = gm.DB.QueryRow(`SELECT status FROM group_members WHERE group_id = ? AND user_id = ?`, g.ID, userID).Scan(&isMember)
			if err != nil {
				if g.CreatorID == userID {
					g.IsMember = "creator"
				} else {
					g.IsMember = ""
				}
			} else {
				g.IsMember = isMember
			}
		} else {
			g.IsMember = ""
		}
		groups = append(groups, &g)
	}
	return groups, nil
}
