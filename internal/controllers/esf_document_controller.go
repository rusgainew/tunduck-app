package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rusgainew/tunduck-app/internal/models"
	repositorypostgres "github.com/rusgainew/tunduck-app/internal/repository/repository_postgres"
	"github.com/rusgainew/tunduck-app/internal/services"
	serviceimpl "github.com/rusgainew/tunduck-app/internal/services/service_impl"
	"github.com/rusgainew/tunduck-app/pkg/apperror"
	"github.com/rusgainew/tunduck-app/pkg/logger"
	"github.com/rusgainew/tunduck-app/pkg/middleware"
	"github.com/rusgainew/tunduck-app/pkg/pagination"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EsfDocumentController struct {
	logger  *logger.Logger
	service services.EsfDocumentService
	db      *gorm.DB
}

func NewEsfDocumentController(app *fiber.App, log *logrus.Logger, db *gorm.DB) {
	// Инициализируем слои
	repo := repositorypostgres.NewEsfDocumentRepositoryPostgres(db, log)
	service := serviceimpl.NewEsfDocumentService(repo, db, log)

	l := logger.New(log)

	controller := &EsfDocumentController{
		logger:  l,
		service: service,
		db:      db,
	}

	l.Info(context.Background(), "EsfDocumentController initialized")
	controller.registerRoutes(app)
}

func (c *EsfDocumentController) registerRoutes(app *fiber.App) {
	esfDocumentGroup := app.Group("/api/esf-documents")

	// Публичные routes (без JWT)
	esfDocumentGroup.Get("/", c.getEsfDocuments)
	esfDocumentGroup.Get("/paginated", c.getEsfDocumentsPaginated)
	esfDocumentGroup.Get("/:id", c.getByEsfDocument)

	// Защищенные routes (с JWT)
	protected := esfDocumentGroup.Group("")
	protected.Use(middleware.JWTMiddleware())
	protected.Post("/", c.createEsfDocument)
	protected.Put("/:id", c.updateEsfDocument)
	protected.Delete("/:id", c.deleteEsfDocument)
}

// getEsfDocuments возвращает все документы ЭСФ
func (c *EsfDocumentController) getEsfDocuments(ctx *fiber.Ctx) error {
	c.logger.Info(ctx.Context(), "Fetching ESF documents")

	orgID, err := c.resolveOrgID(ctx)
	if err != nil {
		c.logger.Warn(ctx.Context(), "Failed to resolve org ID", logrus.Fields{"error": err.Error()})
		appErr := apperror.New(apperror.ErrInvalidRequest, "invalid organization ID")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	documents, err := c.service.GetAllDocuments(ctx.Context(), orgID)
	if err != nil {
		appErr, ok := err.(*apperror.AppError)
		if !ok {
			appErr = apperror.New(apperror.ErrInternal, "failed to fetch documents").WithError(err)
		}
		c.logger.Error(ctx.Context(), "Failed to fetch documents", err, logrus.Fields{"org_id": orgID.String()})
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	c.logger.Debug(ctx.Context(), "Documents fetched successfully", logrus.Fields{
		"org_id": orgID.String(),
		"count":  len(documents),
	})

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    documents,
		"count":   len(documents),
	})
}

