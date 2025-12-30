package pagination

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// FilterParams містить параметри для фільтрації
type FilterParams struct {
	Fields map[string]interface{}
}

// DocumentFilterParams спеціалізована структура для фільтрації документів
type DocumentFilterParams struct {
	Status        string // active, archived
	CreatedAfter  string // ISO 8601 дата
	CreatedBefore string
	Search        string // пошук по назві/опису
}

// OrganizationFilterParams спеціалізована структура для фільтрації організацій
type OrganizationFilterParams struct {
	Status string // active, inactive
	Search string // пошук по назві
}

// UserFilterParams спеціалізована структура для фільтрації користувачів
type UserFilterParams struct {
	Status string // active, inactive
	Search string // пошук по імені/email
	RoleID string
}

// ExtractDocumentFilters витягує фільтри для документів
func ExtractDocumentFilters(ctx *fiber.Ctx) DocumentFilterParams {
	return DocumentFilterParams{
		Status:        ctx.Query("status", ""),
		CreatedAfter:  ctx.Query("created_after", ""),
		CreatedBefore: ctx.Query("created_before", ""),
		Search:        ctx.Query("search", ""),
	}
}

// ExtractOrganizationFilters витягує фільтри для організацій
func ExtractOrganizationFilters(ctx *fiber.Ctx) OrganizationFilterParams {
	return OrganizationFilterParams{
		Status: ctx.Query("status", ""),
		Search: ctx.Query("search", ""),
	}
}

// ExtractUserFilters витягує фільтри для користувачів
func ExtractUserFilters(ctx *fiber.Ctx) UserFilterParams {
	return UserFilterParams{
		Status: ctx.Query("status", ""),
		Search: ctx.Query("search", ""),
		RoleID: ctx.Query("role_id", ""),
	}
}

// HasFilters перевіряє, чи встановлені якісь фільтри
func (f DocumentFilterParams) HasFilters() bool {
	return f.Status != "" || f.CreatedAfter != "" ||
		f.CreatedBefore != "" || strings.TrimSpace(f.Search) != ""
}

// HasFilters перевіряє, чи встановлені якісь фільтри
func (f OrganizationFilterParams) HasFilters() bool {
	return f.Status != "" || strings.TrimSpace(f.Search) != ""
}

// HasFilters перевіряє, чи встановлені якісь фільтри
func (f UserFilterParams) HasFilters() bool {
	return f.Status != "" || strings.TrimSpace(f.Search) != "" || f.RoleID != ""
}

// ToMap конвертує фільтри в map для GORM
func (f DocumentFilterParams) ToMap() map[string]interface{} {
	filters := make(map[string]interface{})

	if f.Status != "" {
		filters["status"] = f.Status
	}
	if f.Search != "" {
		filters["search"] = "%" + strings.TrimSpace(f.Search) + "%"
	}

	return filters
}

// ToMap конвертує фільтри в map для GORM
func (f OrganizationFilterParams) ToMap() map[string]interface{} {
	filters := make(map[string]interface{})

	if f.Status != "" {
		filters["status"] = f.Status
	}
	if f.Search != "" {
		filters["search"] = "%" + strings.TrimSpace(f.Search) + "%"
	}

	return filters
}

// ToMap конвертує фільтри в map для GORM
func (f UserFilterParams) ToMap() map[string]interface{} {
	filters := make(map[string]interface{})

	if f.Status != "" {
		filters["status"] = f.Status
	}
	if f.RoleID != "" {
		filters["role_id"] = f.RoleID
	}
	if f.Search != "" {
		filters["search"] = "%" + strings.TrimSpace(f.Search) + "%"
	}

	return filters
}
