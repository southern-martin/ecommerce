package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTaxOperations(t *testing.T) {
	requireServices(t, map[string]string{
		"auth": AuthBaseURL,
		"tax":  TaxBaseURL,
	})

	user := createTestUser(t)

	t.Run("Calculate_tax", func(t *testing.T) {
		calcBody := map[string]interface{}{
			"amount_cents": 10000,
			"country":      "US",
			"state":        "CA",
			"city":         "Los Angeles",
		}

		resp := httpPostWithUserID(t, TaxBaseURL+"/api/v1/tax/calculate", calcBody, user.AccessToken, user.UserID)
		defer resp.Body.Close()
		require.Contains(t, []int{http.StatusOK, http.StatusCreated}, resp.StatusCode, "failed to calculate tax")

		data := readBodyMap(t, resp)
		assertJSONField(t, data, "tax_cents")

		t.Logf("Tax calculated: %v cents", data["tax_cents"])
	})

	t.Run("Get_tax_rates", func(t *testing.T) {
		resp := httpGetWithUserID(t, TaxBaseURL+"/api/v1/tax/rates?country=US", user.AccessToken, user.UserID)
		defer resp.Body.Close()
		require.Contains(t, []int{http.StatusOK, http.StatusNotFound}, resp.StatusCode)

		t.Log("Retrieved tax rates")
	})

	t.Run("Create_tax_rule", func(t *testing.T) {
		ruleBody := map[string]interface{}{
			"country":  "US",
			"state":    "TX",
			"rate":     8.25,
			"category": "general",
		}

		resp := httpPostWithUserID(t, TaxBaseURL+"/api/v1/tax/rules", ruleBody, user.AccessToken, user.UserID)
		defer resp.Body.Close()
		require.Contains(t, []int{http.StatusOK, http.StatusCreated, http.StatusForbidden}, resp.StatusCode,
			"tax rule creation should succeed or be forbidden for non-admin")

		t.Logf("Create tax rule returned status: %d", resp.StatusCode)
	})
}