// getEsfDocumentsPaginated возвращает документы ЭСФ с пагинацией
func (c *EsfDocumentController) getEsfDocumentsPaginated(ctx *fiber.Ctx) error {
	c.logger.Info(ctx.Context(), "Вибірка документів ЕСФ з пагінацією")

	orgID, err := c.resolveOrgID(ctx)
	if err != nil {
		c.logger.Warn(ctx.Context(), "Не вдалося визначити ID організації", logrus.Fields{"error": err.Error()})
		appErr := apperror.New(apperror.ErrInvalidRequest, "invalid organization ID")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	// Витягуємо параметри пагінації та фільтрації
	paginationParams := pagination.ExtractPaginationParams(ctx)
	filterParams := pagination.ExtractDocumentFilters(ctx)

	documents, totalCount, err := c.service.GetAllDocumentsPaginated(ctx.Context(), orgID, paginationParams, filterParams)
	if err != nil {
		appErr, ok := err.(*apperror.AppError)
		if !ok {
			appErr = apperror.New(apperror.ErrInternal, "failed to fetch documents").WithError(err)
		}
		c.logger.Error(ctx.Context(), "Ошибка вибірки документів", err, logrus.Fields{"org_id": orgID.String()})
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	// Формуємо відповідь з пагінацією
	response := pagination.NewPaginatedResponse(documents, paginationParams.Page, paginationParams.PageSize, totalCount)

	c.logger.Debug(ctx.Context(), "Документи успішно вибрані", logrus.Fields{
		"org_id": orgID.String(),
		"count":  len(documents),
		"total":  totalCount,
		"page":   paginationParams.Page,
	})

	return ctx.Status(http.StatusOK).JSON(response)
}

// getByEsfDocument возвращает документ ЭСФ по ID
func (c *EsfDocumentController) getByEsfDocument(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	c.logger.Debug(ctx.Context(), "Fetching document by ID", logrus.Fields{"doc_id": id})

	orgID, err := c.resolveOrgID(ctx)
	if err != nil {
		c.logger.Warn(ctx.Context(), "Failed to resolve org ID", logrus.Fields{"error": err.Error()})
		appErr := apperror.New(apperror.ErrInvalidRequest, "invalid organization ID")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	docID, err := uuid.Parse(id)
	if err != nil {
		c.logger.Warn(ctx.Context(), "Invalid UUID format", logrus.Fields{"id": id})
		appErr := apperror.New(apperror.ErrInvalidRequest, "invalid document ID format")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	document, err := c.service.GetDocumentByID(ctx.Context(), orgID, docID)
	if err != nil {
		appErr, ok := err.(*apperror.AppError)
		if !ok {
			appErr = apperror.New(apperror.ErrInternal, "failed to fetch document").WithError(err)
		}
		c.logger.Error(ctx.Context(), "Failed to fetch document", err, logrus.Fields{
			"org_id": orgID.String(),
			"doc_id": docID.String(),
		})
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	if document == nil {
		c.logger.Warn(ctx.Context(), "Document not found", logrus.Fields{
			"org_id": orgID.String(),
			"doc_id": docID.String(),
		})
		appErr := apperror.New(apperror.ErrDocumentNotFound, "document not found")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    document,
	})
}

// createEsfDocument создает новый документ ЭСФ
func (c *EsfDocumentController) createEsfDocument(ctx *fiber.Ctx) error {
	c.logger.Info(ctx.Context(), "Creating new ESF document")

	orgID, err := c.resolveOrgID(ctx)
	if err != nil {
		c.logger.Warn(ctx.Context(), "Failed to resolve org ID", logrus.Fields{"error": err.Error()})
		appErr := apperror.New(apperror.ErrInvalidRequest, "invalid organization ID")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	var req models.EsfCreateDocumentRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.logger.Warn(ctx.Context(), "Failed to parse request body", logrus.Fields{"error": err.Error()})
		appErr := apperror.New(apperror.ErrInvalidRequest, "invalid request format")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	// Валидируем запрос
	if err := middleware.ValidateStruct(&req); err != nil {
		c.logger.Warn(ctx.Context(), "Validation failed for create request", logrus.Fields{"error": err.Error()})
		appErr := apperror.New(apperror.ErrValidation, "validation error").WithDetails(err.Error())
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	response, err := c.service.CreateDocument(ctx.Context(), orgID, &req)
	if err != nil {
		appErr, ok := err.(*apperror.AppError)
		if !ok {
			appErr = apperror.New(apperror.ErrInternal, "failed to create document").WithError(err)
		}
		c.logger.Error(ctx.Context(), "Failed to create document", err, logrus.Fields{"org_id": orgID.String()})
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	c.logger.Info(ctx.Context(), "Document created successfully", logrus.Fields{"org_id": orgID.String()})
	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    response,
		"message": "Document created successfully",
	})
}

// updateEsfDocument обновляет документ ЭСФ
func (c *EsfDocumentController) updateEsfDocument(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	c.logger.Info(ctx.Context(), "Updating ESF document", logrus.Fields{"doc_id": id})

	orgID, err := c.resolveOrgID(ctx)
	if err != nil {
		c.logger.Warn(ctx.Context(), "Failed to resolve org ID", logrus.Fields{"error": err.Error()})
		appErr := apperror.New(apperror.ErrInvalidRequest, "invalid organization ID")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	docID, err := uuid.Parse(id)
	if err != nil {
		c.logger.Warn(ctx.Context(), "Invalid UUID format", logrus.Fields{"id": id})
		appErr := apperror.New(apperror.ErrInvalidRequest, "invalid document ID format")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	var req models.EsfEditDocumentRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.logger.Warn(ctx.Context(), "Failed to parse request body", logrus.Fields{"error": err.Error()})
		appErr := apperror.New(apperror.ErrInvalidRequest, "invalid request format")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	req.ID = docID

	// Валидируем запрос
	if err := middleware.ValidateStruct(&req); err != nil {
		c.logger.Warn(ctx.Context(), "Validation failed for update request", logrus.Fields{"error": err.Error()})
		appErr := apperror.New(apperror.ErrValidation, "validation error").WithDetails(err.Error())
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	if err := c.service.UpdateDocument(ctx.Context(), orgID, &req); err != nil {
		appErr, ok := err.(*apperror.AppError)
		if !ok {
			appErr = apperror.New(apperror.ErrInternal, "failed to update document").WithError(err)
		}
		c.logger.Error(ctx.Context(), "Failed to update document", err, logrus.Fields{"org_id": orgID.String(), "doc_id": docID.String()})
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	c.logger.Info(ctx.Context(), "Document updated successfully", logrus.Fields{"org_id": orgID.String(), "doc_id": docID.String()})
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Document updated successfully",
	})
}

// deleteEsfDocument удаляет документ ЭСФ
func (c *EsfDocumentController) deleteEsfDocument(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	c.logger.Info(ctx.Context(), "Deleting ESF document", logrus.Fields{"doc_id": id})

	orgID, err := c.resolveOrgID(ctx)
	if err != nil {
		c.logger.Warn(ctx.Context(), "Failed to resolve org ID", logrus.Fields{"error": err.Error()})
		appErr := apperror.New(apperror.ErrInvalidRequest, "invalid organization ID")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	docID, err := uuid.Parse(id)
	if err != nil {
		c.logger.Warn(ctx.Context(), "Invalid UUID format", logrus.Fields{"id": id})
		appErr := apperror.New(apperror.ErrInvalidRequest, "invalid document ID format")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	if err := c.service.DeleteDocument(ctx.Context(), orgID, docID); err != nil {
		appErr, ok := err.(*apperror.AppError)
		if !ok {
			appErr = apperror.New(apperror.ErrInternal, "failed to delete document").WithError(err)
		}
		c.logger.Error(ctx.Context(), "Failed to delete document", err, logrus.Fields{"org_id": orgID.String(), "doc_id": docID.String()})
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	c.logger.Info(ctx.Context(), "Document deleted successfully", logrus.Fields{"org_id": orgID.String(), "doc_id": docID.String()})
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Document deleted successfully",
	})
}

// resolveOrgID достает идентификатор организации из заголовка X-Org-Id или query orgId.
func (c *EsfDocumentController) resolveOrgID(ctx *fiber.Ctx) (uuid.UUID, error) {
	raw := ctx.Get("X-Org-Id")
	if raw == "" {
		raw = ctx.Query("orgId")
	}
	if raw == "" {
		return uuid.Nil, fmt.Errorf("organization id is required (header X-Org-Id or query orgId)")
	}
	orgID, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid organization id: %w", err)
	}
	return orgID, nil
}
