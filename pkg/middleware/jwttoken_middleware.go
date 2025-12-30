package middleware

import (
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware создает middleware для проверки JWT токенов
func JWTMiddleware() fiber.Handler {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Возвращаем middleware который всегда возвращает ошибку
		return func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Configuration Error",
				"message": "JWT_SECRET environment variable is not set",
			})
		}
	}

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Invalid or expired JWT",
			})
		},
		ContextKey: "user",
	})
}

// GetUserFromToken извлекает информацию о пользователе из JWT токена
func GetUserFromToken(c *fiber.Ctx) (*jwt.MapClaims, error) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return &claims, nil
}

// OptionalJWT создает опциональный JWT middleware (не требует токен)
func OptionalJWT() fiber.Handler {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Возвращаем middleware который всегда продолжает выполнение
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	return func(c *fiber.Ctx) error {
		config := jwtware.Config{
			SigningKey: jwtware.SigningKey{Key: []byte(secret)},
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				// Игнорируем ошибки и продолжаем выполнение
				return c.Next()
			},
			ContextKey: "user",
		}

		handler := jwtware.New(config)
		return handler(c)
	}
}
