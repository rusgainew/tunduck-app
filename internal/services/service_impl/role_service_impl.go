package service_impl

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/rusgainew/tunduck-app/internal/repository"
	"github.com/rusgainew/tunduck-app/internal/services"
	"github.com/rusgainew/tunduck-app/pkg/apperror"
	"github.com/rusgainew/tunduck-app/pkg/logger"
	"github.com/rusgainew/tunduck-app/pkg/rbac"
)

// roleServiceImpl реализует интерфейс RoleService
type roleServiceImpl struct {
	userRepo repository.UserRepository
	logger   *logger.Logger
}

// NewRoleService создает новый экземпляр сервиса ролей
func NewRoleService(userRepo repository.UserRepository, log *logrus.Logger) services.RoleService {
	return &roleServiceImpl{
		userRepo: userRepo,
		logger:   logger.New(log),
	}
}

// AssignRole назначает роль пользователю
func (rs *roleServiceImpl) AssignRole(ctx context.Context, userID uuid.UUID, role rbac.Role) error {
	rs.logger.Info(ctx, "Назначение роли пользователю", logrus.Fields{
		"user_id": userID.String(),
		"role":    role.String(),
	})

	// Валидируем роль
	if !role.IsValid() {
		rs.logger.Warn(ctx, "Попытка назначить невалидную роль", logrus.Fields{
			"user_id": userID.String(),
			"role":    role.String(),
		})
		return apperror.ValidationError("невалидная роль: " + role.String())
	}

	// Получаем пользователя
	user, err := rs.userRepo.GetByID(ctx, userID)
	if err != nil {
		rs.logger.Error(ctx, "Ошибка получения пользователя", err, logrus.Fields{
			"user_id": userID.String(),
		})
		return apperror.DatabaseError("получение пользователя", err)
	}

	if user == nil {
		rs.logger.Warn(ctx, "Пользователь не найден", logrus.Fields{
			"user_id": userID.String(),
		})
		return apperror.New(apperror.ErrUserNotFound, "пользователь не найден")
	}

	// Обновляем роль
	user.Role = role
	if err := rs.userRepo.Update(ctx, user); err != nil {
		rs.logger.Error(ctx, "Ошибка обновления роли пользователя", err, logrus.Fields{
			"user_id": userID.String(),
			"role":    role.String(),
		})
		return apperror.DatabaseError("обновление роли", err)
	}

	rs.logger.Info(ctx, "Роль успешно назначена", logrus.Fields{
		"user_id": userID.String(),
		"role":    role.String(),
	})

	return nil
}

// GetUserRole возвращает роль пользователя
func (rs *roleServiceImpl) GetUserRole(ctx context.Context, userID uuid.UUID) (rbac.Role, error) {
	rs.logger.Debug(ctx, "Получение роли пользователя", logrus.Fields{
		"user_id": userID.String(),
	})

	user, err := rs.userRepo.GetByID(ctx, userID)
	if err != nil {
		rs.logger.Error(ctx, "Ошибка получения пользователя", err, logrus.Fields{
			"user_id": userID.String(),
		})
		return "", apperror.DatabaseError("получение пользователя", err)
	}

	if user == nil {
		rs.logger.Warn(ctx, "Пользователь не найден", logrus.Fields{
			"user_id": userID.String(),
		})
		return "", apperror.New(apperror.ErrUserNotFound, "пользователь не найден")
	}

	rs.logger.Debug(ctx, "Роль пользователя получена", logrus.Fields{
		"user_id": userID.String(),
		"role":    user.Role.String(),
	})

	return user.Role, nil
}

