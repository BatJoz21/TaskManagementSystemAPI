package main

import (
	"github.com/gin-gonic/gin"
	"taskmanagementsystem.localhost/tmsapi/database"
	"taskmanagementsystem.localhost/tmsapi/routes"
)

func main() {
	database.InitDB()
	server := gin.Default()

	routes.RegisteredRoutes(server)

	server.Run(":8080")
}
