package apperror

import (
	"fmt"
	"net/http"
)

// ErrorCode определяет тип ошибки приложения
type ErrorCode string

const (
	// Validation errors
	ErrValidation      ErrorCode = "VALIDATION_ERROR"
	ErrInvalidRequest  ErrorCode = "INVALID_REQUEST"
	ErrFieldValidation ErrorCode = "FIELD_VALIDATION"

	// Authentication errors
	ErrUnauthorized       ErrorCode = "UNAUTHORIZED"
	ErrInvalidToken       ErrorCode = "INVALID_TOKEN"
	ErrExpiredToken       ErrorCode = "EXPIRED_TOKEN"
	ErrInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"

	// Authorization errors
	ErrForbidden    ErrorCode = "FORBIDDEN"
	ErrAccessDenied ErrorCode = "ACCESS_DENIED"

	// Resource errors
	ErrNotFound      ErrorCode = "NOT_FOUND"
	ErrAlreadyExists ErrorCode = "ALREADY_EXISTS"
	ErrConflict      ErrorCode = "CONFLICT"

	// User errors
	ErrUserNotFound     ErrorCode = "USER_NOT_FOUND"
	ErrUserExists       ErrorCode = "USER_ALREADY_EXISTS"
	ErrEmailExists      ErrorCode = "EMAIL_ALREADY_EXISTS"
	ErrUsernameExists   ErrorCode = "USERNAME_ALREADY_EXISTS"
	ErrAccountBlocked   ErrorCode = "ACCOUNT_BLOCKED"
	ErrPasswordMismatch ErrorCode = "PASSWORD_MISMATCH"

	// Document errors
	ErrDocumentNotFound ErrorCode = "DOCUMENT_NOT_FOUND"
	ErrInvalidDocument  ErrorCode = "INVALID_DOCUMENT"

	// Organization errors
	ErrOrgNotFound ErrorCode = "ORGANIZATION_NOT_FOUND"
	ErrOrgExists   ErrorCode = "ORGANIZATION_ALREADY_EXISTS"

	// Database errors
	ErrDatabase        ErrorCode = "DATABASE_ERROR"
	ErrDatabaseTimeout ErrorCode = "DATABASE_TIMEOUT"

	// External service errors
	ErrExternalService ErrorCode = "EXTERNAL_SERVICE_ERROR"

	// Server errors
	ErrInternal    ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrConfigError ErrorCode = "CONFIG_ERROR"
)

// AppError представляет структурированную ошибку приложения
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	HTTPStatus int       `json:"-"`
	Err        error     `json:"-"` // Оригинальная ошибка для логирования
	StackTrace string    `json:"-"`
}

// Error реализует интерфейс error
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// New создает новую ошибку приложения
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: getHTTPStatus(code),
	}
}

// NewWithDetails создает ошибку с дополнительными деталями
func NewWithDetails(code ErrorCode, message, details string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		HTTPStatus: getHTTPStatus(code),
	}
}

// WithError обертывает оригинальную ошибку
func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

// WithDetails добавляет детали к ошибке
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// WithHTTPStatus устанавливает HTTP статус
func (e *AppError) WithHTTPStatus(status int) *AppError {
	e.HTTPStatus = status
	return e
}

// getHTTPStatus возвращает соответствующий HTTP статус для ErrorCode
func getHTTPStatus(code ErrorCode) int {
	switch code {
	// 400 Bad Request
	case ErrValidation, ErrInvalidRequest, ErrFieldValidation,
		ErrPasswordMismatch, ErrInvalidDocument:
		return http.StatusBadRequest

	// 401 Unauthorized
	case ErrUnauthorized, ErrInvalidToken, ErrExpiredToken, ErrInvalidCredentials:
		return http.StatusUnauthorized

	// 403 Forbidden
	case ErrForbidden, ErrAccessDenied:
		return http.StatusForbidden

	// 404 Not Found
	case ErrNotFound, ErrUserNotFound, ErrDocumentNotFound, ErrOrgNotFound:
		return http.StatusNotFound

	// 409 Conflict
	case ErrAlreadyExists, ErrConflict, ErrUserExists, ErrEmailExists,
		ErrUsernameExists, ErrOrgExists, ErrAccountBlocked:
		return http.StatusConflict

	// 500 Internal Server Error
	case ErrDatabase, ErrDatabaseTimeout, ErrExternalService,
		ErrInternal, ErrConfigError:
		return http.StatusInternalServerError

	default:
		return http.StatusInternalServerError
	}
}

// Is проверяет, соответствует ли ошибка другой ошибке
func (e *AppError) Is(target error) bool {
	if e == nil {
		return false
	}
	if other, ok := target.(*AppError); ok {
		return e.Code == other.Code
	}
	return false
}

// ErrorResponse структура для отправки ошибки в HTTP ответе
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// ToResponse преобразует AppError в ErrorResponse
func (e *AppError) ToResponse() *ErrorResponse {
	return &ErrorResponse{
		Code:    string(e.Code),
		Message: e.Message,
		Details: e.Details,
	}
}

// Common error constructors for convenience

func ValidationError(message string) *AppError {
	return New(ErrValidation, message)
}

func UnauthorizedError(message string) *AppError {
	return New(ErrUnauthorized, message)
}

func ForbiddenError(message string) *AppError {
	return New(ErrForbidden, message)
}

func NotFoundError(resource string) *AppError {
	return New(ErrNotFound, fmt.Sprintf("%s not found", resource))
}

func ConflictError(message string) *AppError {
	return New(ErrConflict, message)
}

func InternalError(message string) *AppError {
	return New(ErrInternal, message)
}

func DatabaseError(operation string, err error) *AppError {
	return New(ErrDatabase, fmt.Sprintf("database error during %s", operation)).WithError(err)
}
