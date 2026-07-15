package service

import (
	"nailly-back-end/internal/model"
	"testing"
	"time"

	"gorm.io/gorm"
)

type fakeReportStore struct {
	bookings []model.Booking
	startAt  time.Time
	endAt    time.Time
}

func (f *fakeReportStore) FindBookings(startAt, endAt time.Time) ([]model.Booking, error) {
	f.startAt, f.endAt = startAt, endAt
	return f.bookings, nil
}

func TestReportServiceWeek(t *testing.T) {
	location := time.FixedZone("Asia/Bangkok", 7*60*60)
	now := time.Date(2026, 7, 15, 18, 0, 0, 0, location)
	store := &fakeReportStore{bookings: []model.Booking{
		reportBooking(1, now.AddDate(0, 0, -2), 200, model.BookingStatusCompleted, model.PaymentMethodCash),
		reportBooking(2, now.AddDate(0, 0, -1), 1000, model.BookingStatusCompleted, model.PaymentMethodTransfer),
		reportBooking(3, time.Date(2026, 7, 15, 10, 45, 0, 0, location), 650, model.BookingStatusCompleted, model.PaymentMethodTransfer),
		reportBooking(4, time.Date(2026, 7, 15, 14, 0, 0, 0, location), 450, model.BookingStatusCompleted, model.PaymentMethodCash),
		reportBooking(5, time.Date(2026, 7, 15, 16, 0, 0, 0, location), 300, model.BookingStatusPending, model.PaymentMethodCash),
		reportBooking(6, time.Date(2026, 7, 15, 17, 0, 0, 0, location), 500, model.BookingStatusCancelled, model.PaymentMethodCard),
	}}
	reportService := NewReportService(store)
	reportService.now = func() time.Time { return now }

	report, err := reportService.GetReport("week")
	if err != nil {
		t.Fatalf("GetReport() error = %v", err)
	}
	if report.Summary.TodayRevenue != 1100 || report.Summary.TodayAppointments != 3 {
		t.Fatalf("summary = %+v, want revenue 1100 and appointments 3", report.Summary)
	}
	if report.Summary.AvgPerBill != 550 || report.Summary.TargetPercent != 22 || report.Summary.RevenueChange != 10 {
		t.Fatalf("calculated summary = %+v", report.Summary)
	}
	if len(report.ChartData) != 7 || report.ChartData[0].Day != "จ." || report.ChartData[2].Revenue != 1100 {
		t.Fatalf("chartData = %+v", report.ChartData)
	}
	if len(report.PaymentBreakdown) != 2 || report.PaymentBreakdown[0].Method != model.PaymentMethodTransfer || report.PaymentBreakdown[0].Amount != 650 {
		t.Fatalf("paymentBreakdown = %+v", report.PaymentBreakdown)
	}
	if len(report.RecentTransactions) != 4 || report.RecentTransactions[0].ID != 4 {
		t.Fatalf("recentTransactions = %+v", report.RecentTransactions)
	}
}

func TestReportServiceRejectsInvalidPeriod(t *testing.T) {
	_, err := NewReportService(&fakeReportStore{}).GetReport("year")
	if err == nil {
		t.Fatal("GetReport() error = nil, want invalid period error")
	}
}

func TestReportServiceMonthPractice(t *testing.T) {
	t.Skip("TODO: ฝึกเพิ่ม test เดือนกุมภาพันธ์ปีอธิกสุรทิน และตรวจว่ากราฟมี 29 วัน")
}

func reportBooking(id uint, startAt time.Time, price int, status model.BookingStatus, method model.PaymentMethod) model.Booking {
	return model.Booking{
		Model: gorm.Model{ID: id}, StartAt: startAt, Price: price, Status: status,
		CustomerName: "Customer", ServiceName: "Gel Manicure", PaymentMethod: method,
		Technician: &model.NailTechnician{TechnicianName: "Nุ่น"},
	}
}
