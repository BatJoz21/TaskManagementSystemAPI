package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"taskmanagementsystem.localhost/tmsapi/models"
	"taskmanagementsystem.localhost/tmsapi/utils"
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
		LastName:  registerStruct.LastName,
		Email:     registerStruct.Email,
		Password:  registerStruct.Password,
	}
	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message:": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message:": "New user created!"})
}

func login(context *gin.Context) {
	var userLogin models.LoginUserStruct
	err := context.ShouldBindJSON(&userLogin)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message:": err.Error()})
		return
	}

	user := models.User{
		Email:    userLogin.Email,
		Password: userLogin.Password,
	}
	err = user.ValidateCredentials()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message:": err.Error()})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message:": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message:": "Login successfull", "token": token})
}
