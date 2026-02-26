package domain

import "time"

// ShipmentStatus represents the status of a shipment.
type ShipmentStatus string

const (
	StatusPending      ShipmentStatus = "pending"
	StatusLabelCreated ShipmentStatus = "label_created"
	StatusPickedUp     ShipmentStatus = "picked_up"
	StatusInTransit    ShipmentStatus = "in_transit"
	StatusDelivered    ShipmentStatus = "delivered"
	StatusException    ShipmentStatus = "exception"
)

// AllowedShipmentTransitions defines valid status transitions.
var AllowedShipmentTransitions = map[ShipmentStatus][]ShipmentStatus{
	StatusPending:      {StatusLabelCreated, StatusException},
	StatusLabelCreated: {StatusPickedUp, StatusException},
	StatusPickedUp:     {StatusInTransit, StatusException},
	StatusInTransit:    {StatusDelivered, StatusException},
	StatusException:    {StatusInTransit, StatusDelivered},
}

// CanTransition checks if a status transition is valid.
func CanTransition(from, to ShipmentStatus) bool {
	allowed, ok := AllowedShipmentTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// Carrier represents a shipping carrier.
type Carrier struct {
	Code               string
	Name               string
	IsActive           bool
	SupportedCountries []string
	APIBaseURL         string
	CreatedAt          time.Time
}

// CarrierCredential represents a seller's carrier API credentials.
type CarrierCredential struct {
	ID          string
	SellerID    string
	CarrierCode string
	Credentials string // JSON-encoded credentials
	IsActive    bool
}

// Address represents a shipping address.
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// Shipment represents a shipment.
type Shipment struct {
	ID              string
	OrderID         string
	SellerID        string
	CarrierCode     string
	ServiceCode     string
	TrackingNumber  string
	LabelURL        string
	Status          ShipmentStatus
	Origin          Address
	Destination     Address
	WeightGrams     int
	RateCents       int64
	Currency        string
	Items           []ShipmentItem
	TrackingEvents  []TrackingEvent
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ShipmentItem represents an item within a shipment.
type ShipmentItem struct {
	ID          string
	ShipmentID  string
	ProductID   string
	VariantID   string
	ProductName string
	Quantity    int
}

// TrackingEvent represents a tracking event for a shipment.
type TrackingEvent struct {
	ID          string
	ShipmentID  string
	Status      string
	Description string
	Location    string
	EventAt     time.Time
	CreatedAt   time.Time
}

// ShippingRate represents a rate quote from a carrier.
type ShippingRate struct {
	CarrierCode      string
	ServiceName      string
	RateCents        int64
	Currency         string
	EstimatedDaysMin int
	EstimatedDaysMax int
}
