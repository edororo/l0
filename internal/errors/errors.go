package errors

import "errors"

var (
	ErrInvalidInput = errors.New("Invalid input")
	ErrNotFound     = errors.New("Not found")
	ErrInternal     = errors.New("Internal server error")
)
