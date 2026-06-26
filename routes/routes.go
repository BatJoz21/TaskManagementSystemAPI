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
	taskGroup.GET("tasks/:id/attachment/view", viewAttachmentFile)
	taskGroup.GET("tasks/:id/attachment/download", downloadAttachmentFile)
	taskGroup.PUT("tasks/:id", updateTask)
	taskGroup.PUT("tasks/:id/complete", markTaskComplete)
	taskGroup.PUT("tasks/deleted/:id", restoreTask)
	taskGroup.DELETE("tasks/:id", deleteTask)
	taskGroup.DELETE("tasks/:id/attachment/delete", deleteTaskAttachment)

	adminGroup := server.Group("/admin")
	adminGroup.Use(middlewares.Authenticate)
	adminGroup.Use(middlewares.AdminMiddlewares())
	adminGroup.GET("/dashboard", getDashboardData)
	adminGroup.GET("/users", getUsers)
	adminGroup.GET("/users/:id", getUser)
	adminGroup.PUT("/users/:id", updateUser)
}
