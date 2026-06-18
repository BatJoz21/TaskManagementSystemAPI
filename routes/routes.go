package routes

import "github.com/gin-gonic/gin"

func RegisteredRoutes(server *gin.Engine) {
	server.POST("/tasks", createTask)
	server.GET("/tasks", getTasks)
	server.GET("/tasks/:id", getTaskByID)
	server.PUT("/tasks/:id", updateTask)
	server.DELETE("/tasks/:id", deleteTask)
}
