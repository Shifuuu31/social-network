package sqlite

import (
	"database/sql"

	"social-network/pkg/models"
)

type EventRepository struct {
	DB *sql.DB
}

func (repo *EventRepository) GetAll(groupID string) ([]models.Event, error) {
	events := []models.Event{}
	rows, err := repo.DB.Query("SELECT event_id, created_by, content, title, strftime('%d.%m.%Y', date) FROM event WHERE group_id = ?  ORDER BY date DESC;", groupID)
	if err != nil {
		return events, err
	}
	for rows.Next() {
		var event models.Event
		rows.Scan(&event.ID, &event.AuthorID, &event.Content, &event.Title, &event.Date)
		events = append(events, event)
	}
	return events, nil
}

func (repo *EventRepository) GetData(eventId string) (models.Event, error) {
	row := repo.DB.QueryRow("SELECT title, content, event_id, group_id, strftime('%d.%m.%Y', date), created_by FROM event WHERE event_id = ? ", eventId)
	var event models.Event
	if err := row.Scan(&event.Title, &event.Content, &event.ID, &event.GroupID, &event.Date, &event.AuthorID); err != nil {
		return event, err
	}
	return event, nil
}

func (repo *EventRepository) Save(event models.Event) error {
	stmt, err := repo.DB.Prepare("INSERT INTO event (event_id, group_id, created_by, content, title, date) values (?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(event.ID, event.GroupID, event.AuthorID, event.Content, event.Title, event.Date); err != nil {
		return err
	}
	return nil
}

func (repo *EventRepository) AddParticipant(eventID, userID string) error {
	stmt, err := repo.DB.Prepare("INSERT INTO event_users (event_id, user_id) values (?,?)")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(eventID, userID); err != nil {
		return err
	}
	return nil
}

func (repo *EventRepository) RemoveParticipant(eventID, userID string) error {
	stmt, err := repo.DB.Prepare("DELETE FROM event_users WHERE user_id = ? AND event_id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userID, eventID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *EventRepository) IsParticipating(eventID, userID string) (bool, error) {
	row := repo.DB.QueryRow("SELECT COUNT() FROM event_users WHERE event_id = ? AND  user_id = ?", eventID, userID)
	var participate int
	if err := row.Scan(&participate); err != nil {
		return false, err
	}
	if participate == 1 {
		return true, nil
	} else {
		return false, nil
	}
}
