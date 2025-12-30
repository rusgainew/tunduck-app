package repositorypostgres

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rusgainew/tunduck-app/internal/repository"
	"github.com/rusgainew/tunduck-app/pkg/apperror"
	"github.com/rusgainew/tunduck-app/pkg/entity"
	"github.com/rusgainew/tunduck-app/pkg/logger"
	"github.com/rusgainew/tunduck-app/pkg/pagination"
)

// esfOrganizationPostgres реализует интерфейс EsfOrganizationRepository для PostgreSQL
type esfOrganizationPostgres struct {
	logger *logger.Logger
	db     *gorm.DB
}

// NewEsfOrganizationRepositoryPostgres создает новый экземпляр репозитория организаций
func NewEsfOrganizationRepositoryPostgres(db *gorm.DB, log *logrus.Logger) repository.EsfOrganizationRepository {
	return &esfOrganizationPostgres{
		logger: logger.New(log),
		db:     db,
	}
}

// GetAll возвращает все организации из БД
func (eop *esfOrganizationPostgres) GetAll(ctx context.Context) ([]*entity.EstOrganization, error) {
	eop.logger.Debug(ctx, "Fetching all organizations from database", logrus.Fields{})

	var organizations []*entity.EstOrganization

	if err := eop.db.WithContext(ctx).Find(&organizations).Error; err != nil {
		eop.logger.Error(ctx, "Failed to fetch organizations from database", err, logrus.Fields{})
		return nil, apperror.DatabaseError("fetching organizations", err)
	}

	eop.logger.Debug(ctx, "Organizations fetched successfully", logrus.Fields{"count": len(organizations)})
	return organizations, nil
}

// GetByID возвращает организацию по ID
func (eop *esfOrganizationPostgres) GetByID(ctx context.Context, id string) (*entity.EstOrganization, error) {
	eop.logger.Debug(ctx, "Fetching organization by ID from database", logrus.Fields{"id": id})

	var organization entity.EstOrganization

	if err := eop.db.WithContext(ctx).Where("id = ?", id).First(&organization).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			eop.logger.Debug(ctx, "Organization not found", logrus.Fields{"id": id})
			return nil, nil
		}
		eop.logger.Error(ctx, "Failed to fetch organization by ID from database", err, logrus.Fields{"id": id})
		return nil, apperror.DatabaseError("fetching organization by ID", err)
	}

	eop.logger.Debug(ctx, "Organization fetched successfully", logrus.Fields{"id": id})
	return &organization, nil
}

// Insert создает новую организацию в БД
func (eop *esfOrganizationPostgres) Insert(ctx context.Context, org *entity.EstOrganization) error {
	eop.logger.Debug(ctx, "Inserting organization into database", logrus.Fields{"name": org.Name, "id": org.ID.String()})

	if err := eop.db.WithContext(ctx).Create(org).Error; err != nil {
		eop.logger.Error(ctx, "Failed to insert organization into database", err, logrus.Fields{"name": org.Name})
		return apperror.DatabaseError("inserting organization", err)
	}

	eop.logger.Debug(ctx, "Organization inserted successfully", logrus.Fields{"id": org.ID.String()})
	return nil
}

// Update обновляет данные организации в БД
func (eop *esfOrganizationPostgres) Update(ctx context.Context, org *entity.EstOrganization) error {
	eop.logger.Debug(ctx, "Updating organization in database", logrus.Fields{"id": org.ID.String()})

	if err := eop.db.WithContext(ctx).Save(org).Error; err != nil {
		eop.logger.Error(ctx, "Failed to update organization in database", err, logrus.Fields{"id": org.ID.String()})
		return apperror.DatabaseError("updating organization", err)
	}

	eop.logger.Debug(ctx, "Organization updated successfully", logrus.Fields{"id": org.ID.String()})
	return nil
}

// Delete удаляет организацию из БД по ID
func (eop *esfOrganizationPostgres) Delete(ctx context.Context, id string) error {
	eop.logger.Debug(ctx, "Deleting organization from database", logrus.Fields{"id": id})

	if err := eop.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.EstOrganization{}).Error; err != nil {
		eop.logger.Error(ctx, "Failed to delete organization from database", err, logrus.Fields{"id": id})
		return apperror.DatabaseError("deleting organization", err)
	}

	eop.logger.Debug(ctx, "Organization deleted successfully", logrus.Fields{"id": id})
	return nil
}

