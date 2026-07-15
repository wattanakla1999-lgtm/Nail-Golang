package repository

import (
	"nailly-back-end/internal/model"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) FindAdminByUsername(username string) (model.Admin, error) {
	var admin model.Admin
	err := r.db.Where("username = ?", username).First(&admin).Error
	return admin, err
}

func (r *AuthRepository) CreateAdmin(admin *model.Admin) error {
	return r.db.Create(admin).Error
}

func (r *AuthRepository) SaveAdmin(admin *model.Admin) error {
	return r.db.Save(admin).Error
}
