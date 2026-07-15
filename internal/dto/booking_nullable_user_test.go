package dto

import (
	"encoding/json"
	"nailly-back-end/internal/model"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestCreateBookingRequestAcceptsNullUserID(t *testing.T) {
	var request CreateBookingRequest
	err := json.Unmarshal([]byte(`{
		"userId": null,
		"serviceId": 2,
		"technicianId": 1,
		"startAt": "2026-07-15T10:00:00+07:00",
		"customerName": "Walk-in",
		"customerPhone": "0812345678"
	}`), &request)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if request.UserID != nil {
		t.Fatalf("UserID = %v, want nil", request.UserID)
	}
}

func TestBookingResponseReturnsNullUser(t *testing.T) {
	response := ToBookingResponse(model.Booking{
		Model: gorm.Model{ID: 1}, ServiceID: 2,
		StartAt: time.Now(), EndAt: time.Now().Add(time.Hour),
	})
	if response.UserID != nil || response.User != nil {
		t.Fatalf("response user = (%v, %+v), want nil", response.UserID, response.User)
	}
}
