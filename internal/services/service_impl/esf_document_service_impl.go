package service_impl

import (
	"context"
	"time"

	"github.com/google/uuid"
	models "github.com/rusgainew/tunduck-app/internal/models"
	"github.com/rusgainew/tunduck-app/internal/repository"
	"github.com/rusgainew/tunduck-app/internal/services"
	"github.com/rusgainew/tunduck-app/pkg/apperror"
	"github.com/rusgainew/tunduck-app/pkg/cache"
	"github.com/rusgainew/tunduck-app/pkg/entity"
	"github.com/rusgainew/tunduck-app/pkg/logger"
	"github.com/rusgainew/tunduck-app/pkg/pagination"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type esfDocumentService struct {
	repo         repository.EsfDocumentRepository
	db           *gorm.DB
	logger       *logger.Logger
	cacheManager cache.CacheManager
}

// NewEsfDocumentService создает новый document service с обязательными зависимостями
func NewEsfDocumentService(repo repository.EsfDocumentRepository, db *gorm.DB, log *logrus.Logger) services.EsfDocumentService {
	if db == nil {
		log.Fatal("database connection is required for EsfDocumentService")
		return nil
	}
	return &esfDocumentService{
		repo:         repo,
		db:           db,
		logger:       logger.New(log),
		cacheManager: nil,
	}
}

// SetCacheManager injects the cache manager into the service
func (s *esfDocumentService) SetCacheManager(cacheManager cache.CacheManager) {
	s.cacheManager = cacheManager
}

func (s *esfDocumentService) GetAllDocuments(ctx context.Context, orgID uuid.UUID) ([]models.EsfCreateDocumentRequest, error) {
	s.logger.Info(ctx, "Fetching all documents", logrus.Fields{"org_id": orgID.String()})

	docs, err := s.repo.GetAllDocuments(ctx, orgID)
	if err != nil {
		s.logger.Error(ctx, "Failed to fetch documents", err, logrus.Fields{"org_id": orgID.String()})
		return nil, apperror.DatabaseError("fetching documents", err)
	}

	s.logger.Debug(ctx, "Documents fetched successfully", logrus.Fields{"org_id": orgID.String(), "count": len(docs)})

	result := make([]models.EsfCreateDocumentRequest, len(docs))
	for i := range docs {
		result[i] = s.toModel(&docs[i])
	}
	return result, nil
}

func (s *esfDocumentService) GetDocumentByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.EsfCreateDocumentRequest, error) {
	s.logger.Info(ctx, "Fetching document by ID", logrus.Fields{"org_id": orgID.String(), "doc_id": id.String()})

	// Try to get from cache first
	if s.cacheManager != nil {
		cacheKey := "doc:id:" + id.String()
		if cached, _ := s.cacheManager.Document().Get(ctx, cacheKey); cached != nil {
			s.logger.Debug(ctx, "Document found in cache", logrus.Fields{"doc_id": id.String()})
			doc := cached.(*models.EsfCreateDocumentRequest)
			return doc, nil
		}
	}

	doc, err := s.repo.GetDocumentByID(ctx, orgID, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to fetch document", err, logrus.Fields{"org_id": orgID.String(), "doc_id": id.String()})
		return nil, apperror.DatabaseError("fetching document", err)
	}

	if doc == nil {
		s.logger.Warn(ctx, "Document not found", logrus.Fields{"org_id": orgID.String(), "doc_id": id.String()})
		return nil, apperror.New(apperror.ErrDocumentNotFound, "document not found")
	}

	s.logger.Debug(ctx, "Document fetched successfully", logrus.Fields{"org_id": orgID.String(), "doc_id": id.String()})

	model := s.toModel(doc)

	// Cache the document (30 minutes TTL)
	if s.cacheManager != nil {
		cacheKey := "doc:id:" + id.String()
		_ = s.cacheManager.Document().Set(ctx, cacheKey, &model, 30*time.Minute)
	}

	return &model, nil
}

func (s *esfDocumentService) CreateDocument(ctx context.Context, orgID uuid.UUID, req *models.EsfCreateDocumentRequest) (*models.EsfCreateDocumentResponse, error) {
	s.logger.Info(ctx, "Creating new document", logrus.Fields{"org_id": orgID.String()})

	doc := s.toEntity(req)
	doc.ID = uuid.New()

	if err := s.repo.CreateDocument(ctx, orgID, &doc); err != nil {
		s.logger.Error(ctx, "Failed to create document", err, logrus.Fields{"org_id": orgID.String(), "doc_id": doc.ID.String()})
		return nil, apperror.DatabaseError("creating document", err)
	}

	s.logger.Info(ctx, "Document created successfully", logrus.Fields{"org_id": orgID.String(), "doc_id": doc.ID.String()})

	return &models.EsfCreateDocumentResponse{
		ResponseId:   "success",
		DocumentUuid: doc.ID.String(),
	}, nil
}

