package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Name  string  `json:"name" gorm:"not null"`
	Email string  `json:"email" gorm:"uniqueIndex;not null"`
	Phone *string `json:"phone,omitempty" gorm:"type:varchar(50);uniqueIndex"`
	Age   int     `json:"age"`
}

func (User) TableName() string {
	return "user_dbs"
}
