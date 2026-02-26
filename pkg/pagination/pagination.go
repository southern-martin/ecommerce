package pagination

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

// PaginationParams holds the pagination parameters for a query.
type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// Offset calculates the SQL offset based on page and page size.
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// ParseFromGin extracts pagination parameters from Gin query parameters.
// Defaults: page=1, page_size=20, max page_size=100.
func ParseFromGin(c *gin.Context) PaginationParams {
	page := defaultPage
	pageSize := defaultPageSize

	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}
}

// PaginatedResult is a generic paginated response.
type PaginatedResult[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// NewPaginatedResult creates a PaginatedResult from the given items, total count, and pagination params.
func NewPaginatedResult[T any](items []T, total int64, params PaginationParams) PaginatedResult[T] {
	totalPages := int(math.Ceil(float64(total) / float64(params.PageSize)))
	if totalPages < 1 {
		totalPages = 1
	}

	return PaginatedResult[T]{
		Items:      items,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}
}

// ApplyToGORM applies pagination (limit and offset) to a GORM query.
func ApplyToGORM(db *gorm.DB, params PaginationParams) *gorm.DB {
	return db.Offset(params.Offset()).Limit(params.PageSize)
}
