package errors

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NotFoundError indicates the requested resource was not found.
type NotFoundError struct {
	Code    string
	Message string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewNotFoundError creates a new NotFoundError.
func NewNotFoundError(code, message string) *NotFoundError {
	return &NotFoundError{Code: code, Message: message}
}

// ValidationError indicates invalid input data.
type ValidationError struct {
	Code    string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewValidationError creates a new ValidationError.
func NewValidationError(code, message string) *ValidationError {
	return &ValidationError{Code: code, Message: message}
}

// UnauthorizedError indicates missing or invalid authentication.
type UnauthorizedError struct {
	Code    string
	Message string
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewUnauthorizedError creates a new UnauthorizedError.
func NewUnauthorizedError(code, message string) *UnauthorizedError {
	return &UnauthorizedError{Code: code, Message: message}
}

// ForbiddenError indicates insufficient permissions.
type ForbiddenError struct {
	Code    string
	Message string
}

func (e *ForbiddenError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewForbiddenError creates a new ForbiddenError.
func NewForbiddenError(code, message string) *ForbiddenError {
	return &ForbiddenError{Code: code, Message: message}
}

// ConflictError indicates a resource conflict (e.g., duplicate).
type ConflictError struct {
	Code    string
	Message string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewConflictError creates a new ConflictError.
func NewConflictError(code, message string) *ConflictError {
	return &ConflictError{Code: code, Message: message}
}

// InternalError indicates an unexpected server error.
type InternalError struct {
	Code    string
	Message string
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewInternalError creates a new InternalError.
func NewInternalError(code, message string) *InternalError {
	return &InternalError{Code: code, Message: message}
}

// ToGRPCStatus converts a custom error to a gRPC status.
func ToGRPCStatus(err error) *status.Status {
	if err == nil {
		return status.New(codes.OK, "")
	}

	switch e := err.(type) {
	case *NotFoundError:
		return status.New(codes.NotFound, e.Message)
	case *ValidationError:
		return status.New(codes.InvalidArgument, e.Message)
	case *UnauthorizedError:
		return status.New(codes.Unauthenticated, e.Message)
	case *ForbiddenError:
		return status.New(codes.PermissionDenied, e.Message)
	case *ConflictError:
		return status.New(codes.AlreadyExists, e.Message)
	case *InternalError:
		return status.New(codes.Internal, e.Message)
	default:
		return status.New(codes.Unknown, err.Error())
	}
}

// ToHTTPStatus converts a custom error to an HTTP status code.
func ToHTTPStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err.(type) {
	case *NotFoundError:
		return http.StatusNotFound
	case *ValidationError:
		return http.StatusBadRequest
	case *UnauthorizedError:
		return http.StatusUnauthorized
	case *ForbiddenError:
		return http.StatusForbidden
	case *ConflictError:
		return http.StatusConflict
	case *InternalError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
