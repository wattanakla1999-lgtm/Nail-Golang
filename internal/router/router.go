package router

import (
	"nailly-back-end/internal/handler"
	"nailly-back-end/internal/repository"
	"nailly-back-end/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func New(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	rootHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "nailly-api",
			"status":  "running",
			"version": "v1.0.1",
		})
	}
	healthHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}

	r.GET("/", rootHandler)
	r.HEAD("/", rootHandler)
	r.GET("/health", healthHandler)
	r.HEAD("/health", healthHandler)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	api := r.Group("/api")
	api.GET("/users", userHandler.GetUsers)
	api.GET("/users/email/:email", userHandler.GetUserByEmail)
	api.GET("/users/age/:age", userHandler.GetUsersOlderThan)
	api.GET("/users/:id", userHandler.GetUserByID)
	api.POST("/users", userHandler.CreateUser)
	api.PUT("/users/:id", userHandler.UpdateUser)
	api.DELETE("/users/:id", userHandler.DeleteUser)

	return r
}
