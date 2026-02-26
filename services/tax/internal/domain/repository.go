package domain

import "context"

// TaxZoneRepository defines persistence operations for TaxZone.
type TaxZoneRepository interface {
	Create(ctx context.Context, zone *TaxZone) error
	GetByID(ctx context.Context, id string) (*TaxZone, error)
	GetByLocation(ctx context.Context, countryCode, stateCode string) (*TaxZone, error)
	List(ctx context.Context) ([]*TaxZone, error)
}

// TaxRuleRepository defines persistence operations for TaxRule.
type TaxRuleRepository interface {
	Create(ctx context.Context, rule *TaxRule) error
	GetByID(ctx context.Context, id string) (*TaxRule, error)
	ListByZone(ctx context.Context, zoneID string) ([]*TaxRule, error)
	ListActive(ctx context.Context) ([]*TaxRule, error)
	Update(ctx context.Context, rule *TaxRule) error
	Delete(ctx context.Context, id string) error
	GetByZoneAndCategory(ctx context.Context, zoneID, category string) ([]*TaxRule, error)
}
