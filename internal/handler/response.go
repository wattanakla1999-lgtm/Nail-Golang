package handler

import (
	"nailly-back-end/internal/apperror"

	"github.com/gin-gonic/gin"
)

func respondError(c *gin.Context, err error) {
	appErr := apperror.FromError(err)
	response := gin.H{"error": appErr.Message}
	if appErr.Code != "" {
		response["code"] = appErr.Code
	}
	c.JSON(appErr.Status, response)
}
