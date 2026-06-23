package routes

import (
	"github.com/gin-gonic/gin"
	"taskmanagementsystem.localhost/tmsapi/middlewares"
)

func RegisteredRoutes(server *gin.Engine) {
	server.POST("/signup", signup)
	server.POST("/login", login)

	taskGroup := server.Group("/")
	taskGroup.Use(middlewares.Authenticate)
	taskGroup.POST("tasks", createTask)
	taskGroup.GET("tasks", getTasks)
	taskGroup.GET("tasks/deleted", getDeletedTasks)
	taskGroup.GET("tasks/:id", getTaskByID)
	taskGroup.PUT("tasks/:id", updateTask)
	taskGroup.PUT("tasks/:id/complete", markTaskComplete)
	taskGroup.PUT("tasks/deleted/:id", restoreTask)
	taskGroup.DELETE("tasks/:id", deleteTask)
}
