package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"taskmanagementsystem.localhost/tmsapi/models"
)

func getTasks(context *gin.Context) {
	tasks, err := models.GetAllTasks()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message:": err.Error()})
		return
	}

	context.JSON(http.StatusOK, tasks)
}
