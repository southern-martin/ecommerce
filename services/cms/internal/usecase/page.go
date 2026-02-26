package usecase

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/southern-martin/ecommerce/services/cms/internal/domain"
)

// PageUseCase handles page operations.
type PageUseCase struct {
	pageRepo  domain.PageRepository
	publisher domain.EventPublisher
}

// NewPageUseCase creates a new PageUseCase.
func NewPageUseCase(pageRepo domain.PageRepository, publisher domain.EventPublisher) *PageUseCase {
	return &PageUseCase{
		pageRepo:  pageRepo,
		publisher: publisher,
	}
}

// generateSlug creates a URL-friendly slug from a title.
func generateSlug(title string) string {
	slug := strings.ToLower(title)
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

// CreatePage creates a new CMS page.
func (uc *PageUseCase) CreatePage(ctx context.Context, page *domain.Page) error {
	page.ID = uuid.New().String()
	page.Slug = generateSlug(page.Title)
	page.Status = domain.PageStatusDraft

	if err := uc.pageRepo.Create(ctx, page); err != nil {
		return fmt.Errorf("failed to create page: %w", err)
	}

	_ = uc.publisher.Publish(ctx, "cms.page.created", map[string]interface{}{
		"page_id": page.ID,
		"title":   page.Title,
		"slug":    page.Slug,
	})

	return nil
}

// GetPage retrieves a page by ID.
func (uc *PageUseCase) GetPage(ctx context.Context, id string) (*domain.Page, error) {
	return uc.pageRepo.GetByID(ctx, id)
}

// GetPageBySlug retrieves a page by its slug.
func (uc *PageUseCase) GetPageBySlug(ctx context.Context, slug string) (*domain.Page, error) {
	return uc.pageRepo.GetBySlug(ctx, slug)
}

// ListPublishedPages returns published pages with pagination.
func (uc *PageUseCase) ListPublishedPages(ctx context.Context, page, pageSize int) ([]domain.Page, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.pageRepo.ListPublished(ctx, page, pageSize)
}

// ListAllPages returns all pages with pagination.
func (uc *PageUseCase) ListAllPages(ctx context.Context, page, pageSize int) ([]domain.Page, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.pageRepo.ListAll(ctx, page, pageSize)
}

// UpdatePage updates an existing page.
func (uc *PageUseCase) UpdatePage(ctx context.Context, page *domain.Page) error {
	existing, err := uc.pageRepo.GetByID(ctx, page.ID)
	if err != nil {
		return fmt.Errorf("page not found: %w", err)
	}

	// Regenerate slug if title changed
	if page.Title != "" && page.Title != existing.Title {
		page.Slug = generateSlug(page.Title)
	} else {
		page.Slug = existing.Slug
	}

	if err := uc.pageRepo.Update(ctx, page); err != nil {
		return fmt.Errorf("failed to update page: %w", err)
	}

	_ = uc.publisher.Publish(ctx, "cms.page.updated", map[string]interface{}{
		"page_id": page.ID,
	})

	return nil
}

// DeletePage deletes a page by ID.
func (uc *PageUseCase) DeletePage(ctx context.Context, id string) error {
	if _, err := uc.pageRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("page not found: %w", err)
	}

	if err := uc.pageRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete page: %w", err)
	}

	_ = uc.publisher.Publish(ctx, "cms.page.deleted", map[string]interface{}{
		"page_id": id,
	})

	return nil
}

// PublishPage sets a page's status to published.
func (uc *PageUseCase) PublishPage(ctx context.Context, id string) (*domain.Page, error) {
	page, err := uc.pageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("page not found: %w", err)
	}

	now := time.Now()
	page.Status = domain.PageStatusPublished
	page.PublishedAt = &now

	if err := uc.pageRepo.Update(ctx, page); err != nil {
		return nil, fmt.Errorf("failed to publish page: %w", err)
	}

	_ = uc.publisher.Publish(ctx, "cms.page.published", map[string]interface{}{
		"page_id": page.ID,
		"slug":    page.Slug,
	})

	return page, nil
}
