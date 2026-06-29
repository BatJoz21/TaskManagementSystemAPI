package routes

import (
	"net/http"
	"strconv"

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

	authUser := models.ResponseUserStruct{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      role,
	}
	context.JSON(http.StatusOK, gin.H{"message": "Login successfull", "token": token, "user": authUser})
}

func getUsers(context *gin.Context) {
	sort := context.Query("sort")
	order := context.Query("order")

	users, total, err := models.GetUsers(sort, order)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"users": users, "total": total})
}

func getUser(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := models.GetUser(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, user)
}

func getProfilePicture(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Check if user exist
	user, err := models.GetUser(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	// Check if user has profile picture
	if user.ProfilePicture == nil {
		context.Status(http.StatusNoContent)
		return
	}

	path := utils.GetProfilePicturePath(user.ProfilePicture, id)

	context.File(path)
}

func updateUser(context *gin.Context) {
	// Set up needed variable
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Get User If Exist
	existUser, err := models.GetUser(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Get Data From Form
	var updatedUser = models.UpdateUserStruct{
		FirstName: context.PostForm("first_name"),
		LastName:  context.PostForm("last_name"),
		Role:      context.PostForm("role"),
	}

	// Handle upload user profile
	var profile_picture *string
	file, err := context.FormFile("profile_picture")
	if file != nil {
		profile_picture, err = utils.SaveProfilePicture(file, id, context)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		// Remove old profile picture
		if existUser.ProfilePicture != nil && *existUser.ProfilePicture != "" {
			err = utils.RemoveFileAttachment(existUser.ProfilePicture, utils.ProfilePictureDir, id)
			if err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
		}
	} else {
		profile_picture = existUser.ProfilePicture
	}

	user := models.User{
		ID:             existUser.ID,
		FirstName:      updatedUser.FirstName,
		LastName:       updatedUser.LastName,
		Role:           updatedUser.Role,
		ProfilePicture: profile_picture,
	}

	err = user.Update()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = user.UpdateRole()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User updated", "user": user})
}

func changeUserPassword(context *gin.Context) {
	var userAccount models.ChangePasswordStruct
	err := context.ShouldBindJSON(&userAccount)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if userAccount.NewPassword != userAccount.ConfirmNewPassword {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	user := models.User{
		Email:    userAccount.Email,
		Password: userAccount.CurrentPassword,
	}

	err = user.ValidateCredentials()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorize to change password"})
		return
	}

	user.Password = userAccount.NewPassword
	err = user.ChangePassword()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "password changed", "password": user})
}

func toggleBanUser(context *gin.Context) {
	// Set up needed variable
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Get User If Exist
	existUser, err := models.GetUser(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Set up User's new status
	isBanned := true
	if existUser.Status == nil || *existUser.Status == "" {
		isBanned = false
	}

	user := models.User{
		ID:        existUser.ID,
		FirstName: existUser.FirstName,
		LastName:  existUser.LastName,
		Status:    existUser.Status,
	}

	err = user.ToggleBan(isBanned)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

func softDeleteUser(context *gin.Context) {
	// Set up needed variable
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Get User If Exist
	existUser, err := models.GetUser(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Soft Delete user
	user := models.User{
		ID:        existUser.ID,
		FirstName: existUser.FirstName,
		LastName:  existUser.LastName,
	}

	err = user.SoftDelete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
