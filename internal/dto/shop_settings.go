package dto

import "nailly-back-end/internal/model"

type UpdateShopSettingRequest struct {
	ShopStatus string `json:"shopStatus" binding:"required"`
	OpenTime   string `json:"openTime" binding:"required"`
	CloseTime  string `json:"closeTime" binding:"required"`
	ShopPhone  string `json:"shopPhone" binding:"required"`
}

type ShopSettingResponse struct {
	ShopStatus string `json:"shopStatus"`
	OpenTime   string `json:"openTime"`
	CloseTime  string `json:"closeTime"`
	ShopPhone  string `json:"shopPhone"`
}

func ToShopSettingResponse(setting model.ShopSetting) ShopSettingResponse {
	return ShopSettingResponse{
		ShopStatus: setting.ShopStatus,
		OpenTime:   setting.OpenTime,
		CloseTime:  setting.CloseTime,
		ShopPhone:  setting.ShopPhone,
	}
}
