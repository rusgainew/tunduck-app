package repositorypostgres

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rusgainew/tunduck-app/internal/repository"
	"github.com/rusgainew/tunduck-app/pkg/apperror"
	"github.com/rusgainew/tunduck-app/pkg/entity"
	"github.com/rusgainew/tunduck-app/pkg/logger"
	"github.com/rusgainew/tunduck-app/pkg/pagination"
)

type esfDocumentRepositoryPostgres struct {
	logger  *logger.Logger
	baseDB  *gorm.DB
	dbCache map[string]*gorm.DB
	cacheMu sync.RWMutex
}

func NewEsfDocumentRepositoryPostgres(db *gorm.DB, log *logrus.Logger) repository.EsfDocumentRepository {
	return &esfDocumentRepositoryPostgres{
		baseDB:  db,
		logger:  logger.New(log),
		dbCache: make(map[string]*gorm.DB),
	}
}

// GetAllDocuments возвращает все документы ЭСФ
func (edrp *esfDocumentRepositoryPostgres) GetAllDocuments(ctx context.Context, orgID uuid.UUID) ([]entity.EsfDocument, error) {
	edrp.logger.Debug(ctx, "Fetching all documents from organization database", logrus.Fields{"org_id": orgID.String()})

	var documents []entity.EsfDocument

	orgDB, err := edrp.getOrgDB(ctx, orgID)
	if err != nil {
		edrp.logger.Error(ctx, "Failed to get organization database", err, logrus.Fields{"org_id": orgID.String()})
		return nil, apperror.DatabaseError("getting organization database", err)
	}

	err = orgDB.WithContext(ctx).
		Preload("CatalogEntries").
		Find(&documents).Error

	if err != nil {
		edrp.logger.Error(ctx, "Failed to fetch documents from database", err, logrus.Fields{"org_id": orgID.String()})
		return nil, apperror.DatabaseError("fetching documents", err)
	}

	edrp.logger.Debug(ctx, "Documents fetched successfully", logrus.Fields{"org_id": orgID.String(), "count": len(documents)})
	return documents, nil
}

// GetDocumentByID возвращает документ ЭСФ по ID
func (edrp *esfDocumentRepositoryPostgres) GetDocumentByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*entity.EsfDocument, error) {
	edrp.logger.Debug(ctx, "Fetching document by ID", logrus.Fields{"org_id": orgID.String(), "doc_id": id.String()})

	var document entity.EsfDocument

	orgDB, err := edrp.getOrgDB(ctx, orgID)
	if err != nil {
		edrp.logger.Error(ctx, "Failed to get organization database", err, logrus.Fields{"org_id": orgID.String()})
		return nil, apperror.DatabaseError("getting organization database", err)
	}

	err = orgDB.WithContext(ctx).
		Preload("CatalogEntries").
		Where("id = ?", id).
		First(&document).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			edrp.logger.Debug(ctx, "Document not found", logrus.Fields{"org_id": orgID.String(), "doc_id": id.String()})
			return nil, apperror.New(apperror.ErrDocumentNotFound, "document not found")
		}
		edrp.logger.Error(ctx, "Failed to fetch document by ID", err, logrus.Fields{"org_id": orgID.String(), "doc_id": id.String()})
		return nil, apperror.DatabaseError("fetching document by ID", err)
	}

	edrp.logger.Debug(ctx, "Document fetched successfully", logrus.Fields{"org_id": orgID.String(), "doc_id": id.String()})
	return &document, nil
}

// CreateDocument создает новый документ ЭСФ
func (edrp *esfDocumentRepositoryPostgres) CreateDocument(ctx context.Context, orgID uuid.UUID, doc *entity.EsfDocument) error {
	edrp.logger.Debug(ctx, "Creating document in organization database", logrus.Fields{"org_id": orgID.String(), "doc_id": doc.ID.String()})

	orgDB, err := edrp.getOrgDB(ctx, orgID)
	if err != nil {
		edrp.logger.Error(ctx, "Failed to get organization database", err, logrus.Fields{"org_id": orgID.String()})
		return apperror.DatabaseError("getting organization database", err)
	}

	err = orgDB.WithContext(ctx).Create(doc).Error

	if err != nil {
		edrp.logger.Error(ctx, "Failed to create document in database", err, logrus.Fields{"org_id": orgID.String(), "doc_id": doc.ID.String()})
		return apperror.DatabaseError("creating document", err)
	}

	edrp.logger.Debug(ctx, "Document created successfully", logrus.Fields{"org_id": orgID.String(), "doc_id": doc.ID.String()})
	return nil
}

