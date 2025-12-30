package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/rusgainew/tunduck-app/internal/models"
	"github.com/rusgainew/tunduck-app/pkg/cache"
	"github.com/rusgainew/tunduck-app/pkg/pagination"
)

type EsfOrganizationService interface {
	GetAllOrganizations(ctx context.Context) ([]models.EsfOrganizationModel, error)
	GetOrganizationByID(ctx context.Context, id uuid.UUID) (*models.EsfOrganizationModel, error)
	CreateOrganization(ctx context.Context, org *models.EsfOrganizationModel) (uuid.UUID, string, error)
	UpdateOrganization(ctx context.Context, org *models.EsfOrganizationModel) error
	DeleteOrganization(ctx context.Context, id uuid.UUID) error

	// Пагіновані методи
	GetAllOrganizationsPaginated(ctx context.Context, params pagination.PaginationParams, filters pagination.OrganizationFilterParams) ([]models.EsfOrganizationModel, int64, error)

	// Кеширование
	CacheWarmOrganizations(ctx context.Context) error
	SetCacheManager(cacheManager cache.CacheManager)
}
