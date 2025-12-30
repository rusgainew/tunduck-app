package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/rusgainew/tunduck-app/internal/models"
	"github.com/rusgainew/tunduck-app/pkg/cache"
	"github.com/rusgainew/tunduck-app/pkg/pagination"
)

type EsfDocumentService interface {
	GetAllDocuments(ctx context.Context, orgID uuid.UUID) ([]models.EsfCreateDocumentRequest, error)
	GetDocumentByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.EsfCreateDocumentRequest, error)
	CreateDocument(ctx context.Context, orgID uuid.UUID, doc *models.EsfCreateDocumentRequest) (*models.EsfCreateDocumentResponse, error)
	UpdateDocument(ctx context.Context, orgID uuid.UUID, doc *models.EsfEditDocumentRequest) error
	DeleteDocument(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error

	// Пагіновані методи
	GetAllDocumentsPaginated(ctx context.Context, orgID uuid.UUID, params pagination.PaginationParams, filters pagination.DocumentFilterParams) ([]models.EsfCreateDocumentRequest, int64, error)

	// Cache management
	SetCacheManager(cache.CacheManager)
	CacheWarmDocuments(ctx context.Context, orgID uuid.UUID, limit int) error
}
