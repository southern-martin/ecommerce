package money

import (
	"fmt"
	"math"
)

// Money represents a currency-safe monetary value stored as cents.
type Money struct {
	AmountCents int64  `json:"amount_cents"`
	Currency    string `json:"currency"`
}

// NewMoney creates a new Money value from cents and currency code.
func NewMoney(amountCents int64, currency string) Money {
	return Money{
		AmountCents: amountCents,
		Currency:    currency,
	}
}

// NewMoneyFromFloat creates a new Money value from a float amount (e.g., 19.99) and currency.
func NewMoneyFromFloat(amount float64, currency string) Money {
	cents := int64(math.Round(amount * 100))
	return Money{
		AmountCents: cents,
		Currency:    currency,
	}
}

// Add returns a new Money that is the sum of m and other.
// Panics if currencies do not match.
func (m Money) Add(other Money) Money {
	m.assertSameCurrency(other)
	return Money{
		AmountCents: m.AmountCents + other.AmountCents,
		Currency:    m.Currency,
	}
}

// Subtract returns a new Money that is the difference of m and other.
// Panics if currencies do not match.
func (m Money) Subtract(other Money) Money {
	m.assertSameCurrency(other)
	return Money{
		AmountCents: m.AmountCents - other.AmountCents,
		Currency:    m.Currency,
	}
}

// Multiply returns a new Money with the amount multiplied by the given factor.
// The result is rounded to the nearest cent.
func (m Money) Multiply(factor float64) Money {
	cents := int64(math.Round(float64(m.AmountCents) * factor))
	return Money{
		AmountCents: cents,
		Currency:    m.Currency,
	}
}

// Format returns a human-readable string representation of the money value.
// For example: "USD 19.99" or "EUR 100.00".
func (m Money) Format() string {
	whole := m.AmountCents / 100
	frac := m.AmountCents % 100
	if frac < 0 {
		frac = -frac
	}
	return fmt.Sprintf("%s %d.%02d", m.Currency, whole, frac)
}

// ToFloat returns the amount as a float64.
func (m Money) ToFloat() float64 {
	return float64(m.AmountCents) / 100.0
}

// IsZero returns true if the amount is zero.
func (m Money) IsZero() bool {
	return m.AmountCents == 0
}

// IsPositive returns true if the amount is greater than zero.
func (m Money) IsPositive() bool {
	return m.AmountCents > 0
}

// IsNegative returns true if the amount is less than zero.
func (m Money) IsNegative() bool {
	return m.AmountCents < 0
}

// Equals returns true if both Money values have the same amount and currency.
func (m Money) Equals(other Money) bool {
	return m.AmountCents == other.AmountCents && m.Currency == other.Currency
}

// assertSameCurrency panics if the currencies of the two Money values differ.
func (m Money) assertSameCurrency(other Money) {
	if m.Currency != other.Currency {
		panic(fmt.Sprintf("currency mismatch: %s vs %s", m.Currency, other.Currency))
	}
}
