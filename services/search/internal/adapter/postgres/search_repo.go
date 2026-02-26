package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/southern-martin/ecommerce/services/search/internal/domain"
	"gorm.io/gorm"
)

// SearchRepo implements domain.SearchRepository using PostgreSQL.
type SearchRepo struct {
	db *gorm.DB
}

// NewSearchRepo creates a new SearchRepo.
func NewSearchRepo(db *gorm.DB) *SearchRepo {
	return &SearchRepo{db: db}
}

// Index creates or updates a product in the search index (upsert).
func (r *SearchRepo) Index(ctx context.Context, idx *domain.SearchIndex) error {
	model := ToModel(idx)
	if model.ID == "" {
		model.ID = uuid.New().String()
	}

	// Upsert: try to find existing by product_id, then save
	var existing SearchIndexModel
	result := r.db.WithContext(ctx).Where("product_id = ?", model.ProductID).First(&existing)
	if result.Error == nil {
		// Update existing record
		model.ID = existing.ID
		model.CreatedAt = existing.CreatedAt
		if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
			log.Error().Err(err).Str("product_id", idx.ProductID).Msg("failed to update search index")
			return err
		}
	} else {
		// Create new record
		if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
			log.Error().Err(err).Str("product_id", idx.ProductID).Msg("failed to create search index")
			return err
		}
	}

	log.Info().Str("product_id", idx.ProductID).Msg("product indexed")
	return nil
}

// Delete removes a product from the search index by product ID.
func (r *SearchRepo) Delete(ctx context.Context, productID string) error {
	result := r.db.WithContext(ctx).Where("product_id = ?", productID).Delete(&SearchIndexModel{})
	if result.Error != nil {
		log.Error().Err(result.Error).Str("product_id", productID).Msg("failed to delete from search index")
		return result.Error
	}
	log.Info().Str("product_id", productID).Int64("rows_affected", result.RowsAffected).Msg("product removed from index")
	return nil
}

// Search performs a full-text search using ILIKE on name and description.
func (r *SearchRepo) Search(ctx context.Context, filter domain.SearchFilter) ([]domain.SearchResult, int64, error) {
	query := r.db.WithContext(ctx).Model(&SearchIndexModel{})

	// Text search using ILIKE
	if filter.Query != "" {
		pattern := fmt.Sprintf("%%%s%%", filter.Query)
		query = query.Where("name ILIKE ? OR description ILIKE ?", pattern, pattern)
	}

	// Category filter
	if filter.CategoryID != "" {
		query = query.Where("category_id = ?", filter.CategoryID)
	}

	// Price range filters
	if filter.MinPrice > 0 {
		query = query.Where("price_cents >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query = query.Where("price_cents <= ?", filter.MaxPrice)
	}

	// In-stock filter
	if filter.InStock != nil {
		query = query.Where("in_stock = ?", *filter.InStock)
	}

	// Seller filter
	if filter.SellerID != "" {
		query = query.Where("seller_id = ?", filter.SellerID)
	}

	// Count total results
	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Error().Err(err).Msg("failed to count search results")
		return nil, 0, err
	}

	// Sorting
	orderClause := "created_at DESC"
	if filter.SortBy != "" {
		direction := "ASC"
		if filter.SortOrder == "desc" {
			direction = "DESC"
		}
		switch filter.SortBy {
		case "price":
			orderClause = fmt.Sprintf("price_cents %s", direction)
		case "rating":
			orderClause = fmt.Sprintf("rating %s", direction)
		case "name":
			orderClause = fmt.Sprintf("name %s", direction)
		case "created_at":
			orderClause = fmt.Sprintf("created_at %s", direction)
		}
	}
	query = query.Order(orderClause)

	// Pagination
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Execute query
	var models []SearchIndexModel
	if err := query.Find(&models).Error; err != nil {
		log.Error().Err(err).Msg("failed to execute search query")
		return nil, 0, err
	}

	// Convert to domain results
	results := make([]domain.SearchResult, len(models))
	for i, m := range models {
		results[i] = domain.SearchResult{
			ID:          m.ID,
			ProductID:   m.ProductID,
			Name:        m.Name,
			Slug:        m.Slug,
			Description: m.Description,
			PriceCents:  m.PriceCents,
			Currency:    m.Currency,
			ImageURL:    m.ImageURL,
			SellerID:    m.SellerID,
			CategoryID:  m.CategoryID,
			Rating:      m.Rating,
			ReviewCount: m.ReviewCount,
			InStock:     m.InStock,
			Score:       1.0,
			CreatedAt:   m.CreatedAt,
		}
	}

	return results, total, nil
}

// Suggest returns search suggestions matching the query prefix.
func (r *SearchRepo) Suggest(ctx context.Context, query string, limit int) ([]domain.SearchSuggestion, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	pattern := fmt.Sprintf("%%%s%%", query)
	var models []SearchIndexModel
	if err := r.db.WithContext(ctx).
		Where("name ILIKE ?", pattern).
		Order("rating DESC").
		Limit(limit).
		Find(&models).Error; err != nil {
		log.Error().Err(err).Str("query", query).Msg("failed to get suggestions")
		return nil, err
	}

	suggestions := make([]domain.SearchSuggestion, len(models))
	for i, m := range models {
		suggestions[i] = domain.SearchSuggestion{
			Text:      m.Name,
			Type:      "product",
			ProductID: m.ProductID,
		}
	}

	return suggestions, nil
}
