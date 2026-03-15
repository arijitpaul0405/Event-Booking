package models

import (
	"errors"
	"fmt"
	"time"

	"example.com/event-booking/db"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int64
}

type RegisteredEvent struct {
	ID 		int64
	EventID int64
	UserID 	int64
}

type ResultRegisteredEvent struct {
	EventID 	int64
	UserID 		int64
	Name        string
	Description string
	Location    string
	DateTime    time.Time
}

func (e *Event) Save() error {
	query := `
	INSERT INTO events (name, description, location, dateTime, user_id)
	VALUES (?, ?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)

	if err != nil {
		return err
	}

	e.ID, err = result.LastInsertId()

	return err
}

func GetAllEvents(userid int64) (*[]Event, error) {
	query := `SELECT * FROM events WHERE user_id = ?`
	rows, err := db.DB.Query(query, userid)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events = []Event{}
	var event Event
	for rows.Next() {
		rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
		events = append(events, event)
	}

	return &events, nil
}

func GetEventByID(id, userId int64) (*Event, error) {
	query := "SELECT * FROM events WHERE user_id = ? AND id = ?"
	row := db.DB.QueryRow(query, userId, id)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

	if err != nil {
		return &Event{}, err
	}

	return &event, nil
}

func (e *Event) UpdateByID() error {
	query := `
	UPDATE events
	SET name = ?, description = ?, location = ?, dateTime = ?
	WHERE id = ?
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.ID)

	return err
}

func (e *Event) DeleteByID() error {
	query := `DELETE FROM events WHERE id = ?`
	
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.ID)

	return err
}

func (e *Event) Register(user_id int64) (int64, error) {
	query := "INSERT INTO registrations (event_id, user_id) VALUES (?, ?)"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(e.ID, user_id)

	if err != nil {
		err_msg := "Event already registered!"
		fmt.Printf("%v %v", err_msg, err)
		return 0, errors.New(err_msg)
	}

	var registeration_id int64
	registeration_id, err = result.LastInsertId()
	fmt.Println(registeration_id)

	if err != nil {
		return 0, err
	}

	return registeration_id, err
}

func GetRegistrationByUserID(user_id int64) (*[]ResultRegisteredEvent, error) {
	query := `
	SELECT r.event_id, r.user_id, e.name, e.description, e.location, e.datetime
	FROM registrations r
	INNER JOIN events e ON r.event_id = e.id
	WHERE r.user_id = ?
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(user_id)

	if err != nil {
		return nil, err
	}

	var registered_events = []ResultRegisteredEvent{}
	var registered_event = ResultRegisteredEvent{}
	for rows.Next() {
		err = rows.Scan(
			&registered_event.EventID, &registered_event.UserID, 
			&registered_event.Name, &registered_event.Description,
			&registered_event.Location, &registered_event.DateTime,
		)
		registered_events = append(registered_events, registered_event)
	}

	if err != nil {
		return nil, err
	}

	return &registered_events, nil
}

func GetRegistrationByID(registeration_id int64) (*RegisteredEvent, error) {
	query := "SELECT * FROM registrations WHERE id = ?"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var registered_event RegisteredEvent
	row := stmt.QueryRow(registeration_id)
	err = row.Scan(
		&registered_event.ID, &registered_event.EventID, &registered_event.UserID,
	)

	if err != nil {
		return nil, err
	}

	return &registered_event, nil
}

func CancelRegisteration(registeration_id, user_id int64) error {
	query := "DELETE FROM registrations WHERE id = ? AND user_id = ?"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		fmt.Println(err)
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(registeration_id, user_id)

	affectedRow, err := result.RowsAffected()

	if err != nil {
		fmt.Println(err)
		return err
	}

	if affectedRow == 0 {
		err_msg := fmt.Sprintf("No registration with id %v exists for the given user!", registeration_id)
		return errors.New(err_msg)
	}

	return err
}