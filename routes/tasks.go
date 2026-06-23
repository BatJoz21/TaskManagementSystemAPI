package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"taskmanagementsystem.localhost/tmsapi/models"
)

func createTask(context *gin.Context) {
	var taskDTO models.CreateTaskRequest
	err := context.ShouldBindJSON(&taskDTO)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	task := models.Task{
		UsersID:     context.GetInt64("user_id"),
		Title:       taskDTO.Title,
		Description: taskDTO.Description,
		StatusID:    taskDTO.StatusID,
		DueDate:     taskDTO.DueDate,
		Attachment:  taskDTO.Attachment,
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

func updateTask(context *gin.Context) {
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

	var updatedTask models.UpdateTaskRequest
	err = context.ShouldBindJSON(&updatedTask)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	task := models.Task{
		ID:          taskDTO.ID,
		UsersID:     taskDTO.UsersID,
		Title:       updatedTask.Title,
		Description: updatedTask.Description,
		StatusID:    updatedTask.StatusID,
		DueDate:     updatedTask.DueDate,
		Attachment:  updatedTask.Attachment,
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
