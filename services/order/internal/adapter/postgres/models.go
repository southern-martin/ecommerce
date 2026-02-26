package postgres

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/southern-martin/ecommerce/services/order/internal/domain"
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

// OrderModel is the GORM model for the orders table.
type OrderModel struct {
	ID              string      `gorm:"type:uuid;primaryKey"`
	OrderNumber     string      `gorm:"type:varchar(50);uniqueIndex;not null"`
	BuyerID         string      `gorm:"type:uuid;index;not null"`
	Status          string      `gorm:"type:varchar(20);index;not null;default:'pending'"`
	SubtotalCents   int64       `gorm:"not null;default:0"`
	ShippingCents   int64       `gorm:"not null;default:0"`
	TaxCents        int64       `gorm:"not null;default:0"`
	DiscountCents   int64       `gorm:"not null;default:0"`
	TotalCents      int64       `gorm:"not null;default:0"`
	Currency        string      `gorm:"type:varchar(3);not null;default:'USD'"`
	ShippingAddress AddressJSON `gorm:"type:jsonb"`
	Items           []OrderItemModel  `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	SellerOrders    []SellerOrderModel `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	CreatedAt       time.Time   `gorm:"autoCreateTime"`
	UpdatedAt       time.Time   `gorm:"autoUpdateTime"`
}

// TableName returns the table name for OrderModel.
func (OrderModel) TableName() string {
	return "orders"
}

// OrderItemModel is the GORM model for the order_items table.
type OrderItemModel struct {
	ID             string `gorm:"type:uuid;primaryKey"`
	OrderID        string `gorm:"type:uuid;index;not null"`
	ProductID      string `gorm:"type:uuid;not null"`
	VariantID      *string `gorm:"type:uuid"`
	ProductName    string `gorm:"type:varchar(255);not null"`
	VariantName    string `gorm:"type:varchar(255)"`
	SKU            string `gorm:"type:varchar(100)"`
	Quantity       int    `gorm:"not null;default:1"`
	UnitPriceCents int64  `gorm:"not null;default:0"`
	TotalCents     int64  `gorm:"not null;default:0"`
	SellerID       string `gorm:"type:uuid;index;not null"`
	ImageURL       string `gorm:"type:text"`
}

// TableName returns the table name for OrderItemModel.
func (OrderItemModel) TableName() string {
	return "order_items"
}

// SellerOrderModel is the GORM model for the seller_orders table.
type SellerOrderModel struct {
	ID            string `gorm:"type:uuid;primaryKey"`
	OrderID       string `gorm:"type:uuid;index;not null"`
	SellerID      string `gorm:"type:uuid;index;not null"`
	Status        string `gorm:"type:varchar(20);not null;default:'pending'"`
	SubtotalCents int64  `gorm:"not null;default:0"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

// TableName returns the table name for SellerOrderModel.
func (SellerOrderModel) TableName() string {
	return "seller_orders"
}

// ToDomain converts an OrderModel to a domain Order.
func (m *OrderModel) ToDomain() *domain.Order {
	order := &domain.Order{
		ID:              m.ID,
		OrderNumber:     m.OrderNumber,
		BuyerID:         m.BuyerID,
		Status:          domain.OrderStatus(m.Status),
		SubtotalCents:   m.SubtotalCents,
		ShippingCents:   m.ShippingCents,
		TaxCents:        m.TaxCents,
		DiscountCents:   m.DiscountCents,
		TotalCents:      m.TotalCents,
		Currency:        m.Currency,
		ShippingAddress: domain.Address(m.ShippingAddress),
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}

	for _, item := range m.Items {
		order.Items = append(order.Items, item.ToDomain())
	}

	for _, so := range m.SellerOrders {
		order.SellerOrders = append(order.SellerOrders, *so.ToDomain())
	}

	return order
}

// ToModel converts a domain Order to an OrderModel.
func ToOrderModel(o *domain.Order) *OrderModel {
	model := &OrderModel{
		ID:              o.ID,
		OrderNumber:     o.OrderNumber,
		BuyerID:         o.BuyerID,
		Status:          string(o.Status),
		SubtotalCents:   o.SubtotalCents,
		ShippingCents:   o.ShippingCents,
		TaxCents:        o.TaxCents,
		DiscountCents:   o.DiscountCents,
		TotalCents:      o.TotalCents,
		Currency:        o.Currency,
		ShippingAddress: AddressJSON(o.ShippingAddress),
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
	}

	for _, item := range o.Items {
		model.Items = append(model.Items, ToOrderItemModel(&item))
	}

	return model
}

// ToDomain converts an OrderItemModel to a domain OrderItem.
func (m *OrderItemModel) ToDomain() domain.OrderItem {
	variantID := ""
	if m.VariantID != nil {
		variantID = *m.VariantID
	}
	return domain.OrderItem{
		ID:             m.ID,
		OrderID:        m.OrderID,
		ProductID:      m.ProductID,
		VariantID:      variantID,
		ProductName:    m.ProductName,
		VariantName:    m.VariantName,
		SKU:            m.SKU,
		Quantity:       m.Quantity,
		UnitPriceCents: m.UnitPriceCents,
		TotalCents:     m.TotalCents,
		SellerID:       m.SellerID,
		ImageURL:       m.ImageURL,
	}
}

// ToOrderItemModel converts a domain OrderItem to an OrderItemModel.
func ToOrderItemModel(item *domain.OrderItem) OrderItemModel {
	var variantID *string
	if item.VariantID != "" {
		v := item.VariantID
		variantID = &v
	}
	return OrderItemModel{
		ID:             item.ID,
		OrderID:        item.OrderID,
		ProductID:      item.ProductID,
		VariantID:      variantID,
		ProductName:    item.ProductName,
		VariantName:    item.VariantName,
		SKU:            item.SKU,
		Quantity:       item.Quantity,
		UnitPriceCents: item.UnitPriceCents,
		TotalCents:     item.TotalCents,
		SellerID:       item.SellerID,
		ImageURL:       item.ImageURL,
	}
}

// ToDomain converts a SellerOrderModel to a domain SellerOrder.
func (m *SellerOrderModel) ToDomain() *domain.SellerOrder {
	return &domain.SellerOrder{
		ID:            m.ID,
		OrderID:       m.OrderID,
		SellerID:      m.SellerID,
		Status:        domain.OrderStatus(m.Status),
		SubtotalCents: m.SubtotalCents,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

// ToSellerOrderModel converts a domain SellerOrder to a SellerOrderModel.
func ToSellerOrderModel(so *domain.SellerOrder) *SellerOrderModel {
	return &SellerOrderModel{
		ID:            so.ID,
		OrderID:       so.OrderID,
		SellerID:      so.SellerID,
		Status:        string(so.Status),
		SubtotalCents: so.SubtotalCents,
		CreatedAt:     so.CreatedAt,
		UpdatedAt:     so.UpdatedAt,
	}
}
