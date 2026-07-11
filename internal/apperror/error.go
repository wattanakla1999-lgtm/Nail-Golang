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

type AppError struct {
	Status  int
	Message string
	Err     error
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
