package models

import (
	"time"

	"taskmanagementsystem.localhost/tmsapi/database"
)

type RefreshTokenStruct struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	DeviceName *string    `json:"device_name"`
	TokenHash  string     `json:"token_hash"`
	ExpiresAt  time.Time  `json:"expires_at"`
	RevokedAt  *time.Time `json:"revoked_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

func GetRefreshTokenByHash(hash string) (*RefreshTokenStruct, error) {
	query := `SELECT
		id,
		user_id,
		device_name,
		token_hash,
		expires_at,
		revoked_at,
		created_at
	FROM refresh_tokens WHERE token_hash = ?`
	row := database.DB.QueryRow(query, hash)

	var tokenStruct RefreshTokenStruct
	err := row.Scan(&tokenStruct.ID, &tokenStruct.UserID, &tokenStruct.DeviceName, &tokenStruct.TokenHash, &tokenStruct.ExpiresAt, &tokenStruct.RevokedAt, &tokenStruct.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &tokenStruct, nil
}

func (r *RefreshTokenStruct) Save() error {
	query := `INSERT INTO refresh_tokens(user_id, token_hash, expires_at)
		VALUES (?, ?, ?)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Store refresh token in DB
	expireDate := time.Now().Add(time.Hour * 24 * 7)
	result, err := stmt.Exec(r.UserID, r.TokenHash, expireDate)
	if err != nil {
		return err
	}

	tokenID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	r.ID = tokenID
	r.ExpiresAt = expireDate

	return nil
}
