package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/southern-martin/ecommerce/services/media/internal/domain"
)

// MediaUseCase handles media business logic.
type MediaUseCase struct {
	repo      domain.MediaRepository
	storage   domain.StorageClient
	publisher domain.EventPublisher
}

// NewMediaUseCase creates a new MediaUseCase.
func NewMediaUseCase(repo domain.MediaRepository, storage domain.StorageClient, publisher domain.EventPublisher) *MediaUseCase {
	return &MediaUseCase{
		repo:      repo,
		storage:   storage,
		publisher: publisher,
	}
}

// CreateMediaRequest holds the parameters for creating media metadata.
type CreateMediaRequest struct {
	OwnerID      string `json:"owner_id"`
	OwnerType    string `json:"owner_type"`
	OriginalName string `json:"original_name"`
	ContentType  string `json:"content_type"`
	SizeBytes    int64  `json:"size_bytes"`
}

// CreateMediaResponse holds the created media and its upload URL.
type CreateMediaResponse struct {
	Media     *domain.Media      `json:"media"`
	UploadURL *domain.PresignedURL `json:"upload_url"`
}

// CreateMedia creates media metadata and generates an upload URL.
func (uc *MediaUseCase) CreateMedia(ctx context.Context, req CreateMediaRequest) (*CreateMediaResponse, error) {
	id := uuid.New().String()
	fileName := fmt.Sprintf("%s/%s/%s", req.OwnerType, req.OwnerID, id)

	media := &domain.Media{
		ID:           id,
		OwnerID:      req.OwnerID,
		OwnerType:    req.OwnerType,
		FileName:     fileName,
		OriginalName: req.OriginalName,
		ContentType:  req.ContentType,
		SizeBytes:    req.SizeBytes,
		Status:       domain.MediaStatusPending,
		CreatedAt:    time.Now(),
	}

	if err := uc.repo.Create(ctx, media); err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to create media")
		return nil, err
	}

	uploadURL, err := uc.storage.GenerateUploadURL(ctx, fileName, req.ContentType)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to generate upload URL")
		return nil, err
	}

	_ = uc.publisher.Publish(ctx, "media.created", map[string]string{
		"media_id":   id,
		"owner_id":   req.OwnerID,
		"owner_type": req.OwnerType,
	})

	return &CreateMediaResponse{
		Media:     media,
		UploadURL: uploadURL,
	}, nil
}

// GetMedia retrieves a single media by ID.
func (uc *MediaUseCase) GetMedia(ctx context.Context, id string) (*domain.Media, error) {
	return uc.repo.GetByID(ctx, id)
}

// ListMedia retrieves media by owner with pagination.
func (uc *MediaUseCase) ListMedia(ctx context.Context, ownerID, ownerType string, page, pageSize int) ([]domain.Media, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.repo.ListByOwner(ctx, ownerID, ownerType, page, pageSize)
}

// DeleteMedia deletes media from storage and database.
func (uc *MediaUseCase) DeleteMedia(ctx context.Context, id string) error {
	media, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := uc.storage.DeleteObject(ctx, media.FileName); err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to delete object from storage")
		return err
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to delete media from database")
		return err
	}

	_ = uc.publisher.Publish(ctx, "media.deleted", map[string]string{
		"media_id":   id,
		"owner_id":   media.OwnerID,
		"owner_type": media.OwnerType,
	})

	return nil
}

// GenerateUploadURL generates a presigned upload URL for a given key.
func (uc *MediaUseCase) GenerateUploadURL(ctx context.Context, key, contentType string) (*domain.PresignedURL, error) {
	return uc.storage.GenerateUploadURL(ctx, key, contentType)
}

// GenerateDownloadURL generates a presigned download URL for a media file.
func (uc *MediaUseCase) GenerateDownloadURL(ctx context.Context, id string) (*domain.PresignedURL, error) {
	media, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return uc.storage.GenerateDownloadURL(ctx, media.FileName)
}
