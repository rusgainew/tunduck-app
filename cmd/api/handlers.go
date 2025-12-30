package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rusgainew/tunduck-app/internal/controllers"
	"github.com/rusgainew/tunduck-app/internal/services"
	"github.com/rusgainew/tunduck-app/pkg/container"
	"github.com/rusgainew/tunduck-app/pkg/middleware"
)

// RegisterHandlers регистрирует все handlers и routes приложения
func RegisterHandlers(app *fiber.App, cnt *container.Container, orgDBService services.OrganizationDBService) {
	// Инициализируем Rate Limiter
	rateLimiter := cnt.GetRateLimiter()
	logger := cnt.GetLogrus()

	// Инициализируем контроллеры с зависимостями из контейнера
	// Передаем сервисы из контейнера вместо их создания в контроллерах
	controllers.NewAuthController(app, cnt.GetUserService(), logger, cnt.GetCacheManager())
	controllers.NewEsfDocumentController(app, cnt.GetLogrus(), cnt.GetDatabase())
	controllers.NewEsfOrganizationController(app, cnt.GetLogrus(), cnt.GetDatabase())
	controllers.NewUserController(app, cnt.GetLogrus(), cnt.GetDatabase())

	// Применяем Rate Limiting для публичных endpoints (регистрация, логин)
	// Эти routes переопределяются в auth_controller.go
	// Здесь мы добавляем глобальный rate limit для защиты от DDoS
	app.Use("/api/auth/register", middleware.RateLimitMiddleware(rateLimiter, "public", logger))
	app.Use("/api/auth/login", middleware.RateLimitMiddleware(rateLimiter, "public", logger))

	// Применяем Rate Limiting для health и metrics endpoints (более высокий лимит)
	app.Use("/health", middleware.RateLimitMiddleware(rateLimiter, "health", logger))
	app.Use("/metrics", middleware.RateLimitMiddleware(rateLimiter, "metrics", logger))

	// Применяем Rate Limiting для sensitive endpoints (logout)
	// Используем auth middleware для определения пользователя
	app.Use("/api/auth/logout", middleware.RateLimitAuthMiddleware(rateLimiter, "sensitive", logger))

	// Применяем Rate Limiting для защищенных endpoints с auth middleware
	// Это предотвращает abuse и enumeration атаки
	protected := app.Group("/api")
	protected.Use(middleware.RateLimitAuthMiddleware(rateLimiter, "protected", logger))

	// Все user endpoints защищены rate limiting
	protected.Use("/users", middleware.RateLimitAuthMiddleware(rateLimiter, "protected", logger))

	// Все organization endpoints защищены rate limiting
	protected.Use("/esf-organizations", middleware.RateLimitAuthMiddleware(rateLimiter, "protected", logger))

	// Все document endpoints защищены rate limiting
	protected.Use("/esf-documents", middleware.RateLimitAuthMiddleware(rateLimiter, "protected", logger))
}
