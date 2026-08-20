package apperror

import "errors"

var (
	ErrNotFound = errors.New("resource not found")
	ErrInvalid  = errors.New("invalid request")
)