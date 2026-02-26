package domain

import "context"

// CarrierRepository defines the interface for carrier persistence.
type CarrierRepository interface {
	GetAll(ctx context.Context) ([]Carrier, error)
	GetByCode(ctx context.Context, code string) (*Carrier, error)
	Create(ctx context.Context, carrier *Carrier) error
	Update(ctx context.Context, carrier *Carrier) error
}

// CarrierCredentialRepository defines the interface for carrier credential persistence.
type CarrierCredentialRepository interface {
	GetBySellerAndCarrier(ctx context.Context, sellerID, carrierCode string) (*CarrierCredential, error)
	ListBySeller(ctx context.Context, sellerID string) ([]CarrierCredential, error)
	Create(ctx context.Context, cred *CarrierCredential) error
	Update(ctx context.Context, cred *CarrierCredential) error
	Delete(ctx context.Context, id string) error
}

// ShipmentRepository defines the interface for shipment persistence.
type ShipmentRepository interface {
	GetByID(ctx context.Context, id string) (*Shipment, error)
	GetByOrderID(ctx context.Context, orderID string) ([]Shipment, error)
	GetByTrackingNumber(ctx context.Context, trackingNumber string) (*Shipment, error)
	ListBySeller(ctx context.Context, sellerID string, page, pageSize int) ([]Shipment, int64, error)
	Create(ctx context.Context, shipment *Shipment) error
	Update(ctx context.Context, shipment *Shipment) error
}

// TrackingEventRepository defines the interface for tracking event persistence.
type TrackingEventRepository interface {
	GetByShipmentID(ctx context.Context, shipmentID string) ([]TrackingEvent, error)
	Create(ctx context.Context, event *TrackingEvent) error
}
