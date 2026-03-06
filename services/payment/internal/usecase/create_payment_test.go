package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/payment/internal/domain"
	"github.com/southern-martin/ecommerce/services/payment/internal/infrastructure/stripe"
)

// ---------------------------------------------------------------------------
// Hand-written function-field mocks (shared by all *_test.go in this package)
// ---------------------------------------------------------------------------

// mockPaymentRepo implements domain.PaymentRepository.
type mockPaymentRepo struct {
	createFn       func(ctx context.Context, payment *domain.Payment) error
	getByIDFn      func(ctx context.Context, id string) (*domain.Payment, error)
	getByOrderIDFn func(ctx context.Context, orderID string) (*domain.Payment, error)
	getByStripeIDFn func(ctx context.Context, stripePaymentID string) (*domain.Payment, error)
	updateStatusFn func(ctx context.Context, id string, status domain.PaymentStatus, failureReason string) error
	updateStripeIDFn func(ctx context.Context, id string, stripePaymentID string) error
	listFn         func(ctx context.Context, buyerID string, page, pageSize int) ([]*domain.Payment, int64, error)
}

func (m *mockPaymentRepo) Create(ctx context.Context, payment *domain.Payment) error {
	if m.createFn != nil {
		return m.createFn(ctx, payment)
	}
	return nil
}
func (m *mockPaymentRepo) GetByID(ctx context.Context, id string) (*domain.Payment, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}
func (m *mockPaymentRepo) GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	if m.getByOrderIDFn != nil {
		return m.getByOrderIDFn(ctx, orderID)
	}
	return nil, nil
}
func (m *mockPaymentRepo) GetByStripeID(ctx context.Context, stripePaymentID string) (*domain.Payment, error) {
	if m.getByStripeIDFn != nil {
		return m.getByStripeIDFn(ctx, stripePaymentID)
	}
	return nil, nil
}
func (m *mockPaymentRepo) UpdateStatus(ctx context.Context, id string, status domain.PaymentStatus, failureReason string) error {
	if m.updateStatusFn != nil {
		return m.updateStatusFn(ctx, id, status, failureReason)
	}
	return nil
}
func (m *mockPaymentRepo) UpdateStripeID(ctx context.Context, id string, stripePaymentID string) error {
	if m.updateStripeIDFn != nil {
		return m.updateStripeIDFn(ctx, id, stripePaymentID)
	}
	return nil
}
func (m *mockPaymentRepo) List(ctx context.Context, buyerID string, page, pageSize int) ([]*domain.Payment, int64, error) {
	if m.listFn != nil {
		return m.listFn(ctx, buyerID, page, pageSize)
	}
	return nil, 0, nil
}

// mockWalletRepo implements domain.WalletRepository.
type mockWalletRepo struct {
	getOrCreateFn          func(ctx context.Context, sellerID string) (*domain.SellerWallet, error)
	creditPendingFn        func(ctx context.Context, sellerID string, amountCents int64) error
	movePendingToAvailFn   func(ctx context.Context, sellerID string, amountCents int64) error
	debitAvailableFn       func(ctx context.Context, sellerID string, amountCents int64) error
	createTransactionFn    func(ctx context.Context, tx *domain.WalletTransaction) error
	listTransactionsFn     func(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.WalletTransaction, int64, error)
}

func (m *mockWalletRepo) GetOrCreate(ctx context.Context, sellerID string) (*domain.SellerWallet, error) {
	if m.getOrCreateFn != nil {
		return m.getOrCreateFn(ctx, sellerID)
	}
	return &domain.SellerWallet{SellerID: sellerID}, nil
}
func (m *mockWalletRepo) CreditPending(ctx context.Context, sellerID string, amountCents int64) error {
	if m.creditPendingFn != nil {
		return m.creditPendingFn(ctx, sellerID, amountCents)
	}
	return nil
}
func (m *mockWalletRepo) MovePendingToAvailable(ctx context.Context, sellerID string, amountCents int64) error {
	if m.movePendingToAvailFn != nil {
		return m.movePendingToAvailFn(ctx, sellerID, amountCents)
	}
	return nil
}
func (m *mockWalletRepo) DebitAvailable(ctx context.Context, sellerID string, amountCents int64) error {
	if m.debitAvailableFn != nil {
		return m.debitAvailableFn(ctx, sellerID, amountCents)
	}
	return nil
}
func (m *mockWalletRepo) CreateTransaction(ctx context.Context, tx *domain.WalletTransaction) error {
	if m.createTransactionFn != nil {
		return m.createTransactionFn(ctx, tx)
	}
	return nil
}
func (m *mockWalletRepo) ListTransactions(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.WalletTransaction, int64, error) {
	if m.listTransactionsFn != nil {
		return m.listTransactionsFn(ctx, sellerID, page, pageSize)
	}
	return nil, 0, nil
}

