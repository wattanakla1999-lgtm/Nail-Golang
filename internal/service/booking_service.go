package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"nailly-back-end/internal/apperror"
	"nailly-back-end/internal/model"
	"nailly-back-end/internal/repository"
	"nailly-back-end/pkg/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

type BookingStore interface {
	FindAll(filter repository.BookingFilter, pagination utils.Pagination) ([]model.Booking, int64, error)
	FindByID(id uint) (model.Booking, error)
	FindUserByID(id uint) (model.User, error)
	FindServiceByID(id uint) (model.Service, error)
	FindTechnicianByID(id uint) (model.NailTechnician, error)
	HasTechnicianOverlap(technicianID uint, startAt, endAt time.Time, excludeBookingID uint) (bool, error)
	Create(booking *model.Booking) error
	Update(booking *model.Booking) error
	Delete(booking *model.Booking) error
}

type CreateBookingInput struct {
	UserID        uint
	ServiceID     uint
	TechnicianID  *uint
	StartAt       time.Time
	EndAt         *time.Time
	CustomerName  string
	CustomerPhone string
	Note          string
}

type UpdateBookingInput struct {
	UserID          *uint
	ServiceID       *uint
	TechnicianID    *uint
	TechnicianIDSet bool
	StartAt         *time.Time
	EndAt           *time.Time
	CustomerName    *string
	CustomerPhone   *string
	Note            *string
}

type BookingService struct {
	repo             BookingStore
	bookingNoFactory func() (string, error)
}

func NewBookingService(repo BookingStore) *BookingService {
	return &BookingService{repo: repo, bookingNoFactory: generateBookingNo}
}

func (s *BookingService) GetBookings(filter repository.BookingFilter, pagination utils.Pagination) ([]model.Booking, int64, error) {
	if filter.Status != "" && !model.IsValidBookingStatus(filter.Status) {
		return nil, 0, apperror.BadRequest("invalid booking status", apperror.ErrValidation)
	}
	if filter.DateFrom != nil && filter.DateTo != nil && filter.DateFrom.After(*filter.DateTo) {
		return nil, 0, apperror.BadRequest("dateFrom must be before dateTo", apperror.ErrValidation)
	}
	return s.repo.FindAll(filter, pagination)
}

func (s *BookingService) GetBookingByID(id uint) (model.Booking, error) {
	booking, err := s.repo.FindByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Booking{}, apperror.NotFound("booking not found", err)
	}
	return booking, err
}

func (s *BookingService) CreateBooking(input CreateBookingInput) (model.Booking, error) {
	if input.UserID == 0 || input.ServiceID == 0 || input.StartAt.IsZero() {
		return model.Booking{}, apperror.BadRequest("userId, serviceId and startAt are required", apperror.ErrValidation)
	}
	if strings.TrimSpace(input.CustomerName) == "" || strings.TrimSpace(input.CustomerPhone) == "" {
		return model.Booking{}, apperror.BadRequest("customerName and customerPhone are required", apperror.ErrValidation)
	}

	user, err := s.findUser(input.UserID)
	if err != nil {
		return model.Booking{}, err
	}
	serviceModel, err := s.findService(input.ServiceID)
	if err != nil {
		return model.Booking{}, err
	}
	if serviceModel.ServicePrice < 0 || serviceModel.Duration <= 0 {
		return model.Booking{}, apperror.Internal("service has invalid price or duration", apperror.ErrValidation)
	}

	var technician *model.NailTechnician
	if input.TechnicianID != nil {
		if *input.TechnicianID == 0 {
			return model.Booking{}, apperror.BadRequest("technicianId must be greater than 0", apperror.ErrValidation)
		}
		found, findErr := s.findTechnician(*input.TechnicianID)
		if findErr != nil {
			return model.Booking{}, findErr
		}
		technician = &found
	}

	endAt := input.StartAt.Add(time.Duration(serviceModel.Duration) * time.Minute)
	if input.EndAt != nil {
		endAt = *input.EndAt
	}
	if !input.StartAt.Before(endAt) {
		return model.Booking{}, apperror.BadRequest("startAt must be before endAt", apperror.ErrValidation)
	}
	if err := s.ensureNoOverlap(input.TechnicianID, input.StartAt, endAt, 0); err != nil {
		return model.Booking{}, err
	}

	bookingNo, err := s.bookingNoFactory()
	if err != nil {
		return model.Booking{}, apperror.Internal("could not generate booking number", err)
	}
	booking := model.Booking{
		BookingNo:       bookingNo,
		UserID:          input.UserID,
		ServiceID:       input.ServiceID,
		TechnicianID:    input.TechnicianID,
		StartAt:         input.StartAt,
		EndAt:           endAt,
		CustomerName:    strings.TrimSpace(input.CustomerName),
		CustomerPhone:   strings.TrimSpace(input.CustomerPhone),
		ServiceName:     serviceModel.ServiceName,
		Price:           serviceModel.ServicePrice,
		DurationMinutes: serviceModel.Duration,
		Status:          model.BookingStatusPending,
		Note:            input.Note,
		User:            user,
		Service:         serviceModel,
		Technician:      technician,
	}
	if err := s.repo.Create(&booking); err != nil {
		if errors.Is(err, repository.ErrTechnicianOverlap) {
			return model.Booking{}, apperror.Conflict("technician is already booked for this time", err)
		}
		return model.Booking{}, err
	}
	return booking, nil
}

