package controllers

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	repositorypostgres "github.com/rusgainew/tunduck-app/internal/repository/repository_postgres"
	"github.com/rusgainew/tunduck-app/internal/services"
	serviceimpl "github.com/rusgainew/tunduck-app/internal/services/service_impl"
	"github.com/rusgainew/tunduck-app/pkg/apperror"
	"github.com/rusgainew/tunduck-app/pkg/logger"
	"github.com/rusgainew/tunduck-app/pkg/middleware"
	"github.com/rusgainew/tunduck-app/pkg/rbac"
)

type RoleController struct {
	logger      *logger.Logger
	roleService services.RoleService
	db          *gorm.DB
}

func NewRoleController(app *fiber.App, log *logrus.Logger, db *gorm.DB) {
	// Инициализируем слои
	userRepo := repositorypostgres.NewUserRepositoryPostgres(db, log)
	roleService := serviceimpl.NewRoleService(userRepo, log)

	controller := &RoleController{
		logger:      logger.New(log),
		roleService: roleService,
		db:          db,
	}

	controller.logger.Info(context.Background(), "RoleController инициализирован", logrus.Fields{})
	controller.registerRoutes(app)
}

func (c *RoleController) registerRoutes(app *fiber.App) {
	roleGroup := app.Group("/api/roles")

	// Защищенные маршруты (требуют JWT)
	protected := roleGroup.Group("")
	protected.Use(middleware.JWTMiddleware())

	// Маршруты требующие роль администратора
	admin := protected.Group("")
	admin.Use(rbac.RequireAdminRole())

	admin.Post("/assign", c.assignRole)
	admin.Put("/:user_id", c.updateRole)
	admin.Get("/users", c.listUsersByRole)

	// Маршруты доступные авторизованным пользователям
	protected.Get("/:user_id", c.getUserRole)
	protected.Get("/permissions/:user_id", c.checkUserPermissions)
}

// assignRole назначает роль пользователю (только для админов)
func (c *RoleController) assignRole(ctx *fiber.Ctx) error {
	c.logger.Info(ctx.Context(), "Назначение роли пользователю", logrus.Fields{})

	var req struct {
		UserID uuid.UUID `json:"user_id" validate:"required"`
		Role   string    `json:"role" validate:"required"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		c.logger.Warn(ctx.Context(), "Невалидное тело запроса", logrus.Fields{"error": err.Error()})
		appErr := apperror.ValidationError("невалидное тело запроса")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	// Валидируем роль
	role := rbac.Role(req.Role)
	if !role.IsValid() {
		c.logger.Warn(ctx.Context(), "Невалидная роль", logrus.Fields{"role": req.Role})
		appErr := apperror.ValidationError("невалидная роль: " + req.Role)
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	if err := c.roleService.AssignRole(ctx.Context(), req.UserID, role); err != nil {
		appErr, ok := err.(*apperror.AppError)
		if !ok {
			appErr = apperror.New(apperror.ErrInternal, "ошибка при назначении роли").WithError(err)
		}
		c.logger.Error(ctx.Context(), "Ошибка назначения роли", err, logrus.Fields{
			"user_id": req.UserID.String(),
			"role":    req.Role,
		})
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	c.logger.Info(ctx.Context(), "Роль успешно назначена", logrus.Fields{
		"user_id": req.UserID.String(),
		"role":    req.Role,
	})

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Роль успешно назначена",
		"user_id": req.UserID.String(),
		"role":    req.Role,
	})
}

// updateRole обновляет роль пользователя (только для админов)
func (c *RoleController) updateRole(ctx *fiber.Ctx) error {
	userIDParam := ctx.Params("user_id")
	c.logger.Info(ctx.Context(), "Обновление роли пользователя", logrus.Fields{
		"user_id": userIDParam,
	})

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.logger.Warn(ctx.Context(), "Невалидный формат UUID", logrus.Fields{"user_id": userIDParam})
		appErr := apperror.ValidationError("невалидный формат UUID")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	var req struct {
		Role string `json:"role" validate:"required"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		c.logger.Warn(ctx.Context(), "Невалидное тело запроса", logrus.Fields{"error": err.Error()})
		appErr := apperror.ValidationError("невалидное тело запроса")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	// Валидируем роль
	role := rbac.Role(req.Role)
	if !role.IsValid() {
		c.logger.Warn(ctx.Context(), "Невалидная роль", logrus.Fields{"role": req.Role})
		appErr := apperror.ValidationError("невалидная роль: " + req.Role)
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	if err := c.roleService.UpdateRole(ctx.Context(), userID, role); err != nil {
		appErr, ok := err.(*apperror.AppError)
		if !ok {
			appErr = apperror.New(apperror.ErrInternal, "ошибка при обновлении роли").WithError(err)
		}
		c.logger.Error(ctx.Context(), "Ошибка обновления роли", err, logrus.Fields{
			"user_id": userID.String(),
			"role":    req.Role,
		})
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	c.logger.Info(ctx.Context(), "Роль успешно обновлена", logrus.Fields{
		"user_id": userID.String(),
		"role":    req.Role,
	})

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Роль успешно обновлена",
		"user_id": userID.String(),
		"role":    req.Role,
	})
}

// getUserRole получает роль пользователя
func (c *RoleController) getUserRole(ctx *fiber.Ctx) error {
	userIDParam := ctx.Params("user_id")
	c.logger.Debug(ctx.Context(), "Получение роли пользователя", logrus.Fields{
		"user_id": userIDParam,
	})

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.logger.Warn(ctx.Context(), "Невалидный формат UUID", logrus.Fields{"user_id": userIDParam})
		appErr := apperror.ValidationError("невалидный формат UUID")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	role, err := c.roleService.GetUserRole(ctx.Context(), userID)
	if err != nil {
		appErr, ok := err.(*apperror.AppError)
		if !ok {
			appErr = apperror.New(apperror.ErrInternal, "ошибка при получении роли").WithError(err)
		}
		c.logger.Error(ctx.Context(), "Ошибка получения роли", err, logrus.Fields{
			"user_id": userID.String(),
		})
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	c.logger.Debug(ctx.Context(), "Роль успешно получена", logrus.Fields{
		"user_id": userID.String(),
		"role":    role.String(),
	})

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"user_id": userID.String(),
		"role":    role.String(),
	})
}

