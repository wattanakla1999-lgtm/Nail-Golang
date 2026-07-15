package router

import (
	"nailly-back-end/internal/handler"
	"nailly-back-end/internal/repository"
	"nailly-back-end/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterReportRoutes(api *gin.RouterGroup, db *gorm.DB, requireAdmin gin.HandlerFunc) {
	reportRepository := repository.NewReportRepository(db)
	reportService := service.NewReportService(reportRepository)
	reportHandler := handler.NewReportHandler(reportService)

	api.GET("/reports", requireAdmin, reportHandler.GetReport)
}
