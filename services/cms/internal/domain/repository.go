package domain

import "context"

// BannerRepository defines the interface for banner persistence.
type BannerRepository interface {
	GetByID(ctx context.Context, id string) (*Banner, error)
	ListActive(ctx context.Context, position string) ([]Banner, error)
	ListAll(ctx context.Context, page, pageSize int) ([]Banner, int64, error)
	Create(ctx context.Context, banner *Banner) error
	Update(ctx context.Context, banner *Banner) error
	Delete(ctx context.Context, id string) error
}

// PageRepository defines the interface for page persistence.
type PageRepository interface {
	GetByID(ctx context.Context, id string) (*Page, error)
	GetBySlug(ctx context.Context, slug string) (*Page, error)
	ListPublished(ctx context.Context, page, pageSize int) ([]Page, int64, error)
	ListAll(ctx context.Context, page, pageSize int) ([]Page, int64, error)
	Create(ctx context.Context, pg *Page) error
	Update(ctx context.Context, pg *Page) error
	Delete(ctx context.Context, id string) error
}

// ScheduleRepository defines the interface for content schedule persistence.
type ScheduleRepository interface {
	GetPending(ctx context.Context) ([]ContentSchedule, error)
	Create(ctx context.Context, schedule *ContentSchedule) error
	MarkExecuted(ctx context.Context, id string) error
}
