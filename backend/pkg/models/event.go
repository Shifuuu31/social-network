package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Event struct {
	ID          int       `json:"id"`
	GroupId     int       `json:"group_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	VotesCount  int       `json:"vote_count"`
	UserVote    string    `json:"user_vote"` // User's vote status: "going", "not_going", or ""
	EventTime   time.Time `json:"event_time"`
	CreatedAt   time.Time `json:"created_at"`
}

func (e *Event) Validate() error {
	e.Title = strings.TrimSpace(e.Title)
	if e.Title == "" || len(e.Title) > 100 {
		return errors.New("invalid title")
	}
	e.Description = strings.TrimSpace(e.Description)
	if e.Description == "" || len(e.Description) > 500 {
		return errors.New("invalid description")
	}
	if e.EventTime.Before(time.Now()) {
		return errors.New("event time must be in the future")
	}
	if e.GroupId <= 0 {
		return errors.New("invalid group id")
	}
	return nil
}

type EventModel struct {
	DB *sql.DB
}

// InsertEvent inserts a new event.
func (em *EventModel) Insert(event *Event) error {
	query := `
		INSERT INTO events (group_id, title, description, event_time)
		VALUES (?, ?, ?, ?)
	`
	res, err := em.DB.Exec(query, event.GroupId, event.Title, event.Description, event.EventTime)
	if err != nil {
		return fmt.Errorf("insert event: %w", err)
	}
	id, err := res.LastInsertId()
	if err == nil {
		event.ID = int(id)
	}
	return nil
}

// GetEventByID retrieves an event by ID.
func (em *EventModel) GetEventByID(event *Event) error {
	query := `
		SELECT e.id, e.group_id, e.title, e.description, e.event_time, e.created_at,
		       COUNT(ev.id) AS vote_count
		FROM events e
		LEFT JOIN event_votes ev ON e.id = ev.event_id
		WHERE e.id = ?
		GROUP BY e.id
	`

	if err := em.DB.QueryRow(query, event.ID).Scan(
		&event.ID, &event.GroupId, &event.Title, &event.Description,
		&event.EventTime, &event.CreatedAt, &event.VotesCount,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("event not found")
		}
		return fmt.Errorf("get event by id: %w", err)
	}
	// Note: UserVote is not set here as it requires a specific user context
	return nil
}

// UpdateEvent updates an event.
func (em *EventModel) UpdateEvent(event *Event) error {
	query := `
		UPDATE events
		SET title = ?, description = ?, event_time = ?
		WHERE id = ?
	`
	res, err := em.DB.Exec(query, event.Title, event.Description, event.EventTime, event.ID)
	if err != nil {
		return fmt.Errorf("update event: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no event updated")
	}
	return nil
}

// DeleteEvent deletes an event by ID.
func (em *EventModel) DeleteEvent(eventID int) error {
	res, err := em.DB.Exec(`DELETE FROM events WHERE id = ?`, eventID)
	if err != nil {
		return fmt.Errorf("delete event: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("event not found")
	}
	return nil
}

type EventsPayload struct {
	GroupID    int `json:"group_id"`
	UserID     int `json:"user_id"` // Add user ID for vote status
	Start      int `json:"start"`
	NumOfItems int `json:"n_items"`
}

// GetEventsByGroup returns paginated events for a group.
func (em *EventModel) GetEventsByGroup(payload *EventsPayload) ([]*Event, error) {
	var maxid int
	if payload.Start == -1 {
		if err := em.DB.QueryRow(`SELECT MAX(id) FROM events WHERE group_id = ?`, payload.GroupID).Scan(&maxid); err != nil {
			fmt.Println("Error getting max event id:", err)
			return nil, fmt.Errorf("get max event id: %w", err)
		}

		payload.Start = maxid
		if payload.Start == 0 {
			payload.Start = 2147483647 // max int32 value
		}
	}
	fmt.Println("Fetching events for group:", payload.GroupID, "starting from ID:", payload.Start, "limit:", payload.NumOfItems)

	query := `
		SELECT e.id, e.group_id, e.title, e.description, e.event_time, e.created_at,
		       COUNT(ev.id) AS vote_count
		FROM events e
		LEFT JOIN event_votes ev ON e.id = ev.event_id AND ev.vote = 'going'
		WHERE e.group_id = ? AND e.id <= ?
		GROUP BY e.id
		ORDER BY e.id DESC
		LIMIT ?
	`
	rows, err := em.DB.Query(query, payload.GroupID, payload.Start, payload.NumOfItems)
	if err != nil {
		return nil, fmt.Errorf("get events by group: %w", err)
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.GroupId, &e.Title, &e.Description, &e.EventTime, &e.CreatedAt, &e.VotesCount); err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}

		// Check user's vote status for this event
		if payload.UserID > 0 {
			var userVote string
			err = em.DB.QueryRow(`SELECT vote FROM event_votes WHERE event_id = ? AND user_id = ?`, e.ID, payload.UserID).Scan(&userVote)
			if err != nil {
				// User hasn't voted - this is not an error
				e.UserVote = ""
			} else {
				e.UserVote = userVote
			}
		}

		events = append(events, &e)
	}
	return events, nil
}

// IsEventInGroup checks if event belongs to a group.
func (em *EventModel) IsEventInGroup(eventID, groupID int) error {
	var count int
	if err := em.DB.QueryRow(`SELECT COUNT(*) FROM events WHERE id = ? AND group_id = ?`, eventID, groupID).Scan(&count); err != nil || count <= 0 {
		return fmt.Errorf("check event in group: %w", err)
	}
	return nil
}

// IsUserEventCreator checks if user is creator of group that owns the event.
func (em *EventModel) IsUserEventCreator(eventID, userID int) error {
	query := `
		SELECT COUNT(*)
		FROM events e
		JOIN groups g ON e.group_id = g.id
		WHERE e.id = ? AND g.creator_id = ?
	`
	var count int

	if err := em.DB.QueryRow(query, eventID, userID).Scan(&count); err != nil || count <= 0 {
		return fmt.Errorf("check event creator: %w", err)
	}
	return nil
}

// CountVotesByType returns the count of a specific vote type for an event.
func (em *EventModel) CountVotesByType(event *Event, voteType string) error {
	query := `SELECT COUNT(*) FROM event_votes WHERE event_id = ? AND vote = ?`
	return em.DB.QueryRow(query, event.ID, voteType).Scan(&event.VotesCount) // store count in ID
}
