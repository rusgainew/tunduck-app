package controllers

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/rusgainew/tunduck-app/pkg/apperror"
	"github.com/rusgainew/tunduck-app/pkg/entity"
	"github.com/rusgainew/tunduck-app/pkg/logger"
)

type UserController struct {
	logger *logger.Logger
	db     *gorm.DB
}

func NewUserController(app *fiber.App, log *logrus.Logger, db *gorm.DB) {
	controller := &UserController{
		logger: logger.New(log),
		db:     db,
	}

	controller.logger.Info(context.Background(), "UserController initialized", logrus.Fields{})
	controller.registerRoutes(app)
}

func (c *UserController) registerRoutes(app *fiber.App) {
	userGroup := app.Group("/api/users")

	// Публичные routes (без JWT)
	userGroup.Get("/", c.getAllUsers)
	userGroup.Get("/:id", c.getUserByID)
}

// getAllUsers возвращает всех пользователей с пагинацией
func (c *UserController) getAllUsers(ctx *fiber.Ctx) error {
	c.logger.Info(ctx.Context(), "Fetching all users", logrus.Fields{})

	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var users []entity.User
	var total int64

	// Получаем общее количество пользователей
	if err := c.db.Model(&entity.User{}).Count(&total).Error; err != nil {
		c.logger.Error(ctx.Context(), "Failed to count users", err, logrus.Fields{})
		appErr := apperror.New(apperror.ErrInternal, "failed to count users").WithError(err)
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	// Получаем пользователей с пагинацией
	if err := c.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		c.logger.Error(ctx.Context(), "Failed to fetch users", err, logrus.Fields{})
		appErr := apperror.New(apperror.ErrInternal, "failed to fetch users").WithError(err)
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"data":  users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// getUserByID возвращает пользователя по ID
func (c *UserController) getUserByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	c.logger.Info(ctx.Context(), "Fetching user by ID", logrus.Fields{"id": id})

	var user entity.User
	if err := c.db.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			appErr := apperror.New(apperror.ErrNotFound, "user not found")
			return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
		}
		c.logger.Error(ctx.Context(), "Failed to fetch user", err, logrus.Fields{"id": id})
		appErr := apperror.New(apperror.ErrInternal, "failed to fetch user").WithError(err)
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	return ctx.Status(http.StatusOK).JSON(user)
}
