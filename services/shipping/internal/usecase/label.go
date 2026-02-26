package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/southern-martin/ecommerce/services/shipping/internal/domain"
)

// LabelUseCase handles label generation.
type LabelUseCase struct {
	shipmentRepo domain.ShipmentRepository
	publisher    domain.EventPublisher
}

// NewLabelUseCase creates a new LabelUseCase.
func NewLabelUseCase(shipmentRepo domain.ShipmentRepository, publisher domain.EventPublisher) *LabelUseCase {
	return &LabelUseCase{
		shipmentRepo: shipmentRepo,
		publisher:    publisher,
	}
}

// GenerateLabel generates a shipping label for a shipment (mock implementation).
func (uc *LabelUseCase) GenerateLabel(ctx context.Context, shipmentID, sellerID string) (*domain.Shipment, error) {
	shipment, err := uc.shipmentRepo.GetByID(ctx, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("shipment not found: %w", err)
	}

	if shipment.SellerID != sellerID {
		return nil, fmt.Errorf("unauthorized: shipment belongs to different seller")
	}

	if shipment.Status != domain.StatusPending {
		return nil, fmt.Errorf("label can only be generated for pending shipments")
	}

	// Mock label generation
	trackingNumber := fmt.Sprintf("%s%s", shipment.CarrierCode, uuid.New().String()[:12])
	labelURL := fmt.Sprintf("https://labels.example.com/%s/%s.pdf", shipment.CarrierCode, trackingNumber)

	shipment.TrackingNumber = trackingNumber
	shipment.LabelURL = labelURL
	shipment.Status = domain.StatusLabelCreated

	if err := uc.shipmentRepo.Update(ctx, shipment); err != nil {
		return nil, fmt.Errorf("failed to update shipment: %w", err)
	}

	// Publish event
	_ = uc.publisher.Publish(ctx, "shipping.shipment.updated", map[string]interface{}{
		"shipment_id":     shipment.ID,
		"order_id":        shipment.OrderID,
		"status":          string(shipment.Status),
		"tracking_number": shipment.TrackingNumber,
		"label_url":       shipment.LabelURL,
	})

	return shipment, nil
}
