package repository

import (
	"nailly-back-end/internal/model"

	"gorm.io/gorm"
)

type ShopSettingRepository struct {
	db *gorm.DB
}

func NewShopSettingRepository(db *gorm.DB) *ShopSettingRepository {
	return &ShopSettingRepository{db: db}
}

func (r *ShopSettingRepository) Get() (model.ShopSetting, error) {
	setting := model.DefaultShopSetting()
	err := r.db.Where("id = ?", setting.ID).
		Attrs(model.DefaultShopSetting()).
		FirstOrCreate(&setting).Error
	return setting, err
}

func (r *ShopSettingRepository) Save(setting *model.ShopSetting) error {
	return r.db.Save(setting).Error
}
