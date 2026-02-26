package postgres

import (
	"time"

	"github.com/southern-martin/ecommerce/services/cms/internal/domain"
)

// BannerModel is the GORM model for the banners table.
type BannerModel struct {
	ID             string     `gorm:"type:uuid;primaryKey"`
	Title          string     `gorm:"type:varchar(255);not null"`
	ImageURL       string     `gorm:"type:text;not null"`
	LinkURL        string     `gorm:"type:text"`
	Position       string     `gorm:"type:varchar(50);index"`
	SortOrder      int        `gorm:"default:0"`
	TargetAudience string     `gorm:"type:varchar(100)"`
	StartsAt       time.Time  `gorm:"not null"`
	EndsAt         *time.Time `gorm:"index"`
	IsActive       bool       `gorm:"default:true;index"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
}

func (BannerModel) TableName() string { return "banners" }

func (m *BannerModel) ToDomain() *domain.Banner {
	return &domain.Banner{
		ID:             m.ID,
		Title:          m.Title,
		ImageURL:       m.ImageURL,
		LinkURL:        m.LinkURL,
		Position:       m.Position,
		SortOrder:      m.SortOrder,
		TargetAudience: m.TargetAudience,
		StartsAt:       m.StartsAt,
		EndsAt:         m.EndsAt,
		IsActive:       m.IsActive,
		CreatedAt:      m.CreatedAt,
	}
}

func ToBannerModel(b *domain.Banner) *BannerModel {
	return &BannerModel{
		ID:             b.ID,
		Title:          b.Title,
		ImageURL:       b.ImageURL,
		LinkURL:        b.LinkURL,
		Position:       b.Position,
		SortOrder:      b.SortOrder,
		TargetAudience: b.TargetAudience,
		StartsAt:       b.StartsAt,
		EndsAt:         b.EndsAt,
		IsActive:       b.IsActive,
		CreatedAt:      b.CreatedAt,
	}
}

// PageModel is the GORM model for the pages table.
type PageModel struct {
	ID              string     `gorm:"type:uuid;primaryKey"`
	Title           string     `gorm:"type:varchar(255);not null"`
	Slug            string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	ContentHTML     string     `gorm:"type:text"`
	MetaTitle       string     `gorm:"type:varchar(255)"`
	MetaDescription string     `gorm:"type:text"`
	Status          string     `gorm:"type:varchar(20);default:'draft';index"`
	PublishedAt     *time.Time
	CreatedAt       time.Time  `gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime"`
}

func (PageModel) TableName() string { return "pages" }

func (m *PageModel) ToDomain() *domain.Page {
	return &domain.Page{
		ID:              m.ID,
		Title:           m.Title,
		Slug:            m.Slug,
		ContentHTML:     m.ContentHTML,
		MetaTitle:       m.MetaTitle,
		MetaDescription: m.MetaDescription,
		Status:          domain.PageStatus(m.Status),
		PublishedAt:     m.PublishedAt,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func ToPageModel(p *domain.Page) *PageModel {
	return &PageModel{
		ID:              p.ID,
		Title:           p.Title,
		Slug:            p.Slug,
		ContentHTML:     p.ContentHTML,
		MetaTitle:       p.MetaTitle,
		MetaDescription: p.MetaDescription,
		Status:          string(p.Status),
		PublishedAt:     p.PublishedAt,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

// ContentScheduleModel is the GORM model for the content_schedules table.
type ContentScheduleModel struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	ContentType string    `gorm:"type:varchar(50);not null"`
	ContentID   string    `gorm:"type:uuid;not null"`
	Action      string    `gorm:"type:varchar(50);not null"`
	ScheduledAt time.Time `gorm:"not null;index"`
	Executed    bool      `gorm:"default:false;index"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (ContentScheduleModel) TableName() string { return "content_schedules" }

func (m *ContentScheduleModel) ToDomain() *domain.ContentSchedule {
	return &domain.ContentSchedule{
		ID:          m.ID,
		ContentType: m.ContentType,
		ContentID:   m.ContentID,
		Action:      m.Action,
		ScheduledAt: m.ScheduledAt,
		Executed:    m.Executed,
		CreatedAt:   m.CreatedAt,
	}
}

func ToContentScheduleModel(s *domain.ContentSchedule) *ContentScheduleModel {
	return &ContentScheduleModel{
		ID:          s.ID,
		ContentType: s.ContentType,
		ContentID:   s.ContentID,
		Action:      s.Action,
		ScheduledAt: s.ScheduledAt,
		Executed:    s.Executed,
		CreatedAt:   s.CreatedAt,
	}
}
