package controllers

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	models "github.com/rusgainew/tunduck-app/internal/models"
	repositorypostgres "github.com/rusgainew/tunduck-app/internal/repository/repository_postgres"
	"github.com/rusgainew/tunduck-app/internal/services"
	serviceimpl "github.com/rusgainew/tunduck-app/internal/services/service_impl"
	"github.com/rusgainew/tunduck-app/pkg/apperror"
	"github.com/rusgainew/tunduck-app/pkg/logger"
	"github.com/rusgainew/tunduck-app/pkg/middleware"
	"github.com/rusgainew/tunduck-app/pkg/pagination"
)

type EsfOrganizationController struct {
	logger  *logger.Logger
	service services.EsfOrganizationService
	db      *gorm.DB
}

func NewEsfOrganizationController(app *fiber.App, log *logrus.Logger, db *gorm.DB) {
	// Инициализируем слои
	repo := repositorypostgres.NewEsfOrganizationRepositoryPostgres(db, log)
	service := serviceimpl.NewEsfOrganizationService(repo, log)

	controller := &EsfOrganizationController{
		logger:  logger.New(log),
		service: service,
		db:      db,
	}

	controller.logger.Info(context.Background(), "EsfOrganizationController initialized", logrus.Fields{})
	controller.registerRoutes(app)
}

func (c *EsfOrganizationController) registerRoutes(app *fiber.App) {
	esfOrganizationGroup := app.Group("/api/esf-organizations")

	// Публичные routes (без JWT)
	esfOrganizationGroup.Get("/", c.getEsfOrganizations)
	esfOrganizationGroup.Get("/paginated", c.getEsfOrganizationsPaginated)
	esfOrganizationGroup.Get("/:id", c.getByEsfOrganization)

	// Защищенные routes (с JWT)
	protected := esfOrganizationGroup.Group("")
	protected.Use(middleware.JWTMiddleware())
	protected.Post("/", c.createEsfOrganization)
	protected.Put("/:id", c.updateEsfOrganization)
	protected.Delete("/:id", c.deleteEsfOrganization)
}