// UpdateDocument обновляет существующий документ ЭСФ
func (edrp *esfDocumentRepositoryPostgres) UpdateDocument(ctx context.Context, orgID uuid.UUID, doc *entity.EsfDocument) error {
	edrp.logger.Debug(ctx, "Updating document in organization database", logrus.Fields{"org_id": orgID.String(), "doc_id": doc.ID.String()})

	orgDB, err := edrp.getOrgDB(ctx, orgID)
	if err != nil {
		edrp.logger.Error(ctx, "Failed to get organization database", err, logrus.Fields{"org_id": orgID.String()})
		return apperror.DatabaseError("getting organization database", err)
	}

	err = orgDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Обновляем основной документ
		if err := tx.Model(&entity.EsfDocument{}).
			Where("id = ?", doc.ID).
			Updates(doc).Error; err != nil {
			edrp.logger.Error(ctx, "Failed to update document model", err, logrus.Fields{"doc_id": doc.ID.String()})
			return err
		}

		// Удаляем старые записи CatalogEntries
		if err := tx.Where("document_id = ?", doc.ID).
			Delete(&entity.EsfEntries{}).Error; err != nil {
			edrp.logger.Error(ctx, "Failed to delete old catalog entries", err, logrus.Fields{"doc_id": doc.ID.String()})
			return err
		}

		// Создаем новые записи CatalogEntries
		if len(doc.CatalogEntries) > 0 {
			if err := tx.Create(&doc.CatalogEntries).Error; err != nil {
				edrp.logger.Error(ctx, "Failed to create new catalog entries", err, logrus.Fields{"doc_id": doc.ID.String()})
				return err
			}
		}

		return nil
	})

	if err != nil {
		edrp.logger.Error(ctx, "Failed to update document in database (transaction failed)", err, logrus.Fields{"org_id": orgID.String(), "doc_id": doc.ID.String()})
		return apperror.DatabaseError("updating document", err)
	}

	edrp.logger.Debug(ctx, "Document updated successfully", logrus.Fields{"org_id": orgID.String(), "doc_id": doc.ID.String()})
	return nil
}

// DeleteDocument удаляет документ ЭСФ (soft delete)
func (edrp *esfDocumentRepositoryPostgres) DeleteDocument(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	edrp.logger.Debug(ctx, "Deleting document from organization database", logrus.Fields{"org_id": orgID.String(), "doc_id": id.String()})

	orgDB, err := edrp.getOrgDB(ctx, orgID)
	if err != nil {
		edrp.logger.Error(ctx, "Failed to get organization database", err, logrus.Fields{"org_id": orgID.String()})
		return apperror.DatabaseError("getting organization database", err)
	}

	err = orgDB.WithContext(ctx).Delete(&entity.EsfDocument{}, id).Error

	if err != nil {
		edrp.logger.Error(ctx, "Failed to delete document from database", err, logrus.Fields{"org_id": orgID.String(), "doc_id": id.String()})
		return apperror.DatabaseError("deleting document", err)
	}

	edrp.logger.Debug(ctx, "Document deleted successfully", logrus.Fields{"org_id": orgID.String(), "doc_id": id.String()})
	return nil
}

