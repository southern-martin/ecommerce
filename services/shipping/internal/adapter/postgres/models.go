package postgres

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/southern-martin/ecommerce/services/shipping/internal/domain"
)

// AddressJSON is a GORM-compatible JSONB type for Address.
type AddressJSON domain.Address

// Value implements the driver.Valuer interface for JSONB storage.
func (a AddressJSON) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan implements the sql.Scanner interface for JSONB retrieval.
func (a *AddressJSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan AddressJSON: not a byte slice")
	}
	return json.Unmarshal(bytes, a)
}

// CarrierModel is the GORM model for the carriers table.
type CarrierModel struct {
	Code               string         `gorm:"type:varchar(20);primaryKey"`
	Name               string         `gorm:"type:varchar(100);not null"`
	IsActive           bool           `gorm:"default:true"`
	SupportedCountries pq.StringArray `gorm:"type:text[]"`
	APIBaseURL         string         `gorm:"type:varchar(500)"`
	CreatedAt          time.Time      `gorm:"autoCreateTime"`
}

func (CarrierModel) TableName() string { return "carriers" }

func (m *CarrierModel) ToDomain() *domain.Carrier {
	return &domain.Carrier{
		Code:               m.Code,
		Name:               m.Name,
		IsActive:           m.IsActive,
		SupportedCountries: []string(m.SupportedCountries),
		APIBaseURL:         m.APIBaseURL,
		CreatedAt:          m.CreatedAt,
	}
}

func ToCarrierModel(c *domain.Carrier) *CarrierModel {
	return &CarrierModel{
		Code:               c.Code,
		Name:               c.Name,
		IsActive:           c.IsActive,
		SupportedCountries: pq.StringArray(c.SupportedCountries),
		APIBaseURL:         c.APIBaseURL,
		CreatedAt:          c.CreatedAt,
	}
}

// CarrierCredentialModel is the GORM model for the carrier_credentials table.
type CarrierCredentialModel struct {
	ID          string `gorm:"type:uuid;primaryKey"`
	SellerID    string `gorm:"type:uuid;not null;uniqueIndex:idx_seller_carrier"`
	CarrierCode string `gorm:"type:varchar(20);not null;uniqueIndex:idx_seller_carrier"`
	Credentials string `gorm:"type:jsonb;not null"`
	IsActive    bool   `gorm:"default:true"`
}

func (CarrierCredentialModel) TableName() string { return "carrier_credentials" }

func (m *CarrierCredentialModel) ToDomain() *domain.CarrierCredential {
	return &domain.CarrierCredential{
		ID:          m.ID,
		SellerID:    m.SellerID,
		CarrierCode: m.CarrierCode,
		Credentials: m.Credentials,
		IsActive:    m.IsActive,
	}
}

func ToCarrierCredentialModel(c *domain.CarrierCredential) *CarrierCredentialModel {
	return &CarrierCredentialModel{
		ID:          c.ID,
		SellerID:    c.SellerID,
		CarrierCode: c.CarrierCode,
		Credentials: c.Credentials,
		IsActive:    c.IsActive,
	}
}

// ShipmentModel is the GORM model for the shipments table.
type ShipmentModel struct {
	ID             string      `gorm:"type:uuid;primaryKey"`
	OrderID        string      `gorm:"type:uuid;index;not null"`
	SellerID       string      `gorm:"type:uuid;index;not null"`
	CarrierCode    string      `gorm:"type:varchar(20)"`
	ServiceCode    string      `gorm:"type:varchar(50)"`
	TrackingNumber string      `gorm:"type:varchar(100);index"`
	LabelURL       string      `gorm:"type:text"`
	Status         string      `gorm:"type:varchar(20);default:'pending'"`
	Origin         AddressJSON `gorm:"type:jsonb;not null"`
	Destination    AddressJSON `gorm:"type:jsonb;not null"`
	WeightGrams    int         `gorm:"default:0"`
	RateCents      int64       `gorm:"default:0"`
	Currency       string      `gorm:"type:varchar(3);default:'USD'"`
	Items          []ShipmentItemModel  `gorm:"foreignKey:ShipmentID;constraint:OnDelete:CASCADE"`
	TrackingEvents []TrackingEventModel `gorm:"foreignKey:ShipmentID;constraint:OnDelete:CASCADE"`
	CreatedAt      time.Time   `gorm:"autoCreateTime"`
	UpdatedAt      time.Time   `gorm:"autoUpdateTime"`
}

