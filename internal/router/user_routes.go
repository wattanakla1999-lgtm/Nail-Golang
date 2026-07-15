package router

import (
	"nailly-back-end/internal/handler"
	"nailly-back-end/internal/repository"
	"nailly-back-end/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(api *gin.RouterGroup, db *gorm.DB, requireAdmin gin.HandlerFunc) {
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	users := api.Group("/users")
	users.Use(requireAdmin)
	users.GET("", userHandler.GetUsers)
	users.GET("/email/:email", userHandler.GetUserByEmail)
	users.GET("/age/:age", userHandler.GetUsersOlderThan)
	users.GET("/:id", userHandler.GetUserByID)
	users.POST("", userHandler.CreateUser)
	users.PUT("/:id", userHandler.UpdateUser)
	users.DELETE("/:id", userHandler.DeleteUser)
}
