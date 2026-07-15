package router

import (
	"nailly-back-end/internal/handler"
	"nailly-back-end/internal/repository"
	"nailly-back-end/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterBookingRoutes(api *gin.RouterGroup, db *gorm.DB) {
	bookingRepository := repository.NewBookingRepository(db)
	bookingService := service.NewBookingService(bookingRepository)
	bookingHandler := handler.NewBookingHandler(bookingService)

	bookings := api.Group("/bookings")
	bookings.GET("", bookingHandler.GetBookings)
	bookings.GET("/busy-slots", bookingHandler.GetBusySlots)
	bookings.GET("/:id", bookingHandler.GetBookingByID)
	bookings.POST("", bookingHandler.CreateBooking)
	bookings.PUT("/:id", bookingHandler.UpdateBooking)
	bookings.DELETE("/:id", bookingHandler.DeleteBooking)
	bookings.PATCH("/:id/status", bookingHandler.UpdateBookingStatus)
}
