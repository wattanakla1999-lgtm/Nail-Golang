package dto

import (
	"nailly-back-end/internal/service"
	"strconv"
)

type ReportSummaryResponse struct {
	TodayRevenue      int     `json:"todayRevenue"`
	TodayAppointments int     `json:"todayAppointments"`
	AvgPerBill        int     `json:"avgPerBill"`
	TargetPercent     int     `json:"targetPercent"`
	RevenueChange     float64 `json:"revenueChange"`
}

type ReportChartResponse struct {
	Day          string `json:"day"`
	Revenue      int    `json:"revenue"`
	Appointments int    `json:"appointments"`
}

type ReportPaymentResponse struct {
	Method  string `json:"method"`
	Percent int    `json:"percent"`
	Amount  int    `json:"amount"`
}

type ReportTransactionResponse struct {
	ID       string `json:"id"`
	Customer string `json:"customer"`
	Service  string `json:"service"`
	Amount   int    `json:"amount"`
	Time     string `json:"time"`
	Staff    string `json:"staff"`
	Method   string `json:"method"`
}

type ReportResponse struct {
	Summary            ReportSummaryResponse       `json:"summary"`
	ChartData          []ReportChartResponse       `json:"chartData"`
	PaymentBreakdown   []ReportPaymentResponse     `json:"paymentBreakdown"`
	RecentTransactions []ReportTransactionResponse `json:"recentTransactions"`
}

func ToReportResponse(report service.Report) ReportResponse {
	response := ReportResponse{
		Summary: ReportSummaryResponse{
			TodayRevenue: report.Summary.TodayRevenue, TodayAppointments: report.Summary.TodayAppointments,
			AvgPerBill: report.Summary.AvgPerBill, TargetPercent: report.Summary.TargetPercent,
			RevenueChange: report.Summary.RevenueChange,
		},
		ChartData:          make([]ReportChartResponse, 0, len(report.ChartData)),
		PaymentBreakdown:   make([]ReportPaymentResponse, 0, len(report.PaymentBreakdown)),
		RecentTransactions: make([]ReportTransactionResponse, 0, len(report.RecentTransactions)),
	}
	for _, item := range report.ChartData {
		response.ChartData = append(response.ChartData, ReportChartResponse(item))
	}
	for _, item := range report.PaymentBreakdown {
		response.PaymentBreakdown = append(response.PaymentBreakdown, ReportPaymentResponse{
			Method: string(item.Method), Percent: item.Percent, Amount: item.Amount,
		})
	}
	for _, item := range report.RecentTransactions {
		response.RecentTransactions = append(response.RecentTransactions, ReportTransactionResponse{
			ID: strconv.FormatUint(uint64(item.ID), 10), Customer: item.Customer, Service: item.Service,
			Amount: item.Amount, Time: item.Time, Staff: item.Staff, Method: string(item.Method),
		})
	}
	return response
}
