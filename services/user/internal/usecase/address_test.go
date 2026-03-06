package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pkgerrors "github.com/southern-martin/ecommerce/pkg/errors"
	"github.com/southern-martin/ecommerce/services/user/internal/domain"
)

// ---------------------------------------------------------------------------
// Shared mocks (placed here because address_test.go is first alphabetically)
// ---------------------------------------------------------------------------

// mockUserProfileRepo mocks domain.UserProfileRepository.
type mockUserProfileRepo struct {
	createFn  func(ctx context.Context, profile *domain.UserProfile) error
	getByIDFn func(ctx context.Context, id string) (*domain.UserProfile, error)
	updateFn  func(ctx context.Context, profile *domain.UserProfile) error
}

func (m *mockUserProfileRepo) Create(ctx context.Context, profile *domain.UserProfile) error {
	if m.createFn != nil {
		return m.createFn(ctx, profile)
	}
	return nil
}
func (m *mockUserProfileRepo) GetByID(ctx context.Context, id string) (*domain.UserProfile, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}
func (m *mockUserProfileRepo) Update(ctx context.Context, profile *domain.UserProfile) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, profile)
	}
	return nil
}

// mockAddressRepo mocks domain.AddressRepository.
type mockAddressRepo struct {
	createFn             func(ctx context.Context, addr *domain.Address) error
	getByIDFn            func(ctx context.Context, id string) (*domain.Address, error)
	listByUserIDFn       func(ctx context.Context, userID string) ([]domain.Address, error)
	updateFn             func(ctx context.Context, addr *domain.Address) error
	deleteFn             func(ctx context.Context, id string) error
	countByUserIDFn      func(ctx context.Context, userID string) (int64, error)
	clearDefaultByUserFn func(ctx context.Context, userID string) error
}

func (m *mockAddressRepo) Create(ctx context.Context, addr *domain.Address) error {
	if m.createFn != nil {
		return m.createFn(ctx, addr)
	}
	return nil
}
func (m *mockAddressRepo) GetByID(ctx context.Context, id string) (*domain.Address, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}
func (m *mockAddressRepo) ListByUserID(ctx context.Context, userID string) ([]domain.Address, error) {
	if m.listByUserIDFn != nil {
		return m.listByUserIDFn(ctx, userID)
	}
	return nil, nil
}
func (m *mockAddressRepo) Update(ctx context.Context, addr *domain.Address) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, addr)
	}
	return nil
}
func (m *mockAddressRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}
func (m *mockAddressRepo) CountByUserID(ctx context.Context, userID string) (int64, error) {
	if m.countByUserIDFn != nil {
		return m.countByUserIDFn(ctx, userID)
	}
	return 0, nil
}
func (m *mockAddressRepo) ClearDefaultByUserID(ctx context.Context, userID string) error {
	if m.clearDefaultByUserFn != nil {
		return m.clearDefaultByUserFn(ctx, userID)
	}
	return nil
}

// mockSellerProfileRepo mocks domain.SellerProfileRepository.
type mockSellerProfileRepo struct {
	createFn    func(ctx context.Context, seller *domain.SellerProfile) error
	getByIDFn   func(ctx context.Context, id string) (*domain.SellerProfile, error)
	getByUserFn func(ctx context.Context, userID string) (*domain.SellerProfile, error)
	updateFn    func(ctx context.Context, seller *domain.SellerProfile) error
	listFn      func(ctx context.Context, page, size int) ([]domain.SellerProfile, int64, error)
}

