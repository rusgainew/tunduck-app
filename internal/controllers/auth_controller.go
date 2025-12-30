package controllers

import (
	"context"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rusgainew/tunduck-app/internal/models"
	"github.com/rusgainew/tunduck-app/internal/services"
	"github.com/rusgainew/tunduck-app/pkg/apperror"
	"github.com/rusgainew/tunduck-app/pkg/cache"
	"github.com/rusgainew/tunduck-app/pkg/logger"
	"github.com/rusgainew/tunduck-app/pkg/middleware"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	logger       *logger.Logger
	service      services.UserService
	validate     *validator.Validate
	cacheManager cache.CacheManager
}

// NewAuthController инициализирует контроллер с сервисом из контейнера
func NewAuthController(app *fiber.App, userService services.UserService, log *logrus.Logger, cacheManager cache.CacheManager) {
	l := logger.New(log)

	controller := &AuthController{
		logger:       l,
		service:      userService, // Используем сервис из контейнера
		validate:     validator.New(),
		cacheManager: cacheManager,
	}

	l.Info(context.TODO(), "AuthController initialized")
	controller.registerRoutes(app, log)
}

func (c *AuthController) registerRoutes(app *fiber.App, log *logrus.Logger) {
	authGroup := app.Group("/api/auth")

	// Публичные endpoints
	authGroup.Post("/register", c.register)
	authGroup.Post("/login", c.login)

	// Защищенные endpoints с JWT валидацией и поддержкой blacklist для logout
	protected := authGroup.Group("")
	protected.Use(middleware.JWTBlacklistMiddleware(os.Getenv("JWT_SECRET"), log, c.cacheManager))
	protected.Get("/me", c.getCurrentUser)
	protected.Post("/logout", c.logout)
}

