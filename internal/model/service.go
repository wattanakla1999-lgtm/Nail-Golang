package model

import "gorm.io/gorm"

type Service struct {
	gorm.Model

	ServiceID    string `gorm:"column:service_id;type:varchar(50);not null" json:"serviceId"`
	ServiceName  string `gorm:"type:varchar(255);not null" json:"name"`
	ServicePrice int    `gorm:"not null" json:"price"`
	Duration     int    `gorm:"not null" json:"duration"`
	ServiceImg   string `gorm:"type:varchar(500)" json:"img,omitempty"`
	Popular      bool   `gorm:"default:false" json:"popular"`
	Description  string `gorm:"type:text" json:"description,omitempty"`
}

func (Service) TableName() string {
	return "service_dbs"
}
