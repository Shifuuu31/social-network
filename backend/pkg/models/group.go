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
		SELECT g.id, g.creator_id, g.title, g.description, g.created_at,
		       COUNT(m.id) AS member_count
		FROM groups g
		LEFT JOIN group_members m ON g.id = m.group_id AND m.status = 'member'
		WHERE g.id = ? AND g.creator_id = ?
		GROUP BY g.id
	`

	if err := gm.DB.QueryRow(query, group.ID, group.CreatorID).Scan(
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
	Start      int    `json:"start"`
	NumOfItems int    `json:"n_items"`
	Type       string `json:"type"`
}

func (gm *GroupModel) GetGroups(userID int, filter *GroupsPayload) ([]*Group, error) {
	var (
		query string
		args  []any
	)

	if filter.Start == -1 {
		switch filter.Type {
		case "all":
			if err := gm.DB.QueryRow(`SELECT MAX(id) FROM groups`).Scan(&filter.Start); err != nil {
				return nil, fmt.Errorf("get max group id (all): %w", err)
			}
		case "creator":
			if err := gm.DB.QueryRow(`SELECT MAX(id) FROM groups WHERE creator_id = ?`, userID).Scan(&filter.Start); err != nil {
				return nil, fmt.Errorf("get max group id (creator): %w", err)
			}
		case "member":
			if err := gm.DB.QueryRow(`
				SELECT MAX(g.id) 
				FROM groups g
				INNER JOIN group_members gm ON g.id = gm.group_id
				WHERE gm.user_id = ? AND gm.status = 'member'`, userID).Scan(&filter.Start); err != nil {
				return nil, fmt.Errorf("get max group id (member): %w", err)
			}
		default:
			return nil, fmt.Errorf("invalid group type: %s", filter.Type)
		}
	}

	switch filter.Type {
	case "all":
		query = `
			SELECT g.id, g.creator_id, g.title, g.description, g.image_uuid, g.created_at,
			       COUNT(m.id) AS member_count
			FROM groups g
			LEFT JOIN group_members m ON g.id = m.group_id AND m.status = 'member'
			WHERE g.id <= ?
			GROUP BY g.id
			ORDER BY g.id DESC
			LIMIT ?
		`
		args = []any{filter.Start, filter.NumOfItems}

	case "creator":
		query = `
			SELECT g.id, g.creator_id, g.title, g.description, g.image_uuid, g.created_at,
			       COUNT(m.id) AS member_count
			FROM groups g
			LEFT JOIN group_members m ON g.id = m.group_id AND m.status = 'member'
			WHERE g.creator_id = ? AND g.id <= ?
			GROUP BY g.id
			ORDER BY g.id DESC
			LIMIT ?
		`
		args = []any{userID, filter.Start, filter.NumOfItems}

	case "member":
		query = `
			SELECT g.id, g.creator_id, g.title, g.description, g.image_uuid, g.created_at,
			       COUNT(m2.id) AS member_count
			FROM groups g
			INNER JOIN group_members m1 ON g.id = m1.group_id AND m1.user_id = ? AND m1.status = 'member'
			LEFT JOIN group_members m2 ON g.id = m2.group_id AND m2.status = 'member'
			WHERE g.id <= ?
			GROUP BY g.id
			ORDER BY g.id DESC
			LIMIT ?
		`
		args = []any{userID, filter.Start, filter.NumOfItems}

	default:
		return nil, fmt.Errorf("invalid group type: %s", filter.Type)
	}

	rows, err := gm.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("get groups (%s): %w", filter.Type, err)
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
