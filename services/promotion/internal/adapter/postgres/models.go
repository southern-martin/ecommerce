package postgres

import (
	"time"

	"github.com/lib/pq"
	"github.com/southern-martin/ecommerce/services/promotion/internal/domain"
)

// CouponModel is the GORM model for the coupons table.
type CouponModel struct {
	ID               string         `gorm:"type:uuid;primaryKey"`
	Code             string         `gorm:"type:varchar(100);uniqueIndex;not null"`
	Type             string         `gorm:"type:varchar(20);not null"`
	DiscountValue    int64          `gorm:"not null;default:0"`
	MinOrderCents    int64          `gorm:"not null;default:0"`
	MaxDiscountCents int64          `gorm:"not null;default:0"`
	UsageLimit       int            `gorm:"not null;default:0"`
	UsageCount       int            `gorm:"not null;default:0"`
	PerUserLimit     int            `gorm:"not null;default:0"`
	Scope            string         `gorm:"type:varchar(20);not null;default:'all'"`
	ScopeIDs         pq.StringArray `gorm:"type:text[]"`
	CreatedBy        string         `gorm:"type:varchar(255);not null"`
	StartsAt         time.Time      `gorm:"not null"`
	ExpiresAt        time.Time      `gorm:"not null"`
	IsActive         bool           `gorm:"not null;default:true"`
	CreatedAt        time.Time      `gorm:"autoCreateTime"`
}

// TableName returns the table name for CouponModel.
func (CouponModel) TableName() string {
	return "coupons"
}

// ToDomain converts a CouponModel to a domain Coupon.
func (m *CouponModel) ToDomain() *domain.Coupon {
	return &domain.Coupon{
		ID:               m.ID,
		Code:             m.Code,
		Type:             domain.CouponType(m.Type),
		DiscountValue:    m.DiscountValue,
		MinOrderCents:    m.MinOrderCents,
		MaxDiscountCents: m.MaxDiscountCents,
		UsageLimit:       m.UsageLimit,
		UsageCount:       m.UsageCount,
		PerUserLimit:     m.PerUserLimit,
		Scope:            domain.CouponScope(m.Scope),
		ScopeIDs:         []string(m.ScopeIDs),
		CreatedBy:        m.CreatedBy,
		StartsAt:         m.StartsAt,
		ExpiresAt:        m.ExpiresAt,
		IsActive:         m.IsActive,
		CreatedAt:        m.CreatedAt,
	}
}

// ToCouponModel converts a domain Coupon to a CouponModel.
func ToCouponModel(c *domain.Coupon) *CouponModel {
	return &CouponModel{
		ID:               c.ID,
		Code:             c.Code,
		Type:             string(c.Type),
		DiscountValue:    c.DiscountValue,
		MinOrderCents:    c.MinOrderCents,
		MaxDiscountCents: c.MaxDiscountCents,
		UsageLimit:       c.UsageLimit,
		UsageCount:       c.UsageCount,
		PerUserLimit:     c.PerUserLimit,
		Scope:            string(c.Scope),
		ScopeIDs:         pq.StringArray(c.ScopeIDs),
		CreatedBy:        c.CreatedBy,
		StartsAt:         c.StartsAt,
		ExpiresAt:        c.ExpiresAt,
		IsActive:         c.IsActive,
		CreatedAt:        c.CreatedAt,
	}
}

// CouponUsageModel is the GORM model for the coupon_usages table.
type CouponUsageModel struct {
	ID            string    `gorm:"type:uuid;primaryKey"`
	CouponID      string    `gorm:"type:uuid;index;not null"`
	UserID        string    `gorm:"type:uuid;index;not null"`
	OrderID       string    `gorm:"type:uuid;not null"`
	DiscountCents int64     `gorm:"not null;default:0"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

// TableName returns the table name for CouponUsageModel.
func (CouponUsageModel) TableName() string {
	return "coupon_usages"
}

// ToDomain converts a CouponUsageModel to a domain CouponUsage.
func (m *CouponUsageModel) ToDomain() *domain.CouponUsage {
	return &domain.CouponUsage{
		ID:            m.ID,
		CouponID:      m.CouponID,
		UserID:        m.UserID,
		OrderID:       m.OrderID,
		DiscountCents: m.DiscountCents,
		CreatedAt:     m.CreatedAt,
	}
}

// ToCouponUsageModel converts a domain CouponUsage to a CouponUsageModel.
func ToCouponUsageModel(u *domain.CouponUsage) *CouponUsageModel {
	return &CouponUsageModel{
		ID:            u.ID,
		CouponID:      u.CouponID,
		UserID:        u.UserID,
		OrderID:       u.OrderID,
		DiscountCents: u.DiscountCents,
		CreatedAt:     u.CreatedAt,
	}
}

// FlashSaleModel is the GORM model for the flash_sales table.
type FlashSaleModel struct {
	ID        string              `gorm:"type:uuid;primaryKey"`
	Name      string              `gorm:"type:varchar(255);not null"`
	StartsAt  time.Time           `gorm:"not null"`
	EndsAt    time.Time           `gorm:"not null"`
	IsActive  bool                `gorm:"not null;default:true"`
	Items     []FlashSaleItemModel `gorm:"foreignKey:FlashSaleID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time           `gorm:"autoCreateTime"`
}

