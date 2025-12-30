package repository

import (
	"context"

	"github.com/rusgainew/tunduck-app/pkg/entity"
	"github.com/rusgainew/tunduck-app/pkg/pagination"
)

type EsfOrganizationRepository interface {
	// Define methods for EsfOrganizationRepository here
	GetAll(ctx context.Context) ([]*entity.EstOrganization, error)
	GetByID(ctx context.Context, id string) (*entity.EstOrganization, error)
	Insert(ctx context.Context, org *entity.EstOrganization) error
	Update(ctx context.Context, org *entity.EstOrganization) error
	Delete(ctx context.Context, id string) error
	CreateDatabase(ctx context.Context, dbName string) error

	// Пагіновані методи
	GetAllPaginated(ctx context.Context, params pagination.PaginationParams, filters pagination.OrganizationFilterParams) ([]*entity.EstOrganization, int64, error)
}
