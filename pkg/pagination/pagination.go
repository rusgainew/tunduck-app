package pagination

import (
	"math"

	"github.com/gofiber/fiber/v2"
)

// PaginationParams містить параметри для пагінації
type PaginationParams struct {
	Page     int    `query:"page" default:"1"`
	PageSize int    `query:"page_size" default:"10"`
	Sort     string `query:"sort" default:"created_at"`
	Order    string `query:"order" default:"desc"`
}

// PaginatedResponse містить дані з інформацією про пагінацію
type PaginatedResponse struct {
	Data       interface{}    `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}

// PaginationInfo містить інформацію про пагінацію
type PaginationInfo struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// ExtractPaginationParams витягує параметри пагінації з запиту
func ExtractPaginationParams(ctx *fiber.Ctx) PaginationParams {
	page := ctx.QueryInt("page", 1)
	pageSize := ctx.QueryInt("page_size", 10)
	sort := ctx.Query("sort", "created_at")
	order := ctx.Query("order", "desc")

	// Валідація
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	if order != "asc" && order != "desc" {
		order = "desc"
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
		Order:    order,
	}
}

// GetOffset повертає offset для SQL запиту
func (p PaginationParams) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

// GetLimit повертає limit для SQL запиту
func (p PaginationParams) GetLimit() int {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p.PageSize
}

// NewPaginationInfo створює нову інформацію про пагінацію
func NewPaginationInfo(page, pageSize int, totalItems int64) PaginationInfo {
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
	if totalPages < 1 {
		totalPages = 1
	}

	return PaginationInfo{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// NewPaginatedResponse створює нову пагінвану відповідь
func NewPaginatedResponse(data interface{}, page, pageSize int, totalItems int64) PaginatedResponse {
	return PaginatedResponse{
		Data:       data,
		Pagination: NewPaginationInfo(page, pageSize, totalItems),
	}
}