// TableName returns the table name for FlashSaleModel.
func (FlashSaleModel) TableName() string {
	return "flash_sales"
}

// ToDomain converts a FlashSaleModel to a domain FlashSale.
func (m *FlashSaleModel) ToDomain() *domain.FlashSale {
	fs := &domain.FlashSale{
		ID:        m.ID,
		Name:      m.Name,
		StartsAt:  m.StartsAt,
		EndsAt:    m.EndsAt,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
	}
	for _, item := range m.Items {
		fs.Items = append(fs.Items, *item.ToDomain())
	}
	return fs
}

// ToFlashSaleModel converts a domain FlashSale to a FlashSaleModel.
func ToFlashSaleModel(fs *domain.FlashSale) *FlashSaleModel {
	model := &FlashSaleModel{
		ID:        fs.ID,
		Name:      fs.Name,
		StartsAt:  fs.StartsAt,
		EndsAt:    fs.EndsAt,
		IsActive:  fs.IsActive,
		CreatedAt: fs.CreatedAt,
	}
	for _, item := range fs.Items {
		model.Items = append(model.Items, *ToFlashSaleItemModel(&item))
	}
	return model
}

// FlashSaleItemModel is the GORM model for the flash_sale_items table.
type FlashSaleItemModel struct {
	ID             string `gorm:"type:uuid;primaryKey"`
	FlashSaleID    string `gorm:"type:uuid;index;not null"`
	ProductID      string `gorm:"type:uuid;not null"`
	VariantID      string `gorm:"type:uuid"`
	SalePriceCents int64  `gorm:"not null;default:0"`
	QuantityLimit  int    `gorm:"not null;default:0"`
	SoldCount      int    `gorm:"not null;default:0"`
}

// TableName returns the table name for FlashSaleItemModel.
func (FlashSaleItemModel) TableName() string {
	return "flash_sale_items"
}

// ToDomain converts a FlashSaleItemModel to a domain FlashSaleItem.
func (m *FlashSaleItemModel) ToDomain() *domain.FlashSaleItem {
	return &domain.FlashSaleItem{
		ID:             m.ID,
		FlashSaleID:    m.FlashSaleID,
		ProductID:      m.ProductID,
		VariantID:      m.VariantID,
		SalePriceCents: m.SalePriceCents,
		QuantityLimit:  m.QuantityLimit,
		SoldCount:      m.SoldCount,
	}
}

// ToFlashSaleItemModel converts a domain FlashSaleItem to a FlashSaleItemModel.
func ToFlashSaleItemModel(item *domain.FlashSaleItem) *FlashSaleItemModel {
	return &FlashSaleItemModel{
		ID:             item.ID,
		FlashSaleID:    item.FlashSaleID,
		ProductID:      item.ProductID,
		VariantID:      item.VariantID,
		SalePriceCents: item.SalePriceCents,
		QuantityLimit:  item.QuantityLimit,
		SoldCount:      item.SoldCount,
	}
}

// BundleModel is the GORM model for the bundles table.
type BundleModel struct {
	ID               string         `gorm:"type:uuid;primaryKey"`
	Name             string         `gorm:"type:varchar(255);not null"`
	SellerID         string         `gorm:"type:uuid;index;not null"`
	ProductIDs       pq.StringArray `gorm:"type:text[]"`
	BundlePriceCents int64          `gorm:"not null;default:0"`
	SavingsCents     int64          `gorm:"not null;default:0"`
	IsActive         bool           `gorm:"not null;default:true"`
	CreatedAt        time.Time      `gorm:"autoCreateTime"`
}

// TableName returns the table name for BundleModel.
func (BundleModel) TableName() string {
	return "bundles"
}

// ToDomain converts a BundleModel to a domain Bundle.
func (m *BundleModel) ToDomain() *domain.Bundle {
	return &domain.Bundle{
		ID:               m.ID,
		Name:             m.Name,
		SellerID:         m.SellerID,
		ProductIDs:       []string(m.ProductIDs),
		BundlePriceCents: m.BundlePriceCents,
		SavingsCents:     m.SavingsCents,
		IsActive:         m.IsActive,
		CreatedAt:        m.CreatedAt,
	}
}

// ToBundleModel converts a domain Bundle to a BundleModel.
func ToBundleModel(b *domain.Bundle) *BundleModel {
	return &BundleModel{
		ID:               b.ID,
		Name:             b.Name,
		SellerID:         b.SellerID,
		ProductIDs:       pq.StringArray(b.ProductIDs),
		BundlePriceCents: b.BundlePriceCents,
		SavingsCents:     b.SavingsCents,
		IsActive:         b.IsActive,
		CreatedAt:        b.CreatedAt,
	}
}
