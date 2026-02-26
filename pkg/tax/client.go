package tax

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TaxRate represents a tax rate returned by the tax service.
type TaxRate struct {
	Country  string  `json:"country"`
	State    string  `json:"state"`
	Rate     float64 `json:"rate"`
	Category string  `json:"category"`
}

// TaxCalculation represents the result of a tax calculation.
type TaxCalculation struct {
	SubtotalCents int64   `json:"subtotal_cents"`
	TaxCents      int64   `json:"tax_cents"`
	TotalCents    int64   `json:"total_cents"`
	TaxRate       float64 `json:"tax_rate"`
}

// Client is a gRPC client for the tax service.
type Client struct {
	conn *grpc.ClientConn
	addr string
}

// NewClient creates a new tax service gRPC client.
func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tax service at %s: %w", addr, err)
	}

	return &Client{
		conn: conn,
		addr: addr,
	}, nil
}

// Close closes the gRPC connection.
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// CalculateTax calculates the tax for a given subtotal, country, and state.
func (c *Client) CalculateTax(_ context.Context, subtotalCents int64, country, state, category string) (*TaxCalculation, error) {
	rate := lookupDefaultRate(country, state, category)

	taxCents := int64(float64(subtotalCents) * rate)
	return &TaxCalculation{
		SubtotalCents: subtotalCents,
		TaxCents:      taxCents,
		TotalCents:    subtotalCents + taxCents,
		TaxRate:       rate,
	}, nil
}

// GetTaxRate returns the tax rate for a given country, state, and category.
func (c *Client) GetTaxRate(_ context.Context, country, state, category string) (*TaxRate, error) {
	rate := lookupDefaultRate(country, state, category)
	return &TaxRate{
		Country:  country,
		State:    state,
		Rate:     rate,
		Category: category,
	}, nil
}

// lookupDefaultRate returns a default tax rate based on country.
// This is a placeholder until the actual gRPC tax service is implemented.
func lookupDefaultRate(country, _, _ string) float64 {
	rates := map[string]float64{
		"US": 0.08,
		"GB": 0.20,
		"DE": 0.19,
		"FR": 0.20,
		"JP": 0.10,
		"CA": 0.13,
		"AU": 0.10,
		"IN": 0.18,
	}

	if rate, ok := rates[country]; ok {
		return rate
	}
	return 0.10
}
