package model

import "gorm.io/gorm"

type ShopSetting struct {
	gorm.Model

	ShopStatus string `gorm:"type:varchar(10);not null;default:open;check:shop_status IN ('open','closed')" json:"shopStatus"`
	OpenTime   string `gorm:"type:varchar(5);not null;default:10:00" json:"openTime"`
	CloseTime  string `gorm:"type:varchar(5);not null;default:20:00" json:"closeTime"`
	ShopPhone  string `gorm:"type:varchar(50);not null;default:02-123-4567" json:"shopPhone"`
}

func (ShopSetting) TableName() string {
	return "shop_settings"
}

func DefaultShopSetting() ShopSetting {
	return ShopSetting{
		Model:      gorm.Model{ID: 1},
		ShopStatus: "open",
		OpenTime:   "10:00",
		CloseTime:  "20:00",
		ShopPhone:  "02-123-4567",
	}
}