func (m *mockSellerProfileRepo) Create(ctx context.Context, seller *domain.SellerProfile) error {
	if m.createFn != nil {
		return m.createFn(ctx, seller)
	}
	return nil
}
func (m *mockSellerProfileRepo) GetByID(ctx context.Context, id string) (*domain.SellerProfile, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}
func (m *mockSellerProfileRepo) GetByUserID(ctx context.Context, userID string) (*domain.SellerProfile, error) {
	if m.getByUserFn != nil {
		return m.getByUserFn(ctx, userID)
	}
	return nil, nil
}
func (m *mockSellerProfileRepo) Update(ctx context.Context, seller *domain.SellerProfile) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, seller)
	}
	return nil
}
func (m *mockSellerProfileRepo) List(ctx context.Context, page, size int) ([]domain.SellerProfile, int64, error) {
	if m.listFn != nil {
		return m.listFn(ctx, page, size)
	}
	return nil, 0, nil
}

// mockEventPublisher mocks events.Publisher (NO ctx parameter).
type mockEventPublisher struct {
	publishFn func(subject string, data interface{}) error
}

func (m *mockEventPublisher) Publish(subject string, data interface{}) error {
	if m.publishFn != nil {
		return m.publishFn(subject, data)
	}
	return nil
}

// mockFollowRepo mocks domain.FollowRepository.
type mockFollowRepo struct {
	createFn         func(ctx context.Context, follow *domain.UserFollow) error
	deleteFn         func(ctx context.Context, followerID, sellerID string) error
	listByFollowerFn func(ctx context.Context, followerID string, page, size int) ([]domain.SellerProfile, int64, error)
	countBySellerFn  func(ctx context.Context, sellerID string) (int64, error)
	existsFn         func(ctx context.Context, followerID, sellerID string) (bool, error)
}

func (m *mockFollowRepo) Create(ctx context.Context, follow *domain.UserFollow) error {
	if m.createFn != nil {
		return m.createFn(ctx, follow)
	}
	return nil
}
func (m *mockFollowRepo) Delete(ctx context.Context, followerID, sellerID string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, followerID, sellerID)
	}
	return nil
}
func (m *mockFollowRepo) ListByFollowerID(ctx context.Context, followerID string, page, size int) ([]domain.SellerProfile, int64, error) {
	if m.listByFollowerFn != nil {
		return m.listByFollowerFn(ctx, followerID, page, size)
	}
	return nil, 0, nil
}
func (m *mockFollowRepo) CountBySellerID(ctx context.Context, sellerID string) (int64, error) {
	if m.countBySellerFn != nil {
		return m.countBySellerFn(ctx, sellerID)
	}
	return 0, nil
}
func (m *mockFollowRepo) Exists(ctx context.Context, followerID, sellerID string) (bool, error) {
	if m.existsFn != nil {
		return m.existsFn(ctx, followerID, sellerID)
	}
	return false, nil
}

// mockWishlistRepo mocks domain.WishlistRepository.
type mockWishlistRepo struct {
	createFn       func(ctx context.Context, item *domain.WishlistItem) error
	deleteFn       func(ctx context.Context, userID, productID string) error
	listByUserIDFn func(ctx context.Context, userID string) ([]domain.WishlistItem, error)
	existsFn       func(ctx context.Context, userID, productID string) (bool, error)
}

func (m *mockWishlistRepo) Create(ctx context.Context, item *domain.WishlistItem) error {
	if m.createFn != nil {
		return m.createFn(ctx, item)
	}
	return nil
}
func (m *mockWishlistRepo) Delete(ctx context.Context, userID, productID string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, userID, productID)
	}
	return nil
}
func (m *mockWishlistRepo) ListByUserID(ctx context.Context, userID string) ([]domain.WishlistItem, error) {
	if m.listByUserIDFn != nil {
		return m.listByUserIDFn(ctx, userID)
	}
	return nil, nil
}
func (m *mockWishlistRepo) Exists(ctx context.Context, userID, productID string) (bool, error) {
	if m.existsFn != nil {
		return m.existsFn(ctx, userID, productID)
	}
	return false, nil
}

// ---------------------------------------------------------------------------
// helper
// ---------------------------------------------------------------------------

func strPtr(s string) *string { return &s }

