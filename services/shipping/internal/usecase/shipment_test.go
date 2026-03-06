package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/southern-martin/ecommerce/services/shipping/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShipmentUseCase_CreateShipment(t *testing.T) {
	t.Run("success with all fields", func(t *testing.T) {
		var created *domain.Shipment
		var publishedSubject string
		repo := &mockShipmentRepo{
			createFn: func(ctx context.Context, shipment *domain.Shipment) error {
				created = shipment
				return nil
			},
		}
		pub := &mockEventPublisher{
			publishFn: func(ctx context.Context, subject string, data interface{}) error {
				publishedSubject = subject
				return nil
			},
		}
		uc := NewShipmentUseCase(repo, pub)

		req := CreateShipmentRequest{
			OrderID:     "order-1",
			SellerID:    "seller-1",
			CarrierCode: "fedex",
			ServiceCode: "standard",
			Origin:      domain.Address{Street: "123 Main St", City: "Austin", State: "TX", PostalCode: "78701", Country: "US"},
			Destination: domain.Address{Street: "456 Oak Ave", City: "Dallas", State: "TX", PostalCode: "75201", Country: "US"},
			WeightGrams: 1000,
			RateCents:   1500,
			Currency:    "EUR",
			Items: []domain.ShipmentItem{
				{ProductID: "prod-1", ProductName: "Widget", Quantity: 2},
				{ProductID: "prod-2", ProductName: "Gadget", Quantity: 1},
			},
		}

		result, err := uc.CreateShipment(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, "order-1", result.OrderID)
		assert.Equal(t, "seller-1", result.SellerID)
		assert.Equal(t, "fedex", result.CarrierCode)
		assert.Equal(t, domain.StatusPending, result.Status)
		assert.Equal(t, "EUR", result.Currency)
		assert.Equal(t, int64(1500), result.RateCents)
		assert.Equal(t, 1000, result.WeightGrams)
		assert.Equal(t, result, created)
		assert.Equal(t, "shipping.shipment.created", publishedSubject)
	})

	t.Run("defaults currency to USD", func(t *testing.T) {
		repo := &mockShipmentRepo{
			createFn: func(ctx context.Context, shipment *domain.Shipment) error {
				return nil
			},
		}
		pub := &mockEventPublisher{
			publishFn: func(ctx context.Context, subject string, data interface{}) error {
				return nil
			},
		}
		uc := NewShipmentUseCase(repo, pub)

		req := CreateShipmentRequest{
			OrderID:     "order-1",
			SellerID:    "seller-1",
			CarrierCode: "ups",
			Currency:    "", // empty - should default to USD
		}

		result, err := uc.CreateShipment(context.Background(), req)

		require.NoError(t, err)
		assert.Equal(t, "USD", result.Currency)
	})

	t.Run("items get UUIDs assigned", func(t *testing.T) {
		repo := &mockShipmentRepo{
			createFn: func(ctx context.Context, shipment *domain.Shipment) error {
				return nil
			},
		}
		pub := &mockEventPublisher{
			publishFn: func(ctx context.Context, subject string, data interface{}) error {
				return nil
			},
		}
		uc := NewShipmentUseCase(repo, pub)

		req := CreateShipmentRequest{
			OrderID:     "order-1",
			SellerID:    "seller-1",
			CarrierCode: "fedex",
			Items: []domain.ShipmentItem{
				{ProductID: "prod-1", ProductName: "Widget", Quantity: 1},
				{ProductID: "prod-2", ProductName: "Gadget", Quantity: 3},
			},
		}

		result, err := uc.CreateShipment(context.Background(), req)

		require.NoError(t, err)
		require.Len(t, result.Items, 2)
		for _, item := range result.Items {
			assert.NotEmpty(t, item.ID, "each item should have an assigned ID")
			assert.Equal(t, result.ID, item.ShipmentID, "each item should reference the shipment")
		}
		// IDs should be unique
		assert.NotEqual(t, result.Items[0].ID, result.Items[1].ID)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &mockShipmentRepo{
			createFn: func(ctx context.Context, shipment *domain.Shipment) error {
				return fmt.Errorf("db error")
			},
		}
		pub := &mockEventPublisher{
			publishFn: func(ctx context.Context, subject string, data interface{}) error {
				return nil
			},
		}
		uc := NewShipmentUseCase(repo, pub)

		_, err := uc.CreateShipment(context.Background(), CreateShipmentRequest{
			OrderID:     "order-1",
			SellerID:    "seller-1",
			CarrierCode: "fedex",
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create shipment")
	})
}

func TestShipmentUseCase_GetShipment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := &domain.Shipment{ID: "ship-1", OrderID: "order-1"}
		repo := &mockShipmentRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.Shipment, error) {
				assert.Equal(t, "ship-1", id)
				return expected, nil
			},
		}
		uc := NewShipmentUseCase(repo, nil)

		result, err := uc.GetShipment(context.Background(), "ship-1")

		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("not found", func(t *testing.T) {
		repo := &mockShipmentRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.Shipment, error) {
				return nil, fmt.Errorf("not found")
			},
		}
		uc := NewShipmentUseCase(repo, nil)

		result, err := uc.GetShipment(context.Background(), "nope")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestShipmentUseCase_GetShipmentByTracking(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := &domain.Shipment{ID: "ship-1", TrackingNumber: "FEDEX123"}
		repo := &mockShipmentRepo{
			getByTrackingNumberFn: func(ctx context.Context, trackingNumber string) (*domain.Shipment, error) {
				assert.Equal(t, "FEDEX123", trackingNumber)
				return expected, nil
			},
		}
		uc := NewShipmentUseCase(repo, nil)

		result, err := uc.GetShipmentByTracking(context.Background(), "FEDEX123")

		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestShipmentUseCase_GetShipmentsByOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := []domain.Shipment{
			{ID: "ship-1", OrderID: "order-1"},
			{ID: "ship-2", OrderID: "order-1"},
		}
		repo := &mockShipmentRepo{
			getByOrderIDFn: func(ctx context.Context, orderID string) ([]domain.Shipment, error) {
				assert.Equal(t, "order-1", orderID)
				return expected, nil
			},
		}
		uc := NewShipmentUseCase(repo, nil)

		result, err := uc.GetShipmentsByOrder(context.Background(), "order-1")

		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestShipmentUseCase_ListSellerShipments(t *testing.T) {
	tests := []struct {
		name             string
		page             int
		pageSize         int
		expectedPage     int
		expectedPageSize int
	}{
		{"normal values", 2, 25, 2, 25},
		{"page 0 clamps to 1", 0, 20, 1, 20},
		{"negative page clamps to 1", -5, 20, 1, 20},
		{"pageSize 0 defaults to 20", 1, 0, 1, 20},
		{"pageSize over 100 defaults to 20", 1, 200, 1, 20},
		{"negative pageSize defaults to 20", 1, -1, 1, 20},
		{"pageSize exactly 100 is valid", 1, 100, 1, 100},
		{"pageSize exactly 1 is valid", 1, 1, 1, 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			expected := []domain.Shipment{{ID: "ship-1"}}
			var capturedPage, capturedPageSize int
			repo := &mockShipmentRepo{
				listBySellerFn: func(ctx context.Context, sellerID string, page, pageSize int) ([]domain.Shipment, int64, error) {
					capturedPage = page
					capturedPageSize = pageSize
					return expected, 1, nil
				},
			}
			uc := NewShipmentUseCase(repo, nil)

			result, total, err := uc.ListSellerShipments(context.Background(), "seller-1", tc.page, tc.pageSize)

			require.NoError(t, err)
			assert.Equal(t, expected, result)
			assert.Equal(t, int64(1), total)
			assert.Equal(t, tc.expectedPage, capturedPage, "page should be clamped")
			assert.Equal(t, tc.expectedPageSize, capturedPageSize, "pageSize should be clamped")
		})
	}
}
