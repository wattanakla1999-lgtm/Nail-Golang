package handler

import (
	"nailly-back-end/internal/dto"
	"nailly-back-end/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service *service.DashboardService
}

func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: dashboardService}
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	stats, err := h.service.GetStats()
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.ToDashboardStatsResponse(stats))
}
