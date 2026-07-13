package router

import (
	"nailly-back-end/internal/handler"
	"nailly-back-end/internal/repository"
	"nailly-back-end/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterNailTechnicianRoutes(api *gin.RouterGroup, db *gorm.DB) {
	nailTechnicianRepository := repository.NewNailTechnicianRepository(db)
	nailTechnicianService := service.NewNailTechnicianService(nailTechnicianRepository)
	nailTechnicianHandler := handler.NewNailTechnicianHandler(nailTechnicianService)

	nailTechnicians := api.Group("/nail_technician")
	nailTechnicians.GET("", nailTechnicianHandler.GetNailTechnicians)
	nailTechnicians.GET("/:id", nailTechnicianHandler.GetNailTechnicianByID)
	nailTechnicians.POST("", nailTechnicianHandler.CreateNailTechnician)
	nailTechnicians.PUT("/:id", nailTechnicianHandler.UpdateNailTechnician)
	nailTechnicians.DELETE("/:id", nailTechnicianHandler.DeleteNailTechnician)
}
