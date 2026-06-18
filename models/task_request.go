package models

import (
	"time"
)

type CreateTaskRequest struct {
	UsersID     int64     `json:"users_id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	StatusID    int64     `json:"status_id" binding:"required"`
	DueDate     time.Time `json:"due_date" binding:"required"`
	Attachment  *string   `json:"attachment"`
	TagID       int64     `json:"tag_id" binding:"required"`
}

type UpdateTaskRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StatusID    int64     `json:"status_id"`
	DueDate     time.Time `json:"due_date"`
	Attachment  *string   `json:"attachment"`
	TagID       int64     `json:"tag_id"`
}
