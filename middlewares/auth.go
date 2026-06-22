package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"taskmanagementsystem.localhost/tmsapi/utils"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")
	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message:": "Not authorized."})
	}

	token = strings.TrimPrefix(token, "Bearer ")
	user_id, err := utils.VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message:": "Not authorized."})
	}

	context.Set("user_id", user_id)
	context.Next()
}
