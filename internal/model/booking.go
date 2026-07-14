package model

import (
	"time"

	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusInService BookingStatus = "in_service"
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusNoShow    BookingStatus = "no_show"
)

type Booking struct {
	gorm.Model

	BookingNo       string        `gorm:"type:varchar(50);not null;uniqueIndex" json:"bookingNo"`
	UserID          uint          `gorm:"not null;index" json:"userId"`
	ServiceID       uint          `gorm:"not null;index" json:"serviceId"`
	TechnicianID    *uint         `gorm:"index" json:"technicianId,omitempty"`
	StartAt         time.Time     `gorm:"not null;index" json:"startAt"`
	EndAt           time.Time     `gorm:"not null;check:end_at > start_at" json:"endAt"`
	CustomerName    string        `gorm:"type:varchar(255);not null" json:"customerName"`
	CustomerPhone   string        `gorm:"type:varchar(50);not null" json:"customerPhone"`
	ServiceName     string        `gorm:"type:varchar(255);not null" json:"serviceName"`
	Price           int           `gorm:"not null;check:price >= 0" json:"price"`
	DurationMinutes int           `gorm:"not null;check:duration_minutes > 0" json:"durationMinutes"`
	Status          BookingStatus `gorm:"type:varchar(20);not null;default:pending;index;check:status IN ('pending','confirmed','in_service','completed','cancelled','no_show')" json:"status"`
	Note            string        `gorm:"type:text" json:"note,omitempty"`
	CancelReason    string        `gorm:"type:text" json:"cancelReason,omitempty"`

	User       User            `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	Service    Service         `gorm:"foreignKey:ServiceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	Technician *NailTechnician `gorm:"foreignKey:TechnicianID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
}

func (Booking) TableName() string {
	return "bookings"
}

func IsValidBookingStatus(status BookingStatus) bool {
	switch status {
	case BookingStatusPending, BookingStatusConfirmed, BookingStatusInService,
		BookingStatusCompleted, BookingStatusCancelled, BookingStatusNoShow:
		return true
	default:
		return false
	}
}
