package domain

import "time"

// Media represents a media file entity.
type Media struct {
	ID           string      `json:"id"`
	OwnerID      string      `json:"owner_id"`
	OwnerType    string      `json:"owner_type"`
	FileName     string      `json:"file_name"`
	OriginalName string      `json:"original_name"`
	ContentType  string      `json:"content_type"`
	SizeBytes    int64       `json:"size_bytes"`
	URL          string      `json:"url"`
	ThumbnailURL string      `json:"thumbnail_url"`
	Width        int         `json:"width"`
	Height       int         `json:"height"`
	Status       MediaStatus `json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
}

// MediaStatus represents the processing status of a media file.
type MediaStatus string

const (
	MediaStatusPending   MediaStatus = "pending"
	MediaStatusProcessed MediaStatus = "processed"
	MediaStatusFailed    MediaStatus = "failed"
)

// PresignedURL represents a presigned URL for upload or download.
type PresignedURL struct {
	URL       string    `json:"url"`
	Method    string    `json:"method"`
	ExpiresAt time.Time `json:"expires_at"`
}

// MediaFilter holds filter criteria for listing media.
type MediaFilter struct {
	OwnerID     string `json:"owner_id"`
	OwnerType   string `json:"owner_type"`
	ContentType string `json:"content_type"`
	Page        int    `json:"page"`
	PageSize    int    `json:"page_size"`
}
