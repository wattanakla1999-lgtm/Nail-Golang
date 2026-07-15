package router

import (
	"nailly-back-end/internal/handler"
	"nailly-back-end/internal/repository"
	"nailly-back-end/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterDashboardRoutes(api *gin.RouterGroup, db *gorm.DB, requireAdmin gin.HandlerFunc) {
	dashboardRepository := repository.NewDashboardRepository(db)
	dashboardService := service.NewDashboardService(dashboardRepository)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)

	dashboard := api.Group("/dashboard")
	dashboard.Use(requireAdmin)
	dashboard.GET("/stats", dashboardHandler.GetStats)
}
