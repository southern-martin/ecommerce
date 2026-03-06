package usecase

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/southern-martin/ecommerce/services/shipping/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLabelUseCase_GenerateLabel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		shipment := &domain.Shipment{
			ID:          "ship-1",
			OrderID:     "order-1",
			SellerID:    "seller-1",
			CarrierCode: "fedex",
			Status:      domain.StatusPending,
		}
		var updatedShipment *domain.Shipment
		var publishedSubject string
		var publishedData interface{}
		repo := &mockShipmentRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.Shipment, error) {
				assert.Equal(t, "ship-1", id)
				return shipment, nil
			},
			updateFn: func(ctx context.Context, s *domain.Shipment) error {
				updatedShipment = s
				return nil
			},
		}
		pub := &mockEventPublisher{
			publishFn: func(ctx context.Context, subject string, data interface{}) error {
				publishedSubject = subject
				publishedData = data
				return nil
			},
		}
		uc := NewLabelUseCase(repo, pub)

		result, err := uc.GenerateLabel(context.Background(), "ship-1", "seller-1")

		require.NoError(t, err)
		require.NotNil(t, result)

		// Tracking number format: carrierCode + first 12 chars of UUID
		assert.True(t, strings.HasPrefix(result.TrackingNumber, "fedex"),
			"tracking number should start with carrier code, got: %s", result.TrackingNumber)
		assert.Len(t, result.TrackingNumber, len("fedex")+12,
			"tracking number should be carrierCode + 12 UUID chars")

		// Label URL format
		assert.True(t, strings.HasPrefix(result.LabelURL, "https://labels.example.com/fedex/"),
			"label URL should start with expected prefix, got: %s", result.LabelURL)
		assert.True(t, strings.HasSuffix(result.LabelURL, ".pdf"),
			"label URL should end with .pdf")

		// Status changed
		assert.Equal(t, domain.StatusLabelCreated, result.Status)
		assert.Equal(t, result, updatedShipment)

		// Event published
		assert.Equal(t, "shipping.shipment.updated", publishedSubject)
		dataMap, ok := publishedData.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "ship-1", dataMap["shipment_id"])
		assert.Equal(t, "label_created", dataMap["status"])
		assert.NotEmpty(t, dataMap["tracking_number"])
		assert.NotEmpty(t, dataMap["label_url"])
	})

	t.Run("wrong seller error", func(t *testing.T) {
		repo := &mockShipmentRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.Shipment, error) {
				return &domain.Shipment{
					ID:       "ship-1",
					SellerID: "seller-1",
					Status:   domain.StatusPending,
				}, nil
			},
		}
		uc := NewLabelUseCase(repo, nil)

		result, err := uc.GenerateLabel(context.Background(), "ship-1", "seller-other")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("non-pending shipment error", func(t *testing.T) {
		repo := &mockShipmentRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.Shipment, error) {
				return &domain.Shipment{
					ID:       "ship-1",
					SellerID: "seller-1",
					Status:   domain.StatusLabelCreated, // not pending
				}, nil
			},
		}
		uc := NewLabelUseCase(repo, nil)

		result, err := uc.GenerateLabel(context.Background(), "ship-1", "seller-1")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "pending")
	})

	t.Run("shipment not found error", func(t *testing.T) {
		repo := &mockShipmentRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.Shipment, error) {
				return nil, fmt.Errorf("not found")
			},
		}
		uc := NewLabelUseCase(repo, nil)

		result, err := uc.GenerateLabel(context.Background(), "ship-nope", "seller-1")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "shipment not found")
	})

	t.Run("update repo error", func(t *testing.T) {
		repo := &mockShipmentRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.Shipment, error) {
				return &domain.Shipment{
					ID:          "ship-1",
					SellerID:    "seller-1",
					CarrierCode: "ups",
					Status:      domain.StatusPending,
				}, nil
			},
			updateFn: func(ctx context.Context, s *domain.Shipment) error {
				return fmt.Errorf("db write error")
			},
		}
		uc := NewLabelUseCase(repo, nil)

		result, err := uc.GenerateLabel(context.Background(), "ship-1", "seller-1")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to update shipment")
	})
}
