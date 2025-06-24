package webserver

import "fmt"

type ValidationError struct {
	err error
}

func (e *ValidationError) Error() string {
	return fmt.Errorf("validation error: %w", e.err).Error()
}

func NewValidationError(err error) *ValidationError {
	return &ValidationError{err}
}

type InvalidRequestError struct {
	err error
}

func (e *InvalidRequestError) Error() string {
	return fmt.Errorf("invalid request: %w", e.err).Error()
}

func NewInvalidRequestError(err error) *InvalidRequestError {
	return &InvalidRequestError{err}
}

type NotFoundError struct {
	err error
}

func (e *NotFoundError) Error() string {
	return fmt.Errorf("not found: %w", e.err).Error()
}

func NewNotFoundError(err error) *NotFoundError {
	return &NotFoundError{err}
}
