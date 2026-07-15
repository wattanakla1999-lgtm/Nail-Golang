package apperror

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrValidation     = errors.New("validation error")
	ErrDuplicateEmail = errors.New("email already exists")
)

const CodeBookingTimeOverlap = "BOOKING_TIME_OVERLAP"

type AppError struct {
	Status  int
	Code    string
	Message string
	Err     error
}

func NewWithCode(status int, code, message string, err error) *AppError {
	return &AppError{
		Status:  status,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(status int, message string, err error) *AppError {
	return &AppError{
		Status:  status,
		Message: message,
		Err:     err,
	}
}

func BadRequest(message string, err error) *AppError {
	return New(http.StatusBadRequest, message, err)
}

func NotFound(message string, err error) *AppError {
	return New(http.StatusNotFound, message, err)
}

func Conflict(message string, err error) *AppError {
	return New(http.StatusConflict, message, err)
}

func ConflictWithCode(code, message string, err error) *AppError {
	return NewWithCode(http.StatusConflict, code, message, err)
}

func Internal(message string, err error) *AppError {
	return New(http.StatusInternalServerError, message, err)
}

func FromError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFound("user not found", err)
	}

	return Internal("internal server error", err)
}
