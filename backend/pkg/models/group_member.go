package models

import (
	"database/sql"
	"fmt"
	"time"
)

type GroupMember struct {
	ID         int       `json:"id"`
	GroupID    int       `json:"group_id"`
	UserID     int       `json:"user_id"`
	Status     string    `json:"status"` // "invited", "requested", "member", "declined"
	CreatedAt  time.Time `json:"created_at"`
	PrevStatus string    `json:"prev_status"`
}

// TODO
func (gm *GroupMember) Validate() error {
	return nil
}

type GroupMemberModel struct {
	DB *sql.DB
}

func (gmm *GroupMemberModel) Upsert(member *GroupMember) error {
	query := `
		INSERT INTO group_members (group_id, user_id, prev_status, status)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(group_id, user_id) DO UPDATE SET status=excluded.status
	`

	if _, err := gmm.DB.Exec(query, member.GroupID, member.UserID, member.PrevStatus, member.Status); err != nil {
		return fmt.Errorf("upsert group member: %w", err)
	}
	return nil
}

func (gmm *GroupMemberModel) Delete(member *GroupMember) error {
	query := `
		DELETE FROM group_members
		WHERE group_id = ? AND user_id = ?
	`
	res, err := gmm.DB.Exec(query, member.GroupID, member.UserID)
	if err != nil {
		return fmt.Errorf("delete group member: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (gmm *GroupMemberModel) GetMember(member *GroupMember) error {
	query := `
		SELECT id, group_id, user_id, status, prev_status, created_at
		FROM group_members
		WHERE group_id = ? AND user_id = ?
	`
	row := gmm.DB.QueryRow(query, member.GroupID, member.UserID)
	return row.Scan(&member.ID, &member.GroupID, &member.UserID, &member.Status, &member.PrevStatus, &member.CreatedAt)
}

func (gmm *GroupMemberModel) GetGroupMembers(groupID int) ([]*GroupMember, error) {
	query := `
		SELECT id, group_id, user_id, status, created_at
		FROM group_members
		WHERE group_id = ?
	`
	rows, err := gmm.DB.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("get group members: %w", err)
	}
	defer rows.Close()

	var members []*GroupMember
	for rows.Next() {
		var m GroupMember
		if err := rows.Scan(&m.ID, &m.GroupID, &m.UserID, &m.Status, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan group member: %w", err)
		}
		members = append(members, &m)
	}
	return members, nil
}

func (gmm *GroupMemberModel) IsUserGroupMember(payload *GroupMember) error {
	query := `
		SELECT COUNT(*)
		FROM group_members
		WHERE group_id = ? AND user_id = ?
	`
	var count int

	if err := gmm.DB.QueryRow(query, payload.GroupID, payload.UserID).Scan(&count); err != nil || count <= 0 {
		return fmt.Errorf("check user in group: %w", err)
	}
	return nil
}
