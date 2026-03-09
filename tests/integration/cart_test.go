package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCartOperations(t *testing.T) {
	requireServices(t, map[string]string{
		"auth":    AuthBaseURL,
		"product": ProductBaseURL,
		"cart":    CartBaseURL,
	})

	buyer := createTestUser(t)
	seller := createTestUser(t)
	product := createTestProduct(t, seller)

	var addedProductID string

	t.Run("Add_item_to_cart", func(t *testing.T) {
		cartItem := map[string]interface{}{
			"product_id":   product.ID,
			"product_name": product.Name,
			"variant_id":   "",
			"quantity":     2,
			"price_cents":  product.BasePriceCents,
			"seller_id":    seller.UserID,
		}

		resp := httpPostWithUserID(t, CartBaseURL+"/api/v1/cart/items", cartItem, buyer.AccessToken, buyer.UserID)
		defer resp.Body.Close()
		require.Contains(t, []int{http.StatusOK, http.StatusCreated}, resp.StatusCode, "failed to add item to cart")

		addedProductID = product.ID
		t.Logf("Added product %s to cart", addedProductID)
	})

	t.Run("Get_cart", func(t *testing.T) {
		resp := httpGetWithUserID(t, CartBaseURL+"/api/v1/cart", buyer.AccessToken, buyer.UserID)
		defer resp.Body.Close()
		assertStatus(t, resp, http.StatusOK)

		cartData := readBodyMap(t, resp)
		items, ok := cartData["items"].([]interface{})
		require.True(t, ok, "cart items should be an array")
		assert.GreaterOrEqual(t, len(items), 1, "cart should have at least one item")

		t.Logf("Cart has %d item(s)", len(items))
	})

	t.Run("Update_item_quantity", func(t *testing.T) {
		updateBody := map[string]interface{}{
			"product_id": addedProductID,
			"quantity":   5,
		}

		resp := httpPatch(t, CartBaseURL+"/api/v1/cart/items", updateBody, buyer.AccessToken, buyer.UserID)
		defer resp.Body.Close()
		require.Contains(t, []int{http.StatusOK, http.StatusNoContent}, resp.StatusCode, "failed to update cart item quantity")

		t.Logf("Updated product %s quantity to 5", addedProductID)
	})

	t.Run("Remove_item", func(t *testing.T) {
		removeBody := map[string]interface{}{
			"product_id": addedProductID,
		}

		resp := httpDeleteWithBody(t, CartBaseURL+"/api/v1/cart/items", removeBody, buyer.AccessToken, buyer.UserID)
		defer resp.Body.Close()
		require.Contains(t, []int{http.StatusOK, http.StatusNoContent}, resp.StatusCode, "failed to remove cart item")

		t.Logf("Removed product %s from cart", addedProductID)
	})

	t.Run("Clear_cart", func(t *testing.T) {
		// Add an item back so we have something to clear
		cartItem := map[string]interface{}{
			"product_id":   product.ID,
			"product_name": product.Name,
			"quantity":     1,
			"price_cents":  product.BasePriceCents,
			"seller_id":    seller.UserID,
		}

		addResp := httpPostWithUserID(t, CartBaseURL+"/api/v1/cart/items", cartItem, buyer.AccessToken, buyer.UserID)
		defer addResp.Body.Close()
		require.Contains(t, []int{http.StatusOK, http.StatusCreated}, addResp.StatusCode)

		// Clear the entire cart
		resp := httpDelete(t, CartBaseURL+"/api/v1/cart", buyer.AccessToken, buyer.UserID)
		defer resp.Body.Close()
		require.Contains(t, []int{http.StatusOK, http.StatusNoContent}, resp.StatusCode, "failed to clear cart")

		// Verify cart is empty
		getResp := httpGetWithUserID(t, CartBaseURL+"/api/v1/cart", buyer.AccessToken, buyer.UserID)
		defer getResp.Body.Close()
		assertStatus(t, getResp, http.StatusOK)

		cartData := readBodyMap(t, getResp)
		if items, ok := cartData["items"].([]interface{}); ok {
			assert.Equal(t, 0, len(items), "cart should be empty after clearing")
		}

		_ = fmt.Sprintf("cart cleared")
		t.Log("Cart cleared successfully")
	})
}
