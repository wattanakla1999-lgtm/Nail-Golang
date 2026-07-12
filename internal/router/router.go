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
	keepAliveHandler := func(c *gin.Context) {
		var result int
		if err := db.WithContext(c.Request.Context()).Raw("SELECT 1").Scan(&result).Error; err != nil {
			if c.Request.Method == http.MethodHead {
				c.Status(http.StatusServiceUnavailable)
				return
			}

			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "error",
				"database": "unreachable",
			})
			return
		}

		if c.Request.Method == http.MethodHead {
			c.Status(http.StatusOK)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"database": "active",
		})
	}

	r.GET("/", rootHandler)
	r.HEAD("/", rootHandler)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	api := r.Group("/api")
	api.GET("/keep-alive", keepAliveHandler)
	api.HEAD("/keep-alive", keepAliveHandler)
	api.GET("/users", userHandler.GetUsers)
	api.GET("/users/email/:email", userHandler.GetUserByEmail)
	api.GET("/users/age/:age", userHandler.GetUsersOlderThan)
	api.GET("/users/:id", userHandler.GetUserByID)
	api.POST("/users", userHandler.CreateUser)
	api.PUT("/users/:id", userHandler.UpdateUser)
	api.DELETE("/users/:id", userHandler.DeleteUser)

	return r
}
