package errs

import (
	"errors"
	"fmt"
	"net/http"
)

type ServiceError struct {
	Underlying error
	Code       string
	Message    string
	Status     int
	Details    any
}

func (e *ServiceError) Error() string {
	if e.Underlying != nil {
		return fmt.Sprintf("%s: %v", e.Code, e.Underlying)
	}
	return e.Code
}

func New(code string, status int, message string, opts ...Option) *ServiceError {
	se := &ServiceError{
		Code:    code,
		Message: message,
		Status:  status,
	}

	for _, opt := range opts {
		opt(se)
	}

	return se
}

type Option func(*ServiceError)

func WithDetails(d interface{}) Option {
	return func(e *ServiceError) {
		e.Details = d
	}
}

func WithError(err error) Option {
	return func(e *ServiceError) {
		e.Underlying = err
	}
}

func AsServiceError(err error) *ServiceError {
	var se *ServiceError
	if errors.As(err, &se) {
		return se
	}

	return &ServiceError{
		Underlying: err,
		Code:       "internal_error",
		Message:    "Internal server error",
		Status:     http.StatusInternalServerError,
		Details:    err.Error(),
	}
}