// ---------------------------------------------------------------------------
// Address use-case tests
// ---------------------------------------------------------------------------

func TestCreateAddress_Success_FirstAddressGetsDefault(t *testing.T) {
	var created *domain.Address
	repo := &mockAddressRepo{
		countByUserIDFn: func(_ context.Context, _ string) (int64, error) {
			return 0, nil // no existing addresses
		},
		createFn: func(_ context.Context, addr *domain.Address) error {
			created = addr
			return nil
		},
	}

	uc := NewAddressUseCase(repo, zerolog.Nop())
	addr, err := uc.CreateAddress(context.Background(), "user-1", CreateAddressInput{
		FullName:   "John Doe",
		Phone:      "555-1234",
		Street:     "123 Main St",
		City:       "Springfield",
		PostalCode: "62704",
		Country:    "US",
	})

	require.NoError(t, err)
	assert.True(t, addr.IsDefault, "first address should be default")
	assert.Equal(t, "user-1", created.UserID)
	assert.Equal(t, "John Doe", created.FullName)
}

func TestCreateAddress_Success_NotFirstAddress(t *testing.T) {
	repo := &mockAddressRepo{
		countByUserIDFn: func(_ context.Context, _ string) (int64, error) {
			return 3, nil // already has addresses
		},
		createFn: func(_ context.Context, _ *domain.Address) error { return nil },
	}

	uc := NewAddressUseCase(repo, zerolog.Nop())
	addr, err := uc.CreateAddress(context.Background(), "user-1", CreateAddressInput{
		FullName:   "Jane Doe",
		Phone:      "555-5678",
		Street:     "456 Elm St",
		City:       "Shelbyville",
		PostalCode: "62705",
		Country:    "US",
	})

	require.NoError(t, err)
	assert.False(t, addr.IsDefault, "non-first address should not be default")
}

func TestCreateAddress_MaxAddressesError(t *testing.T) {
	repo := &mockAddressRepo{
		countByUserIDFn: func(_ context.Context, _ string) (int64, error) {
			return 10, nil // already at limit
		},
	}

	uc := NewAddressUseCase(repo, zerolog.Nop())
	_, err := uc.CreateAddress(context.Background(), "user-1", CreateAddressInput{
		FullName:   "X",
		Phone:      "X",
		Street:     "X",
		City:       "X",
		PostalCode: "X",
		Country:    "X",
	})

	require.Error(t, err)
	var valErr *pkgerrors.ValidationError
	assert.ErrorAs(t, err, &valErr)
	assert.Equal(t, "MAX_ADDRESSES", valErr.Code)
}

func TestListAddresses_Success(t *testing.T) {
	expected := []domain.Address{
		{ID: "a1", UserID: "user-1", FullName: "Addr1"},
		{ID: "a2", UserID: "user-1", FullName: "Addr2"},
	}
	repo := &mockAddressRepo{
		listByUserIDFn: func(_ context.Context, _ string) ([]domain.Address, error) {
			return expected, nil
		},
	}

	uc := NewAddressUseCase(repo, zerolog.Nop())
	list, err := uc.ListAddresses(context.Background(), "user-1")

	require.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, "a1", list[0].ID)
}

func TestUpdateAddress_Success(t *testing.T) {
	existing := &domain.Address{
		ID:       "addr-1",
		UserID:   "user-1",
		FullName: "Old Name",
		Phone:    "111",
		Street:   "Old St",
		City:     "Old City",
	}
	repo := &mockAddressRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Address, error) {
			return existing, nil
		},
		updateFn: func(_ context.Context, _ *domain.Address) error { return nil },
	}

	uc := NewAddressUseCase(repo, zerolog.Nop())
	result, err := uc.UpdateAddress(context.Background(), "user-1", "addr-1", UpdateAddressInput{
		FullName: strPtr("New Name"),
		City:     strPtr("New City"),
	})

	require.NoError(t, err)
	assert.Equal(t, "New Name", result.FullName)
	assert.Equal(t, "New City", result.City)
	assert.Equal(t, "111", result.Phone, "nil fields should remain unchanged")
	assert.Equal(t, "Old St", result.Street, "nil fields should remain unchanged")
}

