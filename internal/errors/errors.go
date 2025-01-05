package errors

import (
	"fmt"
)

// AppError represents a custom application error
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("code: %d, message: %s, detail: %s", e.Code, e.Message, e.Detail)
}

// Common error codes
const (
	ErrBadRequest         = 400
	ErrUnauthorized       = 401
	ErrForbidden          = 403
	ErrNotFound           = 404
	ErrInternalServer     = 500
	ErrServiceUnavailable = 503
)

// New creates a new AppError
func New(code int, message string, detail string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

// NewBadRequest creates a new bad request error
func NewBadRequest(message string, detail string) *AppError {
	return New(ErrBadRequest, message, detail)
}

// NewUnauthorized creates a new unauthorized error
func NewUnauthorized(message string, detail string) *AppError {
	return New(ErrUnauthorized, message, detail)
}

// NewInternalServer creates a new internal server error
func NewInternalServer(message string, detail string) *AppError {
	return New(ErrInternalServer, message, detail)
}
