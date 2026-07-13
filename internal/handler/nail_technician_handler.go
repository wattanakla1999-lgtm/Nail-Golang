package handler

import (
	"nailly-back-end/internal/apperror"
	"nailly-back-end/internal/dto"
	"nailly-back-end/internal/repository"
	"nailly-back-end/internal/service"
	"nailly-back-end/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NailTechnicianHandler struct {
	service *service.NailTechnicianService
}

func NewNailTechnicianHandler(service *service.NailTechnicianService) *NailTechnicianHandler {
	return &NailTechnicianHandler{service: service}
}

func (h *NailTechnicianHandler) GetNailTechnicians(c *gin.Context) {
	filter := repository.NailTechnicianFilter{
		TechnicianName: c.Query("technician_name"),
		Phone:          c.Query("phone"),
		Specialty:      c.Query("specialty"),
	}
	pagination := utils.NewPagination(c.DefaultQuery("page", "1"), c.DefaultQuery("limit", "10"))

	technicians, total, err := h.service.GetNailTechnicians(filter, pagination)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Data:  dto.ToNailTechnicianResponses(technicians),
		Page:  pagination.Page,
		Limit: pagination.Limit,
		Total: total,
	})
}

func (h *NailTechnicianHandler) GetNailTechnicianByID(c *gin.Context) {
	technician, err := h.service.GetNailTechnicianByID(c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToNailTechnicianResponse(technician))
}

func (h *NailTechnicianHandler) CreateNailTechnician(c *gin.Context) {
	var input dto.CreateNailTechnicianRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}

	technician, err := h.service.CreateNailTechnician(input.ToModel())
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToNailTechnicianResponse(technician))
}

func (h *NailTechnicianHandler) UpdateNailTechnician(c *gin.Context) {
	var input dto.UpdateNailTechnicianRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, apperror.BadRequest("invalid request body", err))
		return
	}

	technician, err := h.service.UpdateNailTechnician(c.Param("id"), input.ToModel())
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToNailTechnicianResponse(technician))
}

func (h *NailTechnicianHandler) DeleteNailTechnician(c *gin.Context) {
	if err := h.service.DeleteNailTechnician(c.Param("id")); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "nail technician deleted"})
}
