package rbac

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rusgainew/tunduck-app/pkg/apperror"
)

// ContextKey используется для хранения UserContext в контексте Fiber
const ContextKey = "rbac_user_context"

// ExtractUserContext извлекает контекст пользователя из запроса
// Ожидается, что middleware установил это значение в контексте
func ExtractUserContext(ctx *fiber.Ctx) *UserContext {
	userCtx := ctx.Locals(ContextKey)
	if userCtx == nil {
		return nil
	}

	uc, ok := userCtx.(*UserContext)
	if !ok {
		return nil
	}

	return uc
}

// RequireRole создает middleware, который проверяет, имеет ли пользователь требуемую роль
func RequireRole(allowedRoles ...Role) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userCtx := ExtractUserContext(ctx)
		if userCtx == nil {
			appErr := apperror.New(apperror.ErrUnauthorized, "пользователь не авторизован")
			return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
		}

		// Проверяем, есть ли роль пользователя в разрешенных
		allowed := false
		for _, role := range allowedRoles {
			if userCtx.Role == role {
				allowed = true
				break
			}
		}

		if !allowed {
			appErr := apperror.New(apperror.ErrForbidden, "недостаточно прав доступа")
			return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
		}

		return ctx.Next()
	}
}

// RequirePermission создает middleware, который проверяет наличие определенного разрешения
func RequirePermission(permission Permission) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userCtx := ExtractUserContext(ctx)
		if userCtx == nil {
			appErr := apperror.New(apperror.ErrUnauthorized, "пользователь не авторизован")
			return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
		}

		if !userCtx.HasPermission(permission) {
			appErr := apperror.New(apperror.ErrForbidden, "недостаточно прав для выполнения этого действия")
			return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
		}

		return ctx.Next()
	}
}

// RequireAdminRole создает middleware для требования роли администратора
func RequireAdminRole() fiber.Handler {
	return RequireRole(RoleAdmin)
}

// SetUserContext устанавливает контекст пользователя в запросе
// Обычно вызывается из JWT middleware после проверки токена
func SetUserContext(ctx *fiber.Ctx, userID uuid.UUID, role Role) {
	userContext := NewUserContext(userID, role)
	ctx.Locals(ContextKey, userContext)
}
