package service

import (
	"errors"
	"nailly-back-end/internal/apperror"
	"nailly-back-end/internal/model"
	"nailly-back-end/internal/repository"
	"nailly-back-end/pkg/utils"
	"net/http"
	"sort"
	"testing"
	"time"

	"gorm.io/gorm"
)

type fakeBookingStore struct {
	users         map[uint]model.User
	services      map[uint]model.Service
	technicians   map[uint]model.NailTechnician
	bookings      map[uint]model.Booking
	nextBookingID uint
}

func newFakeBookingStore() *fakeBookingStore {
	return &fakeBookingStore{
		users: map[uint]model.User{
			1: {Model: gorm.Model{ID: 1}, Name: "Somying"},
			2: {Model: gorm.Model{ID: 2}, Name: "Mali"},
		},
		services: map[uint]model.Service{
			2: {Model: gorm.Model{ID: 2}, ServiceName: "Gel nails", ServicePrice: 300, Duration: 60},
			3: {Model: gorm.Model{ID: 3}, ServiceName: "Nail art", ServicePrice: 450, Duration: 90},
		},
		technicians: map[uint]model.NailTechnician{
			1: {Model: gorm.Model{ID: 1}, TechnicianName: "Nok"},
			2: {Model: gorm.Model{ID: 2}, TechnicianName: "Ploy"},
		},
		bookings:      make(map[uint]model.Booking),
		nextBookingID: 1,
	}
}

func (f *fakeBookingStore) FindAll(filter repository.BookingFilter, pagination utils.Pagination) ([]model.Booking, int64, error) {
	bookings := make([]model.Booking, 0)
	for _, booking := range f.bookings {
		if booking.DeletedAt.Valid || filter.UserID != nil && booking.UserID != *filter.UserID ||
			filter.ServiceID != nil && booking.ServiceID != *filter.ServiceID ||
			filter.TechnicianID != nil && (booking.TechnicianID == nil || *booking.TechnicianID != *filter.TechnicianID) ||
			filter.Status != "" && booking.Status != filter.Status ||
			filter.DateFrom != nil && booking.StartAt.Before(*filter.DateFrom) ||
			filter.DateTo != nil && booking.StartAt.After(*filter.DateTo) {
			continue
		}
		bookings = append(bookings, f.hydrate(booking))
	}
	sort.Slice(bookings, func(i, j int) bool { return bookings[i].StartAt.Before(bookings[j].StartAt) })
	total := int64(len(bookings))
	start := pagination.Offset
	if start > len(bookings) {
		start = len(bookings)
	}
	end := start + pagination.Limit
	if end > len(bookings) {
		end = len(bookings)
	}
	return bookings[start:end], total, nil
}

func (f *fakeBookingStore) FindByID(id uint) (model.Booking, error) {
	booking, ok := f.bookings[id]
	if !ok || booking.DeletedAt.Valid {
		return model.Booking{}, gorm.ErrRecordNotFound
	}
	return f.hydrate(booking), nil
}

