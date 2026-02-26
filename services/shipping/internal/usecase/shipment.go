package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/southern-martin/ecommerce/services/shipping/internal/domain"
)

// ShipmentUseCase handles shipment operations.
type ShipmentUseCase struct {
	shipmentRepo domain.ShipmentRepository
	publisher    domain.EventPublisher
}

// NewShipmentUseCase creates a new ShipmentUseCase.
func NewShipmentUseCase(shipmentRepo domain.ShipmentRepository, publisher domain.EventPublisher) *ShipmentUseCase {
	return &ShipmentUseCase{
		shipmentRepo: shipmentRepo,
		publisher:    publisher,
	}
}

// CreateShipmentRequest is the input for creating a shipment.
type CreateShipmentRequest struct {
	OrderID     string
	SellerID    string
	CarrierCode string
	ServiceCode string
	Origin      domain.Address
	Destination domain.Address
	WeightGrams int
	RateCents   int64
	Currency    string
	Items       []domain.ShipmentItem
}

// CreateShipment creates a new shipment.
func (uc *ShipmentUseCase) CreateShipment(ctx context.Context, req CreateShipmentRequest) (*domain.Shipment, error) {
	shipmentID := uuid.New().String()

	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}

	// Assign IDs to items
	for i := range req.Items {
		req.Items[i].ID = uuid.New().String()
		req.Items[i].ShipmentID = shipmentID
	}

	shipment := &domain.Shipment{
		ID:          shipmentID,
		OrderID:     req.OrderID,
		SellerID:    req.SellerID,
		CarrierCode: req.CarrierCode,
		ServiceCode: req.ServiceCode,
		Status:      domain.StatusPending,
		Origin:      req.Origin,
		Destination: req.Destination,
		WeightGrams: req.WeightGrams,
		RateCents:   req.RateCents,
		Currency:    currency,
		Items:       req.Items,
	}

	if err := uc.shipmentRepo.Create(ctx, shipment); err != nil {
		return nil, fmt.Errorf("failed to create shipment: %w", err)
	}

	// Publish event
	_ = uc.publisher.Publish(ctx, "shipping.shipment.created", map[string]interface{}{
		"shipment_id":     shipment.ID,
		"order_id":        shipment.OrderID,
		"seller_id":       shipment.SellerID,
		"carrier_code":    shipment.CarrierCode,
		"tracking_number": shipment.TrackingNumber,
	})

	return shipment, nil
}

// GetShipment retrieves a shipment by ID.
func (uc *ShipmentUseCase) GetShipment(ctx context.Context, id string) (*domain.Shipment, error) {
	return uc.shipmentRepo.GetByID(ctx, id)
}

// GetShipmentByTracking retrieves a shipment by tracking number.
func (uc *ShipmentUseCase) GetShipmentByTracking(ctx context.Context, trackingNumber string) (*domain.Shipment, error) {
	return uc.shipmentRepo.GetByTrackingNumber(ctx, trackingNumber)
}

// GetShipmentsByOrder retrieves all shipments for an order.
func (uc *ShipmentUseCase) GetShipmentsByOrder(ctx context.Context, orderID string) ([]domain.Shipment, error) {
	return uc.shipmentRepo.GetByOrderID(ctx, orderID)
}

// ListSellerShipments lists shipments for a seller with pagination.
func (uc *ShipmentUseCase) ListSellerShipments(ctx context.Context, sellerID string, page, pageSize int) ([]domain.Shipment, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.shipmentRepo.ListBySeller(ctx, sellerID, page, pageSize)
}
