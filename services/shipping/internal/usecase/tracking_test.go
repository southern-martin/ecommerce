package usecase

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/southern-martin/ecommerce/services/shipping/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrackingUseCase_AddTrackingEvent(t *testing.T) {
	t.Run("success with valid transition", func(t *testing.T) {
		shipment := &domain.Shipment{
			ID:             "ship-1",
			OrderID:        "order-1",
			SellerID:       "seller-1",
			TrackingNumber: "FEDEX123",
			Status:         domain.StatusPending,
		}
		var createdEvent *domain.TrackingEvent
		var updatedShipment *domain.Shipment
		var publishedSubject string
		var publishedData interface{}

		shipmentRepo := &mockShipmentRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.Shipment, error) {
				assert.Equal(t, "ship-1", id)
				return shipment, nil
			},
			updateFn: func(ctx context.Context, s *domain.Shipment) error {
				updatedShipment = s
				return nil
			},
		}
		trackingRepo := &mockTrackingEventRepo{
			createFn: func(ctx context.Context, event *domain.TrackingEvent) error {
				createdEvent = event
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

		uc := NewTrackingUseCase(shipmentRepo, trackingRepo, pub)

		eventTime := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
		req := AddTrackingEventRequest{
			ShipmentID:  "ship-1",
			Status:      "label_created",
			Description: "Label has been created",
			Location:    "Austin, TX",
			EventAt:     eventTime,
		}

		result, err := uc.AddTrackingEvent(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, "ship-1", result.ShipmentID)
		assert.Equal(t, "label_created", result.Status)
		assert.Equal(t, "Label has been created", result.Description)
		assert.Equal(t, "Austin, TX", result.Location)
		assert.Equal(t, eventTime, result.EventAt)

		// Event was persisted
		assert.Equal(t, result, createdEvent)

		// Shipment status was updated
		assert.Equal(t, domain.StatusLabelCreated, updatedShipment.Status)

		// Event published
		assert.Equal(t, "shipping.shipment.updated", publishedSubject)
		dataMap, ok := publishedData.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "ship-1", dataMap["shipment_id"])
		assert.Equal(t, "label_created", dataMap["status"])
	})

	t.Run("defaults EventAt to now if zero", func(t *testing.T) {
		shipment := &domain.Shipment{
			ID:     "ship-1",
			Status: domain.StatusPickedUp,
		}
		var createdEvent *domain.TrackingEvent

		shipmentRepo := &mockShipmentRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.Shipment, error) {
				return shipment, nil
			},
			updateFn: func(ctx context.Context, s *domain.Shipment) error {
				return nil
			},
		}
		trackingRepo := &mockTrackingEventRepo{
			createFn: func(ctx context.Context, event *domain.TrackingEvent) error {
				createdEvent = event
				return nil
			},
		}
		pub := &mockEventPublisher{
			publishFn: func(ctx context.Context, subject string, data interface{}) error {
				return nil
			},
		}

		uc := NewTrackingUseCase(shipmentRepo, trackingRepo, pub)

		before := time.Now()
		req := AddTrackingEventRequest{
			ShipmentID: "ship-1",
			Status:     "in_transit",
			// EventAt is zero
		}

		result, err := uc.AddTrackingEvent(context.Background(), req)

		require.NoError(t, err)
		assert.False(t, createdEvent.EventAt.IsZero(), "EventAt should be set to now")
		assert.True(t, result.EventAt.After(before) || result.EventAt.Equal(before),
			"EventAt should be at or after the time we recorded before the call")
	})

	t.Run("invalid transition error", func(t *testing.T) {
		shipment := &domain.Shipment{
			ID:     "ship-1",
			Status: domain.StatusPending,
		}
		shipmentRepo := &mockShipmentRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.Shipment, error) {
				return shipment, nil
			},
		}
		uc := NewTrackingUseCase(shipmentRepo, nil, nil)

		req := AddTrackingEventRequest{
			ShipmentID: "ship-1",
			Status:     "delivered", // pending -> delivered is invalid
		}

		result, err := uc.AddTrackingEvent(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid status transition")
		assert.Contains(t, err.Error(), "pending")
		assert.Contains(t, err.Error(), "delivered")
	})

	t.Run("delivered triggers different event subject", func(t *testing.T) {
		shipment := &domain.Shipment{
			ID:             "ship-1",
			OrderID:        "order-1",
			TrackingNumber: "UPS456",
			Status:         domain.StatusInTransit,
		}
		var publishedSubject string

		shipmentRepo := &mockShipmentRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.Shipment, error) {
				return shipment, nil
			},
			updateFn: func(ctx context.Context, s *domain.Shipment) error {
				return nil
			},
		}
		trackingRepo := &mockTrackingEventRepo{
			createFn: func(ctx context.Context, event *domain.TrackingEvent) error {
				return nil
			},
		}
		pub := &mockEventPublisher{
			publishFn: func(ctx context.Context, subject string, data interface{}) error {
				publishedSubject = subject
				return nil
			},
		}

		uc := NewTrackingUseCase(shipmentRepo, trackingRepo, pub)

		req := AddTrackingEventRequest{
			ShipmentID:  "ship-1",
			Status:      "delivered",
			Description: "Package delivered",
			Location:    "Front door",
			EventAt:     time.Now(),
		}

		result, err := uc.AddTrackingEvent(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "shipping.shipment.delivered", publishedSubject)
	})

	t.Run("shipment not found error", func(t *testing.T) {
		shipmentRepo := &mockShipmentRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.Shipment, error) {
				return nil, fmt.Errorf("not found")
			},
		}
		uc := NewTrackingUseCase(shipmentRepo, nil, nil)

		req := AddTrackingEventRequest{
			ShipmentID: "nope",
			Status:     "in_transit",
		}

		result, err := uc.AddTrackingEvent(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "shipment not found")
	})

	t.Run("tracking repo create error", func(t *testing.T) {
		shipment := &domain.Shipment{
			ID:     "ship-1",
			Status: domain.StatusPending,
		}
		shipmentRepo := &mockShipmentRepo{
			getByIDFn: func(ctx context.Context, id string) (*domain.Shipment, error) {
				return shipment, nil
			},
		}
		trackingRepo := &mockTrackingEventRepo{
			createFn: func(ctx context.Context, event *domain.TrackingEvent) error {
				return fmt.Errorf("db write error")
			},
		}
		uc := NewTrackingUseCase(shipmentRepo, trackingRepo, nil)

		req := AddTrackingEventRequest{
			ShipmentID: "ship-1",
			Status:     "label_created",
			EventAt:    time.Now(),
		}

		result, err := uc.AddTrackingEvent(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to create tracking event")
	})
}

func TestTrackingUseCase_GetTrackingEvents(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := []domain.TrackingEvent{
			{ID: "evt-1", ShipmentID: "ship-1", Status: "label_created"},
			{ID: "evt-2", ShipmentID: "ship-1", Status: "picked_up"},
		}
		trackingRepo := &mockTrackingEventRepo{
			getByShipmentIDFn: func(ctx context.Context, shipmentID string) ([]domain.TrackingEvent, error) {
				assert.Equal(t, "ship-1", shipmentID)
				return expected, nil
			},
		}
		uc := NewTrackingUseCase(nil, trackingRepo, nil)

		result, err := uc.GetTrackingEvents(context.Background(), "ship-1")

		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("repo error", func(t *testing.T) {
		trackingRepo := &mockTrackingEventRepo{
			getByShipmentIDFn: func(ctx context.Context, shipmentID string) ([]domain.TrackingEvent, error) {
				return nil, fmt.Errorf("db error")
			},
		}
		uc := NewTrackingUseCase(nil, trackingRepo, nil)

		result, err := uc.GetTrackingEvents(context.Background(), "ship-1")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