func (ShipmentModel) TableName() string { return "shipments" }

func (m *ShipmentModel) ToDomain() *domain.Shipment {
	s := &domain.Shipment{
		ID:             m.ID,
		OrderID:        m.OrderID,
		SellerID:       m.SellerID,
		CarrierCode:    m.CarrierCode,
		ServiceCode:    m.ServiceCode,
		TrackingNumber: m.TrackingNumber,
		LabelURL:       m.LabelURL,
		Status:         domain.ShipmentStatus(m.Status),
		Origin:         domain.Address(m.Origin),
		Destination:    domain.Address(m.Destination),
		WeightGrams:    m.WeightGrams,
		RateCents:      m.RateCents,
		Currency:       m.Currency,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
	for _, item := range m.Items {
		s.Items = append(s.Items, *item.ToDomain())
	}
	for _, evt := range m.TrackingEvents {
		s.TrackingEvents = append(s.TrackingEvents, *evt.ToDomain())
	}
	return s
}

func ToShipmentModel(s *domain.Shipment) *ShipmentModel {
	model := &ShipmentModel{
		ID:             s.ID,
		OrderID:        s.OrderID,
		SellerID:       s.SellerID,
		CarrierCode:    s.CarrierCode,
		ServiceCode:    s.ServiceCode,
		TrackingNumber: s.TrackingNumber,
		LabelURL:       s.LabelURL,
		Status:         string(s.Status),
		Origin:         AddressJSON(s.Origin),
		Destination:    AddressJSON(s.Destination),
		WeightGrams:    s.WeightGrams,
		RateCents:      s.RateCents,
		Currency:       s.Currency,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}
	for _, item := range s.Items {
		model.Items = append(model.Items, *ToShipmentItemModel(&item))
	}
	return model
}

// ShipmentItemModel is the GORM model for the shipment_items table.
type ShipmentItemModel struct {
	ID          string `gorm:"type:uuid;primaryKey"`
	ShipmentID  string `gorm:"type:uuid;index;not null"`
	ProductID   string `gorm:"type:uuid;not null"`
	VariantID   string `gorm:"type:uuid"`
	ProductName string `gorm:"type:varchar(255)"`
	Quantity    int    `gorm:"not null;default:1"`
}

func (ShipmentItemModel) TableName() string { return "shipment_items" }

func (m *ShipmentItemModel) ToDomain() *domain.ShipmentItem {
	return &domain.ShipmentItem{
		ID:          m.ID,
		ShipmentID:  m.ShipmentID,
		ProductID:   m.ProductID,
		VariantID:   m.VariantID,
		ProductName: m.ProductName,
		Quantity:    m.Quantity,
	}
}

func ToShipmentItemModel(i *domain.ShipmentItem) *ShipmentItemModel {
	return &ShipmentItemModel{
		ID:          i.ID,
		ShipmentID:  i.ShipmentID,
		ProductID:   i.ProductID,
		VariantID:   i.VariantID,
		ProductName: i.ProductName,
		Quantity:    i.Quantity,
	}
}

// TrackingEventModel is the GORM model for the tracking_events table.
type TrackingEventModel struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	ShipmentID  string    `gorm:"type:uuid;index;not null"`
	Status      string    `gorm:"type:varchar(50);not null"`
	Description string    `gorm:"type:text"`
	Location    string    `gorm:"type:varchar(200)"`
	EventAt     time.Time `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (TrackingEventModel) TableName() string { return "tracking_events" }

func (m *TrackingEventModel) ToDomain() *domain.TrackingEvent {
	return &domain.TrackingEvent{
		ID:          m.ID,
		ShipmentID:  m.ShipmentID,
		Status:      m.Status,
		Description: m.Description,
		Location:    m.Location,
		EventAt:     m.EventAt,
		CreatedAt:   m.CreatedAt,
	}
}

func ToTrackingEventModel(e *domain.TrackingEvent) *TrackingEventModel {
	return &TrackingEventModel{
		ID:          e.ID,
		ShipmentID:  e.ShipmentID,
		Status:      e.Status,
		Description: e.Description,
		Location:    e.Location,
		EventAt:     e.EventAt,
		CreatedAt:   e.CreatedAt,
	}
}
