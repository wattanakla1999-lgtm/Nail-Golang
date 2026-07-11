package handler

import (
	"nailly-back-end/internal/apperror"

	"github.com/gin-gonic/gin"
)

func respondError(c *gin.Context, err error) {
	appErr := apperror.FromError(err)
	c.JSON(appErr.Status, gin.H{"error": appErr.Message})
}
