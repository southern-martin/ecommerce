package domain

import "time"

// TaxZone represents a geographic tax jurisdiction.
type TaxZone struct {
	ID          string
	CountryCode string
	StateCode   string // optional
	Name        string // "Australia", "California"
}

// TaxRule represents a tax rule applied within a zone.
type TaxRule struct {
	ID        string
	ZoneID    string
	TaxName   string  // "GST", "VAT", "State Sales Tax"
	Rate      float64 // 0.10 = 10%
	Category  string  // product category override, empty = all
	Inclusive bool    // true = tax included in price (AU/EU)
	StartsAt  time.Time
	ExpiresAt *time.Time
	IsActive  bool
}

// TaxCalculationRequest holds the input for a tax calculation.
type TaxCalculationRequest struct {
	Items           []TaxItem
	ShippingAddress TaxAddress
}

// TaxItem represents a single item in a tax calculation request.
type TaxItem struct {
	ProductID  string
	VariantID  string
	Category   string
	PriceCents int64
	Quantity   int
}

// TaxAddress represents a shipping address for tax determination.
type TaxAddress struct {
	CountryCode string
	StateCode   string
	City        string
	PostalCode  string
}

// TaxCalculation holds the result of a tax calculation.
type TaxCalculation struct {
	SubtotalCents  int64
	TaxAmountCents int64
	Breakdown      []TaxBreakdown
}

// TaxBreakdown shows a single tax component in the calculation result.
type TaxBreakdown struct {
	TaxName      string
	Rate         float64
	AmountCents  int64
	Jurisdiction string
}