// mockPayoutRepo implements domain.PayoutRepository.
type mockPayoutRepo struct {
	createFn       func(ctx context.Context, payout *domain.Payout) error
	getByIDFn      func(ctx context.Context, id string) (*domain.Payout, error)
	listBySellerFn func(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.Payout, int64, error)
	updateStatusFn func(ctx context.Context, id string, status domain.PayoutStatus) error
}

func (m *mockPayoutRepo) Create(ctx context.Context, payout *domain.Payout) error {
	if m.createFn != nil {
		return m.createFn(ctx, payout)
	}
	return nil
}
func (m *mockPayoutRepo) GetByID(ctx context.Context, id string) (*domain.Payout, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}
func (m *mockPayoutRepo) ListBySeller(ctx context.Context, sellerID string, page, pageSize int) ([]*domain.Payout, int64, error) {
	if m.listBySellerFn != nil {
		return m.listBySellerFn(ctx, sellerID, page, pageSize)
	}
	return nil, 0, nil
}
func (m *mockPayoutRepo) UpdateStatus(ctx context.Context, id string, status domain.PayoutStatus) error {
	if m.updateStatusFn != nil {
		return m.updateStatusFn(ctx, id, status)
	}
	return nil
}

// mockStripeClient implements stripe.StripeClient.
type mockStripeClient struct {
	createPaymentIntentFn  func(amountCents int64, currency string, metadata map[string]string) (string, string, error)
	confirmPaymentIntentFn func(paymentIntentID string) error
	createRefundFn         func(paymentIntentID string, amountCents int64) (string, error)
	createTransferFn       func(amountCents int64, destinationAccountID string, metadata map[string]string) (string, error)
}

func (m *mockStripeClient) CreatePaymentIntent(amountCents int64, currency string, metadata map[string]string) (string, string, error) {
	if m.createPaymentIntentFn != nil {
		return m.createPaymentIntentFn(amountCents, currency, metadata)
	}
	return "pi_test", "secret_test", nil
}
func (m *mockStripeClient) ConfirmPaymentIntent(paymentIntentID string) error {
	if m.confirmPaymentIntentFn != nil {
		return m.confirmPaymentIntentFn(paymentIntentID)
	}
	return nil
}
func (m *mockStripeClient) CreateRefund(paymentIntentID string, amountCents int64) (string, error) {
	if m.createRefundFn != nil {
		return m.createRefundFn(paymentIntentID, amountCents)
	}
	return "re_test", nil
}
func (m *mockStripeClient) CreateTransfer(amountCents int64, destinationAccountID string, metadata map[string]string) (string, error) {
	if m.createTransferFn != nil {
		return m.createTransferFn(amountCents, destinationAccountID, metadata)
	}
	return "tr_test", nil
}

// mockEventPublisher implements domain.EventPublisher.
type mockEventPublisher struct {
	publishFn func(ctx context.Context, subject string, data interface{}) error
}

func (m *mockEventPublisher) Publish(ctx context.Context, subject string, data interface{}) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, subject, data)
	}
	return nil
}

// Compile-time interface satisfaction checks.
var (
	_ domain.PaymentRepository = (*mockPaymentRepo)(nil)
	_ domain.WalletRepository  = (*mockWalletRepo)(nil)
	_ domain.PayoutRepository  = (*mockPayoutRepo)(nil)
	_ stripe.StripeClient      = (*mockStripeClient)(nil)
	_ domain.EventPublisher    = (*mockEventPublisher)(nil)
)

// ---------------------------------------------------------------------------
// CreatePaymentUseCase tests
// ---------------------------------------------------------------------------

