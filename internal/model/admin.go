package model

import "gorm.io/gorm"

type Admin struct {
	gorm.Model

	Username     string `gorm:"type:varchar(100);not null;uniqueIndex" json:"username"`
	Name         string `gorm:"type:varchar(255);not null" json:"name"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	Role         string `gorm:"type:varchar(20);not null;default:admin" json:"role"`
	Active       bool   `gorm:"not null;default:true" json:"active"`
}

func (Admin) TableName() string {
	return "admins"
}
