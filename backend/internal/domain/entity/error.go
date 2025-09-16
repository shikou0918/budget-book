package entity

import "fmt"

// ValidationError represents a validation error with a custom message
type ValidationError struct {
	Message string
}

// Error returns the formatted validation error message
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", e.Message)
}

// NewValidationError creates a new ValidationError instance with the given message
func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}

// NotFoundError represents an error when a resource is not found
type NotFoundError struct {
	Resource string
	ID       interface{}
}

// Error returns the formatted not found error message
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %v not found", e.Resource, e.ID)
}

// NewNotFoundError creates a new NotFoundError instance with the given resource and ID
func NewNotFoundError(resource string, id interface{}) *NotFoundError {
	return &NotFoundError{Resource: resource, ID: id}
}