// checkUserPermissions проверяет разрешения пользователя
func (c *RoleController) checkUserPermissions(ctx *fiber.Ctx) error {
	userIDParam := ctx.Params("user_id")
	c.logger.Debug(ctx.Context(), "Проверка разрешений пользователя", logrus.Fields{
		"user_id": userIDParam,
	})

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.logger.Warn(ctx.Context(), "Невалидный формат UUID", logrus.Fields{"user_id": userIDParam})
		appErr := apperror.ValidationError("невалидный формат UUID")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	role, err := c.roleService.GetUserRole(ctx.Context(), userID)
	if err != nil {
		appErr, ok := err.(*apperror.AppError)
		if !ok {
			appErr = apperror.New(apperror.ErrInternal, "ошибка при получении роли").WithError(err)
		}
		c.logger.Error(ctx.Context(), "Ошибка получения роли", err, logrus.Fields{
			"user_id": userID.String(),
		})
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	// Получаем все разрешения для роли
	permissions := role.GetAllPermissions()

	c.logger.Debug(ctx.Context(), "Разрешения успешно получены", logrus.Fields{
		"user_id":     userID.String(),
		"role":        role.String(),
		"permissions": len(permissions),
	})

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success":     true,
		"user_id":     userID.String(),
		"role":        role.String(),
		"permissions": permissions,
	})
}

// listUsersByRole получает список пользователей с определенной ролью (только для админов)
func (c *RoleController) listUsersByRole(ctx *fiber.Ctx) error {
	roleParam := ctx.Query("role")
	c.logger.Debug(ctx.Context(), "Получение списка пользователей с ролью", logrus.Fields{
		"role": roleParam,
	})

	role := rbac.Role(roleParam)
	if !role.IsValid() {
		c.logger.Warn(ctx.Context(), "Невалидная роль", logrus.Fields{"role": roleParam})
		appErr := apperror.ValidationError("невалидная роль: " + roleParam)
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	userIDs, err := c.roleService.ListUsersByRole(ctx.Context(), role)
	if err != nil {
		appErr, ok := err.(*apperror.AppError)
		if !ok {
			appErr = apperror.New(apperror.ErrInternal, "ошибка при получении списка пользователей").WithError(err)
		}
		c.logger.Error(ctx.Context(), "Ошибка получения списка пользователей", err, logrus.Fields{
			"role": role.String(),
		})
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	c.logger.Debug(ctx.Context(), "Список пользователей успешно получен", logrus.Fields{
		"role":  role.String(),
		"count": len(userIDs),
	})

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"role":    role.String(),
		"users":   userIDs,
		"count":   len(userIDs),
	})
}
