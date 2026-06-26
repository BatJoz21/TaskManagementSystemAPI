package routes

import (
	"github.com/gin-gonic/gin"
	"taskmanagementsystem.localhost/tmsapi/middlewares"
)

func RegisteredRoutes(server *gin.Engine) {
	server.POST("/signup", signup)
	server.POST("/login", login)

	userGroup := server.Group("/")
	userGroup.Use(middlewares.Authenticate)
	userGroup.POST("tasks", createTask)
	userGroup.GET("tasks", getTasks)
	userGroup.GET("tasks/deleted", getDeletedTasks)
	userGroup.GET("tasks/:id", getTaskByID)
	userGroup.GET("tasks/:id/attachment/view", viewAttachmentFile)
	userGroup.GET("tasks/:id/attachment/download", downloadAttachmentFile)
	userGroup.GET("users/profile/:id", getUser)
	userGroup.GET("users/profile/:id/picture", getProfilePicture)
	userGroup.PUT("tasks/:id", updateTask)
	userGroup.PUT("tasks/:id/complete", markTaskComplete)
	userGroup.PUT("tasks/deleted/:id", restoreTask)
	userGroup.PUT("users/profile/:id", updateUser)
	userGroup.PUT("users/change-password", changeUserPassword)
	userGroup.DELETE("tasks/:id", deleteTask)
	userGroup.DELETE("tasks/:id/attachment/delete", deleteTaskAttachment)

	adminGroup := server.Group("/admin")
	adminGroup.Use(middlewares.Authenticate)
	adminGroup.Use(middlewares.AdminMiddlewares())
	adminGroup.GET("/dashboard", getDashboardData)
	adminGroup.GET("/users", getUsers)
	adminGroup.GET("/users/:id", getUser)
	adminGroup.PUT("/users/:id", updateUser)
}
