package postgres

import (
	"time"

	"github.com/southern-martin/ecommerce/services/tax/internal/domain"
)

// TaxZoneModel is the GORM model for tax_zones table.
type TaxZoneModel struct {
	ID          string `gorm:"type:uuid;primaryKey"`
	CountryCode string `gorm:"type:varchar(2);not null;index:idx_zone_location"`
	StateCode   string `gorm:"type:varchar(10);index:idx_zone_location"`
	Name        string `gorm:"type:varchar(100);not null"`
}

func (TaxZoneModel) TableName() string {
	return "tax_zones"
}

// ToDomain converts the GORM model to a domain entity.
func (m *TaxZoneModel) ToDomain() *domain.TaxZone {
	return &domain.TaxZone{
		ID:          m.ID,
		CountryCode: m.CountryCode,
		StateCode:   m.StateCode,
		Name:        m.Name,
	}
}

// TaxZoneModelFromDomain converts a domain entity to a GORM model.
func TaxZoneModelFromDomain(z *domain.TaxZone) *TaxZoneModel {
	return &TaxZoneModel{
		ID:          z.ID,
		CountryCode: z.CountryCode,
		StateCode:   z.StateCode,
		Name:        z.Name,
	}
}

// TaxRuleModel is the GORM model for tax_rules table.
type TaxRuleModel struct {
	ID        string     `gorm:"type:uuid;primaryKey"`
	ZoneID    string     `gorm:"type:uuid;not null;index:idx_rule_zone"`
	TaxName   string     `gorm:"type:varchar(50);not null"`
	Rate      float64    `gorm:"type:decimal(10,6);not null"`
	Category  string     `gorm:"type:varchar(100);index:idx_rule_zone_category"`
	Inclusive bool       `gorm:"type:boolean;default:false"`
	StartsAt  time.Time  `gorm:"type:timestamptz;not null"`
	ExpiresAt *time.Time `gorm:"type:timestamptz"`
	IsActive  bool       `gorm:"type:boolean;default:true;index:idx_rule_active"`
}

func (TaxRuleModel) TableName() string {
	return "tax_rules"
}

// ToDomain converts the GORM model to a domain entity.
func (m *TaxRuleModel) ToDomain() *domain.TaxRule {
	return &domain.TaxRule{
		ID:        m.ID,
		ZoneID:    m.ZoneID,
		TaxName:   m.TaxName,
		Rate:      m.Rate,
		Category:  m.Category,
		Inclusive:  m.Inclusive,
		StartsAt:  m.StartsAt,
		ExpiresAt: m.ExpiresAt,
		IsActive:  m.IsActive,
	}
}

// TaxRuleModelFromDomain converts a domain entity to a GORM model.
func TaxRuleModelFromDomain(r *domain.TaxRule) *TaxRuleModel {
	return &TaxRuleModel{
		ID:        r.ID,
		ZoneID:    r.ZoneID,
		TaxName:   r.TaxName,
		Rate:      r.Rate,
		Category:  r.Category,
		Inclusive:  r.Inclusive,
		StartsAt:  r.StartsAt,
		ExpiresAt: r.ExpiresAt,
		IsActive:  r.IsActive,
	}
}
