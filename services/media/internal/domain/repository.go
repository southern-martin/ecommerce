package domain

import (
	"context"
	"io"
)

// MediaRepository defines the interface for media persistence.
type MediaRepository interface {
	GetByID(ctx context.Context, id string) (*Media, error)
	ListByOwner(ctx context.Context, ownerID, ownerType string, page, pageSize int) ([]Media, int64, error)
	Create(ctx context.Context, media *Media) error
	Update(ctx context.Context, media *Media) error
	Delete(ctx context.Context, id string) error
}

// StorageClient defines the interface for object storage operations.
type StorageClient interface {
	GenerateUploadURL(ctx context.Context, key, contentType string) (*PresignedURL, error)
	GenerateDownloadURL(ctx context.Context, key string) (*PresignedURL, error)
	DeleteObject(ctx context.Context, key string) error
	UploadFile(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error
	GetPublicURL(key string) string
}
