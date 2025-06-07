package domain

import (
	"errors"
	"fmt"
)

// ErrType represents the type of error that occurred
type ErrType string

const (
	ErrTypeRepository   ErrType = "repository"
	ErrTypeVersion      ErrType = "version"
	ErrTypeService      ErrType = "service"
	ErrTypeValidation   ErrType = "validation"
	ErrTypeNotFound     ErrType = "not_found"
	ErrTypeInvalidInput ErrType = "invalid_input"
)

// ApplicationError represents a structured error with context
type ApplicationError struct {
	Type      ErrType
	Operation string
	Message   string
	Err       error
	Context   map[string]any
}

func (e *ApplicationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s %s: %s: %v", e.Type, e.Operation, e.Message, e.Err)
	}
	return fmt.Sprintf("%s %s: %s", e.Type, e.Operation, e.Message)
}

func (e *ApplicationError) Unwrap() error {
	return e.Err
}

// NewError creates a new ApplicationError
func NewError(errType ErrType, operation, message string, err error) *ApplicationError {
	return &ApplicationError{
		Type:      errType,
		Operation: operation,
		Message:   message,
		Err:       err,
		Context:   make(map[string]any),
	}
}

// WithContext adds context to the error
func (e *ApplicationError) WithContext(key string, value any) *ApplicationError {
	if e.Context == nil {
		e.Context = make(map[string]any)
	}
	e.Context[key] = value
	return e
}

// Common error constructors
func NewRepositoryError(operation, message string, err error) *ApplicationError {
	return NewError(ErrTypeRepository, operation, message, err)
}

func NewVersionError(operation, message string, err error) *ApplicationError {
	return NewError(ErrTypeVersion, operation, message, err)
}

func NewServiceError(operation, message string, err error) *ApplicationError {
	return NewError(ErrTypeService, operation, message, err)
}

func NewValidationError(operation, message string, err error) *ApplicationError {
	return NewError(ErrTypeValidation, operation, message, err)
}

func NewNotFoundError(operation, message string) *ApplicationError {
	return NewError(ErrTypeNotFound, operation, message, nil)
}

func NewInvalidInputError(operation, message string, err error) *ApplicationError {
	return NewError(ErrTypeInvalidInput, operation, message, err)
}

// Error checking utilities
func IsErrorType(err error, errType ErrType) bool {
	var appErr *ApplicationError
	if errors.As(err, &appErr) {
		return appErr.Type == errType
	}
	return false
}

func IsRepositoryError(err error) bool {
	return IsErrorType(err, ErrTypeRepository)
}

func IsVersionError(err error) bool {
	return IsErrorType(err, ErrTypeVersion)
}

func IsServiceError(err error) bool {
	return IsErrorType(err, ErrTypeService)
}

func IsValidationError(err error) bool {
	return IsErrorType(err, ErrTypeValidation)
}

func IsNotFoundError(err error) bool {
	return IsErrorType(err, ErrTypeNotFound)
}

func IsInvalidInputError(err error) bool {
	return IsErrorType(err, ErrTypeInvalidInput)
}