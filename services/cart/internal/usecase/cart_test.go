package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/cart/internal/domain"
)

// --- Mocks ---

type mockCartRepo struct {
	getCartFn    func(ctx context.Context, userID string) (*domain.Cart, error)
	saveCartFn   func(ctx context.Context, cart *domain.Cart) error
	deleteCartFn func(ctx context.Context, userID string) error
}

func (m *mockCartRepo) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	return m.getCartFn(ctx, userID)
}
func (m *mockCartRepo) SaveCart(ctx context.Context, cart *domain.Cart) error {
	return m.saveCartFn(ctx, cart)
}
func (m *mockCartRepo) DeleteCart(ctx context.Context, userID string) error {
	return m.deleteCartFn(ctx, userID)
}

type mockEventPublisher struct {
	publishFn func(ctx context.Context, subject string, payload interface{}) error
}

func (m *mockEventPublisher) Publish(ctx context.Context, subject string, payload interface{}) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, subject, payload)
	}
	return nil
}

func emptyCart(userID string) *domain.Cart {
	return &domain.Cart{UserID: userID, Items: []domain.CartItem{}}
}

func newTestUseCase(repo *mockCartRepo, pub *mockEventPublisher) *CartUseCase {
	logger := zerolog.Nop()
	return NewCartUseCase(repo, pub, logger)
}

// --- AddItem tests ---

func TestAddItem_NewItem(t *testing.T) {
	var savedCart *domain.Cart
	repo := &mockCartRepo{
		getCartFn:  func(_ context.Context, _ string) (*domain.Cart, error) { return emptyCart("u1"), nil },
		saveCartFn: func(_ context.Context, c *domain.Cart) error { savedCart = c; return nil },
	}
	uc := newTestUseCase(repo, &mockEventPublisher{})

	cart, err := uc.AddItem(context.Background(), "u1", domain.CartItem{
		ProductID: "p1", VariantID: "v1", Quantity: 2, PriceCents: 1000,
	})

	require.NoError(t, err)
	require.Len(t, cart.Items, 1)
	assert.Equal(t, "p1", cart.Items[0].ProductID)
	assert.Equal(t, 2, cart.Items[0].Quantity)
	assert.NotNil(t, savedCart)
}

func TestAddItem_IncrementExisting(t *testing.T) {
	existing := &domain.Cart{
		UserID: "u1",
		Items: []domain.CartItem{
			{ProductID: "p1", VariantID: "v1", Quantity: 3, PriceCents: 900, ProductName: "Old"},
		},
	}
	repo := &mockCartRepo{
		getCartFn:  func(_ context.Context, _ string) (*domain.Cart, error) { return existing, nil },
		saveCartFn: func(_ context.Context, _ *domain.Cart) error { return nil },
	}
	uc := newTestUseCase(repo, &mockEventPublisher{})

	cart, err := uc.AddItem(context.Background(), "u1", domain.CartItem{
		ProductID: "p1", VariantID: "v1", Quantity: 2, PriceCents: 1000, ProductName: "Updated",
	})

	require.NoError(t, err)
	require.Len(t, cart.Items, 1)
	assert.Equal(t, 5, cart.Items[0].Quantity)       // 3 + 2
	assert.Equal(t, int64(1000), cart.Items[0].PriceCents) // Updated
	assert.Equal(t, "Updated", cart.Items[0].ProductName)  // Denormalized field updated
}

func TestAddItem_EmptyUserID(t *testing.T) {
	uc := newTestUseCase(&mockCartRepo{}, &mockEventPublisher{})
	_, err := uc.AddItem(context.Background(), "", domain.CartItem{ProductID: "p1", Quantity: 1})
	assert.ErrorIs(t, err, ErrInvalidUserID)
}

func TestAddItem_EmptyProductID(t *testing.T) {
	uc := newTestUseCase(&mockCartRepo{}, &mockEventPublisher{})
	_, err := uc.AddItem(context.Background(), "u1", domain.CartItem{ProductID: "", Quantity: 1})
	assert.ErrorIs(t, err, ErrInvalidProduct)
}

func TestAddItem_ZeroQuantity(t *testing.T) {
	uc := newTestUseCase(&mockCartRepo{}, &mockEventPublisher{})
	_, err := uc.AddItem(context.Background(), "u1", domain.CartItem{ProductID: "p1", Quantity: 0})
	assert.ErrorIs(t, err, ErrInvalidQuantity)
}

