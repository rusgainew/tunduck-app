package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rusgainew/tunduck-app/pkg/apperror"
)

// SuccessResponse представляет успешный ответ API
type SuccessResponse struct {
	Code      int         `json:"code"`
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

// ErrorResponse представляет ошибочный ответ API
type ErrorResponse struct {
	Code      string `json:"code"`
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Details   string `json:"details,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

// Success отправляет успешный ответ
func Success(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(SuccessResponse{
		Code:      statusCode,
		Status:    "success",
		Message:   message,
		Data:      data,
		RequestID: c.Get("X-Request-ID"),
	})
}

// SuccessCreated отправляет 201 Created ответ
func SuccessCreated(c *fiber.Ctx, message string, data interface{}) error {
	return Success(c, fiber.StatusCreated, message, data)
}

// SuccessOK отправляет 200 OK ответ
func SuccessOK(c *fiber.Ctx, message string, data interface{}) error {
	return Success(c, fiber.StatusOK, message, data)
}

// SuccessNoContent отправляет 204 No Content ответ
func SuccessNoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// Error отправляет ошибочный ответ
func Error(c *fiber.Ctx, appErr *apperror.AppError) error {
	response := ErrorResponse{
		Code:      string(appErr.Code),
		Status:    appErr.HTTPStatus,
		Message:   appErr.Message,
		Details:   appErr.Details,
		RequestID: c.Get("X-Request-ID"),
	}

	return c.Status(appErr.HTTPStatus).JSON(response)
}

// BadRequest отправляет 400 Bad Request
func BadRequest(c *fiber.Ctx, message string) error {
	appErr := apperror.New(apperror.ErrValidation, message)
	return Error(c, appErr)
}

// Unauthorized отправляет 401 Unauthorized
func Unauthorized(c *fiber.Ctx, message string) error {
	appErr := apperror.New(apperror.ErrUnauthorized, message)
	return Error(c, appErr)
}

// Forbidden отправляет 403 Forbidden
func Forbidden(c *fiber.Ctx, message string) error {
	appErr := apperror.New(apperror.ErrForbidden, message)
	return Error(c, appErr)
}

// NotFound отправляет 404 Not Found
func NotFound(c *fiber.Ctx, message string) error {
	appErr := apperror.New(apperror.ErrNotFound, message)
	return Error(c, appErr)
}

// Conflict отправляет 409 Conflict
func Conflict(c *fiber.Ctx, message string) error {
	appErr := apperror.New(apperror.ErrConflict, message)
	return Error(c, appErr)
}

// InternalError отправляет 500 Internal Server Error
func InternalError(c *fiber.Ctx, message string) error {
	appErr := apperror.New(apperror.ErrInternal, message)
	return Error(c, appErr)
}
