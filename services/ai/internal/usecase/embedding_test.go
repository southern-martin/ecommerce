package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/ai/internal/domain"
	"github.com/southern-martin/ecommerce/services/ai/internal/infrastructure/aiclient"
)

// ---------------------------------------------------------------------------
// Hand-written function-field mocks
// ---------------------------------------------------------------------------

// --- EmbeddingRepository mock ---

type mockEmbeddingRepo struct {
	getByEntityFn func(ctx context.Context, entityType domain.EntityType, entityID string) (*domain.Embedding, error)
	listByTypeFn  func(ctx context.Context, entityType domain.EntityType, page, pageSize int) ([]domain.Embedding, int64, error)
	createFn      func(ctx context.Context, embedding *domain.Embedding) error
	updateFn      func(ctx context.Context, embedding *domain.Embedding) error
	deleteFn      func(ctx context.Context, id string) error
}

func (m *mockEmbeddingRepo) GetByEntity(ctx context.Context, entityType domain.EntityType, entityID string) (*domain.Embedding, error) {
	if m.getByEntityFn != nil {
		return m.getByEntityFn(ctx, entityType, entityID)
	}
	return nil, errors.New("not found")
}
func (m *mockEmbeddingRepo) ListByType(ctx context.Context, entityType domain.EntityType, page, pageSize int) ([]domain.Embedding, int64, error) {
	if m.listByTypeFn != nil {
		return m.listByTypeFn(ctx, entityType, page, pageSize)
	}
	return nil, 0, nil
}
func (m *mockEmbeddingRepo) Create(ctx context.Context, embedding *domain.Embedding) error {
	if m.createFn != nil {
		return m.createFn(ctx, embedding)
	}
	return nil
}
func (m *mockEmbeddingRepo) Update(ctx context.Context, embedding *domain.Embedding) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, embedding)
	}
	return nil
}
func (m *mockEmbeddingRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

// --- EventPublisher mock ---

type mockAIEventPublisher struct {
	publishFn func(ctx context.Context, subject string, data interface{}) error
}

func (m *mockAIEventPublisher) Publish(ctx context.Context, subject string, data interface{}) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, subject, data)
	}
	return nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func newEmbeddingUseCase(
	embeddingRepo *mockEmbeddingRepo,
	aiClient *aiclient.MockAIClient,
	pub *mockAIEventPublisher,
) *EmbeddingUseCase {
	return NewEmbeddingUseCase(embeddingRepo, aiClient, pub)
}

// ===========================================================================
// GenerateEmbedding tests
// ===========================================================================

func TestGenerateEmbedding_Success(t *testing.T) {
	embRepo := &mockEmbeddingRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")
	pub := &mockAIEventPublisher{}

	var savedEmbedding *domain.Embedding
	embRepo.createFn = func(_ context.Context, e *domain.Embedding) error {
		savedEmbedding = e
		return nil
	}

	uc := newEmbeddingUseCase(embRepo, client, pub)
	emb, err := uc.GenerateEmbedding(context.Background(), GenerateEmbeddingRequest{
		EntityType: domain.EntityTypeProduct,
		EntityID:   "prod-1",
		Text:       "A great wireless mouse with ergonomic design",
	})

	require.NoError(t, err)
	require.NotNil(t, emb)
	assert.NotEmpty(t, emb.ID)
	assert.Equal(t, domain.EntityTypeProduct, emb.EntityType)
	assert.Equal(t, "prod-1", emb.EntityID)
	assert.Equal(t, "mock-v1", emb.ModelVersion)
	assert.Equal(t, 384, emb.Dimensions)
	assert.Len(t, emb.EmbeddingVector, 384)
	assert.NotNil(t, savedEmbedding)
}