func TestAddItem_NegativeQuantity(t *testing.T) {
	uc := newTestUseCase(&mockCartRepo{}, &mockEventPublisher{})
	_, err := uc.AddItem(context.Background(), "u1", domain.CartItem{ProductID: "p1", Quantity: -1})
	assert.ErrorIs(t, err, ErrInvalidQuantity)
}

func TestAddItem_RepoGetError(t *testing.T) {
	repo := &mockCartRepo{
		getCartFn: func(_ context.Context, _ string) (*domain.Cart, error) {
			return nil, errors.New("db down")
		},
	}
	uc := newTestUseCase(repo, &mockEventPublisher{})
	_, err := uc.AddItem(context.Background(), "u1", domain.CartItem{ProductID: "p1", Quantity: 1})
	assert.EqualError(t, err, "db down")
}

// --- RemoveItem tests ---

func TestRemoveItem_Success(t *testing.T) {
	existing := &domain.Cart{
		UserID: "u1",
		Items: []domain.CartItem{
			{ProductID: "p1", VariantID: "v1", Quantity: 2},
			{ProductID: "p2", VariantID: "v2", Quantity: 1},
		},
	}
	repo := &mockCartRepo{
		getCartFn:  func(_ context.Context, _ string) (*domain.Cart, error) { return existing, nil },
		saveCartFn: func(_ context.Context, _ *domain.Cart) error { return nil },
	}
	uc := newTestUseCase(repo, &mockEventPublisher{})

	cart, err := uc.RemoveItem(context.Background(), "u1", "p1", "v1")
	require.NoError(t, err)
	require.Len(t, cart.Items, 1)
	assert.Equal(t, "p2", cart.Items[0].ProductID)
}

func TestRemoveItem_NotFound(t *testing.T) {
	repo := &mockCartRepo{
		getCartFn: func(_ context.Context, _ string) (*domain.Cart, error) { return emptyCart("u1"), nil },
	}
	uc := newTestUseCase(repo, &mockEventPublisher{})
	_, err := uc.RemoveItem(context.Background(), "u1", "p99", "v99")
	assert.ErrorIs(t, err, ErrItemNotFound)
}

func TestRemoveItem_EmptyUserID(t *testing.T) {
	uc := newTestUseCase(&mockCartRepo{}, &mockEventPublisher{})
	_, err := uc.RemoveItem(context.Background(), "", "p1", "v1")
	assert.ErrorIs(t, err, ErrInvalidUserID)
}

// --- UpdateQuantity tests ---

func TestUpdateQuantity_Success(t *testing.T) {
	existing := &domain.Cart{
		UserID: "u1",
		Items:  []domain.CartItem{{ProductID: "p1", VariantID: "v1", Quantity: 2}},
	}
	repo := &mockCartRepo{
		getCartFn:  func(_ context.Context, _ string) (*domain.Cart, error) { return existing, nil },
		saveCartFn: func(_ context.Context, _ *domain.Cart) error { return nil },
	}
	uc := newTestUseCase(repo, &mockEventPublisher{})

	cart, err := uc.UpdateQuantity(context.Background(), "u1", "p1", "v1", 5)
	require.NoError(t, err)
	assert.Equal(t, 5, cart.Items[0].Quantity)
}

func TestUpdateQuantity_ItemNotFound(t *testing.T) {
	repo := &mockCartRepo{
		getCartFn: func(_ context.Context, _ string) (*domain.Cart, error) { return emptyCart("u1"), nil },
	}
	uc := newTestUseCase(repo, &mockEventPublisher{})
	_, err := uc.UpdateQuantity(context.Background(), "u1", "p99", "v99", 5)
	assert.ErrorIs(t, err, ErrItemNotFound)
}

func TestUpdateQuantity_ZeroQuantity(t *testing.T) {
	uc := newTestUseCase(&mockCartRepo{}, &mockEventPublisher{})
	_, err := uc.UpdateQuantity(context.Background(), "u1", "p1", "v1", 0)
	assert.ErrorIs(t, err, ErrInvalidQuantity)
}

// --- GetCart tests ---

func TestGetCart_Success(t *testing.T) {
	expected := &domain.Cart{UserID: "u1", Items: []domain.CartItem{{ProductID: "p1"}}}
	repo := &mockCartRepo{
		getCartFn: func(_ context.Context, _ string) (*domain.Cart, error) { return expected, nil },
	}
	uc := newTestUseCase(repo, &mockEventPublisher{})

	cart, err := uc.GetCart(context.Background(), "u1")
	require.NoError(t, err)
	assert.Equal(t, expected, cart)
}

