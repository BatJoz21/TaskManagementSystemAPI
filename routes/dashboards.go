package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"taskmanagementsystem.localhost/tmsapi/models"
)

func getDashboardData(context *gin.Context) {
	data, err := models.GetDataForDashboard()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, data)
}
