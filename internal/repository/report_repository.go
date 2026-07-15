package repository

import (
	"nailly-back-end/internal/model"
	"time"

	"gorm.io/gorm"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) FindBookings(startAt, endAt time.Time) ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.db.
		Preload("Technician").
		Where("start_at >= ? AND start_at < ?", startAt, endAt).
		Order("start_at DESC").
		Find(&bookings).Error
	return bookings, err
}