func (f *fakeBookingStore) FindUserByID(id uint) (model.User, error) {
	user, ok := f.users[id]
	if !ok {
		return model.User{}, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (f *fakeBookingStore) FindServiceByID(id uint) (model.Service, error) {
	serviceModel, ok := f.services[id]
	if !ok {
		return model.Service{}, gorm.ErrRecordNotFound
	}
	return serviceModel, nil
}

func (f *fakeBookingStore) FindTechnicianByID(id uint) (model.NailTechnician, error) {
	technician, ok := f.technicians[id]
	if !ok {
		return model.NailTechnician{}, gorm.ErrRecordNotFound
	}
	return technician, nil
}

func (f *fakeBookingStore) HasTechnicianOverlap(technicianID uint, startAt, endAt time.Time, excludeBookingID uint) (bool, error) {
	for _, booking := range f.bookings {
		if booking.ID == excludeBookingID || booking.DeletedAt.Valid || booking.TechnicianID == nil ||
			*booking.TechnicianID != technicianID || booking.Status == model.BookingStatusCancelled ||
			booking.Status == model.BookingStatusNoShow {
			continue
		}
		if booking.StartAt.Before(endAt) && booking.EndAt.After(startAt) {
			return true, nil
		}
	}
	return false, nil
}

func (f *fakeBookingStore) Create(booking *model.Booking) error {
	for _, current := range f.bookings {
		if current.BookingNo == booking.BookingNo {
			return errors.New("duplicate booking number")
		}
	}
	booking.ID = f.nextBookingID
	f.nextBookingID++
	f.bookings[booking.ID] = *booking
	return nil
}

func (f *fakeBookingStore) Update(booking *model.Booking) error {
	if _, ok := f.bookings[booking.ID]; !ok {
		return gorm.ErrRecordNotFound
	}
	f.bookings[booking.ID] = *booking
	return nil
}

func (f *fakeBookingStore) Delete(booking *model.Booking) error {
	booking.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	f.bookings[booking.ID] = *booking
	return nil
}

func (f *fakeBookingStore) hydrate(booking model.Booking) model.Booking {
	booking.User = f.users[booking.UserID]
	booking.Service = f.services[booking.ServiceID]
	booking.Technician = nil
	if booking.TechnicianID != nil {
		technician := f.technicians[*booking.TechnicianID]
		booking.Technician = &technician
	}
	return booking
}

func bookingServiceForTest(store *fakeBookingStore) *BookingService {
	bookingService := NewBookingService(store)
	sequence := 0
	bookingService.bookingNoFactory = func() (string, error) {
		sequence++
		return "BK-TEST-" + time.Unix(int64(sequence), 0).Format("150405"), nil
	}
	return bookingService
}

func validCreateBookingInput() CreateBookingInput {
	technicianID := uint(1)
	return CreateBookingInput{
		UserID: 1, ServiceID: 2, TechnicianID: &technicianID,
		StartAt:      time.Date(2026, 7, 15, 10, 0, 0, 0, time.FixedZone("Asia/Bangkok", 7*60*60)),
		CustomerName: "Somying", CustomerPhone: "0812345678", Note: "flower pattern",
	}
}

func TestCreateBookingSuccess(t *testing.T) {
	store := newFakeBookingStore()
	booking, err := bookingServiceForTest(store).CreateBooking(validCreateBookingInput())
	if err != nil {
		t.Fatalf("CreateBooking() error = %v", err)
	}
	if booking.ID == 0 || booking.BookingNo == "" || booking.Status != model.BookingStatusPending {
		t.Fatalf("booking identifiers/status were not generated: %+v", booking)
	}
	if booking.ServiceName != "Gel nails" || booking.Price != 300 || booking.DurationMinutes != 60 {
		t.Fatalf("service snapshot is incorrect: %+v", booking)
	}
	if want := booking.StartAt.Add(60 * time.Minute); !booking.EndAt.Equal(want) {
		t.Fatalf("EndAt = %v, want %v", booking.EndAt, want)
	}
}

func TestCreateBookingForeignKeyNotFound(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(*CreateBookingInput)
		message string
	}{
		{name: "user", prepare: func(input *CreateBookingInput) { input.UserID = 99 }, message: "user not found"},
		{name: "service", prepare: func(input *CreateBookingInput) { input.ServiceID = 99 }, message: "service not found"},
		{name: "technician", prepare: func(input *CreateBookingInput) { id := uint(99); input.TechnicianID = &id }, message: "technician not found"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := validCreateBookingInput()
			test.prepare(&input)
			_, err := bookingServiceForTest(newFakeBookingStore()).CreateBooking(input)
			assertAppError(t, err, http.StatusNotFound, test.message)
		})
	}
}

func TestCreateBookingRejectsInvalidTimeRange(t *testing.T) {
	input := validCreateBookingInput()
	endAt := input.StartAt.Add(-time.Minute)
	input.EndAt = &endAt
	_, err := bookingServiceForTest(newFakeBookingStore()).CreateBooking(input)
	assertAppError(t, err, http.StatusBadRequest, "startAt must be before endAt")
}

