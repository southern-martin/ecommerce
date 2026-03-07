package usecase

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/media/internal/domain"
)

// ---------------------------------------------------------------------------
// Hand-written function-field mocks
// ---------------------------------------------------------------------------

// --- MediaRepository mock ---

type mockMediaRepo struct {
	getByIDFn     func(ctx context.Context, id string) (*domain.Media, error)
	listByOwnerFn func(ctx context.Context, ownerID, ownerType string, page, pageSize int) ([]domain.Media, int64, error)
	createFn      func(ctx context.Context, media *domain.Media) error
	updateFn      func(ctx context.Context, media *domain.Media) error
	deleteFn      func(ctx context.Context, id string) error
}

func (m *mockMediaRepo) GetByID(ctx context.Context, id string) (*domain.Media, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}
func (m *mockMediaRepo) ListByOwner(ctx context.Context, ownerID, ownerType string, page, pageSize int) ([]domain.Media, int64, error) {
	if m.listByOwnerFn != nil {
		return m.listByOwnerFn(ctx, ownerID, ownerType, page, pageSize)
	}
	return nil, 0, nil
}
func (m *mockMediaRepo) Create(ctx context.Context, media *domain.Media) error {
	if m.createFn != nil {
		return m.createFn(ctx, media)
	}
	return nil
}
func (m *mockMediaRepo) Update(ctx context.Context, media *domain.Media) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, media)
	}
	return nil
}
func (m *mockMediaRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

// --- StorageClient mock ---

type mockStorageClient struct {
	generateUploadURLFn   func(ctx context.Context, key, contentType string) (*domain.PresignedURL, error)
	generateDownloadURLFn func(ctx context.Context, key string) (*domain.PresignedURL, error)
	deleteObjectFn        func(ctx context.Context, key string) error
	uploadFileFn          func(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error
	getPublicURLFn        func(key string) string
}

func (m *mockStorageClient) GenerateUploadURL(ctx context.Context, key, contentType string) (*domain.PresignedURL, error) {
	if m.generateUploadURLFn != nil {
		return m.generateUploadURLFn(ctx, key, contentType)
	}
	return &domain.PresignedURL{URL: "https://storage.example.com/upload/" + key, Method: "PUT"}, nil
}
func (m *mockStorageClient) GenerateDownloadURL(ctx context.Context, key string) (*domain.PresignedURL, error) {
	if m.generateDownloadURLFn != nil {
		return m.generateDownloadURLFn(ctx, key)
	}
	return &domain.PresignedURL{URL: "https://storage.example.com/download/" + key, Method: "GET"}, nil
}
func (m *mockStorageClient) DeleteObject(ctx context.Context, key string) error {
	if m.deleteObjectFn != nil {
		return m.deleteObjectFn(ctx, key)
	}
	return nil
}
func (m *mockStorageClient) UploadFile(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	if m.uploadFileFn != nil {
		return m.uploadFileFn(ctx, key, reader, size, contentType)
	}
	return nil
}
func (m *mockStorageClient) GetPublicURL(key string) string {
	if m.getPublicURLFn != nil {
		return m.getPublicURLFn(key)
	}
	return "https://cdn.example.com/" + key
}

// --- EventPublisher mock ---

type mockEventPublisher struct {
	publishFn func(ctx context.Context, subject string, data interface{}) error
}

func (m *mockEventPublisher) Publish(ctx context.Context, subject string, data interface{}) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, subject, data)
	}
	return nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func newMediaUseCase(
	repo *mockMediaRepo,
	storage *mockStorageClient,
	pub *mockEventPublisher,
) *MediaUseCase {
	return NewMediaUseCase(repo, storage, pub)
}

func defaultMediaMocks() (*mockMediaRepo, *mockStorageClient, *mockEventPublisher) {
	return &mockMediaRepo{}, &mockStorageClient{}, &mockEventPublisher{}
}

// ===========================================================================
// CreateMedia tests
// ===========================================================================

