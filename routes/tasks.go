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
		context.JSON(http.StatusBadRequest, gin.H{"message:": err.Error()})
		return
	}
	taskDTO.UsersID = 3

	task := models.Task{
		UsersID:     taskDTO.UsersID,
		Title:       taskDTO.Title,
		Description: taskDTO.Description,
		StatusID:    taskDTO.StatusID,
		DueDate:     taskDTO.DueDate,
		Attachment:  taskDTO.Attachment,
		TagID:       taskDTO.TagID,
	}

	err = task.Save()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message:": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message:": "task created", "task:": task})
}

func getTasks(context *gin.Context) {
	tasks, err := models.GetAllTasks()
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
	task, err := models.GetTaskByID(id)

	context.JSON(http.StatusOK, task)
}

func updateTask(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message:": err.Error()})
		return
	}

	taskDTO, err := models.GetTaskByID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message:": err.Error()})
		return
	}

	var updatedTask models.UpdateTaskRequest
	err = context.ShouldBindJSON(&updatedTask)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message:": err.Error()})
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
		context.JSON(http.StatusBadRequest, gin.H{"message:": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message:": "Task updated", "task:": task})
}

func deleteTask(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message:": err.Error()})
		return
	}

	taskDTO, err := models.GetTaskByID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message:": err.Error()})
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
		context.JSON(http.StatusBadRequest, gin.H{"message:": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message:": "Task deleted"})
}
