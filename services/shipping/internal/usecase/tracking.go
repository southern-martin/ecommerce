package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/southern-martin/ecommerce/services/shipping/internal/domain"
)

// TrackingUseCase handles tracking event operations.
type TrackingUseCase struct {
	shipmentRepo domain.ShipmentRepository
	trackingRepo domain.TrackingEventRepository
	publisher    domain.EventPublisher
}

// NewTrackingUseCase creates a new TrackingUseCase.
func NewTrackingUseCase(
	shipmentRepo domain.ShipmentRepository,
	trackingRepo domain.TrackingEventRepository,
	publisher domain.EventPublisher,
) *TrackingUseCase {
	return &TrackingUseCase{
		shipmentRepo: shipmentRepo,
		trackingRepo: trackingRepo,
		publisher:    publisher,
	}
}

// AddTrackingEventRequest is the input for adding a tracking event.
type AddTrackingEventRequest struct {
	ShipmentID  string
	Status      string
	Description string
	Location    string
	EventAt     time.Time
}

// AddTrackingEvent adds a tracking event and updates shipment status.
func (uc *TrackingUseCase) AddTrackingEvent(ctx context.Context, req AddTrackingEventRequest) (*domain.TrackingEvent, error) {
	shipment, err := uc.shipmentRepo.GetByID(ctx, req.ShipmentID)
	if err != nil {
		return nil, fmt.Errorf("shipment not found: %w", err)
	}

	newStatus := domain.ShipmentStatus(req.Status)
	if !domain.CanTransition(shipment.Status, newStatus) {
		return nil, fmt.Errorf("invalid status transition from %s to %s", shipment.Status, newStatus)
	}

	eventAt := req.EventAt
	if eventAt.IsZero() {
		eventAt = time.Now()
	}

	event := &domain.TrackingEvent{
		ID:          uuid.New().String(),
		ShipmentID:  req.ShipmentID,
		Status:      req.Status,
		Description: req.Description,
		Location:    req.Location,
		EventAt:     eventAt,
	}

	if err := uc.trackingRepo.Create(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to create tracking event: %w", err)
	}

	// Update shipment status
	shipment.Status = newStatus
	if err := uc.shipmentRepo.Update(ctx, shipment); err != nil {
		return nil, fmt.Errorf("failed to update shipment status: %w", err)
	}

	// Publish appropriate event
	eventSubject := "shipping.shipment.updated"
	if newStatus == domain.StatusDelivered {
		eventSubject = "shipping.shipment.delivered"
	}

	_ = uc.publisher.Publish(ctx, eventSubject, map[string]interface{}{
		"shipment_id":     shipment.ID,
		"order_id":        shipment.OrderID,
		"status":          string(newStatus),
		"tracking_number": shipment.TrackingNumber,
		"description":     req.Description,
		"location":        req.Location,
	})

	return event, nil
}

// GetTrackingEvents retrieves all tracking events for a shipment.
func (uc *TrackingUseCase) GetTrackingEvents(ctx context.Context, shipmentID string) ([]domain.TrackingEvent, error) {
	return uc.trackingRepo.GetByShipmentID(ctx, shipmentID)
}
