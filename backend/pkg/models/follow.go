package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type FollowRequest struct {
	ID         int       `json:"id"`
	FromUserID int       `json:"from_user_id"`
	ToUserID   int       `json:"to_user_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type FollowRequestModel struct {
	DB *sql.DB
}

// CanFollow checks if follower can send a follow request or follow directly.
// Return false if already following or a pending request exists.
func (flm *FollowRequestModel) CanFollow(followRequest *FollowRequest) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM follow_request
		WHERE from_user_id = ? AND to_user_id = ? AND status IN ('pending', 'accepted')
	`
	err := flm.DB.QueryRow(query, followRequest.FromUserID, followRequest.ToUserID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// SendFollowRequest inserts a new follow request with status = 'pending'.
func (flm *FollowRequestModel) Insert(followRequest *FollowRequest) error {
	// Prevent duplicate or conflicting requests
	canFollow, err := flm.CanFollow(followRequest)
	if err != nil {
		return err
	}
	if !canFollow {
		return errors.New("follow request already exists or user is already followed")
	}

	query := `
		INSERT INTO follow_request (from_user_id, to_user_id, status)
		VALUES (?, ?, 'pending')
	`
	_, err = flm.DB.Exec(query, followRequest.FromUserID, followRequest.ToUserID)
	return err
}

func (flm *FollowRequestModel) UpdateStatus(followRequest *FollowRequest) error {
	query := `
		UPDATE follow_request
		SET status = ?
		WHERE from_user_id = ? AND to_user_id = ? AND status = 'pending'
	`
	res, err := flm.DB.Exec(query, followRequest.Status, followRequest.FromUserID, followRequest.ToUserID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("no pending follow request to decline")
	}
	return nil
}

// UnfollowUser deletes an accepted follow relationship.
func (flm *FollowRequestModel) Delete(followRequest *FollowRequest) error {
	query := `
		DELETE FROM follow_request
		WHERE from_user_id = ? AND to_user_id = ? AND status = 'accepted'
	`
	res, err := flm.DB.Exec(query, followRequest.FromUserID, followRequest.ToUserID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("no accepted follow record to unfollow")
	}
	return nil
}

// GetFollows returns either followers or following users based on the followType ("followers" or "following").
func (flm *FollowRequestModel) GetFollows(userID int, followType string) (users []User, err error) {
	var query string

	switch followType {
	case "followers":
		query = `
			SELECT u.id, u.email, u.first_name, u.last_name, u.nickname, u.avatar_url
			FROM users u
			JOIN follow_request fr ON fr.from_user_id = u.id
			WHERE fr.to_user_id = ? AND fr.status = 'accepted'
		`
	case "following":
		query = `
			SELECT u.id, u.email, u.first_name, u.last_name, u.nickname, u.avatar_url
			FROM users u
			JOIN follow_request fr ON fr.to_user_id = u.id
			WHERE fr.from_user_id = ? AND fr.status = 'accepted'
		`
	default:
		return nil, fmt.Errorf("invalid followType: must be 'followers' or 'following'")
	}

	rows, err := flm.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Nickname, &u.AvatarURL); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// GetFollowStatus returns "none", "pending", "accepted", or "declined".
func (flm *FollowRequestModel) GetFollowStatus(followRequest *FollowRequest) error {
	query := `
		SELECT status FROM follow_request
		WHERE (from_user_id = ? AND to_user_id = ?) OR id = ? 
	`
	err := flm.DB.QueryRow(query, followRequest.FromUserID, followRequest.ToUserID, followRequest.ID).Scan(&followRequest.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			followRequest.Status = "none"
			return nil
		}
		return err
	}
	return nil
}

// Close Friends
func (f *FollowRequestModel) AddCloseFriend(userID, friendID int) error {
	_, err := f.DB.Exec(`INSERT OR IGNORE INTO close_friends (user_id, friend_id) VALUES (?, ?)`, userID, friendID)
	return err
}

func (f *FollowRequestModel) RemoveCloseFriend(userID, friendID int) error {
	_, err := f.DB.Exec(`DELETE FROM close_friends WHERE user_id = ? AND friend_id = ?`, userID, friendID)
	return err
}

func (f *FollowRequestModel) ListCloseFriends(userID int) ([]*User, error) {
	rows, err := f.DB.Query(`
		SELECT u.id, u.nickname, u.first_name, u.last_name, u.avatar_url
		FROM users u
		JOIN close_friends cf ON cf.friend_id = u.id
		WHERE cf.user_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []*User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Nickname, &u.FirstName, &u.LastName, &u.AvatarURL); err != nil {
			return nil, err
		}
		friends = append(friends, &u)
	}
	return friends, nil
}

func (f *FollowRequestModel) IsCloseFriend(userID, friendID int) (bool, error) {
	row := f.DB.QueryRow(`SELECT 1 FROM close_friends WHERE user_id = ? AND friend_id = ?`, userID, friendID)
	var exists int
	err := row.Scan(&exists)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
