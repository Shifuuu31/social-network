package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type EventVote struct {
	ID        int       `json:"id"`
	EventID   int       `json:"event_id"`
	UserID    int       `json:"user_id"`
	Vote      string    `json:"vote"` // example values: "yes", "no"
	CreatedAt time.Time `json:"created_at"`
}

func (ev *EventVote) Validate() error {
	ev.Vote = strings.ToLower(strings.TrimSpace(ev.Vote))
	if ev.Vote != "going" && ev.Vote != "not_going" {
		return errors.New("vote must be 'going' or 'not_going'")
	}
	if ev.EventID <= 0 {
		return errors.New("invalid event ID")
	}
	if ev.UserID <= 0 {
		return errors.New("invalid user ID")
	}
	return nil
}

type EventVoteModel struct {
	DB *sql.DB
}

// UpsertVote inserts or updates a vote by a user for an event.
func (evm *EventVoteModel) UpsertVote(vote *EventVote) error {
	query := `
		INSERT INTO event_votes (event_id, user_id, vote)
		VALUES (?, ?, ?)
		ON CONFLICT(event_id, user_id)
		DO UPDATE SET vote = excluded.vote
	`

	if _, err := evm.DB.Exec(query, vote.EventID, vote.UserID, vote.Vote); err != nil {
		return fmt.Errorf("upsert vote: %w", err)
	}
	return nil
}

// DeleteVote removes a user's vote on an event.
func (evm *EventVoteModel) DeleteVote(vote *EventVote) error {
	if _, err := evm.DB.Exec(`DELETE FROM event_votes WHERE event_id = ? AND user_id = ?`, vote.EventID, vote.UserID); err != nil {
		return fmt.Errorf("delete vote: %w", err)
	}
	return nil
}

// HasUserVoted checks if a user has already voted on an event.
func (evm *EventVoteModel) HasUserVoted(vote *EventVote) error {
	query := `
		SELECT vote
		FROM event_votes
		WHERE event_id = ? AND user_id = ?
	`
	return evm.DB.QueryRow(query, vote.EventID, vote.UserID).Scan(&vote.Vote)
}

// CountVotesByType returns the count of a specific vote type for an event.
func (evm *EventVoteModel) CountVotesByType(vote *EventVote) error {
	query := `SELECT COUNT(*) FROM event_votes WHERE event_id = ? AND vote = ?`
	return evm.DB.QueryRow(query, vote.EventID, vote.Vote).Scan(&vote.ID) // store count in ID
}

type VotesPayload struct {
	EventID    string `json:"event_id"`
	Start      int    `json:"star"`
	NumOfItems int    `json:"n_items"`
}

// GetVotesByEvent returns a paginated list of votes for a given event.
func (evm *EventVoteModel) GetVotesByEvent(payload *VotesPayload) ([]*EventVote, error) {
	if payload.Start == -1 {
		if err := evm.DB.QueryRow(`SELECT MAX(id) FROM event_votes WHERE event_id = ?`, payload.EventID).Scan(&payload.Start); err != nil {
			return nil, fmt.Errorf("get max vote id: %w", err)
		}
	}

	query := `
		SELECT id, event_id, user_id, vote, created_at
		FROM event_votes
		WHERE event_id = ? AND id <= ?
		ORDER BY id DESC
		LIMIT ?
	`
	rows, err := evm.DB.Query(query, payload.EventID, payload.Start, payload.NumOfItems)
	if err != nil {
		return nil, fmt.Errorf("get votes by event: %w", err)
	}
	defer rows.Close()

	var votes []*EventVote
	for rows.Next() {
		var v EventVote
		if err := rows.Scan(&v.ID, &v.EventID, &v.UserID, &v.Vote, &v.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan vote: %w", err)
		}
		votes = append(votes, &v)
	}
	return votes, nil
}
