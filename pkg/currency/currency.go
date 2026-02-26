package currency

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// Currency represents a supported currency with its ISO 4217 code,
// display symbol, and the number of decimal places.
type Currency struct {
	Code     string // ISO 4217 code: USD, EUR, GBP, JPY, CNY, etc.
	Symbol   string // Display symbol: $, EUR, GBP, JPY, CNY, etc.
	Decimals int    // Number of decimal places (2 for most, 0 for JPY).
}

// SupportedCurrencies maps ISO 4217 currency codes to their definitions.
var SupportedCurrencies = map[string]Currency{
	"USD": {Code: "USD", Symbol: "$", Decimals: 2},
	"EUR": {Code: "EUR", Symbol: "\u20ac", Decimals: 2},
	"GBP": {Code: "GBP", Symbol: "\u00a3", Decimals: 2},
	"JPY": {Code: "JPY", Symbol: "\u00a5", Decimals: 0},
	"CNY": {Code: "CNY", Symbol: "\u00a5", Decimals: 2},
	"CAD": {Code: "CAD", Symbol: "CA$", Decimals: 2},
	"AUD": {Code: "AUD", Symbol: "A$", Decimals: 2},
	"CHF": {Code: "CHF", Symbol: "CHF", Decimals: 2},
	"INR": {Code: "INR", Symbol: "\u20b9", Decimals: 2},
	"BRL": {Code: "BRL", Symbol: "R$", Decimals: 2},
	"KRW": {Code: "KRW", Symbol: "\u20a9", Decimals: 0},
	"MXN": {Code: "MXN", Symbol: "MX$", Decimals: 2},
	"SGD": {Code: "SGD", Symbol: "S$", Decimals: 2},
	"SEK": {Code: "SEK", Symbol: "kr", Decimals: 2},
	"AED": {Code: "AED", Symbol: "AED", Decimals: 2},
}

// Converter handles currency conversion with cached exchange rates.
// Rates are stored relative to USD (i.e., 1 USD = rate units of currency).
type Converter struct {
	rates    map[string]float64
	mu       sync.RWMutex
	cacheExp time.Time
}

// NewConverter creates a new Converter pre-loaded with default exchange
// rates. These defaults are approximate and should be updated with live
// rates via UpdateRates in production.
func NewConverter() *Converter {
	c := &Converter{
		rates: map[string]float64{
			"USD": 1.0,
			"EUR": 0.92,
			"GBP": 0.79,
			"JPY": 149.50,
			"CNY": 7.24,
			"CAD": 1.36,
			"AUD": 1.53,
			"CHF": 0.88,
			"INR": 83.12,
			"BRL": 4.97,
			"KRW": 1325.0,
			"MXN": 17.15,
			"SGD": 1.34,
			"SEK": 10.45,
			"AED": 3.67,
		},
		cacheExp: time.Now().Add(24 * time.Hour),
	}
	return c
}

// Convert converts an amount in cents from one currency to another.
// The amount is in the smallest unit of the source currency (e.g., cents
// for USD, whole yen for JPY). Returns the amount in the smallest unit
// of the target currency.
func (c *Converter) Convert(amountCents int64, from, to string) (int64, error) {
	if from == to {
		return amountCents, nil
	}

	rate, err := c.GetRate(from, to)
	if err != nil {
		return 0, err
	}

	fromCurrency, fromOk := SupportedCurrencies[from]
	toCurrency, toOk := SupportedCurrencies[to]
	if !fromOk || !toOk {
		return 0, fmt.Errorf("unsupported currency pair: %s -> %s", from, to)
	}

	// Convert from source smallest-unit to a base float amount.
	var baseAmount float64
	if fromCurrency.Decimals > 0 {
		baseAmount = float64(amountCents) / math.Pow(10, float64(fromCurrency.Decimals))
	} else {
		baseAmount = float64(amountCents)
	}

	// Apply the exchange rate.
	convertedAmount := baseAmount * rate

	// Convert to target smallest-unit.
	if toCurrency.Decimals > 0 {
		return int64(math.Round(convertedAmount * math.Pow(10, float64(toCurrency.Decimals)))), nil
	}
	return int64(math.Round(convertedAmount)), nil
}

// GetRate returns the exchange rate from one currency to another.
// The rate represents how many units of the target currency equal one
// unit of the source currency.
func (c *Converter) GetRate(from, to string) (float64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	fromRate, fromOk := c.rates[from]
	toRate, toOk := c.rates[to]
	if !fromOk {
		return 0, fmt.Errorf("unsupported source currency: %s", from)
	}
	if !toOk {
		return 0, fmt.Errorf("unsupported target currency: %s", to)
	}

	// Convert via USD: from -> USD -> to.
	// fromRate is "1 USD = fromRate units of from"
	// toRate is "1 USD = toRate units of to"
	// So from -> to = toRate / fromRate
	return toRate / fromRate, nil
}

// UpdateRates replaces the cached exchange rates with the provided rates.
// Rates must be relative to USD (i.e., 1 USD = rate units of currency).
// The cache expiration is reset to 1 hour from now.
func (c *Converter) UpdateRates(rates map[string]float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for code, rate := range rates {
		c.rates[code] = rate
	}
	// Ensure USD is always 1.
	c.rates["USD"] = 1.0
	c.cacheExp = time.Now().Add(1 * time.Hour)
}

// FormatPrice formats an amount (in the smallest currency unit) as a
// human-readable price string with the currency symbol.
//
// Examples:
//
//	FormatPrice(1999, "USD") => "$19.99"
//	FormatPrice(500, "JPY")  => "JPY500" (no decimals)
//	FormatPrice(1050, "EUR") => "EUR10.50"
func (c *Converter) FormatPrice(amountCents int64, currencyCode string) string {
	cur, ok := SupportedCurrencies[currencyCode]
	if !ok {
		// Fallback: treat as 2-decimal currency.
		whole := amountCents / 100
		frac := amountCents % 100
		if frac < 0 {
			frac = -frac
		}
		return fmt.Sprintf("%s%d.%02d", currencyCode, whole, frac)
	}

	if cur.Decimals == 0 {
		return fmt.Sprintf("%s%d", cur.Symbol, amountCents)
	}

	divisor := int64(math.Pow(10, float64(cur.Decimals)))
	whole := amountCents / divisor
	frac := amountCents % divisor
	if frac < 0 {
		frac = -frac
	}

	format := fmt.Sprintf("%%s%%d.%%0%dd", cur.Decimals)
	return fmt.Sprintf(format, cur.Symbol, whole, frac)
}