func (s *esfDocumentService) UpdateDocument(ctx context.Context, orgID uuid.UUID, req *models.EsfEditDocumentRequest) error {
	s.logger.Info(ctx, "Updating document", logrus.Fields{"org_id": orgID.String(), "doc_id": req.ID.String()})

	doc := s.toEntity(&req.EsfCreateDocumentRequest)
	doc.ID = req.ID

	if err := s.repo.UpdateDocument(ctx, orgID, &doc); err != nil {
		s.logger.Error(ctx, "Failed to update document", err, logrus.Fields{"org_id": orgID.String(), "doc_id": req.ID.String()})
		return apperror.DatabaseError("updating document", err)
	}

	// Invalidate cache
	if s.cacheManager != nil {
		cacheKey := "doc:id:" + req.ID.String()
		_ = s.cacheManager.Document().Delete(ctx, cacheKey)
	}

	s.logger.Info(ctx, "Document updated successfully", logrus.Fields{"org_id": orgID.String(), "doc_id": req.ID.String()})
	return nil
}

func (s *esfDocumentService) DeleteDocument(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info(ctx, "Deleting document", logrus.Fields{"org_id": orgID.String(), "doc_id": id.String()})

	if err := s.repo.DeleteDocument(ctx, orgID, id); err != nil {
		s.logger.Error(ctx, "Failed to delete document", err, logrus.Fields{"org_id": orgID.String(), "doc_id": id.String()})
		return apperror.DatabaseError("deleting document", err)
	}

	// Invalidate cache
	if s.cacheManager != nil {
		cacheKey := "doc:id:" + id.String()
		_ = s.cacheManager.Document().Delete(ctx, cacheKey)
	}

	s.logger.Info(ctx, "Document deleted successfully", logrus.Fields{"org_id": orgID.String(), "doc_id": id.String()})
	return nil
}

// CacheWarmDocuments preloads frequently accessed documents
func (s *esfDocumentService) CacheWarmDocuments(ctx context.Context, orgID uuid.UUID, limit int) error {
	if s.cacheManager == nil {
		return nil
	}

	s.logger.Info(ctx, "Starting document cache warming", logrus.Fields{"org_id": orgID.String(), "limit": limit})

	docs, err := s.repo.GetAllDocuments(ctx, orgID)
	if err != nil {
		s.logger.Error(ctx, "Failed to fetch documents for cache warming", err, logrus.Fields{"org_id": orgID.String()})
		return nil // Non-critical operation, don't fail
	}

	// Limit documents to cache
	if limit > 0 && len(docs) > limit {
		docs = docs[:limit]
	}

	// Batch set documents in cache
	batchData := make(map[string]interface{})
	for _, doc := range docs {
		model := s.toModel(&doc)
		cacheKey := "doc:id:" + doc.ID.String()
		batchData[cacheKey] = &model
	}

	if len(batchData) > 0 {
		_ = s.cacheManager.Document().SetMultiple(ctx, batchData, 30*time.Minute)
		s.logger.Info(ctx, "Document cache warming completed", logrus.Fields{
			"org_id": orgID.String(),
			"count":  len(batchData),
		})
	}

	return nil
}

func (s *esfDocumentService) toEntity(m *models.EsfCreateDocumentRequest) entity.EsfDocument {
	// CatalogEntries это models.CatalogEntriesModels которая является []EsfEntriesModel
	entries := make([]entity.EsfEntries, len(m.CatalogEntries))
	for i := 0; i < len(m.CatalogEntries); i++ {
		catalogEntry := m.CatalogEntries[i] // EsfEntriesModel
		entries[i] = entity.EsfEntries{
			UnitClassificationCode: catalogEntry.UnitClassificationCode,
			SalesTaxCode:           catalogEntry.SalesTaxCode,
			CustomsAuthorityCode:   catalogEntry.CustomsAuthorityCode,
			Quantity:               catalogEntry.Quantity,
			Price:                  catalogEntry.Price,
			VatAmount:              catalogEntry.VatAmount,
			SalesTaxAmount:         catalogEntry.SalesTaxAmount,
			AmountWithoutTaxes:     catalogEntry.AmountWithoutTaxes,
			TotalAmount:            catalogEntry.TotalAmount,
		}
	}

	return entity.EsfDocument{
		ForeignName:                    m.ForeignName,
		IsBranchDataSent:               m.IsBranchDataSent,
		IsPriceWithoutTaxes:            m.IsPriceWithoutTaxes,
		AffiliateTin:                   m.AffiliateTin,
		IsIndustry:                     m.IsIndustry,
		OwnedCrmReceiptCode:            m.OwnedCrmReceiptCode,
		OperationTypeCode:              m.OperationTypeCode,
		DeliveryDate:                   m.DeliveryDate,
		DeliveryTypeCode:               m.DeliveryTypeCode,
		IsResident:                     m.IsResident,
		ContractorTin:                  m.ContractorTin,
		SupplierBankAccount:            m.SupplierBankAccount,
		ContractorBankAccount:          m.ContractorBankAccount,
		CurrencyCode:                   m.CurrencyCode,
		CountryCode:                    m.CountryCode,
		CurrencyRate:                   m.CurrencyRate,
		TotalCurrencyValue:             m.TotalCurrencyValue,
		TotalCurrencyValueWithoutTaxes: m.TotalCurrencyValueWithoutTaxes,
		SupplyContractNumber:           m.SupplyContractNumber,
		ContractStartDate:              m.ContractStartDate,
		Comment:                        m.Comment,
		DeliveryCode:                   m.DeliveryCode,
		PaymentCode:                    m.PaymentCode,
		TaxRateVATCode:                 m.TaxRateVATCode,
		CatalogEntries:                 entries,
		OpeningBalances:                m.OpeningBalances,
		AssessedContributionsAmount:    m.AssessedContributionsAmount,
		PaidAmount:                     m.PaidAmount,
		PenaltiesAmount:                m.PenaltiesAmount,
		FinesAmount:                    m.FinesAmount,
		ClosingBalances:                m.ClosingBalances,
		AmountToBePaid:                 m.AmountToBePaid,
		PersonalAccountNumber:          m.PersonalAccountNumber,
	}
}