func TestCreateMedia_Success(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()

	var savedMedia *domain.Media
	repo.createFn = func(_ context.Context, m *domain.Media) error {
		savedMedia = m
		return nil
	}

	uc := newMediaUseCase(repo, storage, pub)
	resp, err := uc.CreateMedia(context.Background(), CreateMediaRequest{
		OwnerID:      "user-1",
		OwnerType:    "product",
		OriginalName: "photo.jpg",
		ContentType:  "image/jpeg",
		SizeBytes:    1024000,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Media)
	require.NotNil(t, resp.UploadURL)
	assert.NotEmpty(t, resp.Media.ID)
	assert.Equal(t, "user-1", resp.Media.OwnerID)
	assert.Equal(t, "product", resp.Media.OwnerType)
	assert.Equal(t, "photo.jpg", resp.Media.OriginalName)
	assert.Equal(t, "image/jpeg", resp.Media.ContentType)
	assert.Equal(t, int64(1024000), resp.Media.SizeBytes)
	assert.Equal(t, domain.MediaStatusPending, resp.Media.Status)
	assert.NotNil(t, savedMedia)
	assert.Contains(t, resp.UploadURL.URL, "upload")
}

func TestCreateMedia_RepoCreateError(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()
	repo.createFn = func(_ context.Context, _ *domain.Media) error {
		return errors.New("db write failed")
	}

	uc := newMediaUseCase(repo, storage, pub)
	_, err := uc.CreateMedia(context.Background(), CreateMediaRequest{
		OwnerID:      "user-1",
		OwnerType:    "product",
		OriginalName: "photo.jpg",
		ContentType:  "image/jpeg",
		SizeBytes:    1024,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "db write failed")
}

func TestCreateMedia_StorageUploadURLError(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()
	repo.createFn = func(_ context.Context, _ *domain.Media) error { return nil }
	storage.generateUploadURLFn = func(_ context.Context, _, _ string) (*domain.PresignedURL, error) {
		return nil, errors.New("storage unavailable")
	}

	uc := newMediaUseCase(repo, storage, pub)
	_, err := uc.CreateMedia(context.Background(), CreateMediaRequest{
		OwnerID:      "user-1",
		OwnerType:    "product",
		OriginalName: "photo.jpg",
		ContentType:  "image/jpeg",
		SizeBytes:    1024,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "storage unavailable")
}

func TestCreateMedia_FileNameFormat(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()

	var savedFileName string
	repo.createFn = func(_ context.Context, m *domain.Media) error {
		savedFileName = m.FileName
		return nil
	}

	uc := newMediaUseCase(repo, storage, pub)
	resp, err := uc.CreateMedia(context.Background(), CreateMediaRequest{
		OwnerID:      "user-1",
		OwnerType:    "product",
		OriginalName: "photo.jpg",
		ContentType:  "image/jpeg",
		SizeBytes:    1024,
	})

	require.NoError(t, err)
	// FileName should be "product/user-1/<uuid>"
	assert.Contains(t, savedFileName, "product/user-1/")
	assert.Equal(t, savedFileName, resp.Media.FileName)
}

func TestCreateMedia_PublishEventFailDoesNotBreak(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()
	repo.createFn = func(_ context.Context, _ *domain.Media) error { return nil }
	pub.publishFn = func(_ context.Context, _ string, _ interface{}) error {
		return errors.New("nats down")
	}

	uc := newMediaUseCase(repo, storage, pub)
	resp, err := uc.CreateMedia(context.Background(), CreateMediaRequest{
		OwnerID:      "user-1",
		OwnerType:    "product",
		OriginalName: "photo.jpg",
		ContentType:  "image/jpeg",
		SizeBytes:    1024,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
}

// ===========================================================================
// GetMedia tests
// ===========================================================================

func TestGetMedia_Success(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()
	expected := &domain.Media{
		ID:          "media-1",
		OwnerID:     "user-1",
		FileName:    "product/user-1/media-1",
		ContentType: "image/jpeg",
	}
	repo.getByIDFn = func(_ context.Context, id string) (*domain.Media, error) {
		assert.Equal(t, "media-1", id)
		return expected, nil
	}

	uc := newMediaUseCase(repo, storage, pub)
	media, err := uc.GetMedia(context.Background(), "media-1")

	require.NoError(t, err)
	assert.Equal(t, expected, media)
}

func TestGetMedia_NotFound(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()
	repo.getByIDFn = func(_ context.Context, _ string) (*domain.Media, error) {
		return nil, errors.New("not found")
	}

	uc := newMediaUseCase(repo, storage, pub)
	_, err := uc.GetMedia(context.Background(), "nonexistent")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// ===========================================================================
// ListMedia tests
// ===========================================================================

func TestListMedia_Success(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()
	expectedMedia := []domain.Media{
		{ID: "media-1", OwnerID: "user-1"},
		{ID: "media-2", OwnerID: "user-1"},
	}
	repo.listByOwnerFn = func(_ context.Context, ownerID, ownerType string, page, pageSize int) ([]domain.Media, int64, error) {
		assert.Equal(t, "user-1", ownerID)
		assert.Equal(t, "product", ownerType)
		assert.Equal(t, 1, page)
		assert.Equal(t, 20, pageSize)
		return expectedMedia, 2, nil
	}

	uc := newMediaUseCase(repo, storage, pub)
	media, total, err := uc.ListMedia(context.Background(), "user-1", "product", 1, 20)

	require.NoError(t, err)
	assert.Len(t, media, 2)
	assert.Equal(t, int64(2), total)
}

func TestListMedia_DefaultPagination(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()

	var capturedPage, capturedPageSize int
	repo.listByOwnerFn = func(_ context.Context, _, _ string, page, pageSize int) ([]domain.Media, int64, error) {
		capturedPage = page
		capturedPageSize = pageSize
		return nil, 0, nil
	}

	uc := newMediaUseCase(repo, storage, pub)

	// page < 1 should default to 1
	_, _, _ = uc.ListMedia(context.Background(), "user-1", "product", 0, 20)
	assert.Equal(t, 1, capturedPage)
	assert.Equal(t, 20, capturedPageSize)

	// pageSize > 100 should default to 20
	_, _, _ = uc.ListMedia(context.Background(), "user-1", "product", 1, 200)
	assert.Equal(t, 20, capturedPageSize)

	// pageSize < 1 should default to 20
	_, _, _ = uc.ListMedia(context.Background(), "user-1", "product", 1, 0)
	assert.Equal(t, 20, capturedPageSize)
}

// ===========================================================================
// DeleteMedia tests
// ===========================================================================

func TestDeleteMedia_Success(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()
	repo.getByIDFn = func(_ context.Context, _ string) (*domain.Media, error) {
		return &domain.Media{
			ID:       "media-1",
			OwnerID:  "user-1",
			FileName: "product/user-1/media-1",
		}, nil
	}

	var deletedKey string
	storage.deleteObjectFn = func(_ context.Context, key string) error {
		deletedKey = key
		return nil
	}

	var deletedID string
	repo.deleteFn = func(_ context.Context, id string) error {
		deletedID = id
		return nil
	}

	uc := newMediaUseCase(repo, storage, pub)
	err := uc.DeleteMedia(context.Background(), "media-1")

	require.NoError(t, err)
	assert.Equal(t, "product/user-1/media-1", deletedKey)
	assert.Equal(t, "media-1", deletedID)
}

func TestDeleteMedia_NotFound(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()
	repo.getByIDFn = func(_ context.Context, _ string) (*domain.Media, error) {
		return nil, errors.New("not found")
	}

	uc := newMediaUseCase(repo, storage, pub)
	err := uc.DeleteMedia(context.Background(), "nonexistent")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestDeleteMedia_StorageDeleteError(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()
	repo.getByIDFn = func(_ context.Context, _ string) (*domain.Media, error) {
		return &domain.Media{ID: "media-1", FileName: "key"}, nil
	}
	storage.deleteObjectFn = func(_ context.Context, _ string) error {
		return errors.New("storage unavailable")
	}

	uc := newMediaUseCase(repo, storage, pub)
	err := uc.DeleteMedia(context.Background(), "media-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "storage unavailable")
}

func TestDeleteMedia_RepoDeleteError(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()
	repo.getByIDFn = func(_ context.Context, _ string) (*domain.Media, error) {
		return &domain.Media{ID: "media-1", FileName: "key"}, nil
	}
	repo.deleteFn = func(_ context.Context, _ string) error {
		return errors.New("db delete failed")
	}

	uc := newMediaUseCase(repo, storage, pub)
	err := uc.DeleteMedia(context.Background(), "media-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "db delete failed")
}

func TestDeleteMedia_PublishEventFailDoesNotBreak(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()
	repo.getByIDFn = func(_ context.Context, _ string) (*domain.Media, error) {
		return &domain.Media{ID: "media-1", OwnerID: "user-1", FileName: "key"}, nil
	}
	pub.publishFn = func(_ context.Context, _ string, _ interface{}) error {
		return errors.New("nats down")
	}

	uc := newMediaUseCase(repo, storage, pub)
	err := uc.DeleteMedia(context.Background(), "media-1")

	require.NoError(t, err)
}

// ===========================================================================
// UploadMedia tests
// ===========================================================================

func TestUploadMedia_Success(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()

	var savedMedia *domain.Media
	repo.createFn = func(_ context.Context, m *domain.Media) error {
		savedMedia = m
		return nil
	}

	var uploadedKey string
	storage.uploadFileFn = func(_ context.Context, key string, _ io.Reader, _ int64, _ string) error {
		uploadedKey = key
		return nil
	}

	reader := bytes.NewReader([]byte("fake file content"))
	uc := newMediaUseCase(repo, storage, pub)
	media, err := uc.UploadMedia(context.Background(), UploadMediaRequest{
		OwnerID:      "user-1",
		OwnerType:    "product",
		OriginalName: "photo.jpg",
		ContentType:  "image/jpeg",
		SizeBytes:    17,
		Reader:       reader,
	})

	require.NoError(t, err)
	require.NotNil(t, media)
	assert.NotEmpty(t, media.ID)
	assert.Equal(t, "user-1", media.OwnerID)
	assert.Equal(t, "product", media.OwnerType)
	assert.Equal(t, "photo.jpg", media.OriginalName)
	assert.Equal(t, "image/jpeg", media.ContentType)
	assert.Equal(t, domain.MediaStatusProcessed, media.Status)
	assert.Contains(t, media.URL, "cdn.example.com")
	assert.NotEmpty(t, uploadedKey)
	assert.NotNil(t, savedMedia)
}

func TestUploadMedia_FileExtensionPreserved(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()
	repo.createFn = func(_ context.Context, _ *domain.Media) error { return nil }

	var uploadedKey string
	storage.uploadFileFn = func(_ context.Context, key string, _ io.Reader, _ int64, _ string) error {
		uploadedKey = key
		return nil
	}

	reader := bytes.NewReader([]byte("content"))
	uc := newMediaUseCase(repo, storage, pub)
	_, err := uc.UploadMedia(context.Background(), UploadMediaRequest{
		OwnerID:      "user-1",
		OwnerType:    "product",
		OriginalName: "document.pdf",
		ContentType:  "application/pdf",
		SizeBytes:    7,
		Reader:       reader,
	})

	require.NoError(t, err)
	// key should end with .pdf
	assert.Contains(t, uploadedKey, ".pdf")
	assert.Contains(t, uploadedKey, "product/user-1/")
}

func TestUploadMedia_StorageUploadError(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()
	storage.uploadFileFn = func(_ context.Context, _ string, _ io.Reader, _ int64, _ string) error {
		return errors.New("storage upload failed")
	}

	reader := bytes.NewReader([]byte("content"))
	uc := newMediaUseCase(repo, storage, pub)
	_, err := uc.UploadMedia(context.Background(), UploadMediaRequest{
		OwnerID:      "user-1",
		OwnerType:    "product",
		OriginalName: "photo.jpg",
		ContentType:  "image/jpeg",
		SizeBytes:    7,
		Reader:       reader,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "storage upload failed")
}

func TestUploadMedia_RepoCreateError(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()
	repo.createFn = func(_ context.Context, _ *domain.Media) error {
		return errors.New("db write failed")
	}

	reader := bytes.NewReader([]byte("content"))
	uc := newMediaUseCase(repo, storage, pub)
	_, err := uc.UploadMedia(context.Background(), UploadMediaRequest{
		OwnerID:      "user-1",
		OwnerType:    "product",
		OriginalName: "photo.jpg",
		ContentType:  "image/jpeg",
		SizeBytes:    7,
		Reader:       reader,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "db write failed")
}

// ===========================================================================
// GenerateUploadURL tests
// ===========================================================================

func TestGenerateUploadURL_Success(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()
	expected := &domain.PresignedURL{URL: "https://storage.example.com/upload/my-key", Method: "PUT"}
	storage.generateUploadURLFn = func(_ context.Context, key, ct string) (*domain.PresignedURL, error) {
		assert.Equal(t, "my-key", key)
		assert.Equal(t, "image/png", ct)
		return expected, nil
	}

	uc := newMediaUseCase(repo, storage, pub)
	url, err := uc.GenerateUploadURL(context.Background(), "my-key", "image/png")

	require.NoError(t, err)
	assert.Equal(t, expected, url)
}

func TestGenerateUploadURL_Error(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()
	storage.generateUploadURLFn = func(_ context.Context, _, _ string) (*domain.PresignedURL, error) {
		return nil, errors.New("storage error")
	}

	uc := newMediaUseCase(repo, storage, pub)
	_, err := uc.GenerateUploadURL(context.Background(), "key", "image/png")

	require.Error(t, err)
}

// ===========================================================================
// GenerateDownloadURL tests
// ===========================================================================

func TestGenerateDownloadURL_Success(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()
	repo.getByIDFn = func(_ context.Context, id string) (*domain.Media, error) {
		return &domain.Media{ID: id, FileName: "product/user-1/media-1"}, nil
	}
	storage.generateDownloadURLFn = func(_ context.Context, key string) (*domain.PresignedURL, error) {
		assert.Equal(t, "product/user-1/media-1", key)
		return &domain.PresignedURL{URL: "https://storage.example.com/download/" + key, Method: "GET"}, nil
	}

	uc := newMediaUseCase(repo, storage, pub)
	url, err := uc.GenerateDownloadURL(context.Background(), "media-1")

	require.NoError(t, err)
	assert.Contains(t, url.URL, "download")
	assert.Equal(t, "GET", url.Method)
}

func TestGenerateDownloadURL_NotFound(t *testing.T) {
	repo, storage, pub := defaultMediaMocks()
	repo.getByIDFn = func(_ context.Context, _ string) (*domain.Media, error) {
		return nil, errors.New("not found")
	}

	uc := newMediaUseCase(repo, storage, pub)
	_, err := uc.GenerateDownloadURL(context.Background(), "nonexistent")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
