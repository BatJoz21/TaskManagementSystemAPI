package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"taskmanagementsystem.localhost/tmsapi/models"
	"taskmanagementsystem.localhost/tmsapi/utils"
)

func signup(context *gin.Context) {
	var registerStruct models.CreateUserStruct
	err := context.ShouldBindJSON(&registerStruct)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": false})
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
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": false})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "New user created!", "status": true})
}

func login(context *gin.Context) {
	// Get login data from login form
	var userLogin models.LoginUserStruct
	err := context.ShouldBindJSON(&userLogin)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user := models.User{
		Email:    userLogin.Email,
		Password: userLogin.Password,
	}

	// Checking credentials
	err = user.ValidateCredentials()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	role, err := user.GetUserRole()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.Email, role, user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Generate Refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.Email, role, user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Hashing token
	hashedToken := utils.HashRefreshToken(refreshToken)

	// Store refresh token in DB
	tokenStruct := models.RefreshTokenStruct{
		UserID:    user.ID,
		TokenHash: hashedToken,
	}
	err = tokenStruct.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	authUser := models.ResponseUserStruct{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      role,
	}
	context.JSON(http.StatusOK, gin.H{
		"message":       "Login successfull",
		"access_token":  token,
		"refresh_token": refreshToken,
		"user":          authUser,
	})
}

func logout(context *gin.Context) {
	var request models.RefreshTokenRequest
	err := context.ShouldBindJSON(&request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "isLoggedOut": false})
		return
	}

	// Get active refresh token from DB
	hashedRefreshedToken := utils.HashRefreshToken(request.RefreshToken)
	refreshToken, err := models.GetRefreshTokenByHash(hashedRefreshedToken)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error(), "isLoggedOut": false})
		return
	}

	// Revoke refresh token
	err = refreshToken.RevokeRefreshToken()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "isLoggedOut": false})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "user logged out", "isLoggedOut": true})
}

func refreshToken(context *gin.Context) {
	var request models.RefreshTokenRequest
	err := context.ShouldBindJSON(&request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	hashedRefreshedToken := utils.HashRefreshToken(request.RefreshToken)

	refreshToken, err := models.GetRefreshTokenByHash(hashedRefreshedToken)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	if refreshToken.RevokedAt != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "refresh token is invalid"})
		return
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		context.JSON(http.StatusUnauthorized, gin.H{"message": " refresh token is expired"})
		return
	}

	user, err := models.GetUserForJWT(refreshToken.UserID)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	accessToken, err := utils.GenerateToken(user.Email, user.Role, user.ID)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Token refreshed", "access_token": accessToken})
}
