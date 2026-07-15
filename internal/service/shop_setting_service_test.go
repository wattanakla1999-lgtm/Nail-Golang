package service

import (
	"errors"
	"nailly-back-end/internal/apperror"
	"nailly-back-end/internal/model"
	"net/http"
	"testing"
)

type fakeShopSettingStore struct {
	setting model.ShopSetting
	saves   int
}

func newFakeShopSettingStore() *fakeShopSettingStore {
	return &fakeShopSettingStore{setting: model.DefaultShopSetting()}
}

func (f *fakeShopSettingStore) Get() (model.ShopSetting, error) {
	return f.setting, nil
}

func (f *fakeShopSettingStore) Save(setting *model.ShopSetting) error {
	f.setting = *setting
	f.saves++
	return nil
}

func TestUpdateShopSettings(t *testing.T) {
	store := newFakeShopSettingStore()
	setting, err := NewShopSettingService(store).UpdateSettings(UpdateShopSettingInput{
		ShopStatus: "closed", OpenTime: "09:30", CloseTime: "19:00", ShopPhone: "081-234-5678",
	})
	if err != nil {
		t.Fatalf("UpdateSettings() error = %v", err)
	}
	if setting.ShopStatus != "closed" || setting.OpenTime != "09:30" || setting.CloseTime != "19:00" || setting.ShopPhone != "081-234-5678" {
		t.Fatalf("setting = %+v", setting)
	}
	if store.saves != 1 {
		t.Fatalf("Save() calls = %d, want 1", store.saves)
	}
}

func TestUpdateShopSettingsValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   UpdateShopSettingInput
		message string
	}{
		{
			name:    "invalid status",
			input:   UpdateShopSettingInput{ShopStatus: "pause", OpenTime: "10:00", CloseTime: "20:00", ShopPhone: "02-123-4567"},
			message: "shopStatus must be open or closed",
		},
		{
			name:    "invalid time format",
			input:   UpdateShopSettingInput{ShopStatus: "open", OpenTime: "10 AM", CloseTime: "20:00", ShopPhone: "02-123-4567"},
			message: "openTime must use HH:MM format",
		},
		{
			name:    "close before open",
			input:   UpdateShopSettingInput{ShopStatus: "open", OpenTime: "20:00", CloseTime: "10:00", ShopPhone: "02-123-4567"},
			message: "openTime must be before closeTime",
		},
		{
			name:    "missing phone",
			input:   UpdateShopSettingInput{ShopStatus: "open", OpenTime: "10:00", CloseTime: "20:00"},
			message: "shopPhone is required",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewShopSettingService(newFakeShopSettingStore()).UpdateSettings(test.input)
			var appErr *apperror.AppError
			if !errors.As(err, &appErr) || appErr.Status != http.StatusBadRequest || appErr.Message != test.message {
				t.Fatalf("UpdateSettings() error = %v, want 400 %q", err, test.message)
			}
		})
	}
}
