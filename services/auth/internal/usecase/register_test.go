package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/southern-martin/ecommerce/services/auth/internal/domain"
)

// --- Mock for domain.UserRepository ---

type mockUserRepo struct {
	createFn             func(ctx context.Context, user *domain.AuthUser) error
	getByIDFn            func(ctx context.Context, id string) (*domain.AuthUser, error)
	getByEmailFn         func(ctx context.Context, email string) (*domain.AuthUser, error)
	getByOAuthProviderFn func(ctx context.Context, provider, providerID string) (*domain.AuthUser, error)
	updateRefreshTokenFn func(ctx context.Context, id, token string) error
	updatePasswordFn     func(ctx context.Context, id, passwordHash string) error
	updateResetTokenFn   func(ctx context.Context, id, token string, exp time.Time) error
	clearResetTokenFn    func(ctx context.Context, id string) error
	updateRoleFn         func(ctx context.Context, id, role string) error
}

func (m *mockUserRepo) Create(ctx context.Context, user *domain.AuthUser) error {
	return m.createFn(ctx, user)
}
func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*domain.AuthUser, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}
func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.AuthUser, error) {
	if m.getByEmailFn != nil {
		return m.getByEmailFn(ctx, email)
	}
	return nil, errors.New("not found")
}
func (m *mockUserRepo) GetByOAuthProvider(ctx context.Context, provider, providerID string) (*domain.AuthUser, error) {
	if m.getByOAuthProviderFn != nil {
		return m.getByOAuthProviderFn(ctx, provider, providerID)
	}
	return nil, errors.New("not found")
}
func (m *mockUserRepo) UpdateRefreshToken(ctx context.Context, id, token string) error {
	if m.updateRefreshTokenFn != nil {
		return m.updateRefreshTokenFn(ctx, id, token)
	}
	return nil
}
func (m *mockUserRepo) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	if m.updatePasswordFn != nil {
		return m.updatePasswordFn(ctx, id, passwordHash)
	}
	return nil
}
func (m *mockUserRepo) UpdateResetToken(ctx context.Context, id, token string, exp time.Time) error {
	if m.updateResetTokenFn != nil {
		return m.updateResetTokenFn(ctx, id, token, exp)
	}
	return nil
}
func (m *mockUserRepo) ClearResetToken(ctx context.Context, id string) error {
	if m.clearResetTokenFn != nil {
		return m.clearResetTokenFn(ctx, id)
	}
	return nil
}
func (m *mockUserRepo) UpdateRole(ctx context.Context, id, role string) error {
	if m.updateRoleFn != nil {
		return m.updateRoleFn(ctx, id, role)
	}
	return nil
}

// --- Register Tests ---

const testJWTSecret = "test-secret-key-for-tests-at-least-32-chars"

func TestRegister_Success(t *testing.T) {
	var createdUser *domain.AuthUser
	repo := &mockUserRepo{
		getByEmailFn: func(_ context.Context, _ string) (*domain.AuthUser, error) {
			return nil, errors.New("not found")
		},
		createFn: func(_ context.Context, user *domain.AuthUser) error {
			user.ID = "generated-id"
			createdUser = user
			return nil
		},
	}
	uc := NewRegisterUseCase(repo, nil, testJWTSecret, 15*time.Minute, 7*24*time.Hour, zerolog.Nop())

	out, err := uc.Execute(context.Background(), RegisterInput{Email: "test@example.com", Password: "password123"})
	require.NoError(t, err)
	assert.Equal(t, "test@example.com", out.Email)
	assert.Equal(t, "buyer", out.Role)
	assert.NotEmpty(t, out.AccessToken)
	assert.NotEmpty(t, out.RefreshToken)
	assert.NotNil(t, createdUser)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(_ context.Context, _ string) (*domain.AuthUser, error) {
			return &domain.AuthUser{ID: "existing", Email: "test@example.com"}, nil
		},
	}
	uc := NewRegisterUseCase(repo, nil, testJWTSecret, 15*time.Minute, 7*24*time.Hour, zerolog.Nop())

	_, err := uc.Execute(context.Background(), RegisterInput{Email: "test@example.com", Password: "password123"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email already registered")
}

func TestRegister_PasswordHashed(t *testing.T) {
	var savedHash string
	repo := &mockUserRepo{
		getByEmailFn: func(_ context.Context, _ string) (*domain.AuthUser, error) { return nil, errors.New("not found") },
		createFn: func(_ context.Context, user *domain.AuthUser) error {
			savedHash = user.PasswordHash
			user.ID = "id-1"
			return nil
		},
	}
	uc := NewRegisterUseCase(repo, nil, testJWTSecret, 15*time.Minute, 7*24*time.Hour, zerolog.Nop())

	_, err := uc.Execute(context.Background(), RegisterInput{Email: "a@b.com", Password: "mypassword"})
	require.NoError(t, err)
	assert.NotEqual(t, "mypassword", savedHash)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(savedHash), []byte("mypassword")))
}

func TestRegister_DefaultRole(t *testing.T) {
	var savedRole string
	repo := &mockUserRepo{
		getByEmailFn: func(_ context.Context, _ string) (*domain.AuthUser, error) { return nil, errors.New("not found") },
		createFn: func(_ context.Context, user *domain.AuthUser) error {
			savedRole = user.Role
			user.ID = "id-1"
			return nil
		},
	}
	uc := NewRegisterUseCase(repo, nil, testJWTSecret, 15*time.Minute, 7*24*time.Hour, zerolog.Nop())

	_, err := uc.Execute(context.Background(), RegisterInput{Email: "a@b.com", Password: "password123"})
	require.NoError(t, err)
	assert.Equal(t, "buyer", savedRole)
}

func TestRegister_CreateRepoError(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(_ context.Context, _ string) (*domain.AuthUser, error) { return nil, errors.New("not found") },
		createFn:     func(_ context.Context, _ *domain.AuthUser) error { return errors.New("db error") },
	}
	uc := NewRegisterUseCase(repo, nil, testJWTSecret, 15*time.Minute, 7*24*time.Hour, zerolog.Nop())

	_, err := uc.Execute(context.Background(), RegisterInput{Email: "a@b.com", Password: "password123"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create user")
}

func TestRegister_TokensGenerated(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailFn: func(_ context.Context, _ string) (*domain.AuthUser, error) { return nil, errors.New("not found") },
		createFn:     func(_ context.Context, u *domain.AuthUser) error { u.ID = "id-1"; return nil },
	}
	uc := NewRegisterUseCase(repo, nil, testJWTSecret, 15*time.Minute, 7*24*time.Hour, zerolog.Nop())

	out, err := uc.Execute(context.Background(), RegisterInput{Email: "a@b.com", Password: "password123"})
	require.NoError(t, err)
	assert.NotEmpty(t, out.AccessToken)
	assert.NotEmpty(t, out.RefreshToken)
	assert.NotEqual(t, out.AccessToken, out.RefreshToken)
}
