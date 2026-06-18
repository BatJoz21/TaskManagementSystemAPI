package models

import (
	"database/sql"
	"time"
)

type GetTaskResponse struct {
	ID          int64        `json:"id"`
	UsersID     int64        `json:"users_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	StatusName  string       `json:"status_name"`
	DueDate     time.Time    `json:"due_date"`
	Attachment  *string      `json:"attachment"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	TagName     string       `json:"tag_name"`
	DeletedAt   sql.NullTime `json:"deleted_at"`
}