func TestCreateBookingRejectsTechnicianOverlap(t *testing.T) {
	store := newFakeBookingStore()
	input := validCreateBookingInput()
	store.bookings[1] = existingBooking(1, *input.TechnicianID, input.StartAt, model.BookingStatusConfirmed)
	store.nextBookingID = 2
	_, err := bookingServiceForTest(store).CreateBooking(input)
	assertAppError(t, err, http.StatusConflict, "technician is already booked for this time")
}

func TestCancelledAndNoShowDoNotBlockTimeSlot(t *testing.T) {
	for _, status := range []model.BookingStatus{model.BookingStatusCancelled, model.BookingStatusNoShow} {
		t.Run(string(status), func(t *testing.T) {
			store := newFakeBookingStore()
			input := validCreateBookingInput()
			store.bookings[1] = existingBooking(1, *input.TechnicianID, input.StartAt, status)
			store.nextBookingID = 2
			if _, err := bookingServiceForTest(store).CreateBooking(input); err != nil {
				t.Fatalf("CreateBooking() error = %v", err)
			}
		})
	}
}

func TestUpdateBookingStatusSuccess(t *testing.T) {
	store := newFakeBookingStore()
	input := validCreateBookingInput()
	store.bookings[1] = existingBooking(1, *input.TechnicianID, input.StartAt, model.BookingStatusPending)
	booking, err := bookingServiceForTest(store).UpdateBookingStatus(1, model.BookingStatusConfirmed, "")
	if err != nil {
		t.Fatalf("UpdateBookingStatus() error = %v", err)
	}
	if booking.Status != model.BookingStatusConfirmed {
		t.Fatalf("Status = %q, want confirmed", booking.Status)
	}
}

func TestGetBookingsPaginationAndFilters(t *testing.T) {
	store := newFakeBookingStore()
	base := time.Date(2026, 7, 15, 9, 0, 0, 0, time.UTC)
	for i := uint(1); i <= 8; i++ {
		booking := existingBooking(i, 1, base.Add(time.Duration(i)*time.Hour), model.BookingStatusConfirmed)
		if i > 6 {
			booking.UserID = 2
		}
		if i == 6 {
			booking.Status = model.BookingStatusCancelled
		}
		store.bookings[i] = booking
	}
	userID := uint(1)
	from := base.Add(2 * time.Hour)
	to := base.Add(6 * time.Hour)
	filter := repository.BookingFilter{
		UserID: &userID, Status: model.BookingStatusConfirmed, DateFrom: &from, DateTo: &to,
	}
	bookings, total, err := bookingServiceForTest(store).GetBookings(filter, utils.Pagination{Page: 2, Limit: 2, Offset: 2})
	if err != nil {
		t.Fatalf("GetBookings() error = %v", err)
	}
	if total != 4 || len(bookings) != 2 || bookings[0].ID != 4 || bookings[1].ID != 5 {
		t.Fatalf("pagination/filter result = ids %v, total %d; want [4 5], total 4", bookingIDs(bookings), total)
	}
}

func existingBooking(id, technicianID uint, startAt time.Time, status model.BookingStatus) model.Booking {
	return model.Booking{
		Model: gorm.Model{ID: id}, BookingNo: "BK-EXISTING", UserID: 1, ServiceID: 2,
		TechnicianID: &technicianID, StartAt: startAt, EndAt: startAt.Add(time.Hour),
		CustomerName: "Customer", CustomerPhone: "0800000000", ServiceName: "Gel nails",
		Price: 300, DurationMinutes: 60, Status: status,
	}
}

func assertAppError(t *testing.T, err error, status int, message string) {
	t.Helper()
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %v, want AppError", err)
	}
	if appErr.Status != status || appErr.Message != message {
		t.Fatalf("AppError = (%d, %q), want (%d, %q)", appErr.Status, appErr.Message, status, message)
	}
}

func bookingIDs(bookings []model.Booking) []uint {
	ids := make([]uint, len(bookings))
	for i, booking := range bookings {
		ids[i] = booking.ID
	}
	return ids
}
