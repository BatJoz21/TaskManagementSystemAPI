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
	taskGroup.GET("tasks/:id", getTaskByID)
	taskGroup.PUT("tasks/:id", updateTask)
	taskGroup.DELETE("tasks/:id", deleteTask)
}
