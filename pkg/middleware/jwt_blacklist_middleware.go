package middleware

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rusgainew/tunduck-app/pkg/apperror"
	"github.com/rusgainew/tunduck-app/pkg/cache"
	"github.com/rusgainew/tunduck-app/pkg/response"
	"github.com/sirupsen/logrus"
)

// JWTBlacklistMiddleware создает JWT middleware с поддержкой blacklist для logout
func JWTBlacklistMiddleware(secret string, logger *logrus.Logger, cacheManager cache.CacheManager) fiber.Handler {
	if secret == "" {
		logger.Fatal("JWT_SECRET is required for JWT middleware")
		return func(c *fiber.Ctx) error {
			return response.Error(c, apperror.New(apperror.ErrInternal, "JWT middleware configuration error"))
		}
	}

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		requestID := c.Get("X-Request-ID")

		// Проверяем Authorization header
		if authHeader == "" {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
			}).Warn("Missing Authorization header")
			return response.Error(c, apperror.New(apperror.ErrUnauthorized, "Authorization header required"))
		}

		// Извлекаем токен
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
			}).Warn("Invalid Authorization header format")
			return response.Error(c, apperror.New(apperror.ErrInvalidToken, "Authorization header must be Bearer token"))
		}

		token := parts[1]

		// Парсим токен с валидацией подписи
		// Важно: используем jwt.MapClaims{} без указателя для совместимости с GetUserIDFromContext
		parsedToken, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !parsedToken.Valid {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
				"error":      err.Error(),
			}).Warn("JWT validation failed")
			return response.Error(c, apperror.New(apperror.ErrInvalidToken, "Invalid or expired JWT token"))
		}

		// Проверяем blacklist
		tokenHash := hashToken(token)
		isBlacklisted, err := cacheManager.Token().Exists(context.Background(), "blacklist:"+tokenHash)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
				"error":      err.Error(),
			}).Warn("Failed to check token blacklist")
			// Не блокируем запрос, если кеш недоступен
		}

		if isBlacklisted {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
			}).Warn("Blacklisted token used")
			return response.Error(c, apperror.New(apperror.ErrInvalidToken, "Token has been revoked"))
		}

		// Сохраняем parsed token (не claims) в контекст для использования в handlers
		// GetUserIDFromContext и GetClaimsFromContext ожидают *jwt.Token
		c.Locals("user", parsedToken)
		c.Locals("token", token)
		c.Locals("token_hash", tokenHash)

		return c.Next()
	}
}

// AddTokenToBlacklist добавляет токен в черный список
func AddTokenToBlacklist(ctx context.Context, token string, expiryTime time.Time, cacheManager cache.CacheManager) error {
	tokenHash := hashToken(token)
	ttl := time.Until(expiryTime)
	if ttl < 0 {
		ttl = time.Minute // минимум 1 минута
	}

	return cacheManager.Token().Set(ctx, "blacklist:"+tokenHash, "revoked", ttl)
}

// IsTokenBlacklisted проверяет находится ли токен в blacklist
func IsTokenBlacklisted(ctx context.Context, token string, cacheManager cache.CacheManager) (bool, error) {
	tokenHash := hashToken(token)
	return cacheManager.Token().Exists(ctx, "blacklist:"+tokenHash)
}

// hashToken создает хеш токена для экономии памяти в кеше
// Предпочтительнее, чем сохранять весь токен в blacklist
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}

// LogoutHandler предоставляет логику для logout с добавлением токена в blacklist
func CreateLogoutHandler(cacheManager cache.CacheManager, logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		token := c.Locals("token").(string)
		userClaims := c.Locals("user").(*jwt.MapClaims)

		// Получаем время экспирации токена
		expiryTime := time.Now().Add(time.Hour) // По умолчанию 1 час
		if exp, ok := (*userClaims)["exp"].(float64); ok {
			expiryTime = time.Unix(int64(exp), 0)
		}

		// Добавляем токен в blacklist
		if err := AddTokenToBlacklist(c.Context(), token, expiryTime, cacheManager); err != nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
			}).Error("Failed to add token to blacklist")
			return response.Error(c, apperror.New(apperror.ErrInternal, "Failed to logout"))
		}

		logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"user_id":    (*userClaims)["user_id"],
		}).Info("User logged out successfully")

		return c.JSON(map[string]string{
			"message": "Logged out successfully",
		})
	}
}
