package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rusgainew/tunduck-app/pkg/entity"
	"github.com/rusgainew/tunduck-app/pkg/pagination"
)

type EsfDocumentRepository interface {
	GetAllDocuments(ctx context.Context, orgID uuid.UUID) ([]entity.EsfDocument, error)
	GetDocumentByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*entity.EsfDocument, error)
	CreateDocument(ctx context.Context, orgID uuid.UUID, doc *entity.EsfDocument) error
	UpdateDocument(ctx context.Context, orgID uuid.UUID, doc *entity.EsfDocument) error
	DeleteDocument(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error

	// Пагіновані методи
	GetAllDocumentsPaginated(ctx context.Context, orgID uuid.UUID, params pagination.PaginationParams, filters pagination.DocumentFilterParams) ([]entity.EsfDocument, int64, error)
}
