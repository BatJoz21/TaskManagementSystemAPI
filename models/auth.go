package models

import (
	"errors"

	"taskmanagementsystem.localhost/tmsapi/database"
	"taskmanagementsystem.localhost/tmsapi/utils"
)

func GetUserForJWT(user_id int64) (*ResponseUserStruct, error) {
	query := `SELECT theusers.email, user_roles.role
	FROM theusers
	JOIN user_roles ON theusers.id = user_roles.user_id
	WHERE theusers.id = ?`
	row := database.DB.QueryRow(query, user_id)

	var userResponse ResponseUserStruct
	err := row.Scan(&userResponse.Email, &userResponse.Role)
	if err != nil {
		return nil, err
	}

	userResponse.ID = user_id
	return &userResponse, nil
}

func (u *User) ValidateCredentials() error {
	query := `SELECT 
		id, 
		first_name, 
		last_name, 
		password 
	FROM theusers 
	WHERE email = ? AND status IS NULL AND deleted_at IS NULL`
	row := database.DB.QueryRow(query, u.Email)

	var retreivedPassword string
	err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &retreivedPassword)
	if err != nil {
		return err
	}

	// Checking password credentials
	isValid := utils.CheckPasswordHash(u.Password, retreivedPassword)
	if !isValid {
		return errors.New("Credentials invalid")
	}

	return nil
}

func (u *User) ChangePassword() error {
	query := `UPDATE theusers SET password = ? WHERE email = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Hashing password
	hashedPass, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(hashedPass, u.Email)

	return err
}
