package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Group struct {
	ID          int
	CreatorID   int
	Title       string
	Description string
	CreatedAt   time.Time
}

type GroupModel struct {
	DB *sql.DB
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

// InsertGroup inserts a new group into the database.
func (gm *GroupModel) InsertGroup(group *Group) error {
	query := `
		INSERT INTO groups (
			creator_id, title, description
		) VALUES (?, ?, ?)
	`

	res, err := gm.DB.Exec(query,
		group.CreatorID,
		group.Title,
		group.Description,
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
