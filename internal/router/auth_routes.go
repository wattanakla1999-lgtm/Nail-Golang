package router

import (
	"nailly-back-end/internal/handler"
	"nailly-back-end/internal/repository"
	"nailly-back-end/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(api *gin.RouterGroup, db *gorm.DB, jwtManager *service.JWTManager, requireAdmin gin.HandlerFunc) {
	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository, jwtManager)
	authHandler := handler.NewAuthHandler(authService)

	auth := api.Group("/auth")
	auth.POST("/login", authHandler.Login)
	auth.GET("/me", requireAdmin, authHandler.Me)
}