// getEsfOrganizations возвращает все организации ЭСФ
func (c *EsfOrganizationController) getEsfOrganizations(ctx *fiber.Ctx) error {
	c.logger.Info(ctx.Context(), "Fetching all ESF organizations", logrus.Fields{})

	organizations, err := c.service.GetAllOrganizations(ctx.Context())
	if err != nil {
		appErr, ok := err.(*apperror.AppError)
		if !ok {
			appErr = apperror.New(apperror.ErrInternal, "failed to fetch organizations").WithError(err)
		}
		c.logger.Error(ctx.Context(), "Failed to fetch organizations", err, logrus.Fields{})
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	return ctx.Status(http.StatusOK).JSON(organizations)
}

// getEsfOrganizationsPaginated возвращает организации ЭСФ с пагинацией
func (c *EsfOrganizationController) getEsfOrganizationsPaginated(ctx *fiber.Ctx) error {
	c.logger.Info(ctx.Context(), "Вибірка організацій ЕСФ з пагінацією", logrus.Fields{})

	// Витягуємо параметри пагінації та фільтрації
	paginationParams := pagination.ExtractPaginationParams(ctx)
	filterParams := pagination.ExtractOrganizationFilters(ctx)

	organizations, totalCount, err := c.service.GetAllOrganizationsPaginated(ctx.Context(), paginationParams, filterParams)
	if err != nil {
		appErr, ok := err.(*apperror.AppError)
		if !ok {
			appErr = apperror.New(apperror.ErrInternal, "не удалось получить организации").WithError(err)
		}
		c.logger.Error(ctx.Context(), "Ошибка вибірки організацій", err, logrus.Fields{})
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	// Формуємо відповідь з пагінацією
	response := pagination.NewPaginatedResponse(organizations, paginationParams.Page, paginationParams.PageSize, totalCount)

	c.logger.Debug(ctx.Context(), "Організації успішно вибрані", logrus.Fields{
		"count": len(organizations),
		"total": totalCount,
		"page":  paginationParams.Page,
	})

	return ctx.Status(http.StatusOK).JSON(response)
}

// createEsfOrganization создает новую организацию ЭСФ
func (c *EsfOrganizationController) createEsfOrganization(ctx *fiber.Ctx) error {
	c.logger.Info(ctx.Context(), "Створення нової ЕСФ організації", logrus.Fields{})

	var req models.EsfOrganizationModel
	if err := ctx.BodyParser(&req); err != nil {
		c.logger.Warn(ctx.Context(), "Невірне тіло запиту", logrus.Fields{"error": err.Error()})
		appErr := apperror.ValidationError("invalid request body")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	// Валидация обязательных полей
	if req.Name == "" {
		c.logger.Warn(ctx.Context(), "Назва організації обов'язкова", logrus.Fields{})
		appErr := apperror.ValidationError("organization name is required")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	id, dbName, err := c.service.CreateOrganization(ctx.Context(), &req)
	if err != nil {
		appErr, ok := err.(*apperror.AppError)
		if !ok {
			appErr = apperror.New(apperror.ErrInternal, "failed to create organization").WithError(err)
		}
		c.logger.Error(ctx.Context(), "Помилка створення організації", err, logrus.Fields{"name": req.Name})
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	c.logger.Info(ctx.Context(), "Організацію успішно створено", logrus.Fields{"id": id.String(), "dbName": dbName})
	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"id":      id,
		"dbName":  dbName,
		"message": "Organization created and database initialized",
	})
}

// getByEsfOrganization возвращает организацию ЭСФ по ID
func (c *EsfOrganizationController) getByEsfOrganization(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	c.logger.Info(ctx.Context(), "Вибірка ЕСФ організації по ID", logrus.Fields{"id": idParam})

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.logger.Warn(ctx.Context(), "Невірний формат UUID", logrus.Fields{"id": idParam})
		appErr := apperror.ValidationError("invalid UUID format")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	organization, err := c.service.GetOrganizationByID(ctx.Context(), id)
	if err != nil {
		appErr, ok := err.(*apperror.AppError)
		if !ok {
			appErr = apperror.New(apperror.ErrInternal, "failed to fetch organization").WithError(err)
		}
		c.logger.Error(ctx.Context(), "Помилка вибірки організації", err, logrus.Fields{"id": id.String()})
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	return ctx.Status(http.StatusOK).JSON(organization)
}

func (c *EsfOrganizationController) updateEsfOrganization(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	c.logger.Info(ctx.Context(), "Updating ESF organization", logrus.Fields{"id": idParam})

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.logger.Warn(ctx.Context(), "Invalid UUID format", logrus.Fields{"id": idParam})
		appErr := apperror.ValidationError("invalid UUID format")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	var req models.EsfOrganizationModel
	if err := ctx.BodyParser(&req); err != nil {
		c.logger.Warn(ctx.Context(), "Invalid request body", logrus.Fields{"error": err.Error()})
		appErr := apperror.ValidationError("invalid request body")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	// Валидация обязательных полей
	if req.Name == "" {
		c.logger.Warn(ctx.Context(), "Organization name is required", logrus.Fields{})
		appErr := apperror.ValidationError("organization name is required")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	req.DBName = idParam
	if err := c.service.UpdateOrganization(ctx.Context(), &req); err != nil {
		appErr, ok := err.(*apperror.AppError)
		if !ok {
			appErr = apperror.New(apperror.ErrInternal, "failed to update organization").WithError(err)
		}
		c.logger.Error(ctx.Context(), "Failed to update organization", err, logrus.Fields{"id": id.String()})
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	c.logger.Info(ctx.Context(), "Organization updated successfully", logrus.Fields{"id": id.String()})
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Organization updated successfully",
	})
}

func (c *EsfOrganizationController) deleteEsfOrganization(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	c.logger.Info(ctx.Context(), "Deleting ESF organization", logrus.Fields{"id": idParam})

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.logger.Warn(ctx.Context(), "Invalid UUID format", logrus.Fields{"id": idParam})
		appErr := apperror.ValidationError("invalid UUID format")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	if err := c.service.DeleteOrganization(ctx.Context(), id); err != nil {
		appErr, ok := err.(*apperror.AppError)
		if !ok {
			appErr = apperror.New(apperror.ErrInternal, "failed to delete organization").WithError(err)
		}
		c.logger.Error(ctx.Context(), "Failed to delete organization", err, logrus.Fields{"id": id.String()})
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	c.logger.Info(ctx.Context(), "Organization deleted successfully", logrus.Fields{"id": id.String()})
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Organization deleted successfully",
	})
}
