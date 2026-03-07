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

// --- GeneratedContentRepository mock ---

type mockGeneratedContentRepo struct {
	getByEntityFn func(ctx context.Context, entityType domain.ContentType, entityID string) (*domain.GeneratedContent, error)
	createFn      func(ctx context.Context, content *domain.GeneratedContent) error
	updateFn      func(ctx context.Context, content *domain.GeneratedContent) error
}

func (m *mockGeneratedContentRepo) GetByEntity(ctx context.Context, entityType domain.ContentType, entityID string) (*domain.GeneratedContent, error) {
	if m.getByEntityFn != nil {
		return m.getByEntityFn(ctx, entityType, entityID)
	}
	return nil, errors.New("not found")
}
func (m *mockGeneratedContentRepo) Create(ctx context.Context, content *domain.GeneratedContent) error {
	if m.createFn != nil {
		return m.createFn(ctx, content)
	}
	return nil
}
func (m *mockGeneratedContentRepo) Update(ctx context.Context, content *domain.GeneratedContent) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, content)
	}
	return nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func newContentUseCase(
	contentRepo *mockGeneratedContentRepo,
	aiClient *aiclient.MockAIClient,
	pub *mockAIEventPublisher,
) *ContentUseCase {
	return NewContentUseCase(contentRepo, aiClient, pub)
}

// ===========================================================================
// GenerateDescription tests
// ===========================================================================

func TestGenerateDescription_Success(t *testing.T) {
	contentRepo := &mockGeneratedContentRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")
	pub := &mockAIEventPublisher{}

	var savedContent *domain.GeneratedContent
	contentRepo.createFn = func(_ context.Context, c *domain.GeneratedContent) error {
		savedContent = c
		return nil
	}

	uc := newContentUseCase(contentRepo, client, pub)
	content, err := uc.GenerateDescription(context.Background(), GenerateDescriptionRequest{
		ProductID:   "prod-1",
		ProductName: "Wireless Mouse",
		Category:    "Electronics",
	})

	require.NoError(t, err)
	require.NotNil(t, content)
	assert.NotEmpty(t, content.ID)
	assert.Equal(t, domain.ContentTypeProductDescription, content.EntityType)
	assert.Equal(t, "prod-1", content.EntityID)
	assert.NotEmpty(t, content.Content)
	assert.Equal(t, "gpt-4", content.Model)
	assert.Greater(t, content.PromptTokens, 0)
	assert.Greater(t, content.CompletionTokens, 0)
	assert.NotNil(t, savedContent)
	// Content should contain the product name
	assert.Contains(t, content.Content, "Wireless Mouse")
}

func TestGenerateDescription_RepoCreateError(t *testing.T) {
	contentRepo := &mockGeneratedContentRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")
	pub := &mockAIEventPublisher{}

	contentRepo.createFn = func(_ context.Context, _ *domain.GeneratedContent) error {
		return errors.New("db write failed")
	}

	uc := newContentUseCase(contentRepo, client, pub)
	_, err := uc.GenerateDescription(context.Background(), GenerateDescriptionRequest{
		ProductID:   "prod-1",
		ProductName: "Widget",
		Category:    "Tools",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to store generated content")
}

func TestGenerateDescription_PublishEventFailDoesNotBreak(t *testing.T) {
	contentRepo := &mockGeneratedContentRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")
	pub := &mockAIEventPublisher{}

	contentRepo.createFn = func(_ context.Context, _ *domain.GeneratedContent) error { return nil }
	pub.publishFn = func(_ context.Context, _ string, _ interface{}) error {
		return errors.New("nats down")
	}

	uc := newContentUseCase(contentRepo, client, pub)
	content, err := uc.GenerateDescription(context.Background(), GenerateDescriptionRequest{
		ProductID:   "prod-1",
		ProductName: "Widget",
		Category:    "Tools",
	})

	require.NoError(t, err)
	require.NotNil(t, content)
}

func TestGenerateDescription_PublishesCorrectSubject(t *testing.T) {
	contentRepo := &mockGeneratedContentRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")
	pub := &mockAIEventPublisher{}

	contentRepo.createFn = func(_ context.Context, _ *domain.GeneratedContent) error { return nil }

	var publishedSubject string
	pub.publishFn = func(_ context.Context, subject string, _ interface{}) error {
		publishedSubject = subject
		return nil
	}

	uc := newContentUseCase(contentRepo, client, pub)
	_, err := uc.GenerateDescription(context.Background(), GenerateDescriptionRequest{
		ProductID:   "prod-1",
		ProductName: "Widget",
		Category:    "Tools",
	})

	require.NoError(t, err)
	assert.Equal(t, "ai.description.generated", publishedSubject)
}

func TestGenerateDescription_TokenCounts(t *testing.T) {
	contentRepo := &mockGeneratedContentRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")
	pub := &mockAIEventPublisher{}
	contentRepo.createFn = func(_ context.Context, _ *domain.GeneratedContent) error { return nil }

	uc := newContentUseCase(contentRepo, client, pub)
	content, err := uc.GenerateDescription(context.Background(), GenerateDescriptionRequest{
		ProductID:   "prod-1",
		ProductName: "Wireless Mouse",
		Category:    "Electronics",
	})

	require.NoError(t, err)
	// promptTokens = (len("Wireless Mouse") + len("Electronics")) / 4
	expectedPromptTokens := (len("Wireless Mouse") + len("Electronics")) / 4
	assert.Equal(t, expectedPromptTokens, content.PromptTokens)
	// completionTokens = len(content.Content) / 4
	expectedCompletionTokens := len(content.Content) / 4
	assert.Equal(t, expectedCompletionTokens, content.CompletionTokens)
}

// ===========================================================================
// GetGeneratedContent tests
// ===========================================================================

func TestGetGeneratedContent_Success(t *testing.T) {
	contentRepo := &mockGeneratedContentRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")
	pub := &mockAIEventPublisher{}

	expected := &domain.GeneratedContent{
		ID:         "content-1",
		EntityType: domain.ContentTypeProductDescription,
		EntityID:   "prod-1",
		Content:    "A wonderful product",
		Model:      "gpt-4",
	}
	contentRepo.getByEntityFn = func(_ context.Context, entityType domain.ContentType, entityID string) (*domain.GeneratedContent, error) {
		assert.Equal(t, domain.ContentTypeProductDescription, entityType)
		assert.Equal(t, "prod-1", entityID)
		return expected, nil
	}

	uc := newContentUseCase(contentRepo, client, pub)
	content, err := uc.GetGeneratedContent(context.Background(), domain.ContentTypeProductDescription, "prod-1")

	require.NoError(t, err)
	assert.Equal(t, expected, content)
}

func TestGetGeneratedContent_NotFound(t *testing.T) {
	contentRepo := &mockGeneratedContentRepo{}
	client := aiclient.NewMockAIClient("http://localhost:8092")
	pub := &mockAIEventPublisher{}

	uc := newContentUseCase(contentRepo, client, pub)
	_, err := uc.GetGeneratedContent(context.Background(), domain.ContentTypeProductDescription, "nonexistent")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