func TestCreatePayment_Success(t *testing.T) {
	var createdPayment *domain.Payment
	var updatedStripeID string
	var publishedSubject string

	repo := &mockPaymentRepo{
		createFn: func(_ context.Context, p *domain.Payment) error {
			createdPayment = p
			return nil
		},
		updateStripeIDFn: func(_ context.Context, _ string, sid string) error {
			updatedStripeID = sid
			return nil
		},
	}
	sc := &mockStripeClient{
		createPaymentIntentFn: func(amount int64, currency string, meta map[string]string) (string, string, error) {
			return "pi_abc123", "secret_xyz", nil
		},
	}
	pub := &mockEventPublisher{
		publishFn: func(_ context.Context, subj string, _ interface{}) error {
			publishedSubject = subj
			return nil
		},
	}

	uc := NewCreatePaymentUseCase(repo, sc, pub)
	out, err := uc.Execute(context.Background(), CreatePaymentInput{
		OrderID:     "order-1",
		BuyerID:     "buyer-1",
		AmountCents: 5000,
		Currency:    "usd",
		Method:      domain.PaymentMethodCard,
	})

	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, "pi_abc123", out.StripePaymentID)
	assert.Equal(t, "secret_xyz", out.ClientSecret)
	assert.Equal(t, "pending", out.Status)
	assert.NotEmpty(t, out.PaymentID)

	// Verify side effects.
	require.NotNil(t, createdPayment)
	assert.Equal(t, "order-1", createdPayment.OrderID)
	assert.Equal(t, "buyer-1", createdPayment.BuyerID)
	assert.Equal(t, int64(5000), createdPayment.AmountCents)
	assert.Equal(t, domain.PaymentStatusPending, createdPayment.Status)

	assert.Equal(t, "pi_abc123", updatedStripeID)
	assert.Equal(t, domain.EventPaymentInitiated, publishedSubject)
}

func TestCreatePayment_DefaultCurrency(t *testing.T) {
	var createdPayment *domain.Payment
	repo := &mockPaymentRepo{
		createFn: func(_ context.Context, p *domain.Payment) error {
			createdPayment = p
			return nil
		},
	}
	uc := NewCreatePaymentUseCase(repo, &mockStripeClient{}, &mockEventPublisher{})

	out, err := uc.Execute(context.Background(), CreatePaymentInput{
		OrderID:     "order-2",
		BuyerID:     "buyer-2",
		AmountCents: 1000,
		Currency:    "", // should default to "usd"
	})

	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, "usd", createdPayment.Currency)
}

func TestCreatePayment_DefaultMethod(t *testing.T) {
	var createdPayment *domain.Payment
	repo := &mockPaymentRepo{
		createFn: func(_ context.Context, p *domain.Payment) error {
			createdPayment = p
			return nil
		},
	}
	uc := NewCreatePaymentUseCase(repo, &mockStripeClient{}, &mockEventPublisher{})

	out, err := uc.Execute(context.Background(), CreatePaymentInput{
		OrderID:     "order-3",
		BuyerID:     "buyer-3",
		AmountCents: 2000,
		Method:      "", // should default to PaymentMethodCard
	})

	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, domain.PaymentMethodCard, createdPayment.Method)
}

func TestCreatePayment_StripeFailure(t *testing.T) {
	var failedStatus domain.PaymentStatus
	repo := &mockPaymentRepo{
		createFn: func(_ context.Context, _ *domain.Payment) error { return nil },
		updateStatusFn: func(_ context.Context, _ string, status domain.PaymentStatus, _ string) error {
			failedStatus = status
			return nil
		},
	}
	sc := &mockStripeClient{
		createPaymentIntentFn: func(_ int64, _ string, _ map[string]string) (string, string, error) {
			return "", "", errors.New("stripe unavailable")
		},
	}
	uc := NewCreatePaymentUseCase(repo, sc, &mockEventPublisher{})

	out, err := uc.Execute(context.Background(), CreatePaymentInput{
		OrderID:     "order-4",
		BuyerID:     "buyer-4",
		AmountCents: 3000,
	})

	require.Error(t, err)
	assert.Nil(t, out)
	assert.Contains(t, err.Error(), "stripe")
	assert.Equal(t, domain.PaymentStatusFailed, failedStatus)
}

func TestCreatePayment_RepoCreateError(t *testing.T) {
	repo := &mockPaymentRepo{
		createFn: func(_ context.Context, _ *domain.Payment) error {
			return errors.New("db connection lost")
		},
	}
	uc := NewCreatePaymentUseCase(repo, &mockStripeClient{}, &mockEventPublisher{})

	out, err := uc.Execute(context.Background(), CreatePaymentInput{
		OrderID:     "order-5",
		BuyerID:     "buyer-5",
		AmountCents: 4000,
	})

	require.Error(t, err)
	assert.Nil(t, out)
	assert.Contains(t, err.Error(), "failed to create payment record")
}