// CreateDatabase создает новую базу данных для организации и применяет миграции
func (eop *esfOrganizationPostgres) CreateDatabase(ctx context.Context, dbName string) error {
	eop.logger.Debug(ctx, "Creating database for organization", logrus.Fields{"dbName": dbName})

	// Выполняем SQL команду создания базы данных
	// Используем Exec вместо параметризованного запроса, так как имя БД не может быть параметром
	sql := fmt.Sprintf("CREATE DATABASE %s", dbName)

	if err := eop.db.WithContext(ctx).Exec(sql).Error; err != nil {
		eop.logger.Error(ctx, "Failed to create database", err, logrus.Fields{"dbName": dbName})
		return apperror.DatabaseError("creating database", err)
	}

	eop.logger.Debug(ctx, "Database created successfully", logrus.Fields{"dbName": dbName})

	// Подготовим параметры подключения к только что созданной базе
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbName, port, sslmode)

	// Подключаемся к новой базе данных
	newDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		eop.logger.Error(ctx, "Failed to connect to new database", err, logrus.Fields{"dbName": dbName})
		return apperror.DatabaseError("connecting to new database", err)
	}

	// Включаем расширение uuid-ossp для новой базы
	if err := newDB.WithContext(ctx).Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		eop.logger.Warn(ctx, "Failed to create uuid-ossp extension in new database", logrus.Fields{"dbName": dbName})
	}
	// не прерываем, так как возможно миграции пройдут, если расширение уже есть/не требуется

	// Применяем миграции для пустых таблиц EsfDocument и EsfEntries
	if err := newDB.WithContext(ctx).AutoMigrate(&entity.EsfDocument{}, &entity.EsfEntries{}); err != nil {
		eop.logger.Error(ctx, "Failed to run migrations in new database", err, logrus.Fields{"dbName": dbName})
		return apperror.DatabaseError("running migrations", err)
	}

	eop.logger.Debug(ctx, "Migrations applied to new database successfully", logrus.Fields{"dbName": dbName})
	return nil
}

// GetAllPaginated возвращает все организации из БД с пагинацией и фильтрацией
func (eop *esfOrganizationPostgres) GetAllPaginated(ctx context.Context, params pagination.PaginationParams, filters pagination.OrganizationFilterParams) ([]*entity.EstOrganization, int64, error) {
	eop.logger.Debug(ctx, "Fetching organizations with pagination", logrus.Fields{
		"page":        params.Page,
		"page_size":   params.PageSize,
		"sort":        params.Sort,
		"order":       params.Order,
		"has_filters": filters.HasFilters(),
	})

	var organizations []*entity.EstOrganization
	var totalCount int64

	query := eop.db.WithContext(ctx)

	// Применяем фильтры
	if filters.Status != "" {
		eop.logger.Debug(ctx, "Applying status filter", logrus.Fields{"status": filters.Status})
		query = query.Where("status = ?", filters.Status)
	}

	if filters.Search != "" {
		eop.logger.Debug(ctx, "Applying search filter", logrus.Fields{"search": filters.Search})
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}

	// Получаем общее количество
	if err := query.Model(&entity.EstOrganization{}).Count(&totalCount).Error; err != nil {
		eop.logger.Error(ctx, "Failed to count organizations", err, logrus.Fields{})
		return nil, 0, apperror.DatabaseError("counting organizations", err)
	}

	// Применяем сортировку и пагинацию
	if err := query.
		Order(params.Sort + " " + params.Order).
		Offset(params.GetOffset()).
		Limit(params.GetLimit()).
		Find(&organizations).Error; err != nil {
		eop.logger.Error(ctx, "Failed to fetch paginated organizations", err, logrus.Fields{
			"page":      params.Page,
			"page_size": params.PageSize,
		})
		return nil, 0, apperror.DatabaseError("fetching paginated organizations", err)
	}

	eop.logger.Debug(ctx, "Organizations fetched successfully", logrus.Fields{
		"count": len(organizations),
		"total": totalCount,
		"page":  params.Page,
	})

	return organizations, totalCount, nil
}
