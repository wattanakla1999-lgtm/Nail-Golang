package dto

import (
	"nailly-back-end/internal/model"
	"time"
)

type CreateNailTechnicianRequest struct {
	TechnicianID    string    `json:"technicianId"`
	TechnicianName  string    `json:"technicianName"`
	Phone           string    `json:"phone,omitempty"`
	ExperienceYears int       `json:"experienceYears"`
	Specialty       string    `json:"specialty,omitempty"`
	ProfileImg      string    `json:"profileImg,omitempty"`
	Active          bool      `json:"active"`
	Bio             string    `json:"bio,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type UpdateNailTechnicianRequest struct {
	TechnicianName  string `json:"technicianName"`
	Phone           string `json:"phone,omitempty"`
	ExperienceYears int    `json:"experienceYears"`
	Specialty       string `json:"specialty,omitempty"`
	ProfileImg      string `json:"profileImg,omitempty"`
	Active          bool   `json:"active"`
	Bio             string `json:"bio,omitempty"`
}

type NailTechnicianResponse struct {
	TechnicianID    string    `json:"technicianId"`
	TechnicianName  string    `json:"technicianName"`
	Phone           string    `json:"phone,omitempty"`
	ExperienceYears int       `json:"experienceYears"`
	Specialty       string    `json:"specialty,omitempty"`
	ProfileImg      string    `json:"profileImg,omitempty"`
	Active          bool      `json:"active"`
	Bio             string    `json:"bio,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func ToNailTechnicianResponse(technician model.NailTechnician) NailTechnicianResponse {
	return NailTechnicianResponse{
		TechnicianID:    technician.TechnicianID,
		TechnicianName:  technician.TechnicianName,
		Phone:           technician.Phone,
		ExperienceYears: technician.ExperienceYears,
		Specialty:       technician.Specialty,
		ProfileImg:      technician.ProfileImg,
		Active:          technician.Active,
		Bio:             technician.Bio,
		CreatedAt:       technician.CreatedAt.In(thailandLocation),
		UpdatedAt:       technician.UpdatedAt.In(thailandLocation),
	}
}

func ToNailTechnicianResponses(technicians []model.NailTechnician) []NailTechnicianResponse {
	responses := make([]NailTechnicianResponse, 0, len(technicians))
	for _, technician := range technicians {
		responses = append(responses, ToNailTechnicianResponse(technician))
	}

	return responses
}

func (r CreateNailTechnicianRequest) ToModel() model.NailTechnician {
	return model.NailTechnician{
		TechnicianID:    r.TechnicianID,
		TechnicianName:  r.TechnicianName,
		Phone:           r.Phone,
		ExperienceYears: r.ExperienceYears,
		Specialty:       r.Specialty,
		ProfileImg:      r.ProfileImg,
		Active:          r.Active,
		Bio:             r.Bio,
	}
}

func (r UpdateNailTechnicianRequest) ToModel() model.NailTechnician {
	return model.NailTechnician{
		TechnicianName:  r.TechnicianName,
		Phone:           r.Phone,
		ExperienceYears: r.ExperienceYears,
		Specialty:       r.Specialty,
		ProfileImg:      r.ProfileImg,
		Active:          r.Active,
		Bio:             r.Bio,
	}
}
