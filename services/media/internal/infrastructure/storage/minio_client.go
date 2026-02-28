package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
	"github.com/southern-martin/ecommerce/services/media/internal/domain"
)

// MinIOClient implements domain.StorageClient using MinIO/S3.
type MinIOClient struct {
	client         *minio.Client
	bucket         string
	endpoint       string
	publicEndpoint string
}

// NewMinIOClient creates a new MinIO storage client and ensures the bucket exists.
func NewMinIOClient(endpoint, publicEndpoint, accessKey, secretKey, bucket, region string, useSSL bool) (*MinIOClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
		Region: region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: region}); err != nil {
			return nil, fmt.Errorf("failed to create bucket %s: %w", bucket, err)
		}
		log.Info().Str("bucket", bucket).Msg("created MinIO bucket")
	}

	if publicEndpoint == "" {
		publicEndpoint = endpoint
	}

	log.Info().Str("endpoint", endpoint).Str("bucket", bucket).Msg("connected to MinIO")
	return &MinIOClient{client: client, bucket: bucket, endpoint: endpoint, publicEndpoint: publicEndpoint}, nil
}

// GenerateUploadURL returns a presigned PUT URL.
func (c *MinIOClient) GenerateUploadURL(ctx context.Context, key, contentType string) (*domain.PresignedURL, error) {
	expiry := 15 * time.Minute
	url, err := c.client.PresignedPutObject(ctx, c.bucket, key, expiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned put URL: %w", err)
	}
	return &domain.PresignedURL{
		URL:       url.String(),
		Method:    "PUT",
		ExpiresAt: time.Now().Add(expiry),
	}, nil
}

// GenerateDownloadURL returns a presigned GET URL.
func (c *MinIOClient) GenerateDownloadURL(ctx context.Context, key string) (*domain.PresignedURL, error) {
	expiry := 1 * time.Hour
	url, err := c.client.PresignedGetObject(ctx, c.bucket, key, expiry, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned get URL: %w", err)
	}
	return &domain.PresignedURL{
		URL:       url.String(),
		Method:    "GET",
		ExpiresAt: time.Now().Add(expiry),
	}, nil
}

// DeleteObject removes an object from the bucket.
func (c *MinIOClient) DeleteObject(ctx context.Context, key string) error {
	return c.client.RemoveObject(ctx, c.bucket, key, minio.RemoveObjectOptions{})
}

// UploadFile uploads a file directly to MinIO.
func (c *MinIOClient) UploadFile(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	_, err := c.client.PutObject(ctx, c.bucket, key, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	return nil
}

// GetPublicURL returns the public URL for an object using the external-facing endpoint.
func (c *MinIOClient) GetPublicURL(key string) string {
	return fmt.Sprintf("http://%s/%s/%s", c.publicEndpoint, c.bucket, key)
}