func TestGetCart_EmptyUserID(t *testing.T) {
	uc := newTestUseCase(&mockCartRepo{}, &mockEventPublisher{})
	_, err := uc.GetCart(context.Background(), "")
	assert.ErrorIs(t, err, ErrInvalidUserID)
}

// --- ClearCart tests ---

func TestClearCart_Success(t *testing.T) {
	var deletedUserID string
	repo := &mockCartRepo{
		deleteCartFn: func(_ context.Context, userID string) error { deletedUserID = userID; return nil },
	}
	uc := newTestUseCase(repo, &mockEventPublisher{})

	err := uc.ClearCart(context.Background(), "u1")
	require.NoError(t, err)
	assert.Equal(t, "u1", deletedUserID)
}

func TestClearCart_EmptyUserID(t *testing.T) {
	uc := newTestUseCase(&mockCartRepo{}, &mockEventPublisher{})
	err := uc.ClearCart(context.Background(), "")
	assert.ErrorIs(t, err, ErrInvalidUserID)
}

// --- MergeCart tests ---

func TestMergeCart_NewItems(t *testing.T) {
	repo := &mockCartRepo{
		getCartFn:  func(_ context.Context, _ string) (*domain.Cart, error) { return emptyCart("u1"), nil },
		saveCartFn: func(_ context.Context, _ *domain.Cart) error { return nil },
	}
	uc := newTestUseCase(repo, &mockEventPublisher{})

	guestItems := []domain.CartItem{
		{ProductID: "p1", VariantID: "v1", Quantity: 2, PriceCents: 1000},
		{ProductID: "p2", VariantID: "v2", Quantity: 1, PriceCents: 500},
	}
	cart, err := uc.MergeCart(context.Background(), "u1", guestItems)
	require.NoError(t, err)
	assert.Len(t, cart.Items, 2)
}

func TestMergeCart_OverlappingItems(t *testing.T) {
	existing := &domain.Cart{
		UserID: "u1",
		Items:  []domain.CartItem{{ProductID: "p1", VariantID: "v1", Quantity: 3, PriceCents: 900}},
	}
	repo := &mockCartRepo{
		getCartFn:  func(_ context.Context, _ string) (*domain.Cart, error) { return existing, nil },
		saveCartFn: func(_ context.Context, _ *domain.Cart) error { return nil },
	}
	uc := newTestUseCase(repo, &mockEventPublisher{})

	guestItems := []domain.CartItem{
		{ProductID: "p1", VariantID: "v1", Quantity: 2, PriceCents: 1000},
	}
	cart, err := uc.MergeCart(context.Background(), "u1", guestItems)
	require.NoError(t, err)
	require.Len(t, cart.Items, 1)
	assert.Equal(t, 5, cart.Items[0].Quantity) // 3 + 2
	assert.Equal(t, int64(1000), cart.Items[0].PriceCents) // Updated
}

func TestMergeCart_SkipsInvalidItems(t *testing.T) {
	repo := &mockCartRepo{
		getCartFn:  func(_ context.Context, _ string) (*domain.Cart, error) { return emptyCart("u1"), nil },
		saveCartFn: func(_ context.Context, _ *domain.Cart) error { return nil },
	}
	uc := newTestUseCase(repo, &mockEventPublisher{})

	guestItems := []domain.CartItem{
		{ProductID: "", Quantity: 2},          // Empty product ID
		{ProductID: "p1", Quantity: 0},        // Zero quantity
		{ProductID: "p2", Quantity: -1},       // Negative quantity
		{ProductID: "p3", Quantity: 1, PriceCents: 500}, // Valid
	}
	cart, err := uc.MergeCart(context.Background(), "u1", guestItems)
	require.NoError(t, err)
	assert.Len(t, cart.Items, 1) // Only the valid item
	assert.Equal(t, "p3", cart.Items[0].ProductID)
}

func TestMergeCart_EmptyUserID(t *testing.T) {
	uc := newTestUseCase(&mockCartRepo{}, &mockEventPublisher{})
	_, err := uc.MergeCart(context.Background(), "", nil)
	assert.ErrorIs(t, err, ErrInvalidUserID)
}
