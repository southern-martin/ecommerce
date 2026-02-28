package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/southern-martin/ecommerce/services/media/internal/domain"
)

// MockStorageClient implements domain.StorageClient with mock responses.
type MockStorageClient struct {
	endpoint string
	bucket   string
}

// NewMockStorageClient creates a new MockStorageClient.
func NewMockStorageClient(endpoint, bucket string) *MockStorageClient {
	return &MockStorageClient{
		endpoint: endpoint,
		bucket:   bucket,
	}
}

// GenerateUploadURL returns a mock presigned upload URL.
func (c *MockStorageClient) GenerateUploadURL(ctx context.Context, key, contentType string) (*domain.PresignedURL, error) {
	return &domain.PresignedURL{
		URL:       fmt.Sprintf("http://%s/%s/%s", c.endpoint, c.bucket, key),
		Method:    "PUT",
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}, nil
}

// GenerateDownloadURL returns a mock presigned download URL.
func (c *MockStorageClient) GenerateDownloadURL(ctx context.Context, key string) (*domain.PresignedURL, error) {
	return &domain.PresignedURL{
		URL:       fmt.Sprintf("http://%s/%s/%s", c.endpoint, c.bucket, key),
		Method:    "GET",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}, nil
}

// DeleteObject is a no-op for the mock client.
func (c *MockStorageClient) DeleteObject(ctx context.Context, key string) error {
	return nil
}

// UploadFile is a no-op for the mock client.
func (c *MockStorageClient) UploadFile(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	return nil
}

// GetPublicURL returns a mock public URL.
func (c *MockStorageClient) GetPublicURL(key string) string {
	return fmt.Sprintf("http://%s/%s/%s", c.endpoint, c.bucket, key)
}
