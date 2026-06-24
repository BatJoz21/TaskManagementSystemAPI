package models

import (
	"time"

	"taskmanagementsystem.localhost/tmsapi/database"
)

type Task struct {
	ID          int64     `json:"id"`
	UsersID     int64     `json:"users_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StatusID    int64     `json:"status_id"`
	DueDate     time.Time `json:"due_date"`
	Attachment  *string   `json:"attachment"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	TagID       int64     `json:"tag_id"`
	DeletedAt   time.Time `json:"deleted_at"`
}

func (t *Task) Save() error {
	query := `INSERT INTO tasks(users_id, title, description, status_id, due_date, attachment, tag_id)
	VALUES (?, ?, ?, ?, ?, ?, ?)`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(t.UsersID, t.Title, t.Description, t.StatusID, t.DueDate, t.Attachment, t.TagID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	t.ID = id
	return err
}

func GetAllTasks(users_id int64, sort, order, status, tag string, isDeleted bool) ([]GetTaskResponse, error) {
	query := `
	SELECT
		tasks.id,
		tasks.users_id,
		tasks.title,
		tasks.description,
		statuses.name AS status_name,
		tasks.due_date,
		tasks.attachment,
		tasks.created_at,
		tasks.updated_at,
		tags.name_tag AS tag_name,
		tasks.deleted_at
	FROM tasks
	JOIN statuses ON tasks.status_id = statuses.id
	JOIN tags ON tasks.tag_id = tags.id
	WHERE tasks.users_id = ?`

	if isDeleted {
		query += ` AND tasks.deleted_at IS NOT NULL`
	} else {
		query += ` AND tasks.deleted_at IS NULL`
	}

	allowedStatus := map[string]bool{"in_progress": true, "cancelled": true, "complete": true}
	if allowedStatus[status] {
		query += ` HAVING status_name = "` + status + `"`
	}
	allowedTags := map[string]bool{"Work": true, "Study": true, "Personal": true, "Urgent": true}
	if allowedTags[tag] {
		query += ` HAVING tag_name = "` + tag + `"`
	}

	allowedSorts := map[string]bool{"id": true, "title": true, "due_date": true}
	allowedOrder := map[string]bool{"ASC": true, "DESC": true}
	if !allowedSorts[sort] {
		sort = "id"
	}
	if !allowedOrder[order] {
		order = "ASC"
	}

	query += ` ORDER BY tasks.` + sort + ` ` + order

	rows, err := database.DB.Query(query, users_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []GetTaskResponse
	for rows.Next() {
		var task GetTaskResponse
		err = rows.Scan(&task.ID, &task.UsersID, &task.Title, &task.Description, &task.StatusName, &task.DueDate,
			&task.Attachment, &task.CreatedAt, &task.UpdatedAt, &task.TagName, &task.DeletedAt)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func GetTaskByID(id, users_id int64) (*GetTaskResponse, error) {
	query := `
	SELECT
		tasks.id,
		tasks.users_id,
		tasks.title,
		tasks.description,
		tasks.status_id,
		statuses.name AS status_name,
		tasks.due_date,
		tasks.attachment,
		tasks.tag_id,
		tags.name_tag AS tag_name,
		tasks.created_at,
		tasks.updated_at,
		tasks.deleted_at
	FROM tasks
	JOIN statuses ON tasks.status_id = statuses.id
	JOIN tags ON tasks.tag_id = tags.id
	WHERE tasks.id = ? && tasks.users_id = ?`
	row := database.DB.QueryRow(query, id, users_id)

	var task GetTaskResponse
	err := row.Scan(&task.ID, &task.UsersID, &task.Title, &task.Description, &task.StatusID, &task.StatusName, &task.DueDate,
		&task.Attachment, &task.TagID, &task.TagName, &task.CreatedAt, &task.UpdatedAt, &task.DeletedAt)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func GetTaskAttachmentByID(id int64) (*GetAttachmentResponse, error) {
	query := `SELECT attachment FROM tasks WHERE id = ? AND deleted_at IS NULL`
	row := database.DB.QueryRow(query, id)

	var responseData GetAttachmentResponse
	err := row.Scan(&responseData.Attachment)
	if err != nil {
		return nil, err
	}

	return &responseData, nil
}

func (t *Task) Update() error {
	query := `UPDATE tasks 
	SET
		title = ?,
		description = ?,
		status_id = ?,
		due_date = ?,
		attachment = ?,
		tag_id = ?
	WHERE id = ? && deleted_at IS NULL`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(t.Title, t.Description, t.StatusID, t.DueDate, t.Attachment, t.TagID, t.ID)

	return err
}

func (t *Task) CompleteTask() error {
	query := `UPDATE tasks SET status_id = 3 WHERE id = ? AND deleted_at IS NULL`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(t.ID)
	return err
}

func (t *Task) RestoreTask() error {
	query := `UPDATE tasks SET deleted_at = NULL WHERE id = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(t.ID)
	return err
}

func (t *Task) Delete() error {
	query := `UPDATE tasks SET deleted_at = ? WHERE id = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(time.Now(), t.ID)
	return err
}

func (t *Task) DeleteAttachment() error {
	query := `UPDATE tasks SET attachment = NULL WHERE id = ?`
	stmt, err := database.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(t.ID)
	return err
}
