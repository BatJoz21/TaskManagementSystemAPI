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
	StatusID    int64        `json:"status_id"`
	StatusName  string       `json:"status_name"`
	DueDate     time.Time    `json:"due_date"`
	Attachment  *string      `json:"attachment"`
	TagID       int64        `json:"tag_id"`
	TagName     string       `json:"tag_name"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at"`
}
