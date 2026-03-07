package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/southern-martin/ecommerce/services/notification/internal/domain"
)

// ---------------------------------------------------------------------------
// Hand-written function-field mocks (shared with preference_test.go)
// ---------------------------------------------------------------------------

// --- NotificationRepository mock ---

type mockNotificationRepo struct {
	getByIDFn      func(ctx context.Context, id string) (*domain.Notification, error)
	listByUserFn   func(ctx context.Context, userID string, filter domain.NotificationFilter) ([]domain.Notification, int64, error)
	createFn       func(ctx context.Context, notification *domain.Notification) error
	updateFn       func(ctx context.Context, notification *domain.Notification) error
	markAsReadFn   func(ctx context.Context, id string) error
	markAllAsReadFn func(ctx context.Context, userID string) error
	countUnreadFn  func(ctx context.Context, userID string) (int64, error)
}

func (m *mockNotificationRepo) GetByID(ctx context.Context, id string) (*domain.Notification, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}
func (m *mockNotificationRepo) ListByUser(ctx context.Context, userID string, filter domain.NotificationFilter) ([]domain.Notification, int64, error) {
	if m.listByUserFn != nil {
		return m.listByUserFn(ctx, userID, filter)
	}
	return nil, 0, nil
}
func (m *mockNotificationRepo) Create(ctx context.Context, notification *domain.Notification) error {
	if m.createFn != nil {
		return m.createFn(ctx, notification)
	}
	return nil
}
func (m *mockNotificationRepo) Update(ctx context.Context, notification *domain.Notification) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, notification)
	}
	return nil
}
func (m *mockNotificationRepo) MarkAsRead(ctx context.Context, id string) error {
	if m.markAsReadFn != nil {
		return m.markAsReadFn(ctx, id)
	}
	return nil
}
func (m *mockNotificationRepo) MarkAllAsRead(ctx context.Context, userID string) error {
	if m.markAllAsReadFn != nil {
		return m.markAllAsReadFn(ctx, userID)
	}
	return nil
}
func (m *mockNotificationRepo) CountUnread(ctx context.Context, userID string) (int64, error) {
	if m.countUnreadFn != nil {
		return m.countUnreadFn(ctx, userID)
	}
	return 0, nil
}

// --- EventPublisher mock ---

type mockNotifEventPub struct {
	publishFn func(ctx context.Context, subject string, data interface{}) error
}

func (m *mockNotifEventPub) Publish(ctx context.Context, subject string, data interface{}) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, subject, data)
	}
	return nil
}

// --- EmailSender mock ---

type mockEmailSender struct {
	sendFn func(to, subject, htmlBody string) error
}

func (m *mockEmailSender) Send(to, subject, htmlBody string) error {
	if m.sendFn != nil {
		return m.sendFn(to, subject, htmlBody)
	}
	return nil
}

// --- PreferenceRepository mock ---

type mockPreferenceRepo struct {
	getByUserFn func(ctx context.Context, userID string) ([]domain.NotificationPreference, error)
	upsertFn    func(ctx context.Context, preference *domain.NotificationPreference) error
}

