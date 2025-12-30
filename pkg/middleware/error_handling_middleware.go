package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rusgainew/tunduck-app/pkg/apperror"
	"github.com/sirupsen/logrus"
)

// ErrorHandlingMiddleware обрабатывает ошибки приложения и преобразует их в HTTP ответы
func ErrorHandlingMiddleware(logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Выполняем следующий handler
		err := c.Next()

		// Если ошибки нет, возвращаем успешный результат
		if err == nil {
			return nil
		}

		// Получаем requestID из контекста для логирования
		requestID := c.Get("X-Request-ID")

		// Преобразуем ошибку в AppError если это возможно
		var appErr *apperror.AppError
		var httpStatus int
		var response fiber.Map

		// Проверяем, является ли ошибка AppError
		if ae, ok := err.(*apperror.AppError); ok {
			appErr = ae
			httpStatus = ae.HTTPStatus

			// Логируем ошибку приложения
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"method":     c.Method(),
				"path":       c.Path(),
				"status":     httpStatus,
				"error_code": appErr.Code,
				"message":    appErr.Message,
				"details":    appErr.Details,
			}).Warn("Application error occurred")

			response = fiber.Map{
				"error":      appErr.Code,
				"message":    appErr.Message,
				"request_id": requestID,
			}

			// Добавляем детали если они есть
			if appErr.Details != "" {
				response["details"] = appErr.Details
			}

		} else if fe, ok := err.(*fiber.Error); ok {
			// Обработка Fiber ошибок
			httpStatus = fe.Code
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"method":     c.Method(),
				"path":       c.Path(),
				"status":     httpStatus,
				"message":    fe.Message,
			}).Warn("HTTP error occurred")

			response = fiber.Map{
				"error":      "HTTP_ERROR",
				"message":    fe.Message,
				"request_id": requestID,
			}

		} else {
			// Обработка неизвестных ошибок
			httpStatus = fiber.StatusInternalServerError
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"method":     c.Method(),
				"path":       c.Path(),
				"status":     httpStatus,
				"error":      err.Error(),
			}).Error("Unknown error occurred")

			response = fiber.Map{
				"error":      "INTERNAL_SERVER_ERROR",
				"message":    "An unexpected error occurred",
				"request_id": requestID,
			}
		}

		// Отправляем ошибку в формате JSON
		return c.Status(httpStatus).JSON(response)
	}
}

// GlobalErrorHandler обрабатывает паники и необработанные ошибки на уровне приложения
func GlobalErrorHandler(logger *logrus.Logger) func(*fiber.Ctx, error) error {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		message := "Internal Server Error"

		requestID := c.Get("X-Request-ID")

		// Пытаемся извлечь AppError
		if appErr, ok := err.(*apperror.AppError); ok {
			code = appErr.HTTPStatus
			message = appErr.Message

			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"method":     c.Method(),
				"path":       c.Path(),
				"status":     code,
				"error_code": appErr.Code,
			}).Error("Global error handler caught application error")

		} else if fe, ok := err.(*fiber.Error); ok {
			code = fe.Code
			message = fe.Message

			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"method":     c.Method(),
				"path":       c.Path(),
				"status":     code,
			}).Warn("Global error handler caught Fiber error")

		} else {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"method":     c.Method(),
				"path":       c.Path(),
				"status":     code,
				"error":      err.Error(),
			}).Error("Global error handler caught unknown error")
		}

		// Не отправляем ответ, если заголовки уже отправлены
		if c.Response().StatusCode() != 0 {
			return nil
		}

		return c.Status(code).JSON(fiber.Map{
			"error":      "ERROR",
			"message":    message,
			"request_id": requestID,
		})
	}
}
