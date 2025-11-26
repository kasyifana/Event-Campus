package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Custom error types
var (
	ErrBadRequest   = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrInternal     = errors.New("internal server error")
)

// AppError represents an application error with HTTP status code
type AppError struct {
	Err        error
	Message    string
	StatusCode int
}

// Error implements error interface
func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

// NewBadRequestError creates a bad request error
func NewBadRequestError(message string) *AppError {
	return &AppError{
		Err:        ErrBadRequest,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Err:        ErrUnauthorized,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Err:        ErrForbidden,
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Err:        ErrNotFound,
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

// NewConflictError creates a conflict error
func NewConflictError(message string) *AppError {
	return &AppError{
		Err:        ErrConflict,
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

// NewInternalError creates an internal server error
func NewInternalError(message string) *AppError {
	return &AppError{
		Err:        ErrInternal,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

// HandleError handles application errors and sends appropriate response
func HandleError(c *gin.Context, err error) {
	var appErr *AppError

	// Check if error is AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.StatusCode, gin.H{
			"success": false,
			"message": "Request failed",
			"error":   appErr.Error(),
		})
		return
	}

	// Handle specific error types
	switch {
	case errors.Is(err, ErrBadRequest):
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Bad request",
			"error":   err.Error(),
		})
	case errors.Is(err, ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   err.Error(),
		})
	case errors.Is(err, ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Forbidden",
			"error":   err.Error(),
		})
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Not found",
			"error":   err.Error(),
		})
	case errors.Is(err, ErrConflict):
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Conflict",
			"error":   err.Error(),
		})
	default:
		// Log the error (in production, use proper logging)
		fmt.Printf("Internal error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal server error",
			"error":   "Something went wrong",
		})
	}
}

// WrapError wraps an error with a message
func WrapError(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}