func (m *mockPreferenceRepo) GetByUser(ctx context.Context, userID string) ([]domain.NotificationPreference, error) {
	if m.getByUserFn != nil {
		return m.getByUserFn(ctx, userID)
	}
	return nil, nil
}
func (m *mockPreferenceRepo) Upsert(ctx context.Context, preference *domain.NotificationPreference) error {
	if m.upsertFn != nil {
		return m.upsertFn(ctx, preference)
	}
	return nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func defaultNotifMocks() (*mockNotificationRepo, *mockNotifEventPub, *mockEmailSender) {
	return &mockNotificationRepo{}, &mockNotifEventPub{}, &mockEmailSender{}
}

func newNotifUseCase(repo *mockNotificationRepo, pub *mockNotifEventPub, email *mockEmailSender) *NotificationUseCase {
	return NewNotificationUseCase(repo, pub, email)
}

// ===========================================================================
// SendNotification (Create) tests
// ===========================================================================

func TestSendNotification_Success(t *testing.T) {
	repo, pub, email := defaultNotifMocks()
	var saved *domain.Notification
	repo.createFn = func(_ context.Context, n *domain.Notification) error {
		saved = n
		return nil
	}
	repo.updateFn = func(_ context.Context, _ *domain.Notification) error { return nil }

	uc := newNotifUseCase(repo, pub, email)
	notif, err := uc.SendNotification(context.Background(), SendNotificationRequest{
		UserID:  "user-1",
		Type:    "order_update",
		Channel: "in_app",
		Subject: "Order shipped",
		Body:    "Your order has been shipped!",
	})

	require.NoError(t, err)
	require.NotNil(t, notif)
	assert.NotEmpty(t, notif.ID)
	assert.Equal(t, "user-1", notif.UserID)
	assert.Equal(t, domain.NotificationType("order_update"), notif.Type)
	assert.Equal(t, domain.NotificationChannel("in_app"), notif.Channel)
	assert.Equal(t, "Order shipped", notif.Subject)
	assert.Equal(t, "Your order has been shipped!", notif.Body)
	assert.Equal(t, domain.StatusSent, notif.Status)
	assert.NotNil(t, notif.SentAt)
	assert.NotNil(t, saved)
}

func TestSendNotification_EmailChannel(t *testing.T) {
	repo, pub, email := defaultNotifMocks()
	repo.createFn = func(_ context.Context, _ *domain.Notification) error { return nil }
	repo.updateFn = func(_ context.Context, _ *domain.Notification) error { return nil }
	var sentTo, sentSubject string
	email.sendFn = func(to, subject, htmlBody string) error {
		sentTo = to
		sentSubject = subject
		return nil
	}

	uc := newNotifUseCase(repo, pub, email)
	_, err := uc.SendNotification(context.Background(), SendNotificationRequest{
		UserID:  "user-1",
		Type:    "order_update",
		Channel: "email",
		Subject: "Order shipped",
		Body:    "<h1>Shipped!</h1>",
		Data:    "user@example.com",
	})

	require.NoError(t, err)
	assert.Equal(t, "user@example.com", sentTo)
	assert.Equal(t, "Order shipped", sentSubject)
}

func TestSendNotification_EmailSendErrorDoesNotFail(t *testing.T) {
	repo, pub, email := defaultNotifMocks()
	repo.createFn = func(_ context.Context, _ *domain.Notification) error { return nil }
	repo.updateFn = func(_ context.Context, _ *domain.Notification) error { return nil }
	email.sendFn = func(_, _, _ string) error {
		return errors.New("SMTP connection failed")
	}

	uc := newNotifUseCase(repo, pub, email)
	notif, err := uc.SendNotification(context.Background(), SendNotificationRequest{
		UserID:  "user-1",
		Type:    "order_update",
		Channel: "email",
		Subject: "Order shipped",
		Body:    "Shipped!",
		Data:    "user@example.com",
	})

	// Email send failure is best-effort; should not fail the operation
	require.NoError(t, err)
	require.NotNil(t, notif)
}

func TestSendNotification_CreateRepoError(t *testing.T) {
	repo, pub, email := defaultNotifMocks()
	repo.createFn = func(_ context.Context, _ *domain.Notification) error {
		return errors.New("db error")
	}

	uc := newNotifUseCase(repo, pub, email)
	_, err := uc.SendNotification(context.Background(), SendNotificationRequest{
		UserID:  "user-1",
		Type:    "order_update",
		Channel: "in_app",
		Subject: "Test",
		Body:    "Test body",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

func TestSendNotification_UpdateRepoError(t *testing.T) {
	repo, pub, email := defaultNotifMocks()
	repo.createFn = func(_ context.Context, _ *domain.Notification) error { return nil }
	repo.updateFn = func(_ context.Context, _ *domain.Notification) error {
		return errors.New("update failed")
	}

	uc := newNotifUseCase(repo, pub, email)
	_, err := uc.SendNotification(context.Background(), SendNotificationRequest{
		UserID:  "user-1",
		Type:    "system",
		Channel: "push",
		Subject: "System update",
		Body:    "System is updating",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "update failed")
}

func TestSendNotification_PublishEventErrorDoesNotFail(t *testing.T) {
	repo, pub, email := defaultNotifMocks()
	repo.createFn = func(_ context.Context, _ *domain.Notification) error { return nil }
	repo.updateFn = func(_ context.Context, _ *domain.Notification) error { return nil }
	pub.publishFn = func(_ context.Context, _ string, _ interface{}) error {
		return errors.New("event bus down")
	}

	uc := newNotifUseCase(repo, pub, email)
	notif, err := uc.SendNotification(context.Background(), SendNotificationRequest{
		UserID:  "user-1",
		Type:    "promotion",
		Channel: "in_app",
		Subject: "Sale!",
		Body:    "Big sale today!",
	})

	// Event publish failure should not cause SendNotification to fail
	require.NoError(t, err)
	require.NotNil(t, notif)
}

func TestSendNotification_SetsStatusQueued(t *testing.T) {
	repo, pub, email := defaultNotifMocks()
	var createdStatus domain.NotificationStatus
	repo.createFn = func(_ context.Context, n *domain.Notification) error {
		createdStatus = n.Status
		return nil
	}
	repo.updateFn = func(_ context.Context, _ *domain.Notification) error { return nil }

	uc := newNotifUseCase(repo, pub, email)
	_, err := uc.SendNotification(context.Background(), SendNotificationRequest{
		UserID:  "user-1",
		Type:    "system",
		Channel: "in_app",
		Subject: "Test",
		Body:    "Test",
	})

	require.NoError(t, err)
	assert.Equal(t, domain.StatusQueued, createdStatus)
}

// ===========================================================================
// ListUserNotifications tests
// ===========================================================================

func TestListUserNotifications_Success(t *testing.T) {
	repo, pub, email := defaultNotifMocks()
	notifications := []domain.Notification{
		{ID: "notif-1", UserID: "user-1", Subject: "Order shipped"},
		{ID: "notif-2", UserID: "user-1", Subject: "Payment received"},
	}
	repo.listByUserFn = func(_ context.Context, userID string, filter domain.NotificationFilter) ([]domain.Notification, int64, error) {
		assert.Equal(t, "user-1", userID)
		return notifications, 2, nil
	}

	uc := newNotifUseCase(repo, pub, email)
	result, total, err := uc.ListUserNotifications(context.Background(), "user-1", domain.NotificationFilter{})

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
}

func TestListUserNotifications_Empty(t *testing.T) {
	repo, pub, email := defaultNotifMocks()
	// Default listByUserFn returns nil, 0, nil

	uc := newNotifUseCase(repo, pub, email)
	result, total, err := uc.ListUserNotifications(context.Background(), "user-1", domain.NotificationFilter{})

	require.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), total)
}

func TestListUserNotifications_RepoError(t *testing.T) {
	repo, pub, email := defaultNotifMocks()
	repo.listByUserFn = func(_ context.Context, _ string, _ domain.NotificationFilter) ([]domain.Notification, int64, error) {
		return nil, 0, errors.New("query failed")
	}

	uc := newNotifUseCase(repo, pub, email)
	_, _, err := uc.ListUserNotifications(context.Background(), "user-1", domain.NotificationFilter{})

	require.Error(t, err)
}

// ===========================================================================
// MarkAsRead tests
// ===========================================================================

func TestMarkAsRead_Success(t *testing.T) {
	repo, pub, email := defaultNotifMocks()
	var markedID string
	repo.markAsReadFn = func(_ context.Context, id string) error {
		markedID = id
		return nil
	}

	uc := newNotifUseCase(repo, pub, email)
	err := uc.MarkAsRead(context.Background(), "notif-1")

	require.NoError(t, err)
	assert.Equal(t, "notif-1", markedID)
}

func TestMarkAsRead_RepoError(t *testing.T) {
	repo, pub, email := defaultNotifMocks()
	repo.markAsReadFn = func(_ context.Context, _ string) error {
		return errors.New("mark failed")
	}

	uc := newNotifUseCase(repo, pub, email)
	err := uc.MarkAsRead(context.Background(), "notif-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "mark failed")
}

// ===========================================================================
// MarkAllAsRead tests
// ===========================================================================

func TestMarkAllAsRead_Success(t *testing.T) {
	repo, pub, email := defaultNotifMocks()
	var markedUserID string
	repo.markAllAsReadFn = func(_ context.Context, userID string) error {
		markedUserID = userID
		return nil
	}

	uc := newNotifUseCase(repo, pub, email)
	err := uc.MarkAllAsRead(context.Background(), "user-1")

	require.NoError(t, err)
	assert.Equal(t, "user-1", markedUserID)
}

func TestMarkAllAsRead_RepoError(t *testing.T) {
	repo, pub, email := defaultNotifMocks()
	repo.markAllAsReadFn = func(_ context.Context, _ string) error {
		return errors.New("bulk update failed")
	}

	uc := newNotifUseCase(repo, pub, email)
	err := uc.MarkAllAsRead(context.Background(), "user-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "bulk update failed")
}

// ===========================================================================
// GetUnreadCount tests
// ===========================================================================

func TestGetUnreadCount_Success(t *testing.T) {
	repo, pub, email := defaultNotifMocks()
	repo.countUnreadFn = func(_ context.Context, userID string) (int64, error) {
		assert.Equal(t, "user-1", userID)
		return 5, nil
	}

	uc := newNotifUseCase(repo, pub, email)
	count, err := uc.GetUnreadCount(context.Background(), "user-1")

	require.NoError(t, err)
	assert.Equal(t, int64(5), count)
}

func TestGetUnreadCount_Zero(t *testing.T) {
	repo, pub, email := defaultNotifMocks()
	// Default countUnreadFn returns 0, nil

	uc := newNotifUseCase(repo, pub, email)
	count, err := uc.GetUnreadCount(context.Background(), "user-1")

	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestGetUnreadCount_RepoError(t *testing.T) {
	repo, pub, email := defaultNotifMocks()
	repo.countUnreadFn = func(_ context.Context, _ string) (int64, error) {
		return 0, errors.New("count failed")
	}

	uc := newNotifUseCase(repo, pub, email)
	_, err := uc.GetUnreadCount(context.Background(), "user-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "count failed")
}
