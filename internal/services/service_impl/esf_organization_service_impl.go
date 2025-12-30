package service_impl

import (
	"context"
	"regexp"
	"strings"
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
)

// esfOrganizationServiceImpl реализует интерфейс EsfOrganizationService
type esfOrganizationServiceImpl struct {
	repo         repository.EsfOrganizationRepository
	logger       *logger.Logger
	cacheManager cache.CacheManager
	cacheHelper  *cache.CacheHelper
}

// NewEsfOrganizationService создает новый экземпляр сервиса организаций
func NewEsfOrganizationService(repo repository.EsfOrganizationRepository, log *logrus.Logger) services.EsfOrganizationService {
	return &esfOrganizationServiceImpl{
		repo:         repo,
		logger:       logger.New(log),
		cacheManager: nil,
		cacheHelper:  nil,
	}
}

// SetCacheManager устанавливает CacheManager для использования кеша
func (s *esfOrganizationServiceImpl) SetCacheManager(cacheManager cache.CacheManager) {
	s.cacheManager = cacheManager
	s.cacheHelper = cache.NewCacheHelper(cacheManager)
}

// GetAllOrganizations возвращает все организации
func (s *esfOrganizationServiceImpl) GetAllOrganizations(ctx context.Context) ([]models.EsfOrganizationModel, error) {
	s.logger.Info(ctx, "Fetching all organizations", logrus.Fields{})

	orgs, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error(ctx, "Failed to fetch organizations", err, logrus.Fields{})
		return nil, apperror.DatabaseError("fetching organizations", err)
	}

	result := make([]models.EsfOrganizationModel, len(orgs))
	for i, org := range orgs {
		result[i] = models.EsfOrganizationModel{
			ID:          org.ID.String(),
			Name:        org.Name,
			Description: org.Description,
			Token:       org.Token,
			DBName:      org.DBName,
		}
	}

	s.logger.Debug(ctx, "Organizations fetched successfully", logrus.Fields{"count": len(result)})
	return result, nil
}

// GetOrganizationByID возвращает организацию по ID
func (s *esfOrganizationServiceImpl) GetOrganizationByID(ctx context.Context, id uuid.UUID) (*models.EsfOrganizationModel, error) {
	s.logger.Info(ctx, "Fetching organization by ID", logrus.Fields{"org_id": id.String()})

	// Проверяем кеш
	if s.cacheManager != nil {
		cacheKey := "id:" + id.String()
		cached, _ := s.cacheManager.Organization().Get(ctx, cacheKey)
		if cached != nil {
			if org, ok := cached.(*entity.EstOrganization); ok {
				s.logger.Debug(ctx, "Organization retrieved from cache", logrus.Fields{"org_id": id.String()})
				return &models.EsfOrganizationModel{
					ID:          org.ID.String(),
					Name:        org.Name,
					Description: org.Description,
					Token:       org.Token,
					DBName:      org.DBName,
				}, nil
			}
		}
	}

	org, err := s.repo.GetByID(ctx, id.String())
	if err != nil {
		s.logger.Error(ctx, "Failed to fetch organization", err, logrus.Fields{"org_id": id.String()})
		return nil, apperror.DatabaseError("fetching organization", err)
	}

	if org == nil {
		s.logger.Warn(ctx, "Organization not found", logrus.Fields{"org_id": id.String()})
		return nil, apperror.New(apperror.ErrOrgNotFound, "organization not found")
	}

	// Кешируем результат на 2 часа
	if s.cacheManager != nil {
		cacheKey := "id:" + org.ID.String()
		_ = s.cacheManager.Organization().Set(ctx, cacheKey, org, 2*time.Hour)
	}

	result := &models.EsfOrganizationModel{
		ID:          org.ID.String(),
		Name:        org.Name,
		Description: org.Description,
		Token:       org.Token,
		DBName:      org.DBName,
	}

	s.logger.Debug(ctx, "Organization fetched successfully", logrus.Fields{"org_id": id.String()})
	return result, nil
}

// CreateOrganization создает новую организацию и отдельную базу данных для неё
func (s *esfOrganizationServiceImpl) CreateOrganization(ctx context.Context, org *models.EsfOrganizationModel) (uuid.UUID, string, error) {
	s.logger.Info(ctx, "Creating new organization", logrus.Fields{"name": org.Name})

	// Валидация
	if org.Name == "" {
		s.logger.Warn(ctx, "Organization name is required", logrus.Fields{})
		return uuid.Nil, "", apperror.ValidationError("organization name is required")
	}

	// Формируем имя базы данных
	dbName := sanitizeDatabaseName(org.Name) + "_db"
	s.logger.Debug(ctx, "Database name generated", logrus.Fields{"dbName": dbName})

	// Создаем базу данных для организации
	if err := s.repo.CreateDatabase(ctx, dbName); err != nil {
		s.logger.Error(ctx, "Failed to create database", err, logrus.Fields{"dbName": dbName})
		return uuid.Nil, "", apperror.DatabaseError("creating database", err)
	}

	// Создаем организацию
	entity := &entity.EstOrganization{
		ID:          uuid.New(),
		Name:        org.Name,
		Description: org.Description,
		Token:       org.Token,
		DBName:      dbName,
	}

	if err := s.repo.Insert(ctx, entity); err != nil {
		s.logger.Error(ctx, "Failed to create organization", err, logrus.Fields{"name": org.Name})
		return uuid.Nil, "", apperror.DatabaseError("creating organization", err)
	}

	s.logger.Info(ctx, "Organization created successfully", logrus.Fields{"id": entity.ID.String(), "dbName": dbName})
	return entity.ID, dbName, nil
}

