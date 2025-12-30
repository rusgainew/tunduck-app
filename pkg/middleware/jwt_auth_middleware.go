package middleware

import (
	"fmt"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rusgainew/tunduck-app/pkg/apperror"
	"github.com/rusgainew/tunduck-app/pkg/response"
	"github.com/sirupsen/logrus"
)

// JWTAuthMiddleware создает улучшенный middleware для проверки JWT токенов
// с логированием и правильной обработкой ошибок
func JWTAuthMiddleware(secret string, logger *logrus.Logger) fiber.Handler {
	if secret == "" {
		logger.Fatal("JWT_SECRET is required for JWT middleware")
		return func(c *fiber.Ctx) error {
			return response.Error(c, apperror.New(apperror.ErrInternal, "JWT middleware configuration error"))
		}
	}

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			requestID := c.Get("X-Request-ID")

			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"method":     c.Method(),
				"path":       c.Path(),
				"error":      err.Error(),
			}).Warn("JWT validation failed")

			appErr := apperror.New(apperror.ErrInvalidToken, "Invalid or expired JWT token")
			return response.Error(c, appErr)
		},
		ContextKey: "user",
	})
}

// JWTOptionalMiddleware создает опциональный JWT middleware
// Если токен присутствует - валидирует его, если нет - продолжает выполнение
func JWTOptionalMiddleware(secret string, logger *logrus.Logger) fiber.Handler {
	if secret == "" {
		logger.Fatal("JWT_SECRET is required for JWT middleware")
		return func(c *fiber.Ctx) error {
			return response.Error(c, apperror.New(apperror.ErrInternal, "JWT middleware configuration error"))
		}
	}

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		requestID := c.Get("X-Request-ID")

		// Если заголовок Authorization отсутствует, продолжаем без токена
		if authHeader == "" {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
			}).Debug("No Authorization header provided, proceeding as anonymous")
			return c.Next()
		}

		// Если заголовок есть, валидируем токен
		config := jwtware.Config{
			SigningKey: jwtware.SigningKey{Key: []byte(secret)},
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				logger.WithFields(logrus.Fields{
					"request_id": requestID,
					"method":     c.Method(),
					"path":       c.Path(),
					"error":      err.Error(),
				}).Warn("JWT validation failed for optional auth")

				appErr := apperror.New(apperror.ErrInvalidToken, "Invalid or expired JWT token")
				return response.Error(c, appErr)
			},
			ContextKey: "user",
		}

		handler := jwtware.New(config)
		return handler(c)
	}
}

// GetUserIDFromContext извлекает ID пользователя из JWT токена в контексте
func GetUserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	user := c.Locals("user")
	if user == nil {
		return uuid.UUID{}, apperror.New(apperror.ErrUnauthorized, "No user in context")
	}

	token, ok := user.(*jwt.Token)
	if !ok {
		// Логируем тип для отладки
		return uuid.UUID{}, apperror.New(apperror.ErrInvalidToken, fmt.Sprintf("Invalid token format in context, got type: %T", user))
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		// Логируем тип Claims для отладки
		return uuid.UUID{}, apperror.New(apperror.ErrInvalidToken, fmt.Sprintf("Invalid token claims, got type: %T", token.Claims))
	}

	// Извлекаем user_id из claims
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.UUID{}, apperror.New(apperror.ErrInvalidToken, "user_id not found in token claims")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.UUID{}, apperror.New(apperror.ErrInvalidToken, fmt.Sprintf("Invalid user_id format: %v", err))
	}

	return userID, nil
}

// GetUsernameFromContext извлекает username из JWT токена в контексте
func GetUsernameFromContext(c *fiber.Ctx) (string, error) {
	user := c.Locals("user")
	if user == nil {
		return "", apperror.New(apperror.ErrUnauthorized, "No user in context")
	}

	token, ok := user.(*jwt.Token)
	if !ok {
		return "", apperror.New(apperror.ErrInvalidToken, "Invalid token format in context")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", apperror.New(apperror.ErrInvalidToken, "Invalid token claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", apperror.New(apperror.ErrInvalidToken, "username not found in token claims")
	}

	return username, nil
}

// GetClaimsFromContext извлекает все claims из JWT токена
func GetClaimsFromContext(c *fiber.Ctx) (jwt.MapClaims, error) {
	user := c.Locals("user")
	if user == nil {
		return nil, apperror.New(apperror.ErrUnauthorized, "No user in context")
	}

	token, ok := user.(*jwt.Token)
	if !ok {
		return nil, apperror.New(apperror.ErrInvalidToken, "Invalid token format in context")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, apperror.New(apperror.ErrInvalidToken, "Invalid token claims")
	}

	return claims, nil
}