// getOrgDB возвращает подключение к БД организации по ее ID, кэшируя соединения.
func (edrp *esfDocumentRepositoryPostgres) getOrgDB(ctx context.Context, orgID uuid.UUID) (*gorm.DB, error) {
	if orgID == uuid.Nil {
		return nil, fmt.Errorf("organization id is required")
	}

	edrp.cacheMu.RLock()
	if cached, ok := edrp.dbCache[orgID.String()]; ok {
		edrp.cacheMu.RUnlock()
		return cached, nil
	}
	edrp.cacheMu.RUnlock()

	var org entity.EstOrganization
	if err := edrp.baseDB.WithContext(ctx).Select("db_name").Where("id = ?", orgID).First(&org).Error; err != nil {
		edrp.logger.Error(ctx, "Failed to fetch organization database name", err, logrus.Fields{"orgID": orgID.String()})
		return nil, apperror.DatabaseError("fetching organization database name", err)
	}

	if org.DBName == "" {
		edrp.logger.Error(ctx, "Organization has empty database name", nil, logrus.Fields{"orgID": orgID.String()})
		return nil, fmt.Errorf("organization %s has empty database name", orgID)
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, org.DBName, port, sslmode)
	orgDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		edrp.logger.Error(ctx, "Failed to connect to organization database", err, logrus.Fields{"dbName": org.DBName})
		return nil, apperror.DatabaseError("connecting to organization database", err)
	}

	edrp.cacheMu.Lock()
	edrp.dbCache[orgID.String()] = orgDB
	edrp.cacheMu.Unlock()

	return orgDB, nil
}

// GetAllDocumentsPaginated возвращает документы ЭСФ с пагинацией и фильтрацией
func (edrp *esfDocumentRepositoryPostgres) GetAllDocumentsPaginated(ctx context.Context, orgID uuid.UUID, params pagination.PaginationParams, filters pagination.DocumentFilterParams) ([]entity.EsfDocument, int64, error) {
	edrp.logger.Debug(ctx, "Fetching documents with pagination", logrus.Fields{
		"org_id":      orgID.String(),
		"page":        params.Page,
		"page_size":   params.PageSize,
		"sort":        params.Sort,
		"order":       params.Order,
		"has_filters": filters.HasFilters(),
	})

	var documents []entity.EsfDocument
	var totalCount int64

	orgDB, err := edrp.getOrgDB(ctx, orgID)
	if err != nil {
		edrp.logger.Error(ctx, "Failed to get organization database", err, logrus.Fields{"org_id": orgID.String()})
		return nil, 0, apperror.DatabaseError("getting organization database", err)
	}

	query := orgDB.WithContext(ctx)

	// Применяем фильтры
	if filters.Status != "" {
		edrp.logger.Debug(ctx, "Applying status filter", logrus.Fields{"status": filters.Status})
		query = query.Where("status = ?", filters.Status)
	}

	if filters.Search != "" {
		edrp.logger.Debug(ctx, "Applying search filter", logrus.Fields{"search": filters.Search})
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}

	if filters.CreatedAfter != "" {
		edrp.logger.Debug(ctx, "Applying created_after filter", logrus.Fields{"created_after": filters.CreatedAfter})
		query = query.Where("created_at >= ?", filters.CreatedAfter)
	}

	if filters.CreatedBefore != "" {
		edrp.logger.Debug(ctx, "Applying created_before filter", logrus.Fields{"created_before": filters.CreatedBefore})
		query = query.Where("created_at <= ?", filters.CreatedBefore)
	}

	// Получаем общее количество
	if err := query.Model(&entity.EsfDocument{}).Count(&totalCount).Error; err != nil {
		edrp.logger.Error(ctx, "Failed to count documents", err, logrus.Fields{"org_id": orgID.String()})
		return nil, 0, apperror.DatabaseError("counting documents", err)
	}

	// Применяем сортировку и пагинацию
	if err := query.
		Preload("CatalogEntries").
		Order(params.Sort + " " + params.Order).
		Offset(params.GetOffset()).
		Limit(params.GetLimit()).
		Find(&documents).Error; err != nil {
		edrp.logger.Error(ctx, "Failed to fetch paginated documents", err, logrus.Fields{
			"org_id":    orgID.String(),
			"page":      params.Page,
			"page_size": params.PageSize,
		})
		return nil, 0, apperror.DatabaseError("fetching paginated documents", err)
	}

	edrp.logger.Debug(ctx, "Documents fetched successfully", logrus.Fields{
		"org_id": orgID.String(),
		"count":  len(documents),
		"total":  totalCount,
		"page":   params.Page,
	})

	return documents, totalCount, nil
}
