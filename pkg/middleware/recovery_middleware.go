package middleware

import (
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// RecoveryMiddleware создает middleware для восстановления после паник
// Перехватывает паники в обработчиках и возвращает ошибку 500
func RecoveryMiddleware(logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				// Логируем паник с полной информацией
				logger.WithFields(logrus.Fields{
					"panic":       r,
					"method":      c.Method(),
					"path":        c.Path(),
					"request_id":  c.Get("X-Request-ID"),
					"stack_trace": string(debug.Stack()),
				}).Error("Panic recovered in request handler")

				// Отправляем ошибку 500
				err = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":      "Internal server error",
					"request_id": c.Get("X-Request-ID"),
				})
			}
		}()

		return c.Next()
	}
}
