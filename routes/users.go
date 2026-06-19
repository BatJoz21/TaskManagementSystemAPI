package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"taskmanagementsystem.localhost/tmsapi/models"
)

func signup(context *gin.Context) {
	var registerStruct models.CreateUserStruct
	err := context.ShouldBindJSON(&registerStruct)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message:": err.Error()})
		return
	}

	user := models.User{
		FirstName: registerStruct.FirstName,
		LastName: registerStruct.LastName,
		Email: registerStruct.Email,
		Password: registerStruct.Password,
	}
	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message:": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message:": "New user created!"})
}
