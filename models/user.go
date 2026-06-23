package models

import (
	"errors"
	"time"

	"taskmanagementsystem.localhost/tmsapi/database"
	"taskmanagementsystem.localhost/tmsapi/utils"
)

type User struct {
	ID             int64     `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Status         string    `json:"status"`
	ProfilePicture string    `json:"profile_picture"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"deleted_at"`
}

func (u *User) Save() error {
	query := `INSERT INTO theusers(first_name, last_name, email, password)
		VALUES (?, ?, ?, ?)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedPass, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.FirstName, u.LastName, u.Email, hashedPass)
	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()
	u.ID = userId
	if err != nil {
		return err
	}

	err = u.SetNewUserRole()
	return err
}

func (u *User) SetNewUserRole() error {
	query := `INSERT INTO user_roles(user_id, role)
		VALUES (?, ?)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.ID, "user")
	return err
}

func (u *User) ValidateCredentials() error {
	query := `SELECT 
		id, 
		first_name, 
		last_name, 
		password 
	FROM theusers 
	WHERE email = ? && status IS NULL`
	row := database.DB.QueryRow(query, u.Email)

	var retreivedPassword string
	err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &retreivedPassword)
	if err != nil {
		return err
	}

	isValid := utils.CheckPasswordHash(u.Password, retreivedPassword)
	if !isValid {
		return errors.New("Credentials invalid")
	}

	return nil
}

func (u *User) GetUserRole() (string, error) {
	query := `SELECT role FROM user_roles WHERE user_id = ?`
	row := database.DB.QueryRow(query, u.ID)

	var retreivedRole string
	err := row.Scan(&retreivedRole)
	if err != nil {
		return "", err
	}

	return retreivedRole, nil
}
