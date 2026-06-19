package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("mysql", "root:@tcp(localhost:3306)/tms?parseTime=true")
	if err != nil {
		panic("Couldn't connect to database: " + err.Error())
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	err = DB.Ping()
	if err != nil {
		panic("Couldn't ping database: " + err.Error())
	}

	createTables()
}

func createTables() {
	createStatusesTable := `CREATE TABLE IF NOT EXISTS statuses (
		id INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(80) NOT NULL
	)`
	_, err := DB.Exec(createStatusesTable)
	if err != nil {
		panic("Failed to create statuses table: " + err.Error())
	}

	createTagsTable := `CREATE TABLE IF NOT EXISTS tags (
		id INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		name_tag VARCHAR(80) NOT NULL
	)`
	_, err = DB.Exec(createTagsTable)
	if err != nil {
		panic("Failed to create tags table: " + err.Error())
	}

	createTasksTable := `CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		users_id INTEGER UNSIGNED NOT NULL,
		title VARCHAR(200) NOT NULL,
		description TEXT NOT NULL,
		status_id INTEGER UNSIGNED NOT NULL,
		due_date DATETIME NOT NULL,
		attachment VARCHAR(250),
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		tag_id INTEGER UNSIGNED NOT NULL,
		deleted_at DATETIME NOT NULL,

		CONSTRAINT user_id_fk
			FOREIGN KEY(users_id)
			REFERENCES users (id)
			ON DELETE CASCADE
			ON UPDATE CASCADE,

		CONSTRAINT status_id_fk
			FOREIGN KEY(status_id)
			REFERENCES statuses (id)
			ON DELETE CASCADE
			ON UPDATE CASCADE,

		CONSTRAINT task_id_fk
			FOREIGN KEY(tag_id)
			REFERENCES tags (id)
			ON DELETE CASCADE
			ON UPDATE CASCADE
	)`
	_, err = DB.Exec(createTasksTable)
	if err != nil {
		panic("Failed to create tasks table: " + err.Error())
	}

	createUsersTable := `CREATE TABLE IF NOT EXISTS theusers (
		id INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		first_name VARCHAR(150) NOT NULL,
		last_name VARCHAR(150) NOT NULL,
		email VARCHAR(225) NOT NULL,
		password VARCHAR(225) NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at DATETIME
	)`
	_, err = DB.Exec(createUsersTable)
	if err != nil {
		panic("Failed to create users table: " + err.Error())
	}

	createUserGroupsTable := `CREATE TABLE IF NOT EXISTS user_roles (
		id INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		user_id INTEGER UNSIGNED NOT NULL,
		role VARCHAR(225) NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

		CONSTRAINT user_roles_user_id_fk
			FOREIGN KEY(user_id)
			REFERENCES theusers (id)
			ON DELETE CASCADE
			ON UPDATE CASCADE
	)`
	_, err = DB.Exec(createUserGroupsTable)
	if err != nil {
		panic("Failed to create user_groups table: " + err.Error())
	}
}
