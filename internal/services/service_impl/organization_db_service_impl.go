package service_impl

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/rusgainew/tunduck-app/pkg/entity"
)

// OrganizationDBServiceImpl реализация сервиса для управления динамическими БД организаций
type OrganizationDBServiceImpl struct {
	mainDB     *gorm.DB
	logger     *logrus.Logger
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbSSLMode  string
}

// NewOrganizationDBService создает новый сервис управления БД организаций
func NewOrganizationDBService(
	mainDB *gorm.DB,
	logger *logrus.Logger,
	dbHost, dbPort, dbUser, dbPassword, dbSSLMode string,
) *OrganizationDBServiceImpl {
	return &OrganizationDBServiceImpl{
		mainDB:     mainDB,
		logger:     logger,
		dbHost:     dbHost,
		dbPort:     dbPort,
		dbUser:     dbUser,
		dbPassword: dbPassword,
		dbSSLMode:  dbSSLMode,
	}
}

// getOrganizationDBName возвращает имя БД для организации
func (s *OrganizationDBServiceImpl) getOrganizationDBName(organizationID uuid.UUID) string {
	// Заменяем дефисы на подчеркивание для корректности имени БД
	return fmt.Sprintf("org_%s", organizationID.String()[:8])
}

// CreateOrganizationDatabase создает отдельную БД для организации с таблицами EsfDocument и EsfEntries
func (s *OrganizationDBServiceImpl) CreateOrganizationDatabase(ctx context.Context, organizationID uuid.UUID) error {
	dbName := s.getOrganizationDBName(organizationID)

	s.logger.WithFields(logrus.Fields{
		"organization_id": organizationID,
		"database_name":   dbName,
	}).Info("Creating organization database")

	// Подключаемся к главной БД для создания новой БД
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=%s",
		s.dbHost, s.dbUser, s.dbPassword, s.dbPort, s.dbSSLMode)

	tempDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		s.logger.WithError(err).Error("Failed to connect to postgres for database creation")
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Создаем новую БД
	if err := tempDB.Exec(fmt.Sprintf("CREATE DATABASE %q", dbName)).Error; err != nil {
		s.logger.WithError(err).WithField("database_name", dbName).Warn("Failed to create database (may already exist)")
		// Не возвращаем ошибку, т.к. БД может уже существовать
	}

	// Подключаемся к новой БД
	newDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		s.dbHost, s.dbUser, s.dbPassword, dbName, s.dbPort, s.dbSSLMode)

	newDB, err := gorm.Open(postgres.Open(newDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		s.logger.WithError(err).Error("Failed to connect to new organization database")
		return fmt.Errorf("failed to connect to organization database: %w", err)
	}

	// Создаем расширение uuid-ossp для генерации UUID
	if err := newDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		s.logger.WithError(err).Warn("Failed to create uuid-ossp extension")
	}

	// Выполняем миграции для новой БД (создаем таблицы EsfDocument и EsfEntries)
	if err := newDB.AutoMigrate(&entity.EsfDocument{}, &entity.EsfEntries{}); err != nil {
		s.logger.WithError(err).Error("Failed to migrate organization database")
		return fmt.Errorf("failed to migrate organization database: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"organization_id": organizationID,
		"database_name":   dbName,
	}).Info("Organization database created successfully")

	return nil
}

// GetOrganizationDatabase получает подключение к БД организации
func (s *OrganizationDBServiceImpl) GetOrganizationDatabase(ctx context.Context, organizationID uuid.UUID) (*gorm.DB, error) {
	dbName := s.getOrganizationDBName(organizationID)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		s.dbHost, s.dbUser, s.dbPassword, dbName, s.dbPort, s.dbSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		s.logger.WithError(err).WithField("organization_id", organizationID).Error("Failed to connect to organization database")
		return nil, fmt.Errorf("failed to connect to organization database: %w", err)
	}

	return db, nil
}

// DeleteOrganizationDatabase удаляет БД организации
func (s *OrganizationDBServiceImpl) DeleteOrganizationDatabase(ctx context.Context, organizationID uuid.UUID) error {
	dbName := s.getOrganizationDBName(organizationID)

	s.logger.WithFields(logrus.Fields{
		"organization_id": organizationID,
		"database_name":   dbName,
	}).Info("Deleting organization database")

	// Подключаемся к главной БД для удаления БД организации
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=%s",
		s.dbHost, s.dbUser, s.dbPassword, s.dbPort, s.dbSSLMode)

	tempDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		s.logger.WithError(err).Error("Failed to connect to postgres for database deletion")
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Удаляем БД
	if err := tempDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %q", dbName)).Error; err != nil {
		s.logger.WithError(err).Error("Failed to delete organization database")
		return fmt.Errorf("failed to delete organization database: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"organization_id": organizationID,
		"database_name":   dbName,
	}).Info("Organization database deleted successfully")

	return nil
}
