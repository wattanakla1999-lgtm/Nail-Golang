package dto

import (
	"nailly-back-end/internal/model"
	"testing"

	"gorm.io/gorm"
)

func TestServiceResponseUsesGormIDAsForeignKey(t *testing.T) {
	response := ToServiceResponse(model.Service{
		Model: gorm.Model{ID: 12}, ServiceID: "SVC-012",
	})
	if response.ID != 12 || response.ServiceID != 12 || response.ServiceCode != "SVC-012" {
		t.Fatalf("response IDs = %+v", response)
	}
}

func TestNailTechnicianResponseUsesGormIDAsForeignKey(t *testing.T) {
	response := ToNailTechnicianResponse(model.NailTechnician{
		Model: gorm.Model{ID: 7}, TechnicianID: "TECH-007",
	})
	if response.ID != 7 || response.TechnicianID != 7 || response.TechnicianCode != "TECH-007" {
		t.Fatalf("response IDs = %+v", response)
	}
}
