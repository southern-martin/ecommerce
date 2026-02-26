package postgres

import (
	"time"

	"github.com/southern-martin/ecommerce/services/media/internal/domain"
)

// MediaModel is the GORM model for the media_files table.
type MediaModel struct {
	ID           string    `gorm:"type:uuid;primaryKey"`
	OwnerID      string    `gorm:"type:uuid;index;not null"`
	OwnerType    string    `gorm:"type:varchar(50);index;not null"`
	FileName     string    `gorm:"type:varchar(500);not null"`
	OriginalName string    `gorm:"type:varchar(500);not null"`
	ContentType  string    `gorm:"type:varchar(100);not null"`
	SizeBytes    int64     `gorm:"default:0"`
	URL          string    `gorm:"type:text"`
	ThumbnailURL string    `gorm:"type:text"`
	Width        int       `gorm:"default:0"`
	Height       int       `gorm:"default:0"`
	Status       string    `gorm:"type:varchar(20);default:'pending'"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

// TableName returns the table name for the MediaModel.
func (MediaModel) TableName() string { return "media_files" }

// ToDomain converts a MediaModel to a domain Media entity.
func (m *MediaModel) ToDomain() *domain.Media {
	return &domain.Media{
		ID:           m.ID,
		OwnerID:      m.OwnerID,
		OwnerType:    m.OwnerType,
		FileName:     m.FileName,
		OriginalName: m.OriginalName,
		ContentType:  m.ContentType,
		SizeBytes:    m.SizeBytes,
		URL:          m.URL,
		ThumbnailURL: m.ThumbnailURL,
		Width:        m.Width,
		Height:       m.Height,
		Status:       domain.MediaStatus(m.Status),
		CreatedAt:    m.CreatedAt,
	}
}

// ToMediaModel converts a domain Media entity to a MediaModel.
func ToMediaModel(m *domain.Media) *MediaModel {
	return &MediaModel{
		ID:           m.ID,
		OwnerID:      m.OwnerID,
		OwnerType:    m.OwnerType,
		FileName:     m.FileName,
		OriginalName: m.OriginalName,
		ContentType:  m.ContentType,
		SizeBytes:    m.SizeBytes,
		URL:          m.URL,
		ThumbnailURL: m.ThumbnailURL,
		Width:        m.Width,
		Height:       m.Height,
		Status:       string(m.Status),
		CreatedAt:    m.CreatedAt,
	}
}
