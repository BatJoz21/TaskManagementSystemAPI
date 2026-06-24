package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"taskmanagementsystem.localhost/tmsapi/models"
	"taskmanagementsystem.localhost/tmsapi/utils"
)

func createTask(context *gin.Context) {
	// Create DTO for Creating Task
	status_id, err := strconv.ParseInt(context.PostForm("status_id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid status_id"})
		return
	}

	tag_id, err := strconv.ParseInt(context.PostForm("tag_id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid tag_id"})
		return
	}

	due_date, err := time.Parse(time.RFC3339, context.PostForm("due_date"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid due_date"})
		return
	}

	taskDTO := models.CreateTaskRequest{
		Title:       context.PostForm("title"),
		Description: context.PostForm("description"),
		StatusID:    status_id,
		DueDate:     due_date,
		TagID:       tag_id,
	}

	// Handle file uploads
	var attachment *string
	file, err := context.FormFile("attachment")
	if err == nil {
		attachment, err = utils.SaveTaskAttachment(file, context)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	// Insert Task data to database
	task := models.Task{
		UsersID:     context.GetInt64("user_id"),
		Title:       taskDTO.Title,
		Description: taskDTO.Description,
		StatusID:    taskDTO.StatusID,
		DueDate:     taskDTO.DueDate,
		Attachment:  attachment,
		TagID:       taskDTO.TagID,
	}
	err = task.Save()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "task created", "task": task})
}

func getTasks(context *gin.Context) {
	// Read query parameters
	sort := context.Query("sort")
	order := context.DefaultQuery("order", "ASC")
	status := context.Query("status")
	tag := context.Query("tag")

	tasks, err := models.GetAllTasks(context.GetInt64("user_id"), sort, order, status, tag, false)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message:": err.Error()})
		return
	}

	context.JSON(http.StatusOK, tasks)
}

func getDeletedTasks(context *gin.Context) {
	// Read query parameters
	sort := context.Query("sort")
	order := context.DefaultQuery("order", "ASC")
	status := context.Query("status")
	tag := context.Query("tag")

	tasks, err := models.GetAllTasks(context.GetInt64("user_id"), sort, order, status, tag, true)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message:": err.Error()})
		return
	}

	context.JSON(http.StatusOK, tasks)
}

func getTaskByID(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message:": err.Error()})
		return
	}
	task, err := models.GetTaskByID(id, context.GetInt64("user_id"))

	context.JSON(http.StatusOK, task)
}

func viewAttachmentFile(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message:": err.Error()})
		return
	}

	attachmentInfo, err := models.GetTaskAttachmentByID(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message:": err.Error()})
		return
	}

	path := utils.UploadRoots + utils.TaskAttachmentDir + strconv.FormatInt(context.GetInt64("user_id"), 10) + "/" + *attachmentInfo.Attachment

	context.Header("Content-Disposition", "inline")
	context.File(path)
}

func downloadAttachmentFile(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message:": err.Error()})
		return
	}

	attachmentInfo, err := models.GetTaskAttachmentByID(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message:": err.Error()})
		return
	}

	path := utils.UploadRoots + utils.TaskAttachmentDir + strconv.FormatInt(context.GetInt64("user_id"), 10) + "/" + *attachmentInfo.Attachment

	context.FileAttachment(
		path,
		*attachmentInfo.Attachment,
	)
}

func updateTask(context *gin.Context) {
	// Set up needed variable
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	status_id, err := strconv.ParseInt(context.PostForm("status_id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid status_id"})
		return
	}

	tag_id, err := strconv.ParseInt(context.PostForm("tag_id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid tag_id"})
		return
	}

	due_date, err := time.Parse(time.RFC3339, context.PostForm("due_date"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid due_date"})
		return
	}

	// Get existed task
	taskDTO, err := models.GetTaskByID(id, context.GetInt64("user_id"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Get Data From Request
	var updatedTask = models.UpdateTaskRequest{
		Title:       context.PostForm("title"),
		Description: context.PostForm("description"),
		StatusID:    status_id,
		DueDate:     due_date,
		TagID:       tag_id,
	}

	// Handle attachment file
	var attachment *string
	file, err := context.FormFile("attachment")
	if file != nil {
		attachment, err = utils.SaveTaskAttachment(file, context)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		if taskDTO.Attachment != nil && *taskDTO.Attachment != "" {
			err = utils.RemoveFileAttachment(taskDTO.Attachment, context.GetInt64("user_id"))
			if err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
		}
	} else {
		attachment = taskDTO.Attachment
	}

	task := models.Task{
		ID:          taskDTO.ID,
		UsersID:     taskDTO.UsersID,
		Title:       updatedTask.Title,
		Description: updatedTask.Description,
		StatusID:    updatedTask.StatusID,
		DueDate:     updatedTask.DueDate,
		Attachment:  attachment,
		TagID:       updatedTask.TagID,
	}
	err = task.Update()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Task updated", "task": task})
}

func markTaskComplete(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	taskDTO, err := models.GetTaskByID(id, context.GetInt64("user_id"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	task := models.Task{
		ID:          taskDTO.ID,
		UsersID:     taskDTO.UsersID,
		Title:       taskDTO.Title,
		Description: taskDTO.Description,
		DueDate:     taskDTO.DueDate,
		Attachment:  taskDTO.Attachment,
	}
	err = task.CompleteTask()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Task marked complete"})
}

func restoreTask(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	taskDTO, err := models.GetTaskByID(id, context.GetInt64("user_id"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	task := models.Task{
		ID:          taskDTO.ID,
		UsersID:     taskDTO.UsersID,
		Title:       taskDTO.Title,
		Description: taskDTO.Description,
		DueDate:     taskDTO.DueDate,
		Attachment:  taskDTO.Attachment,
	}
	err = task.RestoreTask()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Task restored"})
}

func deleteTask(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	taskDTO, err := models.GetTaskByID(id, context.GetInt64("user_id"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	task := models.Task{
		ID:          taskDTO.ID,
		UsersID:     taskDTO.UsersID,
		Title:       taskDTO.Title,
		Description: taskDTO.Description,
		DueDate:     taskDTO.DueDate,
		Attachment:  taskDTO.Attachment,
	}
	err = task.Delete()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

func deleteTaskAttachment(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	taskDTO, err := models.GetTaskByID(id, context.GetInt64("user_id"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if taskDTO.Attachment != nil && *taskDTO.Attachment != "" {
		err = utils.RemoveFileAttachment(taskDTO.Attachment, context.GetInt64("user_id"))
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		task := models.Task{
			ID:          taskDTO.ID,
			UsersID:     taskDTO.UsersID,
			Title:       taskDTO.Title,
			Description: taskDTO.Description,
			DueDate:     taskDTO.DueDate,
			Attachment:  taskDTO.Attachment,
		}
		err = task.DeleteAttachment()
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		context.JSON(http.StatusOK, gin.H{"message": "Task attachment deleted"})
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Task has no attachment"})
	}
}
