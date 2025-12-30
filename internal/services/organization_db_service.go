package services

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrganizationDBService интерфейс для управления динамическими БД организаций
type OrganizationDBService interface {
	// CreateOrganizationDatabase создает отдельную БД для организации
	CreateOrganizationDatabase(ctx context.Context, organizationID uuid.UUID) error

	// GetOrganizationDatabase получает подключение к БД организации
	GetOrganizationDatabase(ctx context.Context, organizationID uuid.UUID) (*gorm.DB, error)

	// DeleteOrganizationDatabase удаляет БД организации
	DeleteOrganizationDatabase(ctx context.Context, organizationID uuid.UUID) error
}
