package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Name  string `json:"name" gorm:"not null"`
	Email string `json:"email" gorm:"uniqueIndex;not null"`
	Age   int    `json:"age"`
}

func (User) TableName() string {
	return "user_dbs"
}
