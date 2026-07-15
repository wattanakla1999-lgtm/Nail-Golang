package service

import (
	"nailly-back-end/internal/model"
	"nailly-back-end/internal/repository"
	"testing"
	"time"

	"gorm.io/gorm"
)

type fakeDashboardStore struct {
	today time.Time
}

func (f *fakeDashboardStore) CountAppointments(startAt, endAt time.Time) (int64, error) {
	if startAt.Year() == f.today.Year() && startAt.YearDay() == f.today.YearDay() {
		return 12, nil
	}
	return 9, nil
}

func (f *fakeDashboardStore) CountCustomers() (int64, error) {
	return 248, nil
}

func (f *fakeDashboardStore) CountNewCustomers(startAt, endAt time.Time) (int64, error) {
	return 5, nil
}

func (f *fakeDashboardStore) CountServices() (int64, error) {
	return 8, nil
}

func (f *fakeDashboardStore) SumCompletedRevenue(startAt, endAt time.Time) (int64, error) {
	return 3200, nil
}

func (f *fakeDashboardStore) FindAppointments(startAt, endAt time.Time, limit int) ([]model.Booking, error) {
	return []model.Booking{
		{
			Model: gorm.Model{ID: 1}, CustomerName: "คุณพิม พิมประภา", ServiceName: "ต่อเล็บเจล",
			StartAt: time.Date(2026, 7, 15, 10, 0, 0, 0, f.today.Location()), Status: model.BookingStatusInService,
		},
		{
			Model: gorm.Model{ID: 2}, CustomerName: "คุณอรัญญา", ServiceName: "ทาสีเจล",
			StartAt: time.Date(2026, 7, 15, 13, 0, 0, 0, f.today.Location()), Status: model.BookingStatusCompleted,
		},
	}, nil
}

func (f *fakeDashboardStore) FindPopularServices(limit int) ([]repository.PopularServiceRow, int64, error) {
	return []repository.PopularServiceRow{
		{Name: "ทาสีเจลมือ", Count: 45},
		{Name: "สปามือ-เท้า", Count: 30},
	}, 100, nil
}

func TestDashboardServiceGetStats(t *testing.T) {
	location := time.FixedZone("Asia/Bangkok", 7*60*60)
	now := time.Date(2026, 7, 15, 18, 0, 0, 0, location)
	dashboardService := NewDashboardService(&fakeDashboardStore{today: now})
	dashboardService.now = func() time.Time { return now }

	stats, err := dashboardService.GetStats()
	if err != nil {
		t.Fatalf("GetStats() error = %v", err)
	}
	if stats.TodayAppointments != 12 || stats.TodayAppointmentsChange != 3 {
		t.Fatalf("appointment stats = %d change %d", stats.TodayAppointments, stats.TodayAppointmentsChange)
	}
	if stats.TotalCustomers != 248 || stats.TotalCustomersChange != "+5 ใหม่" || stats.ActiveServicesCount != 8 || stats.TodayRevenue != 3200 {
		t.Fatalf("summary = %+v", stats)
	}
	if len(stats.Appointments) != 2 || stats.Appointments[0].Status != "active" || stats.Appointments[1].Status != "done" {
		t.Fatalf("appointments = %+v", stats.Appointments)
	}
	if len(stats.PopularServices) != 2 || stats.PopularServices[0].Percent != 45 || stats.PopularServices[0].Rate != 0 {
		t.Fatalf("popularServices = %+v", stats.PopularServices)
	}
}
