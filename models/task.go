package models

import (
	"database/sql"
	"time"

	"taskmanagementsystem.localhost/tmsapi/database"
)

type Task struct {
	ID          uint           `json:"id"`
	UsersID     uint           `json:"users_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	StatusID    uint           `json:"status_id"`
	DueDate     time.Time      `json:"due_date"`
	Attachment  sql.NullString `json:"attachment"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	TagID       uint           `json:"tag_id"`
	DeletedAt   sql.NullTime   `json:"deleted_at"`
}

func GetAllTasks() ([]Task, error) {
	query := `SELECT * FROM tasks`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.UsersID, &task.Title, &task.Description, &task.StatusID, &task.DueDate,
			&task.Attachment, &task.CreatedAt, &task.UpdatedAt, &task.TagID, &task.DeletedAt)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