func TestUpdateAddress_WrongOwner(t *testing.T) {
	repo := &mockAddressRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Address, error) {
			return &domain.Address{ID: "addr-1", UserID: "user-other"}, nil
		},
	}

	uc := NewAddressUseCase(repo, zerolog.Nop())
	_, err := uc.UpdateAddress(context.Background(), "user-1", "addr-1", UpdateAddressInput{
		FullName: strPtr("X"),
	})

	require.Error(t, err)
	var forbErr *pkgerrors.ForbiddenError
	assert.ErrorAs(t, err, &forbErr)
}

func TestDeleteAddress_Success(t *testing.T) {
	deleted := false
	repo := &mockAddressRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Address, error) {
			return &domain.Address{ID: "addr-1", UserID: "user-1"}, nil
		},
		deleteFn: func(_ context.Context, id string) error {
			assert.Equal(t, "addr-1", id)
			deleted = true
			return nil
		},
	}

	uc := NewAddressUseCase(repo, zerolog.Nop())
	err := uc.DeleteAddress(context.Background(), "user-1", "addr-1")

	require.NoError(t, err)
	assert.True(t, deleted)
}

func TestDeleteAddress_WrongOwner(t *testing.T) {
	repo := &mockAddressRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Address, error) {
			return &domain.Address{ID: "addr-1", UserID: "user-other"}, nil
		},
	}

	uc := NewAddressUseCase(repo, zerolog.Nop())
	err := uc.DeleteAddress(context.Background(), "user-1", "addr-1")

	require.Error(t, err)
	var forbErr *pkgerrors.ForbiddenError
	assert.ErrorAs(t, err, &forbErr)
}

func TestSetDefault_Success(t *testing.T) {
	clearCalled := false
	updatedAddr := (*domain.Address)(nil)

	repo := &mockAddressRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Address, error) {
			return &domain.Address{ID: "addr-2", UserID: "user-1", IsDefault: false}, nil
		},
		clearDefaultByUserFn: func(_ context.Context, userID string) error {
			assert.Equal(t, "user-1", userID)
			clearCalled = true
			return nil
		},
		updateFn: func(_ context.Context, addr *domain.Address) error {
			updatedAddr = addr
			return nil
		},
	}

	uc := NewAddressUseCase(repo, zerolog.Nop())
	err := uc.SetDefault(context.Background(), "user-1", "addr-2")

	require.NoError(t, err)
	assert.True(t, clearCalled)
	require.NotNil(t, updatedAddr)
	assert.True(t, updatedAddr.IsDefault)
}

func TestSetDefault_WrongOwner(t *testing.T) {
	repo := &mockAddressRepo{
		getByIDFn: func(_ context.Context, _ string) (*domain.Address, error) {
			return &domain.Address{ID: "addr-2", UserID: "user-other"}, nil
		},
	}

	uc := NewAddressUseCase(repo, zerolog.Nop())
	err := uc.SetDefault(context.Background(), "user-1", "addr-2")

	require.Error(t, err)
	var forbErr *pkgerrors.ForbiddenError
	assert.ErrorAs(t, err, &forbErr)
}

func TestCreateAddress_CountError(t *testing.T) {
	repo := &mockAddressRepo{
		countByUserIDFn: func(_ context.Context, _ string) (int64, error) {
			return 0, errors.New("db error")
		},
	}

	uc := NewAddressUseCase(repo, zerolog.Nop())
	_, err := uc.CreateAddress(context.Background(), "user-1", CreateAddressInput{
		FullName:   "X",
		Phone:      "X",
		Street:     "X",
		City:       "X",
		PostalCode: "X",
		Country:    "X",
	})

	require.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}
