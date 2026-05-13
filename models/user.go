package models

import (
	"REST-API/db"
	"REST-API/utils"
	"errors"
)

// https://github.com/go-playground/validator
type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"-" binding:"required,min=6,max=255"`
}

func (u *User) Save() error {
	query := "INSERT INTO users(email, password) VALUES (?, ?)"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)

	if err != nil {
		return err
	}

	res, err := stmt.Exec(u.Email, hashedPassword)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	u.ID = id
	return err

}

func (u *User) ValidateCredentials() error {
	query := `
	SELECT id, password
	FROM users
	WHERE email = ?
	`
	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string
	err := row.Scan(&u.ID, &retrievedPassword)
	if err != nil {
		return errors.New("credentials invalid")
	}

	passwordIsValid := utils.VerifyPassword(u.Password, retrievedPassword)

	if !passwordIsValid {
		return errors.New("credentials invalid")
	}

	return nil
}