func (s *esfDocumentService) toModel(e *entity.EsfDocument) models.EsfCreateDocumentRequest {
	entries := make([]models.EsfEntriesModel, len(e.CatalogEntries))
	for i, ent := range e.CatalogEntries {
		entries[i] = models.EsfEntriesModel{
			UnitClassificationCode: ent.UnitClassificationCode,
			SalesTaxCode:           ent.SalesTaxCode,
			CustomsAuthorityCode:   ent.CustomsAuthorityCode,
			Quantity:               ent.Quantity,
			Price:                  ent.Price,
			VatAmount:              ent.VatAmount,
			SalesTaxAmount:         ent.SalesTaxAmount,
			AmountWithoutTaxes:     ent.AmountWithoutTaxes,
			TotalAmount:            ent.TotalAmount,
		}
	}

	return models.EsfCreateDocumentRequest{
		ForeignName:                    e.ForeignName,
		IsBranchDataSent:               e.IsBranchDataSent,
		IsPriceWithoutTaxes:            e.IsPriceWithoutTaxes,
		AffiliateTin:                   e.AffiliateTin,
		IsIndustry:                     e.IsIndustry,
		OwnedCrmReceiptCode:            e.OwnedCrmReceiptCode,
		OperationTypeCode:              e.OperationTypeCode,
		DeliveryDate:                   e.DeliveryDate,
		DeliveryTypeCode:               e.DeliveryTypeCode,
		IsResident:                     e.IsResident,
		ContractorTin:                  e.ContractorTin,
		SupplierBankAccount:            e.SupplierBankAccount,
		ContractorBankAccount:          e.ContractorBankAccount,
		CurrencyCode:                   e.CurrencyCode,
		CountryCode:                    e.CountryCode,
		CurrencyRate:                   e.CurrencyRate,
		TotalCurrencyValue:             e.TotalCurrencyValue,
		TotalCurrencyValueWithoutTaxes: e.TotalCurrencyValueWithoutTaxes,
		SupplyContractNumber:           e.SupplyContractNumber,
		ContractStartDate:              e.ContractStartDate,
		Comment:                        e.Comment,
		DeliveryCode:                   e.DeliveryCode,
		PaymentCode:                    e.PaymentCode,
		TaxRateVATCode:                 e.TaxRateVATCode,
		CatalogEntries:                 entries,
		OpeningBalances:                e.OpeningBalances,
		AssessedContributionsAmount:    e.AssessedContributionsAmount,
		PaidAmount:                     e.PaidAmount,
		PenaltiesAmount:                e.PenaltiesAmount,
		FinesAmount:                    e.FinesAmount,
		ClosingBalances:                e.ClosingBalances,
		AmountToBePaid:                 e.AmountToBePaid,
		PersonalAccountNumber:          e.PersonalAccountNumber,
	}
}

// GetAllDocumentsPaginated возвращает документы с пагинацией
func (s *esfDocumentService) GetAllDocumentsPaginated(ctx context.Context, orgID uuid.UUID, params pagination.PaginationParams, filters pagination.DocumentFilterParams) ([]models.EsfCreateDocumentRequest, int64, error) {
	s.logger.Info(ctx, "Fetching documents with pagination", logrus.Fields{
		"org_id":    orgID.String(),
		"page":      params.Page,
		"page_size": params.PageSize,
	})

	docs, totalCount, err := s.repo.GetAllDocumentsPaginated(ctx, orgID, params, filters)
	if err != nil {
		s.logger.Error(ctx, "Failed to fetch paginated documents", err, logrus.Fields{"org_id": orgID.String()})
		return nil, 0, err
	}

	result := make([]models.EsfCreateDocumentRequest, len(docs))
	for i := range docs {
		result[i] = s.toModel(&docs[i])
	}

	s.logger.Debug(ctx, "Paginated documents fetched successfully", logrus.Fields{
		"org_id": orgID.String(),
		"count":  len(result),
		"total":  totalCount,
		"page":   params.Page,
	})

	return result, totalCount, nil
}
