package routes

import "github.com/gin-gonic/gin"

func RegisteredRoutes(server *gin.Engine) {
	server.GET("/tasks", getTasks)
}
