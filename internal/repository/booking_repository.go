package repository

import (
	"errors"
	"nailly-back-end/internal/model"
	"nailly-back-end/pkg/utils"
	"time"

	"gorm.io/gorm"
)

var ErrTechnicianOverlap = errors.New("technician booking time overlaps")

type BookingFilter struct {
	UserID       *uint
	ServiceID    *uint
	TechnicianID *uint
	Status       model.BookingStatus
	DateFrom     *time.Time
	DateTo       *time.Time
}

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func bookingRelations(db *gorm.DB) *gorm.DB {
	return db.Preload("User").Preload("Service").Preload("Technician")
}

func (r *BookingRepository) FindAll(filter BookingFilter, pagination utils.Pagination) ([]model.Booking, int64, error) {
	var bookings []model.Booking
	query := r.db.Model(&model.Booking{})

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.ServiceID != nil {
		query = query.Where("service_id = ?", *filter.ServiceID)
	}
	if filter.TechnicianID != nil {
		query = query.Where("technician_id = ?", *filter.TechnicianID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.DateFrom != nil {
		query = query.Where("start_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("start_at <= ?", *filter.DateTo)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := bookingRelations(query).
		Order("start_at ASC").
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(&bookings).Error
	return bookings, total, err
}

func (r *BookingRepository) FindByID(id uint) (model.Booking, error) {
	var booking model.Booking
	err := bookingRelations(r.db).First(&booking, id).Error
	return booking, err
}

func (r *BookingRepository) FindUserByID(id uint) (model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return user, err
}

func (r *BookingRepository) FindServiceByID(id uint) (model.Service, error) {
	var service model.Service
	err := r.db.First(&service, id).Error
	return service, err
}

func (r *BookingRepository) FindTechnicianByID(id uint) (model.NailTechnician, error) {
	var technician model.NailTechnician
	err := r.db.First(&technician, id).Error
	return technician, err
}

func (r *BookingRepository) HasTechnicianOverlap(technicianID uint, startAt, endAt time.Time, excludeBookingID uint) (bool, error) {
	query := r.db.Model(&model.Booking{}).
		Where("technician_id = ?", technicianID).
		Where("status NOT IN ?", []model.BookingStatus{model.BookingStatusCancelled, model.BookingStatusNoShow}).
		Where("start_at < ? AND end_at > ?", endAt, startAt)
	if excludeBookingID != 0 {
		query = query.Where("id <> ?", excludeBookingID)
	}

	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *BookingRepository) FindBusyBookings(startAt, endAt time.Time, technicianID *uint) ([]model.Booking, error) {
	var bookings []model.Booking
	query := r.db.Model(&model.Booking{}).
		Where("status NOT IN ?", []model.BookingStatus{model.BookingStatusCancelled, model.BookingStatusNoShow}).
		Where("start_at < ? AND end_at > ?", endAt, startAt)
	if technicianID != nil {
		query = query.Where("technician_id = ?", *technicianID)
	}
	err := query.Order("start_at ASC").Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) FindActiveTechnicianIDs() ([]uint, error) {
	var ids []uint
	err := r.db.Model(&model.NailTechnician{}).
		Where("active = ?", true).
		Order("id ASC").
		Pluck("id", &ids).Error
	return ids, err
}

func (r *BookingRepository) Create(booking *model.Booking) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := lockAndCheckOverlap(tx, booking); err != nil {
			return err
		}
		return tx.Omit("User", "Service", "Technician").Create(booking).Error
	})
}

func (r *BookingRepository) Update(booking *model.Booking) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := lockAndCheckOverlap(tx, booking); err != nil {
			return err
		}
		return tx.Omit("User", "Service", "Technician").Save(booking).Error
	})
}

func (r *BookingRepository) Delete(booking *model.Booking) error {
	return r.db.Delete(booking).Error
}

func lockAndCheckOverlap(tx *gorm.DB, booking *model.Booking) error {
	if booking.TechnicianID == nil || booking.Status == model.BookingStatusCancelled || booking.Status == model.BookingStatusNoShow {
		return nil
	}

	// Serializes schedule writes per technician for the duration of this transaction.
	if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", int64(*booking.TechnicianID)).Error; err != nil {
		return err
	}
	var count int64
	err := tx.Model(&model.Booking{}).
		Where("technician_id = ?", *booking.TechnicianID).
		Where("status NOT IN ?", []model.BookingStatus{model.BookingStatusCancelled, model.BookingStatusNoShow}).
		Where("start_at < ? AND end_at > ?", booking.EndAt, booking.StartAt).
		Where("id <> ?", booking.ID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrTechnicianOverlap
	}
	return nil
}