func (s *BookingService) UpdateBooking(id uint, input UpdateBookingInput) (model.Booking, error) {
	booking, err := s.GetBookingByID(id)
	if err != nil {
		return model.Booking{}, err
	}

	recalculateEnd := false
	if input.UserID != nil {
		user, findErr := s.findUser(*input.UserID)
		if findErr != nil {
			return model.Booking{}, findErr
		}
		booking.UserID, booking.User = *input.UserID, user
	}
	if input.ServiceID != nil {
		serviceModel, findErr := s.findService(*input.ServiceID)
		if findErr != nil {
			return model.Booking{}, findErr
		}
		if serviceModel.ServicePrice < 0 || serviceModel.Duration <= 0 {
			return model.Booking{}, apperror.Internal("service has invalid price or duration", apperror.ErrValidation)
		}
		booking.ServiceID = *input.ServiceID
		booking.Service = serviceModel
		booking.ServiceName = serviceModel.ServiceName
		booking.Price = serviceModel.ServicePrice
		booking.DurationMinutes = serviceModel.Duration
		recalculateEnd = true
	}
	if input.TechnicianIDSet {
		booking.TechnicianID = input.TechnicianID
		booking.Technician = nil
		if input.TechnicianID != nil {
			technician, findErr := s.findTechnician(*input.TechnicianID)
			if findErr != nil {
				return model.Booking{}, findErr
			}
			booking.Technician = &technician
		}
	}
	if input.StartAt != nil {
		booking.StartAt = *input.StartAt
		recalculateEnd = true
	}
	if input.EndAt != nil {
		booking.EndAt = *input.EndAt
	} else if recalculateEnd {
		booking.EndAt = booking.StartAt.Add(time.Duration(booking.DurationMinutes) * time.Minute)
	}
	if input.CustomerName != nil {
		booking.CustomerName = strings.TrimSpace(*input.CustomerName)
	}
	if input.CustomerPhone != nil {
		booking.CustomerPhone = strings.TrimSpace(*input.CustomerPhone)
	}
	if input.Note != nil {
		booking.Note = *input.Note
	}
	if booking.CustomerName == "" || booking.CustomerPhone == "" || !booking.StartAt.Before(booking.EndAt) {
		return model.Booking{}, apperror.BadRequest("booking data or time range is invalid", apperror.ErrValidation)
	}
	if booking.Status != model.BookingStatusCancelled && booking.Status != model.BookingStatusNoShow {
		if err := s.ensureNoOverlap(booking.TechnicianID, booking.StartAt, booking.EndAt, booking.ID); err != nil {
			return model.Booking{}, err
		}
	}
	if err := s.repo.Update(&booking); err != nil {
		if errors.Is(err, repository.ErrTechnicianOverlap) {
			return model.Booking{}, apperror.Conflict("technician is already booked for this time", err)
		}
		return model.Booking{}, err
	}
	return booking, nil
}

func (s *BookingService) UpdateBookingStatus(id uint, status model.BookingStatus, cancelReason string) (model.Booking, error) {
	if !model.IsValidBookingStatus(status) {
		return model.Booking{}, apperror.BadRequest("invalid booking status", apperror.ErrValidation)
	}
	booking, err := s.GetBookingByID(id)
	if err != nil {
		return model.Booking{}, err
	}
	if status != model.BookingStatusCancelled && status != model.BookingStatusNoShow {
		if err := s.ensureNoOverlap(booking.TechnicianID, booking.StartAt, booking.EndAt, booking.ID); err != nil {
			return model.Booking{}, err
		}
	}
	booking.Status = status
	booking.CancelReason = cancelReason
	if err := s.repo.Update(&booking); err != nil {
		if errors.Is(err, repository.ErrTechnicianOverlap) {
			return model.Booking{}, apperror.Conflict("technician is already booked for this time", err)
		}
		return model.Booking{}, err
	}
	return booking, nil
}

func (s *BookingService) DeleteBooking(id uint) error {
	booking, err := s.GetBookingByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(&booking)
}

func (s *BookingService) findUser(id uint) (model.User, error) {
	user, err := s.repo.FindUserByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.User{}, apperror.NotFound("user not found", err)
	}
	return user, err
}

func (s *BookingService) findService(id uint) (model.Service, error) {
	serviceModel, err := s.repo.FindServiceByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Service{}, apperror.NotFound("service not found", err)
	}
	return serviceModel, err
}

func (s *BookingService) findTechnician(id uint) (model.NailTechnician, error) {
	technician, err := s.repo.FindTechnicianByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.NailTechnician{}, apperror.NotFound("technician not found", err)
	}
	return technician, err
}

func (s *BookingService) ensureNoOverlap(technicianID *uint, startAt, endAt time.Time, excludeBookingID uint) error {
	if technicianID == nil {
		return nil
	}
	overlaps, err := s.repo.HasTechnicianOverlap(*technicianID, startAt, endAt, excludeBookingID)
	if err != nil {
		return err
	}
	if overlaps {
		return apperror.Conflict("technician is already booked for this time", apperror.ErrValidation)
	}
	return nil
}

func generateBookingNo() (string, error) {
	random := make([]byte, 6)
	if _, err := rand.Read(random); err != nil {
		return "", err
	}
	return fmt.Sprintf("BK-%s-%s", time.Now().Format("20060102"), strings.ToUpper(hex.EncodeToString(random))), nil
}
