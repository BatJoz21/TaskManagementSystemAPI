package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminMiddlewares() gin.HandlerFunc {
	return func(context *gin.Context) {
		role := context.GetString("role")

		if role != "admin" && role != "developer" {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "access denied"})
			return
		}

		context.Next()
	}
}