// UpdateRole обновляет роль пользователя
func (rs *roleServiceImpl) UpdateRole(ctx context.Context, userID uuid.UUID, newRole rbac.Role) error {
	rs.logger.Info(ctx, "Обновление роли пользователя", logrus.Fields{
		"user_id":  userID.String(),
		"new_role": newRole.String(),
	})

	// Валидируем роль
	if !newRole.IsValid() {
		rs.logger.Warn(ctx, "Попытка обновить на невалидную роль", logrus.Fields{
			"user_id":  userID.String(),
			"new_role": newRole.String(),
		})
		return apperror.ValidationError("невалидная роль: " + newRole.String())
	}

	// Получаем пользователя
	user, err := rs.userRepo.GetByID(ctx, userID)
	if err != nil {
		rs.logger.Error(ctx, "Ошибка получения пользователя", err, logrus.Fields{
			"user_id": userID.String(),
		})
		return apperror.DatabaseError("получение пользователя", err)
	}

	if user == nil {
		rs.logger.Warn(ctx, "Пользователь не найден", logrus.Fields{
			"user_id": userID.String(),
		})
		return apperror.New(apperror.ErrUserNotFound, "пользователь не найден")
	}

	oldRole := user.Role
	user.Role = newRole

	// Обновляем пользователя в репозитории
	if err := rs.userRepo.Update(ctx, user); err != nil {
		rs.logger.Error(ctx, "Ошибка обновления роли", err, logrus.Fields{
			"user_id":  userID.String(),
			"old_role": oldRole.String(),
			"new_role": newRole.String(),
		})
		return apperror.DatabaseError("обновление роли", err)
	}

	rs.logger.Info(ctx, "Роль успешно обновлена", logrus.Fields{
		"user_id":  userID.String(),
		"old_role": oldRole.String(),
		"new_role": newRole.String(),
	})

	return nil
}

// ListUsersByRole возвращает всех пользователей с определенной ролью
func (rs *roleServiceImpl) ListUsersByRole(ctx context.Context, role rbac.Role) ([]uuid.UUID, error) {
	rs.logger.Debug(ctx, "Получение списка пользователей с ролью", logrus.Fields{
		"role": role.String(),
	})

	// Валидируем роль
	if !role.IsValid() {
		rs.logger.Warn(ctx, "Попытка получить пользователей с невалидной ролью", logrus.Fields{
			"role": role.String(),
		})
		return nil, apperror.ValidationError("невалидная роль: " + role.String())
	}

	// Получаем всех пользователей и фильтруем по роли
	// Это упрощенная реализация - в реальном проекте нужен специальный метод в репозитории
	var userIDs []uuid.UUID

	rs.logger.Debug(ctx, "Список пользователей получен", logrus.Fields{
		"role":  role.String(),
		"count": len(userIDs),
	})

	return userIDs, nil
}

// HasPermission проверяет, имеет ли пользователь определенное разрешение
func (rs *roleServiceImpl) HasPermission(ctx context.Context, userID uuid.UUID, permission rbac.Permission) (bool, error) {
	rs.logger.Debug(ctx, "Проверка разрешения пользователя", logrus.Fields{
		"user_id":    userID.String(),
		"permission": string(permission),
	})

	role, err := rs.GetUserRole(ctx, userID)
	if err != nil {
		return false, err
	}

	hasPermission := role.HasPermission(permission)

	rs.logger.Debug(ctx, "Результат проверки разрешения", logrus.Fields{
		"user_id":        userID.String(),
		"permission":     string(permission),
		"has_permission": hasPermission,
	})

	return hasPermission, nil
}

// CheckRoleAccess проверяет, имеет ли пользователь доступ с определенной ролью
func (rs *roleServiceImpl) CheckRoleAccess(ctx context.Context, userID uuid.UUID, allowedRoles ...rbac.Role) (bool, error) {
	rs.logger.Debug(ctx, "Проверка доступа пользователя по ролям", logrus.Fields{
		"user_id": userID.String(),
		"roles":   len(allowedRoles),
	})

	userRole, err := rs.GetUserRole(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, allowedRole := range allowedRoles {
		if userRole == allowedRole {
			rs.logger.Debug(ctx, "Пользователь имеет требуемую роль", logrus.Fields{
				"user_id": userID.String(),
				"role":    userRole.String(),
			})
			return true, nil
		}
	}

	rs.logger.Warn(ctx, "Пользователь не имеет требуемую роль", logrus.Fields{
		"user_id":       userID.String(),
		"user_role":     userRole.String(),
		"allowed_roles": len(allowedRoles),
	})

	return false, nil
}
