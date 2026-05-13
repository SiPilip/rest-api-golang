package models

import (
	"REST-API/db"
	"database/sql"
	"errors"
	"time"
)

type Event struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" binding:"required,min=3,max=80"`
	Description string    `json:"description" binding:"required,max=255"`
	Location    string    `json:"location" binding:"required"`
	DateTime    time.Time `json:"datetime" binding:"required"`
	UserID      int64     `json:"user_id"`
}


func (e *Event) Save() error {
	query := `
	INSERT INTO events(name, description, location, datetime, user_id) 
	VALUES(?, ?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	
	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	e.ID = id
	return err
} 

func GetAllEvents(page, limit int, search string) ([]Event, int64, error) {
	// Hitung offset
	offset := (page - 1) * limit

	// Hitung total datanya
	var total int64
	countQuery := "SELECT COUNT(*) FROM events"
	if search != ""{
		countQuery +=  " WHERE name LIKE ? OR description LIKE ? OR location like ?"
		searchParam := "%" + search + "%"
		err := db.DB.QueryRow(countQuery, searchParam, searchParam, searchParam).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
	} else {
		err := db.DB.QueryRow(countQuery).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
	}
	
	// Query dengan data LIMIT dan OFFSET
	query := "SELECT id, name, description, location, datetime, user_id FROM events"
	var rows *sql.Rows
	var err error

	if search != ""{
		query += " WHERE name LIKE ? OR description LIKE ? OR location LIKE ? ORDER BY id DESC LIMIT ? OFFSET ?"
		searchParams:= "%" + search + "%"
		rows, err =  db.DB.Query(query, searchParams, searchParams, searchParams, limit, offset)
	} else {
		query += " ORDER BY id DESC LIMIT ? OFFSET ?"
		rows, err = db.DB.Query(query, limit, offset)
	}
	
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
		if err != nil {
			return nil, 0, err
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return events, total, nil
}
	// rows, err := db.DB.Query(query)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()

	// var events []Event

	// for rows.Next() {
	// 	var event Event
	// 	err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	events = append(events, event)
	// }

	// if err = rows.Err(); err != nil {
	// 	return nil, err
	// }
	// return events, nil


func GetEventByID(id int64) (*Event, error) {
	query := "SELECT id, name, description, location, datetime, user_id FROM events WHERE id = ?"
	row := db.DB.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

	if err != nil {
			return nil, err
	}

	return &event, nil
}

func (event *Event) Update() error {
	query := `
	UPDATE events
	SET name = ?, description = ?, location = ?, datetime = ?
	WHERE id = ?
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
			return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.ID)
	return err
}

func (event *Event) Delete() error {
	query := `
	DELETE FROM events
	WHERE id = ?
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
			return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.ID)
	return err
}

func (e *Event) Register(userId int64) error {
	query := `
	INSERT INTO registrations(event_id, user_id)
	VALUES (?, ?)
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)
	return err
}

func (e Event) CancelRegistration(userId int64) error {
	query := `
	DELETE FROM registrations WHERE event_id = ? AND user_id = ?
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(e.ID, userId)

	if err != nil {
    return err
	}

	// Check affected row
	rowsAffected, err := result.RowsAffected()
	if err != nil {
			return err
	}
	if rowsAffected == 0 {
			return errors.New("registration not found") 
	}
	return nil
}
