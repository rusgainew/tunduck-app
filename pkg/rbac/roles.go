package rbac

// Role представляет роль пользователя в системе
type Role string

const (
	RoleAdmin  Role = "admin"  // Администратор - полный доступ
	RoleUser   Role = "user"   // Обычный пользователь - ограниченный доступ
	RoleViewer Role = "viewer" // Просмотр - только чтение
)

// IsValid проверяет, валидна ли роль
func (r Role) IsValid() bool {
	return r == RoleAdmin || r == RoleUser || r == RoleViewer
}

// String возвращает строковое представление роли
func (r Role) String() string {
	return string(r)
}

// Permission представляет разрешение для выполнения действия
type Permission string

const (
	// Права для организаций
	PermissionCreateOrganization Permission = "create:organization"
	PermissionReadOrganization   Permission = "read:organization"
	PermissionUpdateOrganization Permission = "update:organization"
	PermissionDeleteOrganization Permission = "delete:organization"

	// Права для документов
	PermissionCreateDocument Permission = "create:document"
	PermissionReadDocument   Permission = "read:document"
	PermissionUpdateDocument Permission = "update:document"
	PermissionDeleteDocument Permission = "delete:document"

	// Права для пользователей
	PermissionCreateUser Permission = "create:user"
	PermissionReadUser   Permission = "read:user"
	PermissionUpdateUser Permission = "update:user"
	PermissionDeleteUser Permission = "delete:user"

	// Права для ролей
	PermissionAssignRole Permission = "assign:role"
	PermissionViewRoles  Permission = "view:roles"
)

// RolePermissions определяет какие разрешения есть у каждой роли
var RolePermissions = map[Role][]Permission{
	RoleAdmin: {
		// Администратор имеет все права
		PermissionCreateOrganization, PermissionReadOrganization, PermissionUpdateOrganization, PermissionDeleteOrganization,
		PermissionCreateDocument, PermissionReadDocument, PermissionUpdateDocument, PermissionDeleteDocument,
		PermissionCreateUser, PermissionReadUser, PermissionUpdateUser, PermissionDeleteUser,
		PermissionAssignRole, PermissionViewRoles,
	},
	RoleUser: {
		// Обычный пользователь может читать и создавать
		PermissionReadOrganization,
		PermissionCreateDocument, PermissionReadDocument, PermissionUpdateDocument,
		PermissionReadUser, PermissionViewRoles,
	},
	RoleViewer: {
		// Viewer может только читать
		PermissionReadOrganization,
		PermissionReadDocument,
		PermissionReadUser,
		PermissionViewRoles,
	},
}

// HasPermission проверяет, есть ли у роли определенное разрешение
func (r Role) HasPermission(permission Permission) bool {
	perms, exists := RolePermissions[r]
	if !exists {
		return false
	}

	for _, p := range perms {
		if p == permission {
			return true
		}
	}
	return false
}

// GetAllPermissions возвращает все разрешения для роли
func (r Role) GetAllPermissions() []Permission {
	if perms, exists := RolePermissions[r]; exists {
		return perms
	}
	return []Permission{}
}
