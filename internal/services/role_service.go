package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/rusgainew/tunduck-app/pkg/rbac"
)

// RoleService определяет интерфейс для управления ролями пользователей
type RoleService interface {
	// AssignRole назначает роль пользователю
	AssignRole(ctx context.Context, userID uuid.UUID, role rbac.Role) error

	// GetUserRole возвращает роль пользователя
	GetUserRole(ctx context.Context, userID uuid.UUID) (rbac.Role, error)

	// UpdateRole обновляет роль пользователя
	UpdateRole(ctx context.Context, userID uuid.UUID, newRole rbac.Role) error

	// ListUsersByRole возвращает всех пользователей с определенной ролью
	ListUsersByRole(ctx context.Context, role rbac.Role) ([]uuid.UUID, error)

	// HasPermission проверяет, имеет ли пользователь определенное разрешение
	HasPermission(ctx context.Context, userID uuid.UUID, permission rbac.Permission) (bool, error)

	// CheckRoleAccess проверяет, имеет ли пользователь доступ с определенной ролью
	CheckRoleAccess(ctx context.Context, userID uuid.UUID, allowedRoles ...rbac.Role) (bool, error)
}