func TestGenerateEmbedding_RepoCreateError(t *testing.T) {
	embRepo := &mockEmbeddingRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")
	pub := &mockAIEventPublisher{}

	embRepo.createFn = func(_ context.Context, _ *domain.Embedding) error {
		return errors.New("db write failed")
	}

	uc := newEmbeddingUseCase(embRepo, client, pub)
	_, err := uc.GenerateEmbedding(context.Background(), GenerateEmbeddingRequest{
		EntityType: domain.EntityTypeProduct,
		EntityID:   "prod-1",
		Text:       "Some text",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to store embedding")
}

func TestGenerateEmbedding_VectorDimensions(t *testing.T) {
	embRepo := &mockEmbeddingRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")
	pub := &mockAIEventPublisher{}
	embRepo.createFn = func(_ context.Context, _ *domain.Embedding) error { return nil }

	uc := newEmbeddingUseCase(embRepo, client, pub)
	emb, err := uc.GenerateEmbedding(context.Background(), GenerateEmbeddingRequest{
		EntityType: domain.EntityTypeCategory,
		EntityID:   "cat-1",
		Text:       "Electronics category",
	})

	require.NoError(t, err)
	assert.Equal(t, 384, len(emb.EmbeddingVector))
	assert.Equal(t, 384, emb.Dimensions)

	// Verify vector values are in [-1, 1] range
	for _, v := range emb.EmbeddingVector {
		assert.GreaterOrEqual(t, v, -1.0)
		assert.LessOrEqual(t, v, 1.0)
	}
}

func TestGenerateEmbedding_PublishEventFailDoesNotBreak(t *testing.T) {
	embRepo := &mockEmbeddingRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")
	pub := &mockAIEventPublisher{}

	embRepo.createFn = func(_ context.Context, _ *domain.Embedding) error { return nil }
	pub.publishFn = func(_ context.Context, _ string, _ interface{}) error {
		return errors.New("nats down")
	}

	uc := newEmbeddingUseCase(embRepo, client, pub)
	emb, err := uc.GenerateEmbedding(context.Background(), GenerateEmbeddingRequest{
		EntityType: domain.EntityTypeProduct,
		EntityID:   "prod-1",
		Text:       "Some text",
	})

	require.NoError(t, err)
	require.NotNil(t, emb)
}

func TestGenerateEmbedding_PublishesCorrectSubject(t *testing.T) {
	embRepo := &mockEmbeddingRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")
	pub := &mockAIEventPublisher{}

	embRepo.createFn = func(_ context.Context, _ *domain.Embedding) error { return nil }

	var publishedSubject string
	pub.publishFn = func(_ context.Context, subject string, _ interface{}) error {
		publishedSubject = subject
		return nil
	}

	uc := newEmbeddingUseCase(embRepo, client, pub)
	_, err := uc.GenerateEmbedding(context.Background(), GenerateEmbeddingRequest{
		EntityType: domain.EntityTypeProduct,
		EntityID:   "prod-1",
		Text:       "Some text",
	})

	require.NoError(t, err)
	assert.Equal(t, "ai.embedding.ready", publishedSubject)
}

// ===========================================================================
// GetEmbedding tests
// ===========================================================================

func TestGetEmbedding_Success(t *testing.T) {
	embRepo := &mockEmbeddingRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")
	pub := &mockAIEventPublisher{}

	expected := &domain.Embedding{
		ID:         "emb-1",
		EntityType: domain.EntityTypeProduct,
		EntityID:   "prod-1",
		Dimensions: 384,
	}
	embRepo.getByEntityFn = func(_ context.Context, entityType domain.EntityType, entityID string) (*domain.Embedding, error) {
		assert.Equal(t, domain.EntityTypeProduct, entityType)
		assert.Equal(t, "prod-1", entityID)
		return expected, nil
	}

	uc := newEmbeddingUseCase(embRepo, client, pub)
	emb, err := uc.GetEmbedding(context.Background(), domain.EntityTypeProduct, "prod-1")

	require.NoError(t, err)
	assert.Equal(t, expected, emb)
}

func TestGetEmbedding_NotFound(t *testing.T) {
	embRepo := &mockEmbeddingRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")
	pub := &mockAIEventPublisher{}

	uc := newEmbeddingUseCase(embRepo, client, pub)
	_, err := uc.GetEmbedding(context.Background(), domain.EntityTypeProduct, "nonexistent")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// ===========================================================================
// SearchSimilar tests
// ===========================================================================

func TestSearchSimilar_Success(t *testing.T) {
	embRepo := &mockEmbeddingRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")
	pub := &mockAIEventPublisher{}

	uc := newEmbeddingUseCase(embRepo, client, pub)
	results, err := uc.SearchSimilar(context.Background(), domain.EntityTypeProduct, "prod-1", 5)

	require.NoError(t, err)
	assert.Len(t, results, 5)

	// Scores should generally decrease
	for _, r := range results {
		assert.NotEmpty(t, r.EntityID)
		assert.Greater(t, r.Score, 0.0)
	}
}

func TestSearchSimilar_DefaultLimit(t *testing.T) {
	embRepo := &mockEmbeddingRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")
	pub := &mockAIEventPublisher{}

	uc := newEmbeddingUseCase(embRepo, client, pub)
	results, err := uc.SearchSimilar(context.Background(), domain.EntityTypeProduct, "prod-1", 0)

	require.NoError(t, err)
	assert.Len(t, results, 10) // default limit
}

func TestSearchSimilar_NegativeLimit(t *testing.T) {
	embRepo := &mockEmbeddingRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")
	pub := &mockAIEventPublisher{}

	uc := newEmbeddingUseCase(embRepo, client, pub)
	results, err := uc.SearchSimilar(context.Background(), domain.EntityTypeProduct, "prod-1", -5)

	require.NoError(t, err)
	assert.Len(t, results, 10) // default limit when negative
}
