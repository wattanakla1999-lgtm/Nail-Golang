package handler

import (
	"nailly-back-end/internal/dto"
	"nailly-back-end/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	service *service.ReportService
}

func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{service: reportService}
}

func (h *ReportHandler) GetReport(c *gin.Context) {
	report, err := h.service.GetReport(c.DefaultQuery("period", "week"))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.ToReportResponse(report))
}
