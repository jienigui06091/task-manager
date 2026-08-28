package apperror

import "net/http"

type AppError struct {
	Code       int
	HTTPStatus int
	Message    string
	Err        error
}

func New(code, httpStatus int, message string) *AppError {
	return &AppError{
		Code:       code,
		HTTPStatus: httpStatus,
		Message:    message,
	}
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

var (
	ErrNotFound = &AppError{
		Code:       40401,
		HTTPStatus: http.StatusNotFound,
		Message:    "resource not found",
	}

	ErrUnauthorized = &AppError{
		Code:       40101,
		HTTPStatus: http.StatusUnauthorized,
		Message:    "unauthorized",
	}

	ErrForbidden = &AppError{
		Code:       40301,
		HTTPStatus: http.StatusForbidden,
		Message:    "forbidden",
	}

	ErrBadRequest = &AppError{
		Code:       40001,
		HTTPStatus: http.StatusBadRequest,
		Message:    "bad request",
	}

	ErrInvalid = ErrBadRequest
)
