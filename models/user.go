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
	ProfilePicture *string   `json:"profile_picture"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"deleted_at"`
}

var allowedUserSorts = map[string]string{"firstName": "first_name", "lastName": "last_name", "email": "email"}

func (u *User) Save() error {
	query := `INSERT INTO theusers(first_name, last_name, email, password)
		VALUES (?, ?, ?, ?)`
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

	result, err := stmt.Exec(u.FirstName, u.LastName, u.Email, hashedPass)
	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()
	u.ID = userId
	if err != nil {
		return err
	}

	// Give new user a role
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

func GetUsers(sort, order string) (*[]ResponseUserStruct, int, error) {
	query := `SELECT
		theusers.id,
		theusers.first_name,
		theusers.last_name,
		theusers.email,
		user_roles.role
	FROM theusers
	JOIN user_roles ON theusers.id = user_roles.user_id`

	// Filtering
	if value, ok := allowedUserSorts[sort]; ok {
		sort = value
	} else {
		sort = "id"
	}
	if !allowedOrder[order] {
		order = "ASC"
	}
	query += ` ORDER BY theusers.` + sort + ` ` + order

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, 0, err
	}

	var users []ResponseUserStruct
	for rows.Next() {
		var user ResponseUserStruct
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Role)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return &users, 0, nil
}

func GetUser(id int64) (*UserProfileStruct, error) {
	query := `SELECT
		theusers.id,
		theusers.first_name,
		theusers.last_name,
		theusers.email,
		theusers.status,
		theusers.profile_picture,
		user_roles.role
	FROM theusers
	JOIN user_roles ON theusers.id = user_roles.user_id
	WHERE theusers.id = ?`
	row := database.DB.QueryRow(query, id)

	var user UserProfileStruct
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Status, &user.ProfilePicture, &user.Role)
	if err != nil {
		return nil, err
	}

	return &user, nil
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

func (u *User) Update() error {
	query := `UPDATE theusers
	SET
		first_name = ?,
		last_name = ?,
		profile_picture = ?
	WHERE id = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.FirstName, u.LastName, u.ProfilePicture, u.ID)

	return err
}

func (u *User) UpdateRole() error {
	query := `UPDATE user_roles
	SET
		role = ?
	WHERE user_id = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.Role, u.ID)

	return err
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
