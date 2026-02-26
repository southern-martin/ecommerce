package currency

import (
	"github.com/gin-gonic/gin"
)

const (
	// contextKeyCurrency is the gin context key used to store the resolved currency.
	contextKeyCurrency = "x_currency"

	// headerCurrency is the HTTP header used to specify the desired currency.
	headerCurrency = "X-Currency"

	// defaultCurrency is the fallback currency when none is specified.
	defaultCurrency = "USD"
)

// GinMiddleware returns a gin.HandlerFunc that extracts the X-Currency header
// from incoming requests and stores the validated currency code in the gin
// context. If the header is missing or contains an unsupported currency code,
// it defaults to USD.
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		currencyCode := c.GetHeader(headerCurrency)

		if currencyCode == "" {
			currencyCode = defaultCurrency
		}

		// Validate against supported currencies.
		if _, ok := SupportedCurrencies[currencyCode]; !ok {
			currencyCode = defaultCurrency
		}

		c.Set(contextKeyCurrency, currencyCode)
		c.Next()
	}
}

// GetCurrency retrieves the current currency code from the gin context.
// Returns "USD" if no currency has been set.
func GetCurrency(c *gin.Context) string {
	if val, exists := c.Get(contextKeyCurrency); exists {
		if s, ok := val.(string); ok && s != "" {
			return s
		}
	}
	return defaultCurrency
}
