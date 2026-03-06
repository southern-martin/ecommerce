package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubtotalCents_EmptyCart(t *testing.T) {
	cart := &Cart{UserID: "user-1", Items: []CartItem{}}
	assert.Equal(t, int64(0), cart.SubtotalCents())
}

func TestSubtotalCents_SingleItem(t *testing.T) {
	cart := &Cart{
		UserID: "user-1",
		Items: []CartItem{
			{ProductID: "p1", PriceCents: 1999, Quantity: 2},
		},
	}
	assert.Equal(t, int64(3998), cart.SubtotalCents())
}

func TestSubtotalCents_MultipleItems(t *testing.T) {
	cart := &Cart{
		UserID: "user-1",
		Items: []CartItem{
			{ProductID: "p1", PriceCents: 1000, Quantity: 3},
			{ProductID: "p2", PriceCents: 2500, Quantity: 1},
			{ProductID: "p3", PriceCents: 500, Quantity: 10},
		},
	}
	// 3000 + 2500 + 5000 = 10500
	assert.Equal(t, int64(10500), cart.SubtotalCents())
}

func TestTotalItems_EmptyCart(t *testing.T) {
	cart := &Cart{UserID: "user-1", Items: []CartItem{}}
	assert.Equal(t, 0, cart.TotalItems())
}

func TestTotalItems_MultipleItems(t *testing.T) {
	cart := &Cart{
		UserID: "user-1",
		Items: []CartItem{
			{ProductID: "p1", Quantity: 3},
			{ProductID: "p2", Quantity: 1},
			{ProductID: "p3", Quantity: 5},
		},
	}
	assert.Equal(t, 9, cart.TotalItems())
}

func TestFindItem_Found(t *testing.T) {
	cart := &Cart{
		UserID: "user-1",
		Items: []CartItem{
			{ProductID: "p1", VariantID: "v1"},
			{ProductID: "p2", VariantID: "v2"},
			{ProductID: "p3", VariantID: "v3"},
		},
	}
	assert.Equal(t, 1, cart.FindItem("p2", "v2"))
}

func TestFindItem_NotFound(t *testing.T) {
	cart := &Cart{
		UserID: "user-1",
		Items: []CartItem{
			{ProductID: "p1", VariantID: "v1"},
		},
	}
	assert.Equal(t, -1, cart.FindItem("p99", "v99"))
}

func TestFindItem_MatchesBothIDs(t *testing.T) {
	cart := &Cart{
		UserID: "user-1",
		Items: []CartItem{
			{ProductID: "p1", VariantID: "v1"},
			{ProductID: "p1", VariantID: "v2"},
		},
	}
	// Same productID but different variantID should not match
	assert.Equal(t, -1, cart.FindItem("p1", "v99"))
	// Exact match on second item
	assert.Equal(t, 1, cart.FindItem("p1", "v2"))
}
