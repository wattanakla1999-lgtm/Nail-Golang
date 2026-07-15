package service

import (
	"math"
	"nailly-back-end/internal/apperror"
	"nailly-back-end/internal/model"
	"sort"
	"strconv"
	"time"
)

const dailyRevenueTarget = 5000

type ReportStore interface {
	FindBookings(startAt, endAt time.Time) ([]model.Booking, error)
}

type ReportSummary struct {
	TodayRevenue      int
	TodayAppointments int
	AvgPerBill        int
	TargetPercent     int
	RevenueChange     float64
}

type ReportChartItem struct {
	Day          string
	Revenue      int
	Appointments int
}

type ReportPaymentItem struct {
	Method  model.PaymentMethod
	Percent int
	Amount  int
}

type ReportTransaction struct {
	ID       uint
	Customer string
	Service  string
	Amount   int
	Time     string
	Staff    string
	Method   model.PaymentMethod
}

type Report struct {
	Summary            ReportSummary
	ChartData          []ReportChartItem
	PaymentBreakdown   []ReportPaymentItem
	RecentTransactions []ReportTransaction
}

type ReportService struct {
	repo ReportStore
	now  func() time.Time
}

func NewReportService(repo ReportStore) *ReportService {
	return &ReportService{repo: repo, now: time.Now}
}

func (s *ReportService) GetReport(period string) (Report, error) {
	if period == "" {
		period = "week"
	}
	if period != "week" && period != "month" {
		return Report{}, apperror.BadRequest("period must be week or month", apperror.ErrValidation)
	}

	now := s.now().In(time.Local)
	todayStart := startOfDay(now)
	tomorrowStart := todayStart.AddDate(0, 0, 1)
	yesterdayStart := todayStart.AddDate(0, 0, -1)
	chartStart, chartEnd := reportRange(now, period)
	queryStart := chartStart
	if yesterdayStart.Before(queryStart) {
		queryStart = yesterdayStart
	}

	bookings, err := s.repo.FindBookings(queryStart, chartEnd)
	if err != nil {
		return Report{}, err
	}

	today := filterBookings(bookings, todayStart, tomorrowStart)
	yesterday := filterBookings(bookings, yesterdayStart, todayStart)
	todayRevenue, todayCompleted := completedRevenue(today)
	yesterdayRevenue, _ := completedRevenue(yesterday)

	return Report{
		Summary: ReportSummary{
			TodayRevenue:      todayRevenue,
			TodayAppointments: countAppointments(today),
			AvgPerBill:        average(todayRevenue, todayCompleted),
			TargetPercent:     percentage(todayRevenue, dailyRevenueTarget),
			RevenueChange:     percentChange(todayRevenue, yesterdayRevenue),
		},
		ChartData:          buildChart(bookings, chartStart, chartEnd, period),
		PaymentBreakdown:   buildPaymentBreakdown(today, todayRevenue),
		RecentTransactions: buildRecentTransactions(bookings, chartStart, chartEnd, 5),
	}, nil
}

func reportRange(now time.Time, period string) (time.Time, time.Time) {
	day := startOfDay(now)
	if period == "month" {
		start := time.Date(day.Year(), day.Month(), 1, 0, 0, 0, 0, day.Location())
		return start, start.AddDate(0, 1, 0)
	}

	daysSinceMonday := (int(day.Weekday()) + 6) % 7
	start := day.AddDate(0, 0, -daysSinceMonday)
	return start, start.AddDate(0, 0, 7)
}

func startOfDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}

func filterBookings(bookings []model.Booking, startAt, endAt time.Time) []model.Booking {
	filtered := make([]model.Booking, 0)
	for _, booking := range bookings {
		if !booking.StartAt.Before(startAt) && booking.StartAt.Before(endAt) {
			filtered = append(filtered, booking)
		}
	}
	return filtered
}

func completedRevenue(bookings []model.Booking) (revenue, count int) {
	for _, booking := range bookings {
		if booking.Status == model.BookingStatusCompleted {
			revenue += booking.Price
			count++
		}
	}
	return revenue, count
}

func countAppointments(bookings []model.Booking) int {
	count := 0
	for _, booking := range bookings {
		if booking.Status != model.BookingStatusCancelled && booking.Status != model.BookingStatusNoShow {
			count++
		}
	}
	return count
}

func average(total, count int) int {
	if count == 0 {
		return 0
	}
	return total / count
}

func percentage(value, total int) int {
	if total == 0 {
		return 0
	}
	return value * 100 / total
}

func percentChange(current, previous int) float64 {
	if previous == 0 {
		if current > 0 {
			return 100
		}
		return 0
	}
	change := float64(current-previous) / float64(previous) * 100
	return math.Round(change*10) / 10
}

func buildChart(bookings []model.Booking, startAt, endAt time.Time, period string) []ReportChartItem {
	items := make([]ReportChartItem, 0)
	thaiDays := []string{"อา.", "จ.", "อ.", "พ.", "พฤ.", "ศ.", "ส."}
	for day := startAt; day.Before(endAt); day = day.AddDate(0, 0, 1) {
		nextDay := day.AddDate(0, 0, 1)
		dayBookings := filterBookings(bookings, day, nextDay)
		revenue, _ := completedRevenue(dayBookings)
		label := thaiDays[day.Weekday()]
		if period == "month" {
			label = strconv.Itoa(day.Day())
		}
		items = append(items, ReportChartItem{Day: label, Revenue: revenue, Appointments: countAppointments(dayBookings)})
	}
	return items
}

func buildPaymentBreakdown(bookings []model.Booking, totalRevenue int) []ReportPaymentItem {
	amounts := map[model.PaymentMethod]int{}
	for _, booking := range bookings {
		if booking.Status == model.BookingStatusCompleted {
			method := booking.PaymentMethod
			if !model.IsValidPaymentMethod(method) {
				method = model.PaymentMethodCash
			}
			amounts[method] += booking.Price
		}
	}

	methods := []model.PaymentMethod{model.PaymentMethodTransfer, model.PaymentMethodCash, model.PaymentMethodCard}
	items := make([]ReportPaymentItem, 0, len(amounts))
	for _, method := range methods {
		if amount := amounts[method]; amount > 0 {
			items = append(items, ReportPaymentItem{Method: method, Percent: percentage(amount, totalRevenue), Amount: amount})
		}
	}
	return items
}

func buildRecentTransactions(bookings []model.Booking, startAt, endAt time.Time, limit int) []ReportTransaction {
	completed := filterBookings(bookings, startAt, endAt)
	sort.Slice(completed, func(i, j int) bool { return completed[i].StartAt.After(completed[j].StartAt) })
	items := make([]ReportTransaction, 0, limit)
	for _, booking := range completed {
		if booking.Status != model.BookingStatusCompleted {
			continue
		}
		staff := ""
		if booking.Technician != nil {
			staff = booking.Technician.TechnicianName
		}
		method := booking.PaymentMethod
		if !model.IsValidPaymentMethod(method) {
			method = model.PaymentMethodCash
		}
		items = append(items, ReportTransaction{
			ID: booking.ID, Customer: booking.CustomerName, Service: booking.ServiceName,
			Amount: booking.Price, Time: booking.StartAt.In(time.Local).Format("15:04"), Staff: staff, Method: method,
		})
		if len(items) == limit {
			break
		}
	}
	return items
}
