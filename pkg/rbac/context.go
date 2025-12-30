package rbac

import (
	"github.com/google/uuid"
)

// UserContext содержит информацию о текущем пользователе для проверки разрешений
type UserContext struct {
	UserID uuid.UUID
	Role   Role
}

// NewUserContext создает новый контекст пользователя
func NewUserContext(userID uuid.UUID, role Role) *UserContext {
	return &UserContext{
		UserID: userID,
		Role:   role,
	}
}

// HasPermission проверяет, есть ли у пользователя определенное разрешение
func (uc *UserContext) HasPermission(permission Permission) bool {
	if uc == nil || !uc.Role.IsValid() {
		return false
	}
	return uc.Role.HasPermission(permission)
}

// IsAdmin проверяет, является ли пользователь администратором
func (uc *UserContext) IsAdmin() bool {
	return uc != nil && uc.Role == RoleAdmin
}

// IsUser проверяет, является ли пользователь обычным пользователем
func (uc *UserContext) IsUser() bool {
	return uc != nil && uc.Role == RoleUser
}

// IsViewer проверяет, является ли пользователь просмотрщиком
func (uc *UserContext) IsViewer() bool {
	return uc != nil && uc.Role == RoleViewer
}
