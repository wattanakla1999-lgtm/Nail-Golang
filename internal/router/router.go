package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func New(db *gorm.DB, allowOrigin string) *gin.Engine {
	r := gin.Default()
	r.Use(corsMiddleware(allowOrigin))

	// Root
	r.GET("/", rootHandler)
	r.HEAD("/", rootHandler)

	api := r.Group("/api")

	// Keep-alive
	api.GET("/keep-alive", keepAliveHandler(db))
	api.HEAD("/keep-alive", keepAliveHandler(db))

	// Users
	RegisterUserRoutes(api, db)

	// Services
	RegisterServiceRoutes(api, db)

	// Nail Technicians
	RegisterNailTechnicianRoutes(api, db)

	// Bookings
	RegisterBookingRoutes(api, db)

	return r
}

func rootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service": "nailly-api",
		"status":  "running",
		"version": "v1.0.1",
	})
}

func keepAliveHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
}

func corsMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", allowOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