// sanitizeDatabaseName очищает имя от недопустимых символов для имени БД
// Оставляет только буквы, цифры и подчеркивания, приводит к нижнему регистру
func sanitizeDatabaseName(name string) string {
	// Приводим к нижнему регистру
	name = strings.ToLower(name)

	// Заменяем пробелы и дефисы на подчеркивания
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")

	// Удаляем все символы кроме букв, цифр и подчеркиваний
	reg := regexp.MustCompile(`[^a-z0-9_]+`)
	name = reg.ReplaceAllString(name, "")

	// Если имя начинается с цифры, добавляем префикс
	if len(name) > 0 && name[0] >= '0' && name[0] <= '9' {
		name = "db_" + name
	}

	return name
}

// UpdateOrganization обновляет данные организации
func (s *esfOrganizationServiceImpl) UpdateOrganization(ctx context.Context, org *models.EsfOrganizationModel) error {
	s.logger.Info(ctx, "Updating organization", logrus.Fields{"name": org.Name})

	// Валидация
	if org.Name == "" {
		s.logger.Warn(ctx, "Organization name is required", logrus.Fields{})
		return apperror.ValidationError("organization name is required")
	}

	// Конвертируем model в entity для обновления
	entity := &entity.EstOrganization{
		Name:        org.Name,
		Description: org.Description,
		Token:       org.Token,
		DBName:      org.DBName,
	}

	// Обновляем в репозитории
	if err := s.repo.Update(ctx, entity); err != nil {
		s.logger.Error(ctx, "Failed to update organization", err, logrus.Fields{"name": org.Name})
		return apperror.DatabaseError("updating organization", err)
	}

	s.logger.Info(ctx, "Organization updated successfully", logrus.Fields{"name": org.Name})
	return nil
}

// DeleteOrganization удаляет организацию по ID
func (s *esfOrganizationServiceImpl) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	s.logger.Info(ctx, "Deleting organization", logrus.Fields{"id": id.String()})

	if err := s.repo.Delete(ctx, id.String()); err != nil {
		s.logger.Error(ctx, "Failed to delete organization", err, logrus.Fields{"id": id.String()})
		return apperror.DatabaseError("deleting organization", err)
	}

	s.logger.Info(ctx, "Organization deleted successfully", logrus.Fields{"id": id.String()})
	return nil
}

// GetAllOrganizationsPaginated возвращает организации с пагинацией
func (s *esfOrganizationServiceImpl) GetAllOrganizationsPaginated(ctx context.Context, params pagination.PaginationParams, filters pagination.OrganizationFilterParams) ([]models.EsfOrganizationModel, int64, error) {
	s.logger.Info(ctx, "Fetching organizations with pagination", logrus.Fields{
		"page":      params.Page,
		"page_size": params.PageSize,
	})

	orgs, totalCount, err := s.repo.GetAllPaginated(ctx, params, filters)
	if err != nil {
		s.logger.Error(ctx, "Failed to fetch paginated organizations", err, logrus.Fields{})
		return nil, 0, err
	}

	result := make([]models.EsfOrganizationModel, len(orgs))
	for i, org := range orgs {
		result[i] = models.EsfOrganizationModel{
			ID:          org.ID.String(),
			Name:        org.Name,
			Description: org.Description,
			Token:       org.Token,
			DBName:      org.DBName,
		}
	}

	s.logger.Debug(ctx, "Paginated organizations fetched successfully", logrus.Fields{
		"count": len(result),
		"total": totalCount,
		"page":  params.Page,
	})

	return result, totalCount, nil
}

// CacheWarmOrganizations предварительно загружает организации в кеш
func (s *esfOrganizationServiceImpl) CacheWarmOrganizations(ctx context.Context) error {
	if s.cacheManager == nil {
		return nil
	}

	s.logger.Info(ctx, "Starting cache warming for organizations")

	orgs, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error(ctx, "Failed to get organizations for cache warming", err)
		return err
	}

	// Подготавливаем данные для пакетного кеширования
	batchData := make(map[string]interface{})
	for _, org := range orgs {
		batchData["id:"+org.ID.String()] = org
		batchData["name:"+org.Name] = org
	}

	// Кешируем все сразу (более эффективно)
	if err := s.cacheManager.Organization().SetMultiple(ctx, batchData, 2*time.Hour); err != nil {
		s.logger.Error(ctx, "Failed to warm organizations cache", err)
		return err
	}

	s.logger.Info(ctx, "Organizations cache warming completed", logrus.Fields{"count": len(orgs)})
	return nil
}
