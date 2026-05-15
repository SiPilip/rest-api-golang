package models

import (
	"REST-API/db"
	"time"
)

type RefreshToken struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func SaveRefreshToken(userId int64, token string, expiresAt time.Time) error {
	query := `
	INSERT INTO refresh_tokens(user_id, token, expires_at)
	VALUES (?, ?, ?)
	`
	_, err := db.DB.Exec(query, userId, token, expiresAt)
	return err
}

func GetRefreshToken(token string) (*RefreshToken, error) {
	query := `
	SELECT id, user_id, token, expires_at
	FROM refresh_tokens WHERE token = ?
	`
	row := db.DB.QueryRow(query, token)

	var rt RefreshToken
	err := row.Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func DeleteRefreshToken(token string) error {
	_, err := db.DB.Exec(`
	DELETE FROM refresh_tokens
	WHERE token = ?
	`, token)
	return err
}

func DeleteAllUserRefreshTokens(userId int64) error {
	_, err := db.DB.Exec(`
	DELETE FROM refresh_tokens
	WHERE user_id = ?
	`, userId)
	return err
}