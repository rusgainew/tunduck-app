package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestIDMiddleware создает middleware для добавления уникального ID каждому запросу
// Используется для трассировки запросов через логи и ответы
func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Пытаемся получить RequestID из заголовка запроса
		requestID := c.Get("X-Request-ID")

		// Если не найден, генерируем новый
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Добавляем RequestID в контекст для использования в обработчиках
		c.Locals("request_id", requestID)

		// Добавляем RequestID в заголовок ответа
		c.Set("X-Request-ID", requestID)

		return c.Next()
	}
}