// @Summary Регистрация нового пользователя
// @Description Регистрация с проверкой уникальности логина и email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Данные для регистрации"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/auth/register [post]
func (c *AuthController) register(ctx *fiber.Ctx) error {
	var req models.RegisterRequest

	if err := ctx.BodyParser(&req); err != nil {
		c.logger.Warn(ctx.Context(), "Failed to parse register request", logrus.Fields{"error": err.Error()})
		return ctx.Status(fiber.StatusBadRequest).JSON(apperror.New(apperror.ErrInvalidRequest, "invalid request format").ToResponse())
	}

	// Валидация
	if err := c.validate.Struct(req); err != nil {
		c.logger.Warn(ctx.Context(), "Validation failed for register request", logrus.Fields{"error": err.Error()})
		return ctx.Status(fiber.StatusBadRequest).JSON(apperror.New(apperror.ErrValidation, "validation error").WithDetails(err.Error()).ToResponse())
	}

	// Проверка совпадения паролей
	if req.Password != req.ConfirmPassword {
		c.logger.Warn(ctx.Context(), "Password mismatch during registration")
		return ctx.Status(fiber.StatusBadRequest).JSON(apperror.New(apperror.ErrPasswordMismatch, "passwords do not match").ToResponse())
	}

	response, err := c.service.Register(ctx.Context(), &req)
	if err != nil {
		appErr, ok := err.(*apperror.AppError)
		if !ok {
			appErr = apperror.New(apperror.ErrInternal, "registration failed").WithError(err)
		}
		c.logger.Error(ctx.Context(), "Registration failed", err, logrus.Fields{"username": req.Username})
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	c.logger.Info(ctx.Context(), "User registered successfully", logrus.Fields{"username": req.Username})
	return ctx.Status(fiber.StatusCreated).JSON(response)
}

// @Summary Вход в систему
// @Description Аутентификация пользователя по логину и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Данные для входа"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/auth/login [post]
func (c *AuthController) login(ctx *fiber.Ctx) error {
	var req models.LoginRequest

	if err := ctx.BodyParser(&req); err != nil {
		c.logger.Warn(ctx.Context(), "Failed to parse login request", logrus.Fields{"error": err.Error()})
		return ctx.Status(fiber.StatusBadRequest).JSON(apperror.New(apperror.ErrInvalidRequest, "invalid request format").ToResponse())
	}

	// Валидация
	if err := c.validate.Struct(req); err != nil {
		c.logger.Warn(ctx.Context(), "Validation failed for login request", logrus.Fields{"error": err.Error()})
		return ctx.Status(fiber.StatusBadRequest).JSON(apperror.New(apperror.ErrValidation, "validation error").WithDetails(err.Error()).ToResponse())
	}

	response, err := c.service.Login(ctx.Context(), &req)
	if err != nil {
		appErr, ok := err.(*apperror.AppError)
		if !ok {
			appErr = apperror.New(apperror.ErrInternal, "login failed").WithError(err)
		}
		c.logger.Warn(ctx.Context(), "Login failed", logrus.Fields{"username": req.Username})
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	c.logger.Info(ctx.Context(), "User logged in successfully", logrus.Fields{"username": req.Username})
	return ctx.Status(fiber.StatusOK).JSON(response)
}

// @Summary Получить текущего пользователя
// @Description Возвращает информацию о текущем авторизованном пользователе
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.UserInfo
// @Failure 401 {object} models.ErrorResponse
// @Router /api/auth/me [get]
func (c *AuthController) getCurrentUser(ctx *fiber.Ctx) error {
	c.logger.Debug(ctx.Context(), "Getting current user info")

	// Извлекаем user_id из JWT токена используя helper функцию
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		c.logger.Warn(ctx.Context(), "Failed to extract user ID from token", logrus.Fields{"error": err.Error()})
		appErr := err.(*apperror.AppError)
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	// Извлекаем все claims из токена
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		c.logger.Warn(ctx.Context(), "Failed to extract claims from token", logrus.Fields{"error": err.Error()})
		appErr := err.(*apperror.AppError)
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	userInfo := &models.UserInfo{
		ID:       userID,
		Username: claims["username"].(string),
		Email:    claims["email"].(string),
		FullName: claims["full_name"].(string),
	}

	c.logger.Debug(ctx.Context(), "Current user info retrieved", logrus.Fields{"user_id": userID})
	return ctx.Status(fiber.StatusOK).JSON(userInfo)
}

// @Summary Выход из системы
// @Description Логаут пользователя и добавление токена в blacklist
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} models.ErrorResponse
// @Router /api/auth/logout [post]
func (c *AuthController) logout(ctx *fiber.Ctx) error {
	c.logger.Debug(ctx.Context(), "Logout request received")

	// Извлекаем user_id из JWT токена
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		c.logger.Warn(ctx.Context(), "Failed to extract user ID from token", logrus.Fields{"error": err.Error()})
		appErr := err.(*apperror.AppError)
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	// Получаем токен из контекста (установлен middleware'ом)
	token, ok := ctx.Locals("token").(string)
	if !ok {
		c.logger.Warn(ctx.Context(), "Failed to extract token from context")
		appErr := apperror.New(apperror.ErrInvalidToken, "Failed to extract token")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	// Получаем claims для получения времени экспирации
	claims, err := middleware.GetClaimsFromContext(ctx)
	if err != nil {
		c.logger.Warn(ctx.Context(), "Failed to extract claims from token", logrus.Fields{"error": err.Error()})
		appErr := err.(*apperror.AppError)
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	// Получаем время экспирации токена
	var expiryTime int64 = 3600 // По умолчанию 1 час
	if exp, ok := claims["exp"].(float64); ok {
		expiryTime = int64(exp)
	}

	// Добавляем токен в blacklist
	if c.cacheManager != nil {
		if err := middleware.AddTokenToBlacklist(ctx.Context(), token, time.Unix(expiryTime, 0), c.cacheManager); err != nil {
			c.logger.Warn(ctx.Context(), "Failed to add token to blacklist", logrus.Fields{"error": err.Error()})
			// Логируем, но не блокируем logout, если кеш недоступен
		}
	}

	c.logger.Info(ctx.Context(), "User logged out successfully", logrus.Fields{"user_id": userID})
	return ctx.Status(fiber.StatusOK).JSON(map[string]string{
		"message": "Logged out successfully",
	})
}
