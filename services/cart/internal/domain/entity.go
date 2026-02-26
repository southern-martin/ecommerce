package domain

import "time"

// Cart represents a user's shopping cart.
type Cart struct {
	UserID    string     `json:"user_id"`
	Items     []CartItem `json:"items"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// CartItem represents a single item in the cart.
type CartItem struct {
	ProductID   string `json:"product_id"`
	VariantID   string `json:"variant_id"`
	ProductName string `json:"product_name"`
	VariantName string `json:"variant_name"`
	SKU         string `json:"sku"`
	PriceCents  int64  `json:"price_cents"`
	Quantity    int    `json:"quantity"`
	ImageURL    string `json:"image_url"`
	SellerID    string `json:"seller_id"`
}

// SubtotalCents returns the total price of all items in the cart in cents.
func (c *Cart) SubtotalCents() int64 {
	var total int64
	for _, item := range c.Items {
		total += item.PriceCents * int64(item.Quantity)
	}
	return total
}

// TotalItems returns the total number of items (sum of quantities) in the cart.
func (c *Cart) TotalItems() int {
	var total int
	for _, item := range c.Items {
		total += item.Quantity
	}
	return total
}

// FindItem returns the index of an item matching productID and variantID, or -1 if not found.
func (c *Cart) FindItem(productID, variantID string) int {
	for i, item := range c.Items {
		if item.ProductID == productID && item.VariantID == variantID {
			return i
		}
	}
	return -1
}
