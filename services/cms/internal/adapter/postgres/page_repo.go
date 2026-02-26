package postgres

import (
	"context"

	"github.com/southern-martin/ecommerce/services/cms/internal/domain"
	"gorm.io/gorm"
)

// PageRepo implements domain.PageRepository.
type PageRepo struct {
	db *gorm.DB
}

// NewPageRepo creates a new PageRepo.
func NewPageRepo(db *gorm.DB) *PageRepo {
	return &PageRepo{db: db}
}

func (r *PageRepo) GetByID(ctx context.Context, id string) (*domain.Page, error) {
	var model PageModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *PageRepo) GetBySlug(ctx context.Context, slug string) (*domain.Page, error) {
	var model PageModel
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *PageRepo) ListPublished(ctx context.Context, page, pageSize int) ([]domain.Page, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&PageModel{}).Where("status = ?", string(domain.PageStatusPublished)).Count(&total)

	var models []PageModel
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Where("status = ?", string(domain.PageStatusPublished)).
		Order("published_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	pages := make([]domain.Page, len(models))
	for i, m := range models {
		pages[i] = *m.ToDomain()
	}
	return pages, total, nil
}

func (r *PageRepo) ListAll(ctx context.Context, page, pageSize int) ([]domain.Page, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&PageModel{}).Count(&total)

	var models []PageModel
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	pages := make([]domain.Page, len(models))
	for i, m := range models {
		pages[i] = *m.ToDomain()
	}
	return pages, total, nil
}

func (r *PageRepo) Create(ctx context.Context, pg *domain.Page) error {
	model := ToPageModel(pg)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *PageRepo) Update(ctx context.Context, pg *domain.Page) error {
	return r.db.WithContext(ctx).Model(&PageModel{}).Where("id = ?", pg.ID).Updates(map[string]interface{}{
		"title":            pg.Title,
		"slug":             pg.Slug,
		"content_html":     pg.ContentHTML,
		"meta_title":       pg.MetaTitle,
		"meta_description": pg.MetaDescription,
		"status":           string(pg.Status),
		"published_at":     pg.PublishedAt,
	}).Error
}

func (r *PageRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&PageModel{}).Error
}
