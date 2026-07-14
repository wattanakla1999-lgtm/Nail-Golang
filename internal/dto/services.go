package dto

import (
	"nailly-back-end/internal/model"
	"time"
)

type CreateServiceRequest struct {
	ServiceID    string    `json:"serviceId"`
	ServiceName  string    `gorm:"type:varchar(255);not null" json:"serviceName"`
	ServicePrice int       `gorm:"not null" json:"servicePrice"`
	Duration     int       `gorm:"not null" json:"duration"`
	ServiceImg   string    `gorm:"type:varchar(500)" json:"img,omitempty"`
	Popular      bool      `gorm:"default:false" json:"popular"`
	Description  string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type UpdateServiceRequest struct {
	ServiceName  string `gorm:"type:varchar(255);not null" json:"serviceName"`
	ServicePrice int    `gorm:"not null" json:"servicePrice"`
	Duration     int    `gorm:"not null" json:"duration"`
	ServiceImg   string `gorm:"type:varchar(500)" json:"img,omitempty"`
	Popular      bool   `gorm:"default:false" json:"popular"`
	Description  string `gorm:"type:text" json:"description,omitempty"`
}

type ServiceResponse struct {
	ID           uint      `json:"id"`
	ServiceID    uint      `json:"serviceId"`
	ServiceCode  string    `json:"serviceCode,omitempty"`
	ServiceName  string    `json:"serviceName"`
	ServicePrice int       `json:"servicePrice"`
	Duration     int       `json:"duration"`
	ServiceImg   string    `json:"img,omitempty"`
	Popular      bool      `json:"popular"`
	Description  string    `json:"description,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func ToServiceResponse(service model.Service) ServiceResponse {
	return ServiceResponse{
		ID:           service.ID,
		ServiceID:    service.ID,
		ServiceCode:  service.ServiceID,
		ServiceName:  service.ServiceName,
		ServicePrice: service.ServicePrice,
		Duration:     service.Duration,
		ServiceImg:   service.ServiceImg,
		Popular:      service.Popular,
		Description:  service.Description,
		CreatedAt:    service.CreatedAt.In(thailandLocation),
		UpdatedAt:    service.UpdatedAt.In(thailandLocation),
	}
}

func ToServiceResponses(services []model.Service) []ServiceResponse {
	responses := make([]ServiceResponse, 0, len(services))
	for _, service := range services {
		responses = append(responses, ToServiceResponse(service))
	}

	return responses
}

func (r CreateServiceRequest) ToModel() model.Service {
	return model.Service{
		ServiceID:    r.ServiceID,
		ServiceName:  r.ServiceName,
		ServicePrice: r.ServicePrice,
		Duration:     r.Duration,
		ServiceImg:   r.ServiceImg,
		Popular:      r.Popular,
		Description:  r.Description,
	}
}

func (r UpdateServiceRequest) ToModel() model.Service {
	return model.Service{
		ServiceName:  r.ServiceName,
		ServicePrice: r.ServicePrice,
		Duration:     r.Duration,
		ServiceImg:   r.ServiceImg,
		Popular:      r.Popular,
		Description:  r.Description,
	}
}
